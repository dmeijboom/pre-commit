# Pre-commit

### Introduction

Run checks before `git commit` so that your code doesn't suck.
Note that the name `pre-commit` already exists so I might change the name of the project later on (any suggestions?).


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