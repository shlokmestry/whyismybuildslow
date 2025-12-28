# WhyIsMyBuildSlow 🐌

A CLI tool that explains *why* your builds are slow with timing analysis,
bottleneck classification, and a live animated TUI.

## Features
- Detects idle build gaps
- Classifies causes (network, cache, Docker)
- Animated terminal UI
- Headless mode for CI (`--no-ui`)

## When should I use this?

Use `whyismybuildslow` when:
- Your build feels slow but logs don’t explain why
- CI builds are slower than local runs
- Dependency installs take unpredictable time
- You want to understand *idle* or *waiting* time, not just total duration

You run your build exactly the same way — just prefix it with `whyismybuildslow`.

## Usage

```bash
whyismybuildslow run -- npm install
whyismybuildslow run --no-ui -- sleep 4
