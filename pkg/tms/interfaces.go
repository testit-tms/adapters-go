package tms

type hasSteps interface {
	getSteps() []stepresult
	addStep(step stepresult)
}

type hasStatus interface {
	addStatus(status string)
}

type hasAttachments interface {
	addAttachments(a string)
}

type hasErrorFields interface {
	addMessage(message string)
	addTrace(trace string)
	addStatus(status string)
}