package tms

import (
	"fmt"
	"runtime/debug"
	"testing"
	"time"

	"github.com/jtolds/gls"
)

func BeforeTest(t *testing.T, m StepMetadata, f func()) {
	testPhaseObject := getCurrentTestPhaseObject(t)
	if testPhaseObject.test != nil {
		logger.Error("cannot add before to test after test has been started")
	}
	before := newBefore(m)
	testPhaseObject.before = before

	defer func() {
		panicObject := recover()

		before.completedOn = time.Now()
		if before.status == "" {
			before.status = getTestStatus(t)
		}
		before.duration = before.completedOn.UnixMilli() - before.startedOn.UnixMilli()

		if panicObject != nil {
			t.Fail()
			before.message = fmt.Sprintf("%+v", panicObject)
			before.trace = string(debug.Stack())
			before.status = failed
		}
	}()
	ctxMgr.SetValues(gls.Values{
		testResultKey:   before,
		nodeKey:         before,
		testInstanceKey: t,
	}, f)
}

func newBefore(m StepMetadata) *fixture {
	before := &fixture{
		description: m.Description,
		startedOn:   time.Now(),
		parameters:  m.Parameters,
		name:        m.Name,
	}

	if before.name == "" {
		before.name = "Before"
	}

	return before
}
