use std::path::PathBuf;
use anyhow::Result;
use git2::{Repository, Status, StatusOptions};

pub fn list_changes() -> Result<Vec<PathBuf>> {
    let repo = Repository::open(".")?;
    let mut opts = StatusOptions::new();

    opts.include_untracked(true);
    opts.exclude_submodules(true);

    let mut changes = vec![];
    let statuses = repo.statuses(Some(&mut opts))?;

    for entry in statuses.iter().filter(|e| e.status() != git2::Status::CURRENT) {
        if !matches!(entry.status(), Status::INDEX_NEW | Status::INDEX_MODIFIED | Status::INDEX_DELETED | Status::INDEX_RENAMED | Status::INDEX_TYPECHANGE) {
            continue;
        }

        if let Some(path) = entry.head_to_index().unwrap().new_file().path() {
            changes.push(path.to_path_buf());
        }
    }

    Ok(changes)
}