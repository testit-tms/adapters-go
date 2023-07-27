package tms

import (
	"runtime/debug"
	"testing"
)

func Fail(err error) {
	addError(err, "Failed", false)
}

func FailNow(err error) {
	addError(err, "Failed", true)
}

func addError(err error, status string, now bool) {
	manipulateOnObjectFromCtx(
		testResultKey,
		func(test interface{}) {
			testStatusDetails := test.(*testResult)
			testStatusDetails.trace = string(debug.Stack())
			testStatusDetails.message = err.Error()
			testStatusDetails.status = status
		})
	manipulateOnObjectFromCtx(
		nodeKey,
		func(node interface{}) {
			n := node.(hasStatus)
			n.addStatus(status)
		})
	manipulateOnObjectFromCtx(
		testInstanceKey,
		func(testInstance interface{}) {
			testInstance.(*testing.T).Fail()
			if now {
				panic(err)
			}
		})
}
