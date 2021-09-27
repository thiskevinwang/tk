<p align="center">
  <h1 align="center">🧰</h1>
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

## Configuration

`tk` will create a yaml configuration file at `<homedir>/.tk/config.yaml`.

```bash
cat ~/.tk/config.yaml

# loglevel: warn
# repos: []
```
