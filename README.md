# WhyIsMyBuildSlow 🐌

A CLI tool that explains **why** your builds are slow using timing analysis,
bottleneck classification, and a live animated terminal UI.

You run your build exactly the same way — just prefix it with  
`whyismybuildslow`.

---

## When should I use this?

Use **whyismybuildslow** when:

- Your build feels slow but logs don’t explain why
- CI builds are slower than local runs
- Dependency installs take unpredictable time
- You want to understand *idle* or *waiting* time, not just total duration

This tool helps surface *where time is actually lost* during a build.

---

## Features

- Detects idle gaps during builds
- Classifies bottlenecks (network, cache, Docker)
- Animated terminal UI for humans
- Headless mode for CI (`--no-ui`)
- Machine-readable output for automation (`--json`)

---

## Usage

Run any command exactly as you normally would — just prefix it:

```bash
whyismybuildslow run -- npm install
whyismybuildslow run --no-ui -- sleep 4
whyismybuildslow run --json -- sleep 2
