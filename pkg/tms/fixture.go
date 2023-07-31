package tms

import "time"

type fixture struct {
	name          string
	description   string
	childrenSteps []step
	status        string
	startedOn     time.Time
	completedOn   time.Time
	duration      int64
	attachments   []string
	parameters    map[string]interface{}
	message       string
	trace         string
}

func (b *fixture) getSteps() []step {
	return b.childrenSteps
}

func (b *fixture) addStep(step step) {
	b.childrenSteps = append(b.childrenSteps, step)
}

func (b *fixture) addStatus(status string) {
	b.status = status
}

func (b *fixture) addAttachments(a string) {
	b.attachments = append(b.attachments, a)
}
