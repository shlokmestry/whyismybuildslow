# WhyIsMyBuildSlow 🐌

<<<<<<< HEAD
A CLI tool that explains *why* your builds are slow with timing analysis,
bottleneck classification, and a live animated TUI.
=======
A CLI tool that explains *why* your builds are slow — using timing analysis,
bottleneck classification, and a live animated terminal UI.

You run your build exactly the same way — just prefix it with
`whyismybuildslow`.

---
>>>>>>> 4b4ee98 (docs: update README)

## Features

<<<<<<< HEAD
## When should I use this?

Use `whyismybuildslow` when:
- Your build feels slow but logs don’t explain why
- CI builds are slower than local runs
- Dependency installs take unpredictable time
- You want to understand *idle* or *waiting* time, not just total duration

You run your build exactly the same way — just prefix it with `whyismybuildslow`.

## Usage
=======
- Detects idle gaps during builds
- Classifies bottlenecks (network, cache, Docker)
- Animated terminal UI for humans
- Headless mode for CI (`--no-ui`)
- Machine-readable output for automation (`--json`)

---

## Installation
>>>>>>> 4b4ee98 (docs: update README)

```bash
go install github.com/shlokmestry/whyismybuildslow/cmd/whyismybuildslow@v1.0.2
