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
