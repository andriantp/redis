use anyhow::Result;
use futures_util::StreamExt;
use redis::AsyncCommands;

use super::redis::Repository;

impl Repository {
    pub async fn publish(
        &self,
        channel: &str,
        message: &str,
    ) -> Result<()> {
        let mut conn = self.client.get_multiplexed_async_connection().await?;

        conn.publish::<_, _, ()>(
            channel,
            message,
        )
        .await?;

        Ok(())
    }

    pub async fn subscribe(
        &self,
        channel: &str,
    ) -> Result<()> {
        let mut pubsub = self.client.get_async_pubsub().await?;

        pubsub.subscribe(channel).await?;

        println!("Subscribed to \"{}\"", channel);

        let mut messages = pubsub.on_message();

        while let Some(message) = messages.next().await {
            let payload: String = message.get_payload()?;

            println!(
                "Received message from \"{}\": {}",
                message.get_channel_name(),
                payload,
            );
        }

        Ok(())
    }
}