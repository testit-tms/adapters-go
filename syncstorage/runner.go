package syncstorage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

const (
	syncStorageVersion        = "v0.1.18"
	syncStorageRepoURL        = "https://github.com/testit-tms/sync-storage-public/releases/download/"
	syncStorageStartupTimeout = 30 * time.Second
	defaultPort               = "49152"
)

// Runner manages the sync-storage process lifecycle and worker coordination.
type Runner struct {
	TestRunID    string
	Port         string
	BaseURL      string
	PrivateToken string

	workerPID           string
	isMaster            bool
	isAlreadyInProgress bool
	isRunning           bool
	isExternal          bool

	process *os.Process
	client  *Client
	logger  *slog.Logger
	mu      sync.Mutex
}

// NewRunner creates a new SyncStorageRunner.
func NewRunner(testRunID, port, baseURL, privateToken string, logger *slog.Logger) *Runner {
	if port == "" {
		port = defaultPort
	}
	return &Runner{
		TestRunID:    testRunID,
		Port:         port,
		BaseURL:      baseURL,
		PrivateToken: privateToken,
		workerPID:    fmt.Sprintf("worker-%d-%d", os.Getpid(), time.Now().UnixMilli()),
		logger:       logger,
	}
}

// Start starts the sync-storage service and registers the worker.
func (r *Runner) Start() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info("Starting Sync Storage service")

	if r.isRunning {
		r.logger.Info("SyncStorage already running")
		return true
	}

	baseURL := fmt.Sprintf("http://localhost:%s", r.Port)
	r.client = NewClient(baseURL)

	// Check if already running externally.
	if r.isAlreadyRunning() {
		r.logger.Info("SyncStorage already started, connecting to existing one", "port", r.Port)
		r.isRunning = true
		r.isExternal = true
		r.registerWorker()
		return true
	}

	// Download and start the binary.
	execPath, err := r.prepareExecutable()
	if err != nil {
		r.logger.Error("Failed to prepare sync-storage executable", "error", err)
		return false
	}

	if err := r.startProcess(execPath); err != nil {
		r.logger.Error("Failed to start sync-storage process", "error", err)
		return false
	}

	if !r.waitForStartup(syncStorageStartupTimeout) {
		r.logger.Error("Sync Storage did not start within timeout")
		return false
	}

	r.isRunning = true
	r.logger.Info("SyncStorage started successfully", "port", r.Port)

	// Small delay like Java/Python implementations.
	time.Sleep(2 * time.Second)
	r.registerWorker()

	return true
}

// IsMaster returns whether this worker is the master coordinator.
func (r *Runner) IsMaster() bool {
	return r.isMaster
}

// IsRunning returns whether sync storage is running.
func (r *Runner) IsRunning() bool {
	return r.isRunning
}

// IsAlreadyInProgress returns the in-progress flag.
func (r *Runner) IsAlreadyInProgress() bool {
	return r.isAlreadyInProgress
}

// SetAlreadyInProgress sets the in-progress flag.
func (r *Runner) SetAlreadyInProgress(v bool) {
	r.isAlreadyInProgress = v
}

// WorkerPID returns the worker identifier.
func (r *Runner) WorkerPID() string {
	return r.workerPID
}

// SendInProgressTestResult sends a test result to sync storage (master only).
func (r *Runner) SendInProgressTestResult(model TestResultCutModel) bool {
	if !r.isMaster {
		r.logger.Debug("Not master worker, skipping sending test result to Sync Storage")
		return false
	}
	if r.isAlreadyInProgress {
		r.logger.Debug("Test already in progress, skipping duplicate send")
		return false
	}

	r.logger.Debug("Sending in-progress test result to Sync Storage")

	if r.client == nil {
		r.logger.Error("Sync storage client not initialized")
		return false
	}

	if err := r.client.SendInProgressTestResult(r.TestRunID, model); err != nil {
		r.logger.Warn("Failed to send test result to Sync Storage", "error", err)
		return false
	}

	r.isAlreadyInProgress = true
	r.logger.Debug("Successfully sent test result to Sync Storage")
	return true
}

// SetWorkerStatus sets the worker status (e.g. "in_progress", "completed").
func (r *Runner) SetWorkerStatus(status string) {
	if !r.isRunning || r.client == nil {
		return
	}

	req := SetWorkerStatusRequest{
		PID:       r.workerPID,
		Status:    status,
		TestRunID: r.TestRunID,
	}
	if err := r.client.SetWorkerStatus(req); err != nil {
		r.logger.Error("Error setting worker status", "error", err)
		return
	}
	r.logger.Info("Successfully set worker status", "pid", r.workerPID, "status", status)
}

// --- internal helpers ---

func (r *Runner) isAlreadyRunning() bool {
	if r.client == nil {
		return false
	}
	return r.client.Health() == nil
}

func (r *Runner) registerWorker() {
	if r.client == nil {
		r.logger.Error("Cannot register worker: client not initialized")
		return
	}

	resp, err := r.client.Register(RegisterRequest{
		PID:       r.workerPID,
		TestRunID: r.TestRunID,
	})
	if err != nil {
		r.logger.Error("Error registering worker", "error", err)
		return
	}

	r.isMaster = resp.IsMaster
	if r.isMaster {
		r.logger.Info("Master worker registered", "pid", r.workerPID)
	} else {
		r.logger.Info("Worker registered", "pid", r.workerPID)
	}
}

func (r *Runner) waitForStartup(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if r.isAlreadyRunning() {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

func (r *Runner) startProcess(execPath string) error {
	args := []string{execPath}
	if r.TestRunID != "" {
		args = append(args, "--testRunId", r.TestRunID)
	}
	if r.Port != "" {
		args = append(args, "--port", r.Port)
	}
	if r.BaseURL != "" {
		args = append(args, "--baseURL", r.BaseURL)
	}
	if r.PrivateToken != "" {
		args = append(args, "--privateToken", r.PrivateToken)
	}

	r.logger.Info("Starting SyncStorage process", "command", args)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = filepath.Dir(execPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start sync-storage: %w", err)
	}

	r.process = cmd.Process
	return nil
}

func (r *Runner) prepareExecutable() (string, error) {
	fileName := r.executableFileName()
	cacheDir := filepath.Join("build", ".caches")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", fmt.Errorf("create cache dir: %w", err)
	}

	targetPath := filepath.Join(cacheDir, fileName)

	if _, err := os.Stat(targetPath); err == nil {
		r.logger.Info("Using existing sync-storage executable", "path", targetPath)
		_ = os.Chmod(targetPath, 0o755)
		return targetPath, nil
	}

	r.logger.Info("Downloading sync-storage from GitHub Releases")
	downloadURL := r.downloadURL(fileName)
	if err := r.downloadFile(downloadURL, targetPath); err != nil {
		return "", err
	}
	return targetPath, nil
}

func (r *Runner) executableFileName() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	var osPart string
	switch osName {
	case "windows":
		osPart = "windows"
	case "darwin":
		osPart = "darwin"
	default:
		osPart = "linux"
	}

	var archPart string
	switch arch {
	case "arm64", "aarch64":
		archPart = "arm64"
	default:
		archPart = "amd64"
	}

	name := fmt.Sprintf("syncstorage-%s-%s_%s", syncStorageVersion, osPart, archPart)
	if osPart == "windows" {
		name += ".exe"
	}
	return name
}

func (r *Runner) downloadURL(fileName string) string {
	return fmt.Sprintf("%s%s/%s", syncStorageRepoURL, syncStorageVersion, fileName)
}

func (r *Runner) downloadFile(url, targetPath string) error {
	r.logger.Info("Downloading file", "url", url, "target", targetPath)

	resp, err := http.Get(url) //nolint:gosec // URL is constructed from known constants
	if err != nil {
		return fmt.Errorf("download sync-storage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download sync-storage: HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(targetPath) //nolint:gosec // path is constructed internally
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	if runtime.GOOS != "windows" {
		_ = os.Chmod(targetPath, 0o755)
	}

	r.logger.Info("File downloaded successfully", "path", targetPath)
	return nil
}
