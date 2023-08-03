package tms

import "testing"

const (
	passed  = "Passed"
	failed  = "Failed"
	skipped = "Skipped"
)

func getTestStatus(t *testing.T) string {
	if t.Failed() {
		return failed
	} else if t.Skipped() {
		return skipped
	} else {
		return passed
	}
}
