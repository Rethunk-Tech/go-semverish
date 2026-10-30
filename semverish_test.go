package semverish

import (
	"slices"
	"testing"
)

func TestCompare(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.0", "1.0.0", 0},
		{"2.9.1", "2.10.0", -1},
		{"1.0.0", "1.0.0-beta", 1},
		{"1.0.0-beta", "1.0.0-beta.2", -1},
		{"1.0.0-beta.2", "1.0.0-beta.10", -1},
		{"1.0.0-alpha", "1.0.0-1", 1},
		{"1.0.0-RC1", "1.0.0-rc1", 0},
		{"v1.2.3+build5", "1.2.3", 0},
		{"1.2.3.4", "1.2.3", 1},
		{"3.0.0-unofficial.1-pathoschild", "3.0.0", -1},
		// SMAPI SemanticVersionTests.CompareTo
		{"0.5.7", "0.5.7", 0},
		{"1.0", "1.0", 0},
		{"1.0-beta", "1.0-beta", 0},
		{"1.0-beta.10", "1.0-beta.10", 0},
		{"1.0-beta", "1.0-beta   ", 0},
		{"1.0-beta+build.001", "1.0-beta+build.001", 0},
		{"1.0-beta+build.001", "1.0-beta+build.006", 0},
		{"0.5.7", "0.5.8", -1},
		{"1.0", "1.1", -1},
		{"1.0-beta", "1.0", -1},
		{"1.0-beta", "1.0-beta.2", -1},
		{"1.0-beta.1", "1.0-beta.2", -1},
		{"1.0-beta.2", "1.0-beta.10", -1},
		{"1.0-beta-2", "1.0-beta-10", -1},
		{"1.0-unofficial.1", "1.0-beta.1", -1},
		{"0.5.8", "0.5.7", 1},
		{"1.1", "1.0", 1},
		{"1.0", "1.0-beta", 1},
		{"1.0-beta.2", "1.0-beta", 1},
		{"1.0-beta.2", "1.0-beta.1", 1},
		{"1.0-beta.10", "1.0-beta.2", 1},
		{"1.0-beta-10", "1.0-beta-2", 1},
		// departures from SMAPI
		{"1.0-5", "1.0-10a", -1},
		{"1.0-3000000000", "1.0-9", 1},
		{"1.0-Beta.2", "1.0-beta.1", 1},
		{"1.0-Unofficial.1", "1.0-unofficial.1", 0},
		{"V1.0", "1.0", 0},
		// measured defects
		{"1.0.0-unofficial.10-pathoschild", "1.0.0-unofficial.9-pathoschild", 1},
		{"1.0.0-beta", "1.0.0-unofficial.1-pathoschild", 1},
		{"1.0-99999999999999999999", "1.0-1", 1},
	}
	for _, c := range cases {
		got, ok := Compare(c.a, c.b)
		if !ok || got != c.want {
			t.Errorf("Compare(%q, %q) = %d, %v; want %d", c.a, c.b, got, ok, c.want)
		}
		if rev, _ := Compare(c.b, c.a); rev != -c.want {
			t.Errorf("Compare(%q, %q) = %d, not the mirror", c.b, c.a, rev)
		}
	}
	for _, bad := range []string{"", "1", "abc", "1.x", "1.0-a..b", "1.0-.a", "1.0-beta.", "1.0--1"} {
		if _, ok := Compare(bad, "1.0"); ok {
			t.Errorf("%q should not parse", bad)
		}
	}
}

func TestCompareDotNet(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want int
	}{
		{"1.0", "1.0.0", -1},
		{"1.0.0", "1.0.0.0", -1},
		{"1.0.1", "1.0.0.9", 1},
		{"01.2", "1.2", 0},
		{" 1.2 ", "1.2", 0},
		{"+1.2", "1.2", 0},
		{"1. 2", "1.2", 0},
		{"2.0", "10.0", -1},
		{"1.-0", "1.0", 0},
		{"1.2.3.-0", "1.2.3.0", 0},
	} {
		got, ok := CompareDotNet(c.a, c.b)
		if !ok || got != c.want {
			t.Errorf("CompareDotNet(%q, %q) = %d, %v; want %d", c.a, c.b, got, ok, c.want)
		}
		if rev, _ := CompareDotNet(c.b, c.a); rev != -c.want {
			t.Errorf("CompareDotNet(%q, %q) = %d, not the mirror", c.b, c.a, rev)
		}
	}
	for _, bad := range []string{"", "1", "1.0.0-beta", "v1.0", "1.2.3.4.5", "3000000000.0", "-1.0", "1.-1", "1..2", "1.2\u00a0", "1.2\u0085"} {
		if _, ok := CompareDotNet(bad, "1.0"); ok {
			t.Errorf("%q should not parse", bad)
		}
	}
}

func TestParseAndString(t *testing.T) {
	for _, c := range []struct {
		in   string
		want Version
		out  string
	}{
		{"1.2", Version{Major: 1, Minor: 2}, "1.2.0"},
		{" V1.2.3.4-RC.1+b.7 ", Version{1, 2, 3, 4, "RC.1", "b.7"}, "1.2.3.4-RC.1+b.7"},
		{"01.02-unofficial.10-pathoschild", Version{Major: 1, Minor: 2, Pre: "unofficial.10-pathoschild"}, "1.2.0-unofficial.10-pathoschild"},
		{"0.0.0", Version{}, "0.0.0"},
	} {
		got, ok := Parse(c.in)
		if !ok || got != c.want || got.String() != c.out {
			t.Errorf("Parse(%q) = %+v, %v, %q; want %+v, %q", c.in, got, ok, got.String(), c.want, c.out)
		}
	}
	// SMAPI's legacy manifest object prints M.m.p[-build].
	if got := (Version{Major: 4, Minor: 5, Patch: 2, Pre: "beta"}).String(); got != "4.5.2-beta" {
		t.Errorf("String = %q", got)
	}
}

// A signed or oversized identifier once wrapped around in the comparison and made the order intransitive.
func TestPrereleaseOrderIsTransitive(t *testing.T) {
	for _, s := range []string{"1.0--9223372036854775808", "1.0-9223372036854775808-"} {
		if _, ok := Parse(s); ok {
			t.Errorf("%q should not parse", s)
		}
	}
	vs := []string{"1.0-99999999999999999999", "1.0-1", "1.0-0"}
	slices.SortFunc(vs, func(a, b string) int { c, _ := Compare(a, b); return c })
	if want := []string{"1.0-0", "1.0-1", "1.0-99999999999999999999"}; !slices.Equal(vs, want) {
		t.Errorf("sorted = %v", vs)
	}
}

func orderLaws(t *testing.T, cmpFn func(a, b string) (int, bool), a, b, c string) {
	t.Helper()
	ss := [3]string{a, b, c}
	var r [3][3]int
	parsed := true
	for i, x := range ss {
		for j, y := range ss {
			v, ok := cmpFn(x, y)
			if !ok {
				if v != 0 {
					t.Fatalf("(%q, %q) = %d with ok false", x, y, v)
				}
				parsed = false
				continue
			}
			r[i][j] = v
		}
	}
	if !parsed {
		return
	}
	for i := range ss {
		if r[i][i] != 0 {
			t.Fatalf("(%q, %q) = %d, not 0", ss[i], ss[i], r[i][i])
		}
		for j := range ss {
			if r[i][j] != -r[j][i] {
				t.Fatalf("(%q, %q) = %d, reverse %d", ss[i], ss[j], r[i][j], r[j][i])
			}
			for k := range ss {
				if r[i][j] <= 0 && r[j][k] <= 0 && r[i][k] > 0 {
					t.Fatalf("%q <= %q <= %q but %q > %q", ss[i], ss[j], ss[k], ss[i], ss[k])
				}
			}
		}
	}
}

func FuzzCompare(f *testing.F) {
	f.Add("1.0", "1.0.0-beta", "v1.0.0-unofficial.1-x+b")
	f.Add("1.0-beta-10", "1.0-Beta.2", "1.0-99999999999999999999")
	f.Add("1.0--1", "1.0-a..b", "x")
	f.Fuzz(func(t *testing.T, a, b, c string) {
		orderLaws(t, Compare, a, b, c)
		for _, s := range []string{a, b, c} {
			v, ok := Parse(s)
			if !ok {
				continue
			}
			w, ok := Parse(v.String())
			if !ok || v.Compare(w) != 0 {
				t.Fatalf("Parse(%q).String() = %q reads back as %+v, %v", s, v.String(), w, ok)
			}
		}
	})
}

func FuzzCompareDotNet(f *testing.F) {
	f.Add("1.0", "1.0.0", "1.0.0.0")
	f.Add(" +1.2", "1.-0", "1.2.3.-0")
	f.Add("3000000000.0", "1..2", "v1.0")
	f.Fuzz(func(t *testing.T, a, b, c string) {
		orderLaws(t, CompareDotNet, a, b, c)
	})
}

func TestEqual(t *testing.T) {
	for _, tc := range []struct {
		a, b string
		want bool
	}{
		{"1.0", "1.0.0", true},
		{"v1.2", "1.2.0+build", true},
		{"1.0", "1.1", false},
		{"not-a-version", "not-a-version", true},
		{"", "", true},
		{"not-a-version", "other", false},
		{"1.0", "one", false},
	} {
		if got := Equal(tc.a, tc.b); got != tc.want {
			t.Errorf("Equal(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}
