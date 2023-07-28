package tms

type hasSteps interface {
	getSteps() []step
	addStep(step step)
}

type hasStatus interface {
	addStatus(status string)
}

type hasAttachments interface {
	addAttachments(a string)
}