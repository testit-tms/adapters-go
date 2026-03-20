package tms

import (
	"time"

	"github.com/testit-tms/adapters-go/models"
	"github.com/testit-tms/adapters-go/syncstorage"
	"golang.org/x/exp/slog"
)

const inProgressLiteral = "InProgress"

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

	// Sync Storage integration: if master and no in-progress → send to sync-storage
	if isSyncStorageActive() && isMasterAndNoInProgress() {
		ok := tr.onMasterNoAlreadyInProgress()
		if ok {
			return ""
		}
		// Fallback to normal processing on failure.
	}

	id, err := client.writeTest(*tr)
	if err != nil {
		logger.Error("error writing test result", "error", err, slog.String("op", op))
	}

	return id
}

// onMasterNoAlreadyInProgress handles the sync-storage master logic:
// send a cut model to sync-storage, mark as InProgress, write to TestIT.
func (tr *TestResult) onMasterNoAlreadyInProgress() bool {
	const op = "tms.TestResult.onMasterNoAlreadyInProgress"

	cutModel := syncstorage.TestResultCutModel{
		AutoTestExternalID: tr.externalId,
		StatusCode:         tr.status,
		StartedOn:          tr.startedOn.Format(time.RFC3339),
	}

	logger.Debug("Sending in-progress test result to Sync Storage",
		"externalId", tr.externalId, "status", tr.status, slog.String("op", op))

	if !syncRunner.SendInProgressTestResult(cutModel) {
		return false
	}

	// Write the test result with InProgress status to TestIT.
	originalStatus := tr.status
	tr.status = inProgressLiteral

	id, err := client.writeTest(*tr)
	if err != nil {
		logger.Warn("Error writing InProgress test result, falling back to normal processing",
			"error", err, slog.String("op", op))
		tr.status = originalStatus
		syncRunner.SetAlreadyInProgress(false)
		return false
	}

	_ = id
	return true
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
		logger.Error("failed to update test result", "error", err, slog.String("op", op))
	}
}

// isSyncStorageActive checks if sync-storage runner is initialized and running.
func isSyncStorageActive() bool {
	return syncRunner != nil && syncRunner.IsRunning()
}

// isMasterAndNoInProgress checks if current worker is master and no test is in progress.
func isMasterAndNoInProgress() bool {
	return syncRunner.IsMaster() && !syncRunner.IsAlreadyInProgress()
}
