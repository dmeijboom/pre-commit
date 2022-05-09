use std::path::PathBuf;
use std::time::Duration;

use colored::Colorize;
use anyhow::Result;
use globset::Glob;
use indicatif::ProgressBar;
use tokio::time;

use crate::config::Check;
use crate::process::Process;

pub struct Task {
    check: Check,
    completed: bool,
    process: Option<Process>,
}

pub struct Runner {
    success: usize,
    changes: Vec<PathBuf>,
    tasks: Vec<Task>,
}

impl Runner {
    pub fn new(changes: Vec<PathBuf>) -> Self {
        Self {
            success: 0,
            tasks: vec![],
            changes,
        }
    }

    pub async fn run(&mut self, pb: ProgressBar, checks: Vec<Check>, run_all: bool) -> Result<bool> {
        self.prepare(checks, run_all)?;

        let task_count = self.tasks
            .iter()
            .filter(|t| t.process.is_some())
            .count();

        pb.set_length(task_count as u64);

        loop {
            let mut dirty = false;
            let tasks = self.tasks
                .iter_mut()
                .filter(|t| !t.completed);

            for task in tasks {
                dirty = true;

                match &mut task.process {
                    None => {
                        task.completed = true;

                        pb.println(format!("{} {}", format!("✓ {}", task.check.name).green(), "[skipped]".white()));
                    }
                    Some(process) => {
                        if let Ok(Some(status)) = process.try_wait() {
                            task.completed = true;
                            pb.inc(1);

                            if status.success() {
                                self.success += 1;

                                pb.println(format!("✓ {}", task.check.name).green().to_string());
                            } else {
                                pb.println(format!("✗ {}", task.check.name).red().to_string());

                                if let Some(process) = task.process.take() {
                                    if let Ok((stdout, stderr)) = process.output().await {
                                        pb.println(stdout.yellow().to_string());
                                        pb.println(stderr.yellow().to_string());
                                    }
                                }
                            }
                        }
                    }
                };
            }

            if !dirty {
                break;
            }

            time::sleep(Duration::from_millis(10)).await;
        }

        pb.finish_and_clear();

        println!("\n{}/{} checks passed", self.success, task_count);

        Ok(self.success == task_count)
    }

    fn prepare(&mut self, checks: Vec<Check>, run_all: bool) -> Result<()> {
        for check in checks {
            if !run_all && !check.when.is_empty() {
                if self.is_skipped(&check)? {
                    self.tasks.push(Task {
                        check,
                        completed: false,
                        process: None,
                    });

                    continue;
                }
            }

            let process = Process::spawn(&check.cmd)?;

            self.tasks.push(Task {
                check,
                completed: false,
                process: Some(process),
            });
        }

        Ok(())
    }

    fn is_skipped(&self, check: &Check) -> Result<bool> {
        for cond in check.when.iter() {
            if let Some(glob_str) = &cond.glob {
                let glob = Glob::new(glob_str)?
                    .compile_matcher();

                if !self.changes.iter().any(|p| glob.is_match(p)) {
                    return Ok(true);
                }
            }
        }

        Ok(false)
    }
}