use std::fs;

use anyhow::Result;
use redis::{Client, TlsCertificates};

#[derive(Clone)]
pub struct Setting {
    pub addr: String,
    pub username: String,
    pub password: String,
    pub path_cert: String,
}

pub struct Repository {
    #[allow(dead_code)]
    setting: Setting,
    pub client: Client,
}

impl Repository {
    pub async fn new(setting: Setting) -> Result<Self> {
        let root_cert = fs::read(&setting.path_cert)?;

        let client = Client::build_with_tls(
            format!(
                "rediss://{}:{}@{}",
                setting.username,
                setting.password,
                setting.addr,
            ),
            TlsCertificates {
                client_tls: None,
                root_cert: Some(root_cert),
            },
        )?;

        Ok(Self { setting, client })
    }
}