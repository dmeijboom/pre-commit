# Pre-commit

### Introduction

Run checks before `git commit` so that your code doesn't suck. Note that the name `pre-commit` already exists so I might
change the name of the project later on (any suggestions?).

### Goals

- Fast
- Customizable
- Easy to use

## Usage

Create a file in the root of the project called `pre-commit.json` and add some checks:

```json
{
  "checks": [
    {
      "name": "Go Tests",
      "cmd": "go test ./...",
      "when": [
        {
          "glob": "*.go"
        }
      ]
    }
  ]
}
```

## Configuration

| Key | Type | Description |
|---|---|---|
|checks|`[]Check`|A list of checks to run before commit|

### Type: `Check`

| Key | Type | Description | Example |
|---|---|---|---|
|name|`string`|Display name of the check|Unittests|
|cmd|`string`|Command to run|go test ./...|
|when|`[]When`|Conditions to test|-|

### Type: `When`

Note that all conditions are validated against the changed files in the Git staging area.

| Key | Type | Description | Example |
|---|---|---|---|
|glob|`string`|Glob pattern (note that it will only be applied to the filename)|`*.go`|
|dir|`string`|Directory (exact match)|`src/`|
