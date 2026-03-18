package tms

import (
	"time"

	"github.com/jtolds/gls"
	"github.com/pkg/errors"
	"github.com/testit-tms/adapters-go/models"
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
		step.CompletedOn = time.Now()
		step.Duration = step.CompletedOn.UnixMilli() - step.StartedOn.UnixMilli()
		manipulateOnObjectFromCtx(
			testInstanceKey,
			func(testInstance interface{}) {
				if panicObject != nil {
					fail(errors.Errorf("%+v", panicObject))
				}
			})
		if step.Status == "" {
			step.Status = models.Passed
		}
		manipulateOnObjectFromCtx(nodeKey, func(currentStepObj interface{}) {
			hasStep := currentStepObj.(hasSteps)
			hasStep.addStep(*step)

			hasStatus := currentStepObj.(hasStatus)
			hasStatus.addStatus(step.Status)
		})
	}()

	ctxMgr.SetValues(gls.Values{nodeKey: step}, f)
}

func newStep(m StepMetadata) *StepResult {
	step := &StepResult{
		Description: m.Description,
		StartedOn:   time.Now(),
		Parameters:  m.Parameters,
		Name:        m.Name,
	}

	if step.Name == "" {
		step.Name = "Step"
	}

	return step
}
