# Security

## Supported versions

The latest tagged release. There are no long-lived release branches, so a fix ships in the next tag
rather than being backported.

## Reporting a vulnerability

Use GitHub's **private vulnerability reporting** on this repository (Security → Report a
vulnerability). It stays private until a fix ships.

Please do not open a public issue for a suspected vulnerability first.

## What is in scope

- Every function takes version strings from remote data (mod pages, tags, manifests). No input may panic, allocate more than a few times its own length, or make the order intransitive; `FuzzCompare` and `FuzzCompareDotNet` check the order laws and `Parse`'s round trip.
- A number or prerelease identifier too large for its type must be refused or compared by value, never wrapped.

## What is out of scope

- Deciding what an unparseable version means (`ok == false`) is the caller's policy.
- A version that is ordered but not a real release, such as `0.0.0` or a hand-typed leading zero, is accepted on purpose.
