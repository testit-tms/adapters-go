package syncstorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a simple HTTP client for the Sync Storage service.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Sync Storage client.
func NewClient(port string) *Client {
	return &Client{
		baseURL: fmt.Sprintf("http://localhost:%s", port),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RegisterRequest represents a worker registration request.
type RegisterRequest struct {
	PID       string `json:"pid"`
	TestRunID string `json:"testRunId"`
}

// RegisterResponse represents a worker registration response.
type RegisterResponse struct {
	IsMaster bool `json:"is_master"`
}

// SetWorkerStatusRequest represents a request to set worker status.
type SetWorkerStatusRequest struct {
	PID       string `json:"pid"`
	Status    string `json:"status"`
	TestRunID string `json:"testRunId"`
}

// TestResultCutModel represents a cut-down test result for sync storage.
type TestResultCutModel struct {
	AutoTestExternalID string `json:"autoTestExternalId"`
	StatusCode         string `json:"statusCode"`
	StartedOn          string `json:"startedOn"`
}

// HealthCheck checks if the sync storage service is running.
func (c *Client) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}
	return nil
}

// Register registers a worker with the sync storage service.
func (c *Client) Register(req RegisterRequest) (*RegisterResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal register request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to register worker: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("register failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var result RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode register response: %w", err)
	}

	return &result, nil
}

// SetWorkerStatus sets the status of a worker.
func (c *Client) SetWorkerStatus(req SetWorkerStatusRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal set worker status request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/set_worker_status", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to set worker status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("set worker status failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// SendInProgressTestResult sends an in-progress test result to sync storage.
func (c *Client) SendInProgressTestResult(testRunID string, model TestResultCutModel) error {
	body, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("failed to marshal test result: %w", err)
	}

	url := fmt.Sprintf("%s/in_progress_test_result?testRunId=%s", c.baseURL, testRunID)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send in-progress test result: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("send in-progress test result failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
