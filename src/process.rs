use std::process::{ExitStatus, Stdio};
use anyhow::Result;
use tokio::process::{Child, Command};

pub struct Process {
    child: Child,
}

impl Process {
    pub fn spawn(command: &str) -> Result<Self> {
        let child = Command::new("/bin/bash")
            .arg("-c")
            .arg(command)
            .stdin(Stdio::null())
            .stdout(Stdio::piped())
            .stderr(Stdio::piped())
            .spawn()?;

        Ok(Self { child })
    }

    pub fn try_wait(&mut self) -> Result<Option<ExitStatus>> {
        Ok(self.child.try_wait()?)
    }

    pub async fn output(self) -> Result<(String, String)> {
        let output = self.child.wait_with_output().await?;

        Ok((
            String::from_utf8(output.stdout)?,
            String::from_utf8(output.stderr)?,
        ))
    }
}