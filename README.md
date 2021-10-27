<p align="center">
  <img alt="toolkit logo" src="https://emojipedia-us.s3.dualstack.us-west-1.amazonaws.com/thumbs/320/twitter/282/toolbox_1f9f0.png" height="96" />
  <h3 align="center">Toolkit</h3>
  <p align="center">Repetitive. Chores. Automated.</p>
</p>
<p align="center">
  <a href="https://github.com/thiskevinwang/tk/releases">
    <img alt="latest release" src="https://img.shields.io/github/v/release/thiskevinwang/tk"/>
  </a>
  <!-- <a href="https://github.com/thiskevinwang/tk/tags">
    <img alt="latest tag" src="https://img.shields.io/github/v/tag/thiskevinwang/tk"/>
  </a> -->
  <img alt="go version" src="https://img.shields.io/github/go-mod/go-version/thiskevinwang/tk"/>
</p>

## What is Toolkit?

Toolkit is a CLI tool for myself to avoid doing some repetitive and tedious daily chores. These include, but are not limited to:

- Logging into the AWS Dashboard
- Selecting from a handful of Github repos to open in the browser

... more tasks to come!

## Installation

### Homebrew

```sh
brew tap thiskevinwang/tk https://github.com/thiskevinwang/tk
brew install tk
tk
```

## Usage

### `tk aws [federated-identity-name]`

Launch AWS console in your browser

### `tk repo`

Select a Github repository to open in your browser

### `tk tree [path]`

List contents of directories in a tree-like format.

## Configuration

`tk` will create a yaml configuration file at `<homedir>/.tk/config.yaml`.

```bash
cat ~/.tk/config.yaml

# loglevel: warn
# repos: []
```

## Benchmark

`go test ./cmd -bench=.`
