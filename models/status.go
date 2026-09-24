package models

import "testing"

const (
	Passed  = "Passed"
	Failed  = "Failed"
	Skipped = "Skipped"
)

func GetTestStatus(t *testing.T) string {
	if t.Failed() {
		return Failed
	} else if t.Skipped() {
		return Skipped
	} else {
		return Passed
	}
}

// MergeStatus keeps the worse outcome so a later Passed cannot overwrite Failed.
// Priority: Failed > Skipped > Passed > empty.
func MergeStatus(current, incoming string) string {
	if rankStatus(incoming) > rankStatus(current) {
		return incoming
	}
	return current
}

func rankStatus(s string) int {
	switch s {
	case Failed:
		return 3
	case Skipped:
		return 2
	case Passed:
		return 1
	default:
		return 0
	}
}
