package tms

import (
	"log"
	"time"
)

type testResult struct {
	externalId  string
	displayName string
	title       string
	description string
	labels      []string
	className   string
	nameSpace   string
	steps       []step
	links       []Link
	resultLinks []Link
	workItemIds []string
	status      string
	message     string
	trace       string
	startedOn   time.Time
	completedOn time.Time
	duration    int64
}

func (tr *testResult) write() {
	err := client.writeTest(*tr)
	if err != nil {
		log.Printf("Error writing test result: %v", err)
	}
}

func (tr *testResult) addStatus(v string) {
	tr.status = v
}

func (tr *testResult) addStep(step step) {
	tr.steps = append(tr.steps, step)
}

func (tr *testResult) getSteps() []step {
	return tr.steps
}
