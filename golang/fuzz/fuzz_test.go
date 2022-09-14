package fuzz

import (
	"testing"
)

func FuzzDiv(f *testing.F) {
	testcases := []struct {
		a, b int
	}{
		{10, 2},
		{5, 3},
		{-6, 3},
		{-6, -3},
	}
	for _, tc := range testcases {
		f.Add(tc.a, tc.b) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, a, b int) {
		q := Div(a, b)
		if q != a/b {
			t.Errorf("Before: %q, after: %q", q, a/b)
		}
	})
}
