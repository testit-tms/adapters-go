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
	befores     []step
	steps       []step
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

func (tr *testResult) addBefore(step step) {
	tr.befores = append(tr.befores, step)
}

func (tr *testResult) getSteps() []step {
	return tr.steps
}

func (tr *testResult) addAttachments(a string) {
	tr.attachments = append(tr.attachments, a)
}
