# Changelog

Notable changes to `go-semverish`. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and versions follow
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-10-29

### Added

- `Compare`, `Parse` and `Version.Compare` order SMAPI-style versions: two to four numbers, `.` and `-` prerelease separators, `unofficial` lowest.
- `Version.String` writes a parsed version back as `M.m.p[.platform][-pre][+build]`.
- `CompareDotNet` orders versions the way .NET's `System.Version` does.

[Unreleased]: https://github.com/Rethunk-Tech/go-semverish/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/Rethunk-Tech/go-semverish/releases/tag/v0.1.0
