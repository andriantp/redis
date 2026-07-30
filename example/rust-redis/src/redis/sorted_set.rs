use anyhow::Result;
use redis::AsyncCommands;

use super::redis::Repository;

impl Repository {
    pub async fn zadd(&self,key: &str,score: f64,member: &str) -> Result<()> {
        let mut conn = self.client.get_multiplexed_async_connection().await?;
        conn.zadd::<_, _, _, ()>(
            key,
            member,
            score,
        )
        .await?;
        Ok(())
    }

    pub async fn zrevrange(&self,key: &str,start: isize,stop: isize) -> Result<Vec<String>> {
        let mut conn = self.client.get_multiplexed_async_connection().await?;
        let ranking: Vec<String> = conn
            .zrevrange(key, start, stop)
            .await?;

        Ok(ranking)
    }
}

