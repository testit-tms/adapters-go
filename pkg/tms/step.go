package tms

import (
	"testing"
	"time"

	"github.com/jtolds/gls"
	"github.com/pkg/errors"
)

// TODO: rename to stepResult and move to separate file
type step struct {
	name          string
	description   string
	childrenSteps []step
	status        string
	startedOn     time.Time
	completedOn   time.Time
	duration      int64
}

type StepMetadata struct {
	Name        string
	Description string
}

func (s *step) getSteps() []step {
	return s.childrenSteps
}

func (s *step) addStep(step step) {
	s.childrenSteps = append(s.childrenSteps, step)
}

func (s *step) addStatus(status string) {
	s.status = status
}

// TODO: try to use StepMetadata as a pointer
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
					Fail(errors.Errorf("%+v", panicObject))
				}
				if testInstance.(*testing.T).Failed() ||
					panicObject != nil {
					if step.status == ""  {
						step.status = failed
					}
				}
			})
		if step.status == "" {
			step.status = passed
		}
		manipulateOnObjectFromCtx(nodeKey, func(currentStepObj interface{}) {
			currentStep := currentStepObj.(hasSteps)
			currentStep.addStep(*step)
		})

		if panicObject != nil {
			panic(panicObject)
		}
	}()

	ctxMgr.SetValues(gls.Values{nodeKey: step}, f)
}

func newStep(m StepMetadata) *step {
	step := &step{
		description: m.Description,
		startedOn:   time.Now(),
	}

	if m.Name != "" {
		step.name = m.Name
	} else {
		step.name = "Step"
	}

	return step
}
