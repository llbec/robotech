use anyhow::{bail, Context, Result};
use chrono::Utc;
use serde::Deserialize;
use serde_json::Value;
use sqlx::PgPool;
use std::collections::{HashMap, HashSet};
use tokio::sync::RwLock;

#[derive(Clone, Debug)]
pub struct MarketInfo {
    pub market: String,
    pub base_asset: String,
    pub quote_asset: String,
}

#[derive(Default, Debug)]
pub struct MarketCatalog {
    perps: HashSet<String>,
    spots: HashMap<String, MarketInfo>,
}

impl MarketCatalog {
    pub fn resolve(&self, coin: &str) -> Result<(String, MarketInfo)> {
        if let Some(info) = self.spots.get(coin) {
            return Ok(("SPOT".into(), info.clone()));
        }
        if self.perps.contains(coin) || coin.contains(':') {
            return Ok((
                "PERPETUAL".into(),
                MarketInfo {
                    market: coin.into(),
                    base_asset: coin.split(':').next_back().unwrap_or(coin).into(),
                    quote_asset: "USDC".into(),
                },
            ));
        }
        bail!("unknown Hyperliquid market: {coin}")
    }

    fn from_values(perp: &Value, spot: &Value) -> Result<Self> {
        let perp: PerpMeta = serde_json::from_value(perp.clone()).context("invalid perp meta")?;
        let spot: SpotMeta = serde_json::from_value(spot.clone()).context("invalid spot meta")?;
        let token_names: HashMap<u32, String> = spot
            .tokens
            .into_iter()
            .map(|token| (token.index, token.name))
            .collect();
        let mut spots = HashMap::new();
        for market in spot.universe {
            if market.tokens.len() != 2 {
                continue;
            }
            let Some(base) = token_names.get(&market.tokens[0]).cloned() else {
                continue;
            };
            let Some(quote) = token_names.get(&market.tokens[1]).cloned() else {
                continue;
            };
            let info = MarketInfo {
                market: market.name.clone(),
                base_asset: base,
                quote_asset: quote,
            };
            spots.insert(format!("@{}", market.index), info.clone());
            spots.insert(market.name, info);
        }
        Ok(Self {
            perps: perp
                .universe
                .into_iter()
                .map(|market| market.name)
                .collect(),
            spots,
        })
    }
}

#[derive(Deserialize)]
struct PerpMeta {
    universe: Vec<PerpMarket>,
}

#[derive(Deserialize)]
struct PerpMarket {
    name: String,
}

#[derive(Deserialize)]
struct SpotMeta {
    tokens: Vec<SpotToken>,
    universe: Vec<SpotMarket>,
}

#[derive(Deserialize)]
struct SpotToken {
    name: String,
    index: u32,
}

#[derive(Deserialize)]
struct SpotMarket {
    name: String,
    tokens: Vec<u32>,
    index: u32,
}

pub struct MetadataStore {
    catalog: RwLock<MarketCatalog>,
}

impl MetadataStore {
    pub fn empty() -> Self {
        Self {
            catalog: RwLock::new(MarketCatalog::default()),
        }
    }

    pub async fn replace_and_persist(&self, pool: &PgPool, perp: Value, spot: Value) -> Result<()> {
        let catalog = MarketCatalog::from_values(&perp, &spot)?;
        sqlx::query(
            "INSERT INTO market_metadata_versions (observed_at, perp_payload, spot_payload) VALUES ($1,$2,$3)",
        )
        .bind(Utc::now())
        .bind(perp)
        .bind(spot)
        .execute(pool)
        .await?;
        *self.catalog.write().await = catalog;
        Ok(())
    }

    pub async fn load_latest(&self, pool: &PgPool) -> Result<bool> {
        let row: Option<(Value, Value)> = sqlx::query_as(
            "SELECT perp_payload,spot_payload FROM market_metadata_versions ORDER BY observed_at DESC,id DESC LIMIT 1",
        )
        .fetch_optional(pool)
        .await?;
        let Some((perp, spot)) = row else {
            return Ok(false);
        };
        *self.catalog.write().await = MarketCatalog::from_values(&perp, &spot)?;
        Ok(true)
    }

    pub async fn resolve(&self, coin: &str) -> Result<(String, MarketInfo)> {
        self.catalog.read().await.resolve(coin)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json::json;

    #[test]
    fn resolves_perp_and_spot_by_versioned_metadata() {
        let catalog = MarketCatalog::from_values(
            &json!({"universe":[{"name":"BTC"}]}),
            &json!({"tokens":[{"name":"USDC","index":0},{"name":"HYPE","index":150}],"universe":[{"name":"HYPE/USDC","tokens":[150,0],"index":107}]}),
        )
        .unwrap();
        assert_eq!(catalog.resolve("BTC").unwrap().0, "PERPETUAL");
        let (_, spot) = catalog.resolve("@107").unwrap();
        assert_eq!(spot.base_asset, "HYPE");
        assert_eq!(spot.quote_asset, "USDC");
    }
}
