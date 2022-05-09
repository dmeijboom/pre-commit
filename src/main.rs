use std::process::exit;
use clap::Parser;
use indicatif::{ProgressBar, ProgressStyle};

mod config;
mod process;
mod git;
mod runner;

use config::Config;
use runner::Runner;

#[derive(Debug, Parser)]
#[clap(about = "Make sure your code is OK before shipping")]
struct Args {
    #[clap(subcommand)]
    cmd: Cmd,
}

#[derive(Debug, Parser)]
enum Cmd {
    #[clap(about = "Run the checks")]
    Run {
        #[clap(long, help = "Run all checks (ignore when conditions)")]
        all: bool
    },
    #[clap(about = "Install git-hook and example config")]
    Install,
}

#[tokio::main]
async fn main() {
    let args = Args::parse();

    match args.cmd {
        Cmd::Run { all } => {
            let config = Config::from_file("pre-commit.json")
                .expect("failed to read config");
            let changes = git::list_changes().expect("failed to list Git changes");
            let style = ProgressStyle::default_bar()
                .template(
                    "[{elapsed_precise}] {bar:40.cyan/blue} {pos:>7}/{len:7} {msg}",
                )
                .progress_chars("##-");
            let pb = ProgressBar::new(config.checks.len() as u64).with_style(style);

            pb.tick();
            pb.enable_steady_tick(1000);

            let mut runner = Runner::new(changes);
            let succeeded = runner.run(pb, config.checks, all)
                .await
                .expect("failed to run checks");

            if !succeeded {
                exit(1);
            }
        }
        Cmd::Install => {
            panic!("not implemented");
        }
    }
}
