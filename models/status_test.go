package models

import "testing"

func TestMergeStatus(t *testing.T) {
	cases := []struct {
		current, incoming, want string
	}{
		{"", Passed, Passed},
		{Passed, Failed, Failed},
		{Failed, Passed, Failed},
		{Failed, Skipped, Failed},
		{Passed, Skipped, Skipped},
		{Skipped, Passed, Skipped},
		{Skipped, Failed, Failed},
		{Passed, Passed, Passed},
		{Failed, Failed, Failed},
	}
	for _, tc := range cases {
		if got := MergeStatus(tc.current, tc.incoming); got != tc.want {
			t.Fatalf("MergeStatus(%q, %q) = %q, want %q", tc.current, tc.incoming, got, tc.want)
		}
	}
}
