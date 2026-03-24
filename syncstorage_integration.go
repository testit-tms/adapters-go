package tms

// onRunningStarted notifies sync-storage that test execution has started.
func onRunningStarted() {
	if syncStorageRunner == nil || !syncStorageRunner.IsRunning() {
		return
	}
	syncStorageRunner.SetWorkerStatus("in_progress")
}

// onBlockCompleted notifies sync-storage that a test block has completed.
func onBlockCompleted() {
	if syncStorageRunner == nil || !syncStorageRunner.IsRunning() {
		return
	}
	syncStorageRunner.SetWorkerStatus("completed")
}
