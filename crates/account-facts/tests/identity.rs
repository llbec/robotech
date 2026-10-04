use account_facts::fact_id;
use shared_types::Network;
#[test]
fn identity_fixed_vector_and_namespace_isolation() {
    let address = "0x0000000000000000000000000000000000000001";
    let id = fact_id(&Network::Mainnet, address, "ETH", 9007199254740993);
    assert_eq!(
        id,
        "hl_fill_v1_7ad47a033ec0a3acbb533af39503356bd718c37d77243de259a5577794929c7c"
    );
    assert_ne!(
        id,
        fact_id(&Network::Testnet, address, "ETH", 9007199254740993)
    );
    assert_ne!(
        id,
        fact_id(&Network::Mainnet, address, "BTC", 9007199254740993)
    );
}
