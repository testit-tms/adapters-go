package tms

import (
	"time"

	"github.com/testit-tms/adapters-go/v2/models"
)

type StepResult struct {
	Name          string
	Description   string
	ChildrenSteps []StepResult
	Status        string
	StartedOn     time.Time
	CompletedOn   time.Time
	Duration      int64
	Attachments   []string
	Parameters    map[string]interface{}
}

func (s *StepResult) getSteps() []StepResult {
	return s.ChildrenSteps
}

func (s *StepResult) addStep(step StepResult) {
	s.ChildrenSteps = append(s.ChildrenSteps, step)
}

func (s *StepResult) addStatus(status string) {
	s.Status = models.MergeStatus(s.Status, status)
}

func (s *StepResult) addAttachments(a string) {
	s.Attachments = append(s.Attachments, a)
}
