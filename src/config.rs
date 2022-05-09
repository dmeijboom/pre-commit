use std::fs;
use std::path::{Path, PathBuf};

use anyhow::Result;
use serde::Deserialize;

#[derive(Debug, Clone, Deserialize)]
pub struct When {
    pub glob: Option<String>,
    pub dir: Option<PathBuf>,
}

#[derive(Debug, Clone, Deserialize)]
pub struct Check {
    pub name: String,
    pub cmd: String,
    #[serde(default)]
    pub when: Vec<When>,
}

#[derive(Debug, Deserialize)]
pub struct Config {
    pub checks: Vec<Check>,
}

impl Config {
    pub fn from_file(path: impl AsRef<Path>) -> Result<Self> {
        let contents = fs::read_to_string(path.as_ref())?;

        Ok(serde_json::from_str(&contents)?)
    }
}