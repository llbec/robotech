use trade_log::realtime::{ReconnectSchedule, WebsocketConfig};
#[test]
fn reconnect_does_not_reset_on_a_brief_success_and_limits_flapping() {
    let c = WebsocketConfig::default();
    let mut s = ReconnectSchedule::default();
    for (i, expected) in [5000, 10000, 20000, 40000, 60000, 60000]
        .into_iter()
        .enumerate()
    {
        s.stable(59, true, &c);
        assert_eq!(s.wait_ms(i as u64 * 1000, &c, 0), expected);
    }
    s.stable(60, false, &c);
    assert_eq!(s.wait_ms(0, &c, 0), 60000);
    s.stable(60, true, &c);
    assert_eq!(s.wait_ms(0, &c, 0), 5000);
    for time in 0..10 {
        s.attempted(time * 1000);
    }
    assert_eq!(s.rate_wait_ms(10000), 50000);
    assert_eq!(s.rate_wait_ms(60000), 0);
}
#[test]
fn jitter_is_bounded_and_legacy_configuration_is_http_only() {
    let c: WebsocketConfig = serde_json::from_str("{}").unwrap();
    assert!(!c.enabled);
    assert!(c.validate(180).is_ok());
    for random in [0, 1, 1000, u64::MAX] {
        let mut s = ReconnectSchedule::default();
        for base in [5000, 10000, 20000, 40000, 60000] {
            let actual = s.wait_ms(0, &c, random);
            assert!((base..=(base + base / 5).min(60000)).contains(&actual));
        }
    }
    for invalid in [
        serde_json::json!({"reconnect_base_seconds":0}),
        serde_json::json!({"reconnect_max_seconds":61}),
        serde_json::json!({"reconnect_jitter_percent":21}),
        serde_json::json!({"reconnect_reset_after_seconds":59}),
        serde_json::json!({"ping_interval_seconds":50}),
        serde_json::json!({"max_pending_bytes":1}),
    ] {
        assert!(
            serde_json::from_value::<WebsocketConfig>(invalid)
                .unwrap()
                .validate(180)
                .is_err()
        );
    }
}

#[test]
fn rebuilding_a_worker_preserves_pending_backoff_and_the_connection_window() {
    let c = WebsocketConfig::default();
    let mut schedule = ReconnectSchedule::default();
    schedule.attempted(1000);
    assert_eq!(schedule.wait_ms(1000, &c, 0), 5000);
    assert_eq!(schedule.pending_wait_ms(2000), 4000);
    assert_eq!(schedule.pending_wait_ms(6000), 0);
    for at in 2..=10 {
        schedule.attempted(at * 1000);
    }
    assert_eq!(schedule.pending_wait_ms(11000), 50000);
    schedule.stable(60, true, &c);
    assert_eq!(schedule.pending_wait_ms(11000), 50000);
}
