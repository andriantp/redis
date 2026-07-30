use std::collections::HashMap;

use anyhow::{anyhow, Result};
use redis::AsyncCommands;

use super::redis::Repository;

impl Repository {
    pub async fn hset(&self,key: &str,values: &HashMap<String, String>) -> Result<()> {
        let mut conn = self.client.get_multiplexed_async_connection().await?;

        conn.hset_multiple::<_, _, _, ()>(
            key,
            &values
                .iter()
                .map(|(k, v)| (k.as_str(), v.as_str()))
                .collect::<Vec<_>>(),
        )
        .await?;

        Ok(())
    }

    pub async fn hget(&self,key: &str, field: &str) -> Result<String> {
        let mut conn = self.client.get_multiplexed_async_connection().await?;
        let result: Option<String> = conn.hget(key, field).await?;
        match result {
            Some(value) => Ok(value),
            None => Err(anyhow!("field not found")),
        }
    }

    pub async fn hgetall(&self,key: &str) -> Result<HashMap<String, String>> {
        let mut conn = self.client.get_multiplexed_async_connection().await?;
        let result: HashMap<String, String> =
            conn.hgetall(key).await?;
        Ok(result)
    }
}

