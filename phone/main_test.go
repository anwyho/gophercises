package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNormalize(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"123-456-7890", "1234567890"},
		{"(123)456-7892", "1234567892"},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			got := normalize(tc.input)
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
