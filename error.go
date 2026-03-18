package tms

import (
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/testit-tms/adapters-go/models"
)

func fail(err error) {
	addError(err, models.Failed, false)
}

func addError(err error, status string, now bool) {
	manipulateOnObjectFromCtx(
		testResultKey,
		func(test interface{}) {
			testStatusDetails := test.(hasErrorFields)
			testStatusDetails.addTrace(string(debug.Stack()))
			testStatusDetails.addMessage(fmt.Sprintf("%+v", err))
			testStatusDetails.addStatus(status)
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
