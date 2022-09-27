package fuzz

import (
	"testing"
)

func TestDiv(t *testing.T) {
	testcases := []struct {
		a, b, want int
	}{
		{10, 2, 5},
		{5, 3, 1},
		{-6, 3, -2},
		{-6, -3, 2},
	}
	for _, tc := range testcases {
		result := Div(tc.a, tc.b)
		if Div(tc.a, tc.b) != tc.want {
			t.Errorf("Div: %q, want %q", result, tc.want)
		}
	}
}

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
