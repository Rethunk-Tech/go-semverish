package semverish

import (
	"cmp"
	"regexp"
	"strconv"
	"strings"
)

var pattern = regexp.MustCompile(`^[vV]?(\d+)\.(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:-([0-9A-Za-z]+(?:[.-][0-9A-Za-z]+)*))?(?:\+([0-9A-Za-z.-]+))?$`)

// Version is a parsed version. Missing numbers are 0.
type Version struct {
	Major, Minor, Patch int
	// Platform is SMAPI's fourth number.
	Platform int
	// Pre and Build are written without their leading '-' and '+'. Build never affects order.
	Pre, Build string
}

// Parse reads a version: two to four numbers, an optional leading v or V, an optional -prerelease whose
// identifiers are separated by '.' or '-', and optional +build, with surrounding whitespace trimmed. Leading
// zeros are accepted.
func Parse(s string) (Version, bool) {
	m := pattern.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return Version{}, false
	}
	var nums [4]int
	for i := range nums {
		if m[i+1] == "" {
			continue
		}
		n, err := strconv.Atoi(m[i+1])
		if err != nil {
			return Version{}, false
		}
		nums[i] = n
	}
	return Version{nums[0], nums[1], nums[2], nums[3], m[5], m[6]}, true
}

// String writes "M.m.p[.platform][-pre][+build]" without a v. The platform number is left out when 0.
func (v Version) String() string {
	var b strings.Builder
	b.WriteString(strconv.Itoa(v.Major))
	b.WriteByte('.')
	b.WriteString(strconv.Itoa(v.Minor))
	b.WriteByte('.')
	b.WriteString(strconv.Itoa(v.Patch))
	if v.Platform != 0 {
		b.WriteByte('.')
		b.WriteString(strconv.Itoa(v.Platform))
	}
	if v.Pre != "" {
		b.WriteByte('-')
		b.WriteString(v.Pre)
	}
	if v.Build != "" {
		b.WriteByte('+')
		b.WriteString(v.Build)
	}
	return b.String()
}

// Compare orders v against w: -1, 0 or 1. A release outranks its prereleases, build metadata is ignored, and
// an "unofficial" identifier ranks below any other.
func (v Version) Compare(w Version) int {
	if c := cmp.Or(
		cmp.Compare(v.Major, w.Major),
		cmp.Compare(v.Minor, w.Minor),
		cmp.Compare(v.Patch, w.Patch),
		cmp.Compare(v.Platform, w.Platform),
	); c != 0 {
		return c
	}
	switch {
	case v.Pre == "" && w.Pre == "":
		return 0
	case v.Pre == "":
		return 1
	case w.Pre == "":
		return -1
	}
	x, y := splitPre(v.Pre), splitPre(w.Pre)
	for i := range min(len(x), len(y)) {
		if c := comparePre(x[i], y[i]); c != 0 {
			return c
		}
	}
	return cmp.Compare(len(x), len(y))
}

// Compare parses both strings and orders them as [Version.Compare] does. ok is false, and the order 0, when
// either does not parse.
func Compare(a, b string) (c int, ok bool) {
	x, okA := Parse(a)
	y, okB := Parse(b)
	if !okA || !okB {
		return 0, false
	}
	return x.Compare(y), true
}

func splitPre(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool { return r == '.' || r == '-' })
}

// comparePre orders two prerelease identifiers: "unofficial" lowest, then numeric ones by value (any length),
// which rank below alphanumeric ones, then text ignoring case.
func comparePre(a, b string) int {
	a, b = strings.ToLower(a), strings.ToLower(b)
	if a == b {
		return 0
	}
	switch {
	case b == "unofficial":
		return 1
	case a == "unofficial":
		return -1
	}
	na, nb := isDigits(a), isDigits(b)
	switch {
	case na && nb:
		a, b = strings.TrimLeft(a, "0"), strings.TrimLeft(b, "0")
		return cmp.Or(cmp.Compare(len(a), len(b)), strings.Compare(a, b))
	case na:
		return -1
	case nb:
		return 1
	}
	return strings.Compare(a, b)
}

func isDigits(s string) bool {
	return !strings.ContainsFunc(s, func(r rune) bool { return r < '0' || r > '9' })
}

// Equal reports whether a and b are the same version: both parse and compare equal ("1.2" and "1.2.0"), or,
// when either does not parse, the strings are identical.
func Equal(a, b string) bool {
	if c, ok := Compare(a, b); ok {
		return c == 0
	}
	return a == b
}
