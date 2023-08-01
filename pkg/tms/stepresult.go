package tms

import "time"

type stepresult struct {
	name          string
	description   string
	childrenSteps []stepresult
	status        string
	startedOn     time.Time
	completedOn   time.Time
	duration      int64
	attachments   []string
	parameters    map[string]interface{}
}

func (s *stepresult) getSteps() []stepresult {
	return s.childrenSteps
}

func (s *stepresult) addStep(step stepresult) {
	s.childrenSteps = append(s.childrenSteps, step)
}

func (s *stepresult) addStatus(status string) {
	s.status = status
}

func (s *stepresult) addAttachments(a string) {
	s.attachments = append(s.attachments, a)
}
