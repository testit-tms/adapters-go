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
	setups      []step
	steps       []step
	teardowns   []step
	links       []Link
	resultLinks []Link
	attachments []string
	workItemIds []string
	parameters  map[string]interface{}
	status      string
	message     string
	trace       string
	startedOn   time.Time
	completedOn time.Time
	duration    int64
}

func (tr *testResult) write() string {
	id, err := client.writeTest(*tr)
	if err != nil {
		log.Printf("Error writing test result: %v", err)
	}

	return id
}

func (tr *testResult) addStatus(v string) {
	tr.status = v
}

func (tr *testResult) addStep(step step) {
	tr.steps = append(tr.steps, step)
}

func (tr *testResult) addBefore(step step) {
	tr.setups = append(tr.setups, step)
}

func (tr *testResult) addAfter(step step) {
	tr.teardowns = append(tr.teardowns, step)
}

func (tr *testResult) getSteps() []step {
	return tr.steps
}

func (tr *testResult) addAttachments(a string) {
	tr.attachments = append(tr.attachments, a)
}
