package tms

import (
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
			logger.Error("Cannot add after to test before test has been started")
		}

		testPhaseObject.after = after

		if testPhaseObject.after != nil {
			testPhaseObject.test.addAfter(step{
				name:          testPhaseObject.after.name,
				description:   testPhaseObject.after.description,
				status:        testPhaseObject.after.status,
				startedOn:     testPhaseObject.after.startedOn,
				completedOn:   testPhaseObject.after.completedOn,
				duration:      testPhaseObject.after.duration,
				attachments:   testPhaseObject.after.attachments,
				parameters:    testPhaseObject.after.parameters,
				childrenSteps: testPhaseObject.after.childrenSteps,
			})

			err := client.updateTest(*testPhaseObject.test)
			if err != nil {
				logger.Error("Failed to update test: %s", err)
			}

			err = client.updateTestResult(testPhaseObject.resultID, *testPhaseObject.test)		
			if err != nil {
				logger.Error("Failed to update test result: %s", err)
			}
		}
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
