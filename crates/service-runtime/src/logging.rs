use crate::config::LoggingConfig;

pub fn init(config: &LoggingConfig) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let level: tracing::Level = config.level.parse()?;
    let subscriber = tracing_subscriber::fmt()
        .with_max_level(level)
        .with_writer(std::io::stderr)
        .with_ansi(false)
        .with_target(false);
    if config.format == "json" {
        subscriber.json().try_init()?;
    } else {
        subscriber.try_init()?;
    }
    Ok(())
}
