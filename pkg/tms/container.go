package tms

import "testing"

type testPhaseContainer struct {
	before *before
	test    *testResult
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
