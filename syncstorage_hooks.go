package tms

// OnRunningStarted should be called when a test execution block starts.
// It sets the worker status to "in_progress" in Sync Storage.
func OnRunningStarted() {
	if !isSyncStorageActive() {
		return
	}
	syncRunner.SetWorkerStatus("in_progress")
}

// OnBlockCompleted should be called when a test execution block finishes.
// It sets the worker status to "completed" in Sync Storage.
func OnBlockCompleted() {
	if !isSyncStorageActive() {
		return
	}
	syncRunner.SetWorkerStatus("completed")
}

// ResetInProgressFlag resets the sync-storage in-progress flag.
// Should be called at the start of each new test to allow a new test result
// to be sent to sync-storage.
func ResetInProgressFlag() {
	if !isSyncStorageActive() {
		return
	}
	syncRunner.SetAlreadyInProgress(false)
}
