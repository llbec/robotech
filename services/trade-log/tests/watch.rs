use trade_log::watch::{AddAccount, BudgetConfig, MonitoringConfig};
#[test]
fn account_validation_and_resource_limits() {
    let mut a = AddAccount {
        account: "0xABCDEF0000000000000000000000000000000000".into(),
        start_time: "2026-01-01T00:00:00Z".into(),
    };
    a.validate().unwrap();
    assert_eq!(a.account, "0xabcdef0000000000000000000000000000000000");
    a.start_time = "2999-01-01T00:00:00Z".into();
    assert_eq!(a.validate().unwrap_err().code, "INVALID_ARGUMENT");
    assert!(
        serde_json::from_str::<AddAccount>(r#"{"account":"x","start_time":"x","enabled":true}"#)
            .is_err()
    );
    assert!(MonitoringConfig::default().validate().is_ok());
    assert!(BudgetConfig::default().validate().is_ok());
    for max in [0, 11] {
        assert!(
            MonitoringConfig {
                max_active_accounts: max,
                ..Default::default()
            }
            .validate()
            .is_err()
        )
    }
    assert!(
        BudgetConfig {
            http_weight_per_minute: 119,
            ..Default::default()
        }
        .validate()
        .is_err()
    );
    assert!(
        BudgetConfig {
            websocket_connections_per_minute: 31,
            ..Default::default()
        }
        .validate()
        .is_err()
    );
}
#[test]
fn conservative_source_weights() {
    use protocol_api::QueryKind;
    use trade_log::acquisition::budget::*;
    assert_eq!(reservation_weight(QueryKind::UserFillsByTime), 120);
    assert_eq!(response_weight(QueryKind::UserFills, 2000), 120);
    assert_eq!(response_weight(QueryKind::Meta, 2000), 20);
    assert_eq!(response_weight(QueryKind::UserFills, 0), 20);
}
