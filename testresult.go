package tms

import (
	"time"

	"golang.org/x/exp/slog"
)

type testResult struct {
	externalId  string
	displayName string
	title       string
	description string
	labels      []string
	className   string
	nameSpace   string
	setups      []stepresult
	steps       []stepresult
	teardowns   []stepresult
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
	externalKey string
}

func (tr *testResult) addStatus(v string) {
	tr.status = v
}

func (tr *testResult) addStep(step stepresult) {
	tr.steps = append(tr.steps, step)
}

func (tr *testResult) addBefore(step stepresult) {
	tr.setups = append(tr.setups, step)
}

func (tr *testResult) addAfter(step stepresult) {
	tr.teardowns = append(tr.teardowns, step)
}

func (tr *testResult) getSteps() []stepresult {
	return tr.steps
}

func (tr *testResult) addAttachments(a string) {
	tr.attachments = append(tr.attachments, a)
}

func (tr *testResult) addMessage(message string) {
	tr.message = message
}

func (tr *testResult) addTrace(trace string) {
	tr.trace = trace
}

func (tr *testResult) write() string {
	const op = "tms.testResult.write"
	id, err := client.writeTest(*tr)
	if err != nil {
		logger.Error("error writing test result", "error", err, slog.String("op", op))
	}

	return id
}

func (tr *testResult) update(resultID string) {
	const op = "tms.testResult.update"
	err := client.updateTest(*tr)
	if err != nil {
		logger.Error("failed to update test", "error", err, slog.String("op", op))
	}

	//
	err = client.updateTestResult(resultID, *tr)
	if err != nil {
		logger.Error("failed to update test result", "error", err, slog.String("op", op))
	}
}
