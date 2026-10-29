<h1 align="center">go-semverish</h1>

<div align="center">

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.27.3+-00ADD8?logo=go&logoColor=white)](go.mod)

</div>

---

Lenient version comparison for mod-loader ecosystems: loose semver as SMAPI and BepInEx read it, and .NET System.Version ordering. Package `semverish`, module `github.com/Rethunk-Tech/go-semverish`.

Mod authors and mod managers write versions loosely: `v1.2`, `1.0.0.1`, `1.0-beta-2`, `2.0-unofficial.3-author`. `semverish` orders them the way SMAPI does, never panics on a string it cannot read, and says so (`ok == false`) instead of guessing. `CompareDotNet` orders them the way .NET's `System.Version` does, which BepInEx 5 uses. Standard library only.

## Quick start

```bash
go get github.com/Rethunk-Tech/go-semverish
```

A worked example and the contracts: [HUMANS.md](HUMANS.md).

## Features

- `Compare` orders two strings, and `Parse` plus `Version.Compare` order without re-parsing.
- Two to four numbers, a missing number is 0, an optional `v` or `V`, build metadata ignored.
- SMAPI's prerelease rules: `.` and `-` both separate identifiers, `unofficial` ranks lowest.
- Prerelease numbers of any length compare by value, so a hostile string cannot break transitivity.
- `CompareDotNet` follows `System.Version`: `1.0` is below `1.0.0`, and `1.0.0-beta` is not a version.

## Documentation

| Doc | Audience |
| --- | --- |
| [HUMANS.md](HUMANS.md) | Running and using the package |
| [AGENTS.md](AGENTS.md) | Internals, invariants, and why |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Process — commits, tests, docs |
| [SECURITY.md](SECURITY.md) | Vulnerability disclosure |
| [CHANGELOG.md](CHANGELOG.md) | Release notes |

## License

MIT — see [LICENSE](LICENSE). © 2026 Rethunk.Tech, LLC.
