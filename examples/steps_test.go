package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testit-tms/adapters-go/pkg/tms"
)

func TestSteps_Success(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "steps success",
		},
		func() {
			tms.Step(
				tms.StepMetadata{
					Name:        "step 1",
					Description: "step 1 description",
				},
				func() {
					tms.Step(tms.StepMetadata{
						Name:        "step 1.1",
						Description: "step 1.1 description",
					}, func() {
						tms.Step(tms.StepMetadata{}, func() {
							assert.True(t, true)
						})
						assert.True(t, true)
					})
					assert.True(t, true)
				},
			)
			tms.Step(
				tms.StepMetadata{
					Name:        "step 2",
					Description: "step 2 description",
				},
				func() {
					tms.Step(tms.StepMetadata{
						Name:        "step 2.1",
						Description: "step 2.1 description",
					}, func() {
						tms.Step(tms.StepMetadata{}, func() {
							assert.True(t, true)
						})
						assert.True(t, true)
					})
					assert.True(t, true)
				},
			)
		})
}

func TestSteps_Failed(t *testing.T) {
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "steps failed",
		},
		func() {
			tms.Step(
				tms.StepMetadata{
					Name:        "step 1",
					Description: "step 1 description",
				},
				func() {
					tms.Step(tms.StepMetadata{
						Name:        "step 1.1",
						Description: "step 1.1 description",
					}, func() {
						tms.Step(tms.StepMetadata{}, func() {
							assert.True(t, false)
						})
						assert.True(t, true)
					})
					assert.True(t, true)
				},
			)
			tms.Step(
				tms.StepMetadata{
					Name:        "step 2",
					Description: "step 2 description",
				},
				func() {
					tms.Step(tms.StepMetadata{
						Name:        "step 2.1",
						Description: "step 2.1 description",
					}, func() {
						tms.Step(tms.StepMetadata{}, func() {
							assert.True(t, true)
						})
						assert.True(t, false)
					})
					assert.True(t, true)
				},
			)
		})
}
