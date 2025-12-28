# WhyIsMyBuildSlow 🐌

A CLI tool that explains *why* your builds are slow with timing analysis,
bottleneck classification, and a live animated TUI.

## Features
- Detects idle build gaps
- Classifies causes (network, cache, Docker)
- Animated terminal UI
- Headless mode for CI (`--no-ui`)

## Usage

```bash
whyismybuildslow run -- npm install
whyismybuildslow run --no-ui -- sleep 4
