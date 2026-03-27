package tms

import (
	"time"

	"github.com/testit-tms/adapters-go/models"
	"golang.org/x/exp/slog"
)

type TestResult struct {
	externalId  string
	displayName string
	title       string
	description string
	labels      []string
	tags        []string
	className   string
	nameSpace   string
	setups      []StepResult
	steps       []StepResult
	teardowns   []StepResult
	links       []models.Link
	resultLinks []models.Link
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

func (tr *TestResult) addStatus(v string) {
	tr.status = v
}

func (tr *TestResult) addStep(step StepResult) {
	tr.steps = append(tr.steps, step)
}

func (tr *TestResult) addBefore(step StepResult) {
	tr.setups = append(tr.setups, step)
}

func (tr *TestResult) addAfter(step StepResult) {
	tr.teardowns = append(tr.teardowns, step)
}

func (tr *TestResult) getSteps() []StepResult {
	return tr.steps
}

func (tr *TestResult) addAttachments(a string) {
	tr.attachments = append(tr.attachments, a)
}

func (tr *TestResult) addMessage(message string) {
	tr.message = message
}

func (tr *TestResult) addTrace(trace string) {
	tr.trace = trace
}

func (tr *TestResult) write() string {
	const op = "tms.TestResult.write"

	// Sync Storage integration: if master and not already in progress,
	// send cut result to sync-storage and write with InProgress status
	if syncStorageRunner != nil && syncStorageRunner.IsRunning() &&
		syncStorageRunner.IsMaster() && !syncStorageRunner.IsAlreadyInProgress() {

		startedOnStr := tr.startedOn.UTC().Format(time.RFC3339)
		if syncStorageRunner.SendInProgressTestResult(tr.externalId, tr.status, startedOnStr) {
			// Write to TMS with InProgress status
			originalStatus := tr.status
			tr.status = "InProgress"
			id, err := client.writeTest(*tr)
			if err != nil {
				logger.Error("error writing in-progress test result", "error", err, slog.String("op", op))
				// Fallback: restore status and write normally
				tr.status = originalStatus
				syncStorageRunner.SetAlreadyInProgress(false)
			} else {
				return id
			}
		}
	}

	// If importRealtime is false, buffer the result for later flush
	if !cfg.ImportRealtime {
		logger.Debug("Buffering test result (importRealtime=false)",
			"externalId", tr.externalId,
			slog.String("op", op))
		addPendingResult(tr)
		return ""
	}

	// Normal realtime write path
	id, err := client.writeTest(*tr)
	if err != nil {
		logger.Error("error writing test result", "error", err, slog.String("op", op))
	}

	return id
}

func (tr *TestResult) update(resultID string) {
	const op = "tms.TestResult.update"
	err := client.updateTest(*tr)
	if err != nil {
		logger.Error("failed to update test", "error", err, slog.String("op", op))
	}

	//
	err = client.updateTestResult(resultID, *tr)
	if err != nil {
		logger.Error(
			"failed to update test result",
			"error", err,
			"resultId", resultID,
			"externalId", tr.externalId,
			slog.String("op", op),
		)
	}
}
