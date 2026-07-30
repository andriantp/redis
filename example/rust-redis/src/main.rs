mod redis;

use std::collections::HashMap;
use std::env;

use anyhow::{anyhow, Result};

use redis::{Repository, Setting};

#[tokio::main]
async fn main() -> Result<()> {
    let command = env::args()
        .nth(1)
        .ok_or_else(|| anyhow!(r#"Usage:
            cargo run -- <command>
                Commands:
                set
                setex
                get
                ttl
                hset
                hget
                hgetall
                publish
                subscribe
                leaderboard"#))?;

    let setting = Setting {
        addr: "localhost:6380".to_string(),
        username: "cache-service".to_string(),
        password: "cache123".to_string(),
        path_cert: "../../docker/certs/ca.crt".to_string(),
    };

    let repo = Repository::new(setting).await?;

    match command.as_str() {

        // Strings
        "set" => {
            repo.set(
                "cache:strings:title",
                "Exploring Redis Features",
            )
            .await?;

            println!("OK");
        }

        "setex" => {
            repo.setex(
                "cache:strings:temporary",
                "Hello Redis",
                60,
            )
            .await?;

            println!("OK");
        }

        "get" => {
            let value = repo
                .get("cache:strings:title")
                .await?;

            println!("{}", value);
        }

        "ttl" => {
            let ttl = repo
                .ttl("cache:strings:temporary")
                .await?;

            println!("{}", ttl);
        }

        // Hashes
        "hset" => {
            let mut profile = HashMap::new();

            profile.insert(
                "name".to_string(),
                "Andrian".to_string(),
            );

            profile.insert(
                "email".to_string(),
                "andrian@example.com".to_string(),
            );

            profile.insert(
                "role".to_string(),
                "admin".to_string(),
            );

            repo.hset(
                "cache:user:1001",
                &profile,
            )
            .await?;

            println!("OK");
        }

        "hget" => {
            let value = repo
                .hget(
                    "cache:user:1001",
                    "name",
                )
                .await?;

            println!("{}", value);
        }

        "hgetall" => {
            let profile = repo
                .hgetall("cache:user:1001")
                .await?;

            println!("{:?}", profile);
        }

        // Pub/Sub
        "publish" => {
            repo.publish(
                "cache:notifications",
                "Hello Redis!",
            )
            .await?;

            println!("Message published");
        }

        "subscribe" => {
            repo.subscribe(
                "cache:notifications",
            )
            .await?;
        }

        // Sorted Sets
        "leaderboard" => {
            repo.zadd(
                "cache:leaderboard",
                1500.0,
                "Alice",
            )
            .await?;

            repo.zadd(
                "cache:leaderboard",
                3200.0,
                "Bob",
            )
            .await?;

            repo.zadd(
                "cache:leaderboard",
                2100.0,
                "Charlie",
            )
            .await?;

            repo.zadd(
                "cache:leaderboard",
                1800.0,
                "David",
            )
            .await?;

            let ranking = repo
                .zrevrange(
                    "cache:leaderboard",
                    0,
                    2,
                )
                .await?;

            println!("{:?}", ranking);
        }

        _ => {
            return Err(anyhow!(
                "Unknown command: {}",
                command
            ));
        }
    }

    Ok(())
}

