package syncstorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a simple HTTP client for the sync-storage API.
// Replaces the auto-generated api_client_syncstorage from the Python adapter
// with direct HTTP calls to the endpoints actually used by the runner.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new sync-storage API client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Health checks if the sync-storage service is running.
func (c *Client) Health() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed: status %d", resp.StatusCode)
	}
	return nil
}

// Register registers a worker with the sync-storage service.
func (c *Client) Register(req RegisterRequest) (*RegisterResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal register request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("register worker: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("register worker: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode register response: %w", err)
	}
	return &result, nil
}

// SendInProgressTestResult sends a test result to sync storage.
func (c *Client) SendInProgressTestResult(testRunID string, model TestResultCutModel) error {
	body, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("marshal test result: %w", err)
	}

	url := fmt.Sprintf("%s/in_progress_test_result?testRunId=%s", c.baseURL, testRunID)
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("send in-progress test result: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("send in-progress test result: status %d, body: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

// SetWorkerStatus sets the worker status in sync storage.
func (c *Client) SetWorkerStatus(req SetWorkerStatusRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal set worker status request: %w", err)
	}

	resp, err := c.httpClient.Post(c.baseURL+"/set_worker_status", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("set worker status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("set worker status: status %d, body: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
