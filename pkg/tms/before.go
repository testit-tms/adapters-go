package tms

import (
	"fmt"
	"runtime/debug"
	"testing"
	"time"

	"github.com/jtolds/gls"
)

type before struct {
	name          string
	description   string
	childrenSteps []step
	status        string
	startedOn     time.Time
	completedOn   time.Time
	duration      int64
	attachments   []string
	parameters    map[string]interface{}
	message       string
	trace         string
}

func (b *before) getSteps() []step {
	return b.childrenSteps
}

func (b *before) addStep(step step) {
	b.childrenSteps = append(b.childrenSteps, step)
}

func (b *before) addStatus(status string) {
	b.status = status
}

func (b *before) addAttachments(a string) {
	b.attachments = append(b.attachments, a)
}

func BeforeTest(t *testing.T, m StepMetadata, f func()) {
	testPhaseObject := getCurrentTestPhaseObject(t)
	if testPhaseObject.test != nil {
		logger.Error("Cannot add before to test after test has been started")
	}
	before := newBefore(m)
	testPhaseObject.before = before

	defer func() {
		panicObject := recover()

		before.completedOn = time.Now()
		before.status = getTestStatus(t)
		before.duration = before.completedOn.UnixMilli() - before.startedOn.UnixMilli()

		if panicObject != nil {
			t.Fail()
			before.message = fmt.Sprintf("%+v", panicObject)
			before.trace = string(debug.Stack())
			before.status = failed

			panic(panicObject)
		}
	}()
	ctxMgr.SetValues(gls.Values{
		testResultKey:   before,
		nodeKey:         before,
		testInstanceKey: t,
	}, f)
}

func newBefore(m StepMetadata) *before {
	before := &before{
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
