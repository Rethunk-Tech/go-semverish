# AGENTS.md

Package `semverish` orders version strings the way SMAPI does (two to four numbers, `.` and `-` prerelease separators, `unofficial` lowest) and, separately, the way .NET's `System.Version` does (`CompareDotNet`). Mod managers and update checkers import it. It imports only the standard library and no consumer package.

## Gate

`gate` is the command. Each override in `.gate.toml` is there because a runner
cannot resolve that check from PATH.

## Shape

- `Parse` is the only reader of SMAPI-style strings. `Compare` goes through `Parse` and `Version.Compare`, so a string and a parsed value never order differently.
- Prerelease numbers compare by trimmed length, then bytes, never by `Atoi`: a wrapped integer made the order intransitive. `FuzzCompare` holds transitivity.
- An empty prerelease identifier fails the regexp, so `Version.Pre` never holds one and splitting it on `.` and `-` is exact.
- `Build` is kept only so `String` can write it back; no comparison reads it.
- `CompareDotNet` is a separate function, not a mode: the two disagree on a missing part (0 against -1), and a flag would make every caller choose.
- The regexp, `comparePre` and the .NET parser stay unexported: the API is the contract, not the grammar's spelling.
- Leniency (leading zeros, `0.0.0`, `V`) is deliberate: Nexus and GitHub versions are hand-typed and never passed SMAPI's parser. Rejecting them would make a version unorderable for no gain.

##| File | Holds |
| --- | --- |
| `semverish.go` | `Version`, `Parse`, `String`, `Compare`, the prerelease comparator |
| `dotnet.go` | `CompareDotNet` and its parser |
