package tms

import (
	"time"

	"github.com/jtolds/gls"
	"github.com/pkg/errors"
)

type StepMetadata struct {
	Name        string
	Description string
	Parameters  map[string]interface{}
}

func Step(m StepMetadata, f func()) {
	step := newStep(m)

	defer func() {
		panicObject := recover()
		step.completedOn = time.Now()
		step.duration = step.completedOn.UnixMilli() - step.startedOn.UnixMilli()
		manipulateOnObjectFromCtx(
			testInstanceKey,
			func(testInstance interface{}) {
				if panicObject != nil {
					fail(errors.Errorf("%+v", panicObject))
				}
			})
		if step.status == "" {
			step.status = passed
		}
		manipulateOnObjectFromCtx(nodeKey, func(currentStepObj interface{}) {
			hasStep := currentStepObj.(hasSteps)
			hasStep.addStep(*step)

			hasStatus := currentStepObj.(hasStatus)
			hasStatus.addStatus(step.status)
		})
	}()

	ctxMgr.SetValues(gls.Values{nodeKey: step}, f)
}

func newStep(m StepMetadata) *stepresult {
	step := &stepresult{
		description: m.Description,
		startedOn:   time.Now(),
		parameters:  m.Parameters,
		name:        m.Name,
	}

	if step.name == "" {
		step.name = "Step"
	}

	return step
}
