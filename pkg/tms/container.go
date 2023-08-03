package tms

import "testing"

type testPhaseContainer struct {
	before   *fixture
	test     *testResult
	resultID string
	after    *fixture
}

func getCurrentTestPhaseObject(t *testing.T) *testPhaseContainer {
	var currentPhaseObject *testPhaseContainer
	if phaseContainer, ok := testPhaseObjects[t.Name()]; ok {
		currentPhaseObject = phaseContainer
	} else {
		currentPhaseObject = &testPhaseContainer{}
		testPhaseObjects[t.Name()] = currentPhaseObject
	}

	return currentPhaseObject
}
