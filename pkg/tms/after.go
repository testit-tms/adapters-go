package tms

import (
	"fmt"
	"runtime/debug"
	"testing"
	"time"

	"github.com/jtolds/gls"
)

func AfterTest(t *testing.T, m StepMetadata, f func()) {
	after := newAfter(m)

	defer func() {
		after.completedOn = time.Now()
		after.duration = after.completedOn.UnixMilli() - after.startedOn.UnixMilli()
		after.status = getTestStatus(t)

		testPhaseObject := getCurrentTestPhaseObject(t)
		if testPhaseObject.test == nil {
			logger.Error("cannot add after to test before test has been started")
		}

		tr := testPhaseObject.test

		panicObject := recover()
		if panicObject != nil {
			t.Fail()
			after.status = failed
			tr.status = failed
			tr.message = fmt.Sprintf("%+v", panicObject)
			tr.trace = string(debug.Stack())
		}

		tr.addAfter(after.convertToStepResult())
		tr.update(testPhaseObject.resultID)
	}()
	ctxMgr.SetValues(gls.Values{
		testResultKey:   after,
		nodeKey:         after,
		testInstanceKey: t,
	}, f)
}

func newAfter(m StepMetadata) *fixture {
	after := &fixture{
		description: m.Description,
		startedOn:   time.Now(),
		parameters:  m.Parameters,
		name:        m.Name,
	}

	if after.name == "" {
		after.name = "After"
	}

	return after
}
