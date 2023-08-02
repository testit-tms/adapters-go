package examples

import (
	"testing"

	"github.com/testit-tms/adapters-go/pkg/tms"
)

func TestParameters_success(t *testing.T) {
	tests := []struct {
		name           string
		parameters     map[string]interface{}
		stepName       string
		stepParameters map[string]interface{}
		expValue       bool
	}{
		{
			name: "add parameters success",
			parameters: map[string]interface{}{
				"param1": "value1",
				"param2": 15,
			},
			stepName: "step1",
			stepParameters: map[string]interface{}{
				"param1": "value1",
				"param2": 15,
			},
			expValue: true,
		},
		{
			name: "add parameters failed",
			parameters: map[string]interface{}{
				"param1": "value1",
				"param2": 15,
			},
			stepName: "step1",
			stepParameters: map[string]interface{}{
				"param1": "value1",
				"param2": 15,
			},
			expValue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tms.Test(t, tms.TestMetadata{
				DisplayName: tt.name,
				Parameters:  tt.parameters,
			}, func() {
				tms.Step(
					tms.StepMetadata{
						Name:       tt.stepName,
						Parameters: tt.stepParameters,
					}, func() {
						tms.True(t, tt.expValue)
					})
			})
		})
	}
}
