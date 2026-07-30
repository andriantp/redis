
use anyhow::{anyhow, Result};
use redis::AsyncCommands;

use super::redis::Repository;

impl Repository {
    pub async fn get(&self, key: &str) -> Result<String> {
        let mut conn = self.client.get_multiplexed_async_connection().await?;

        let result: Option<String> = conn.get(key).await?;

        match result {
            Some(value) => Ok(value),
            None => Err(anyhow!("key not found")),
        }
    }

    pub async fn set(&self,key: &str, value: &str) -> Result<()> {
        let mut conn = self.client.get_multiplexed_async_connection().await?;
        conn.set::<_, _, ()>(key, value).await?;
        Ok(())
    }

    pub async fn setex(&self,key: &str,value: &str,expiration: u64,) -> Result<()> {
        let mut conn = self.client.get_multiplexed_async_connection().await?;
        conn.set_ex::<_, _, ()>(
            key,
            value,
            expiration,
        )
        .await?;

        Ok(())
    }

    pub async fn ttl(&self, key: &str) -> Result<i64> {
        let mut conn = self.client.get_multiplexed_async_connection().await?;
        let ttl: i64 = conn.ttl(key).await?;
        Ok(ttl)
    }
}

