# HUMANS.md

Run and use `go-semverish`. Internals: [AGENTS.md](AGENTS.md).

## Quick start

```bash
go get github.com/Rethunk-Tech/go-semverish
```

```go
package main

import (
	"fmt"
	"slices"

	"github.com/Rethunk-Tech/go-semverish"
)

func main() {
	c, ok := semverish.Compare("1.0-beta-10", "1.0-beta-2")
	fmt.Println(c, ok) // 1 true

	tags := []string{"v1.2", "1.10.0", "1.2.0-unofficial.1-x"}
	slices.SortStableFunc(tags, func(a, b string) int {
		c, _ := semverish.Compare(a, b)
		return c
	})
	fmt.Println(tags)

	v, _ := semverish.Parse("V2.0.1-rc.1+build5")
	fmt.Println(v.String()) // 2.0.1-rc.1+build5

	d, ok := semverish.CompareDotNet("1.0", "1.0.0")
	fmt.Println(d, ok) // -1 true
}
```

## Contracts

| Behaviour | Guarantee |
| --- | --- |
| `Equal` | True when both strings parse and compare equal (`1.0` and `1.0.0`); when either does not parse, true only for identical strings |
| `Compare`, `CompareDotNet` | `ok` is false and the order 0 when either string does not parse; the caller picks the policy |
| `Parse` | Accepts two to four numbers, `v` or `V`, surrounding whitespace, leading zeros and `0.0.0`; refuses one number, five numbers, a number past `int`, and an empty prerelease identifier (`1.0-a..b`, `1.0-beta.`) |
| `Version.Compare` | A total order: release above its prereleases, `Build` ignored, identifiers compared ignoring case, numeric identifiers by value below alphanumeric ones, `unofficial` lowest |
| `Version.String` | `M.m.p[.platform][-pre][+build]`, no `v`; a zero `Platform` is omitted, so `Parse(v.String())` orders equal to `v` |
| `CompareDotNet` | Two to four plain numbers; a missing part ranks below 0 (`1.0` < `1.0.0` < `1.0.0.0`); a prefix, a suffix, one number or a part past int32 is not a version |
| Panics | Never, on any string |

### Where SMAPI's order is not followed

| Input | Here | SMAPI |
| --- | --- | --- |
| `1.0-5` against `1.0-10a` | `<` (numeric below alphanumeric) | `>` (compared as text) |
| `1.0-3000000000` against `1.0-9` | `>` (by value) | `<` (does not fit int32, so text) |
| `1.0-Beta.2` against `1.0-beta.1` | `>` (comparison continues past a case-only tie) | `=` |
| `1.0-beta.` | does not parse | parses |
| `V1.0`, `01.2`, `0.0.0` | parse | `V1.0` and `0.0.0` refused |

### `CompareDotNet` against .NET's `System.Version`

Measured on .NET 10; BepInEx 5 on Mono follows the same reference-source rules (`Version.cs`).

| Input | Result |
| --- | --- |
| `1.0` against `1.0.0` | `-1`: an unset part is -1 |
| `1.0.0` against `1.0.0.0` | `-1` |
| `1.0.0-beta`, `v1.0`, `1`, `1.2.3.4.5`, `3000000000.0` | do not parse |
| `01.2`, ` 1.2 `, `+1.2`, `1. 2` | parse, equal to `1.2` |

## Requirements

Go 1.27.3 or newer, plus `lefthook` and `gitleaks` on PATH if you want the git hooks.

```bash
lefthook install
```

That is per clone; git does not commit hooks. Pre-commit scans staged changes for secrets, and
pre-push runs `gate`.

## Verify

```bash
gate
```
