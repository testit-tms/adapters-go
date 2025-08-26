package examples

import (
	"testing"

	tms "github.com/testit-tms/adapters-go"
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
					tms.True(t, true)
				})
			tms.True(t, true)
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
					tms.True(t, true)
				})
			tms.True(t, true)
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
					tms.True(t, true)
				})
			tms.True(t, true)
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
					tms.True(t, false)
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
					tms.True(t, true)
				})
			tms.True(t, true)
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
					tms.True(t, true)
				})
		},
	)
}
