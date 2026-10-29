package semverish

import (
	"slices"
	"strconv"
	"strings"
)

// CompareDotNet orders two versions the way .NET's System.Version does, which BepInEx 5 uses: two to four plain
// numbers, a part left out ranks below 0 (so 1.0 < 1.0.0), and anything else (a suffix, a leading v, one number) is
// not a version. ok is false, and the order 0, when either is not one.
func CompareDotNet(a, b string) (c int, ok bool) {
	x, okA := parseDotNet(a)
	y, okB := parseDotNet(b)
	if !okA || !okB {
		return 0, false
	}
	return slices.Compare(x[:], y[:]), true
}

func parseDotNet(s string) (v [4]int, ok bool) {
	parts := strings.Split(s, ".")
	if len(parts) < 2 || len(parts) > 4 {
		return v, false
	}
	for i := range v {
		v[i] = -1
		if i >= len(parts) {
			continue
		}
		n, err := strconv.ParseInt(strings.Trim(parts[i], " \t\n\v\f\r"), 10, 32)
		if err != nil || n < 0 {
			return v, false
		}
		v[i] = int(n)
	}
	return v, true
}
