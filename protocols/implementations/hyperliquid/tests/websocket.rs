use futures_util::{SinkExt, StreamExt};
use hyperliquid::{
    parser::HyperliquidParser,
    websocket::{WebsocketSource, decode, subscription},
};
use serde_json::{Value, json};
use trade_log::{
    collection::content::semantic_hash,
    parsing::{ParseContext, ProtocolParser},
    realtime::{RealtimeSource, StreamEvent},
};
const ACCOUNT: &str = "0x0000000000000000000000000000000000000001";
#[test]
fn envelope_modes_account_and_control_messages() {
    for (bytes, expected) in [
        (
            include_bytes!("../../../../tests/fixtures/v0.4/snapshot.json").as_slice(),
            "SNAPSHOT",
        ),
        (
            include_bytes!("../../../../tests/fixtures/v0.4/update.json").as_slice(),
            "LIVE_UPDATE",
        ),
        (
            include_bytes!("../../../../tests/fixtures/v0.4/unknown.json").as_slice(),
            "UNKNOWN",
        ),
    ] {
        let StreamEvent::Data(m) = decode(bytes, ACCOUNT).unwrap() else {
            panic!()
        };
        assert_eq!(m.mode, expected);
        assert_eq!(m.body, bytes);
        assert!(
            serde_json::from_slice::<Value>(&m.fills)
                .unwrap()
                .is_array()
        );
    }
    assert!(matches!(
        decode(
            &serde_json::to_vec(
                &json!({"channel":"subscriptionResponse","data":subscription(ACCOUNT)})
            )
            .unwrap(),
            ACCOUNT
        )
        .unwrap(),
        StreamEvent::Subscribed
    ));
    assert!(
        decode(
            include_bytes!("../../../../tests/fixtures/v0.4/malformed.json"),
            ACCOUNT
        )
        .is_err()
    );
    assert!(
        decode(
            include_bytes!("../../../../tests/fixtures/v0.4/update.json"),
            "0x0000000000000000000000000000000000000002"
        )
        .is_err()
    );
    assert!(matches!(
        decode(b"{\"channel\":\"pong\"}", ACCOUNT).unwrap(),
        StreamEvent::Pong
    ));
    assert!(matches!(
        decode(b"{\"channel\":\"other\"}", ACCOUNT).unwrap(),
        StreamEvent::Ignore
    ));
}
#[test]
fn semantic_comparison_preserves_business_conflicts_and_nulls() {
    let original: Value = serde_json::from_str(include_str!(
        "../../../../tests/fixtures/hyperliquid/fills.json"
    ))
    .unwrap();
    let parse = |fills: Value| {
        HyperliquidParser
            .parse(ParseContext {
                network: &shared_types::Network::Mainnet,
                account: ACCOUNT,
                raw_log_id: "raw_test",
                fills: &serde_json::to_vec(&fills).unwrap(),
                meta: include_bytes!("../../../../tests/fixtures/hyperliquid/meta.json"),
                spot_meta: include_bytes!("../../../../tests/fixtures/hyperliquid/spot-meta.json"),
            })
            .unwrap()
    };
    let before = parse(json!([original[0]]));
    let mut optional = original[0].clone();
    optional["irrelevantOptionalField"] = json!(true);
    optional["px"] = json!(format!("{}0", original[0]["px"].as_str().unwrap()));
    let after = parse(json!([optional]));
    assert_eq!(
        semantic_hash(&before.trades[0]).unwrap(),
        semantic_hash(&after.trades[0]).unwrap()
    );
    assert_eq!(parse(json!([original[0], optional])).duplicate_records, 1);
    let mut changed = before.trades[0].clone();
    changed.payload.fee = "999".into();
    assert_ne!(
        semantic_hash(&before.trades[0]).unwrap(),
        semantic_hash(&changed).unwrap()
    );
    changed = before.trades[0].clone();
    changed.payload.reported_realized_pnl = None;
    let no_value = semantic_hash(&changed).unwrap();
    changed.payload.reported_realized_pnl = Some("0".into());
    assert_ne!(no_value, semantic_hash(&changed).unwrap());
}
#[tokio::test]
async fn actual_socket_subscription_and_application_heartbeat() {
    let listener = tokio::net::TcpListener::bind("127.0.0.1:0").await.unwrap();
    let address = listener.local_addr().unwrap();
    let server = tokio::spawn(async move {
        let (stream, _) = listener.accept().await.unwrap();
        let mut ws = tokio_tungstenite::accept_async(stream).await.unwrap();
        let request = ws.next().await.unwrap().unwrap();
        assert_eq!(
            serde_json::from_str::<Value>(request.to_text().unwrap()).unwrap(),
            subscription(ACCOUNT)
        );
        ws.send(tokio_tungstenite::tungstenite::Message::Text(
            json!({"channel":"subscriptionResponse","data":subscription(ACCOUNT)})
                .to_string()
                .into(),
        ))
        .await
        .unwrap();
        let ping = ws.next().await.unwrap().unwrap();
        assert_eq!(ping.to_text().unwrap(), "{\"method\":\"ping\"}");
        ws.send(tokio_tungstenite::tungstenite::Message::Text(
            "{\"channel\":\"pong\"}".into(),
        ))
        .await
        .unwrap();
    });
    let mut connection = WebsocketSource::with_endpoint(format!("ws://{address}"))
        .connect(ACCOUNT, 1024)
        .await
        .unwrap();
    assert!(matches!(
        connection.next().await.unwrap(),
        StreamEvent::Subscribed
    ));
    connection.ping().await.unwrap();
    assert!(matches!(
        connection.next().await.unwrap(),
        StreamEvent::Pong
    ));
    server.await.unwrap();
}

#[tokio::test]
async fn actual_socket_rejects_oversize_and_archives_malformed_events() {
    for (text, max_bytes, oversize) in [
        ("x".repeat(2048), 1024, true),
        (
            include_str!("../../../../tests/fixtures/v0.4/malformed.json").to_owned(),
            1024,
            false,
        ),
    ] {
        let listener = tokio::net::TcpListener::bind("127.0.0.1:0").await.unwrap();
        let address = listener.local_addr().unwrap();
        let server = tokio::spawn(async move {
            let (stream, _) = listener.accept().await.unwrap();
            let mut socket = tokio_tungstenite::accept_async(stream).await.unwrap();
            let _ = socket.next().await;
            socket
                .send(tokio_tungstenite::tungstenite::Message::Text(text.into()))
                .await
                .unwrap();
            let _ = socket.next().await;
        });
        let mut connection = WebsocketSource::with_endpoint(format!("ws://{address}"))
            .connect(ACCOUNT, max_bytes)
            .await
            .unwrap();
        if oversize {
            assert_eq!(
                connection.next().await.err().unwrap().code,
                "MESSAGE_TOO_LARGE"
            );
        } else {
            assert!(matches!(
                connection.next().await.unwrap(),
                StreamEvent::Rejected(_, _)
            ));
        }
        drop(connection);
        server.await.unwrap();
    }
}
#[tokio::test]
#[ignore = "explicit official Hyperliquid WSS acceptance; requires internet access"]
async fn official_wss_subscribes_and_answers_application_ping() {
    let source = WebsocketSource::new(&shared_types::Network::Mainnet);
    let mut connection = tokio::time::timeout(
        std::time::Duration::from_secs(15),
        source.connect(ACCOUNT, 16777216),
    )
    .await
    .unwrap()
    .unwrap();
    tokio::time::timeout(std::time::Duration::from_secs(15), async {
        loop {
            if matches!(connection.next().await.unwrap(), StreamEvent::Subscribed) {
                break;
            }
        }
    })
    .await
    .unwrap();
    connection.ping().await.unwrap();
    tokio::time::timeout(std::time::Duration::from_secs(15), async {
        loop {
            if matches!(connection.next().await.unwrap(), StreamEvent::Pong) {
                break;
            }
        }
    })
    .await
    .unwrap();
    connection.close().await;
}
