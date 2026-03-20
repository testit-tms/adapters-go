package syncstorage

// RegisterRequest represents the worker registration request.
type RegisterRequest struct {
	PID       string `json:"pid"`
	TestRunID string `json:"testRunId"`
}

// RegisterResponse represents the worker registration response.
type RegisterResponse struct {
	IsMaster bool `json:"is_master"`
}

// SetWorkerStatusRequest represents a request to set worker status.
type SetWorkerStatusRequest struct {
	PID       string `json:"pid"`
	Status    string `json:"status"`
	TestRunID string `json:"testRunId"`
}

// TestResultCutModel represents a lightweight test result sent to sync storage.
type TestResultCutModel struct {
	AutoTestExternalID string `json:"auto_test_external_id"`
	StatusCode         string `json:"status_code"`
	StartedOn          string `json:"started_on,omitempty"`
}
