package semverish

import (
	"strings"
	"testing"
)

func FuzzParse(f *testing.F) {
	for _, s := range []string{
		"1.0",
		"1.0.0",
		"1.2.3",
		"1.2.3.4",
		"v1.2.3",
		"V1.0",
		" V1.2.3.4-RC.1+b.7 ",
		"1.2.3-beta",
		"1.0-beta",
		"1.0-beta.10",
		"1.0-beta+build.001",
		"1.2.3-unofficial",
		"01.02-unofficial.10-pathoschild",
		"3.0.0-unofficial.1-pathoschild",
		"1.0.0-unofficial.10-pathoschild",
		"1.0-99999999999999999999",
		"0.0.0",
		"0.5.7",
		"  1.2.3  ",
		"1.0-beta   ",
		"",
		"1",
		"abc",
		"1.x",
		"1.0-a..b",
		"1.0-.a",
		"1.0-beta.",
		"1.0--1",
		"1.2.3-",
		"1.2.3+",
	} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if len(strings.TrimSpace(s)) > MaxLen {
			if _, ok := Parse(s); ok {
				t.Fatalf("Parse accepted over MaxLen")
			}
			return
		}
		v, ok := Parse(s)
		if !ok {
			return
		}
		out := v.String()
		v2, ok := Parse(out)
		if !ok {
			t.Fatalf("Parse(%q) ok but Parse(String=%q) failed", s, out)
		}
		if v.Compare(v2) != 0 || v2 != v {
			t.Fatalf("round-trip %q -> %q: %+v vs %+v", s, out, v, v2)
		}
	})
}
