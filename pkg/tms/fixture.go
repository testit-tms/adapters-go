package tms

import "time"

type fixture struct {
	name          string
	description   string
	childrenSteps []stepresult
	status        string
	startedOn     time.Time
	completedOn   time.Time
	duration      int64
	attachments   []string
	parameters    map[string]interface{}
	message       string
	trace         string
}

func (b *fixture) getSteps() []stepresult {
	return b.childrenSteps
}

func (b *fixture) addStep(step stepresult) {
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

func (b *fixture) convertToStepResult() stepresult {
	return stepresult{
		name:          b.name,
		description:   b.description,
		childrenSteps: b.childrenSteps,
		status:        b.status,
		startedOn:     b.startedOn,
		completedOn:   b.completedOn,
		duration:      b.duration,
		attachments:   b.attachments,
		parameters:    b.parameters,
	}
}
