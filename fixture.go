package tms

import (
	"time"
)

type fixture struct {
	name          string
	description   string
	childrenSteps []StepResult
	status        string
	startedOn     time.Time
	completedOn   time.Time
	duration      int64
	attachments   []string
	parameters    map[string]interface{}
	message       string
	trace         string
}

func (b *fixture) getSteps() []StepResult {
	return b.childrenSteps
}

func (b *fixture) addStep(step StepResult) {
	b.childrenSteps = append(b.childrenSteps, step)
}

func (b *fixture) addStatus(status string) {
	b.status = status
}

func (b *fixture) addAttachments(a string) {
	b.attachments = append(b.attachments, a)
}

func (b *fixture) addMessage(message string) {
	b.message = message
}

func (b *fixture) addTrace(trace string) {
	b.trace = trace
}

func (b *fixture) convertToStepResult() StepResult {
	return StepResult{
		Name:          b.name,
		Description:   b.description,
		ChildrenSteps: b.childrenSteps,
		Status:        b.status,
		StartedOn:     b.startedOn,
		CompletedOn:   b.completedOn,
		Duration:      b.duration,
		Attachments:   b.attachments,
		Parameters:    b.parameters,
	}
}
