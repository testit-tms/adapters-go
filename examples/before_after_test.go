package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testit-tms/adapters-go/pkg/tms"
)

func TestFixture_success(t *testing.T) {
	tms.BeforeTest(t,
		tms.StepMetadata{
			Name:        "before test",
			Description: "before test description",
		},
		func() {
			tms.Step(
				tms.StepMetadata{
					Name:        "step1",
					Description: "step1 description",
				}, func() {
					assert.True(t, true)
				})
		})
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "fixture success",
		},
		func() {
			tms.Step(
				tms.StepMetadata{
					Name:        "step1",
					Description: "step1 description",
				}, func() {
					assert.True(t, true)
				})
			assert.True(t, true)
		})
	tms.AfterTest(t,
		tms.StepMetadata{
			Name:        "after test",
			Description: "after test description",
		},
		func() {
			tms.Step(
				tms.StepMetadata{
					Name:        "step1",
					Description: "step1 description",
				}, func() {
					assert.True(t, true)
				})
		},
	)
}

func TestFixture_failed(t *testing.T) {
	tms.BeforeTest(t,
		tms.StepMetadata{
			Name:        "before test",
			Description: "before test description",
		},
		func() {
			tms.Step(
				tms.StepMetadata{
					Name:        "step1",
					Description: "step1 description",
				}, func() {
					assert.True(t, false)
				})
		})
	tms.Test(t,
		tms.TestMetadata{
			DisplayName: "fixture failed",
		},
		func() {
			tms.Step(
				tms.StepMetadata{
					Name:        "step1",
					Description: "step1 description",
				}, func() {
					assert.True(t, true)
				})
			assert.True(t, true)
		})
	tms.AfterTest(t,
		tms.StepMetadata{
			Name:        "after test",
			Description: "after test description",
		},
		func() {
			tms.Step(
				tms.StepMetadata{
					Name:        "step1",
					Description: "step1 description",
				}, func() {
					assert.True(t, true)
				})
		},
	)
}
