package syncstorage

import (
	"fmt"
	"io"
	"os"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

const (
	syncStorageVersion        = "v0.3.0"
	syncStorageRepoURL        = "https://github.com/testit-tms/sync-storage-public/releases/download/"
	defaultPort               = "49152"
	startupTimeout            = 30 * time.Second
	startupCheckInterval      = 1 * time.Second
	postStartupDelay          = 2 * time.Second
)

// Runner manages the lifecycle of the Sync Storage process and worker coordination.
type Runner struct {
	testRunID    string
	port         string
	baseURL      string
	privateToken string

	workerPID           string
	isMaster            bool
	isAlreadyInProgress bool
	isRunning           bool
	isExternal          bool

	process *exec.Cmd
	client  *Client
	logger  *slog.Logger
	mu      sync.Mutex
}

// NewRunner creates a new SyncStorage runner.
func NewRunner(testRunID, port, baseURL, privateToken string, logger *slog.Logger) *Runner {
	if port == "" {
		port = defaultPort
	}

	workerPID := fmt.Sprintf("worker-%d-%d", os.Getpid(), time.Now().UnixMilli())

	return &Runner{
		testRunID:    testRunID,
		port:         port,
		baseURL:      baseURL,
		privateToken: privateToken,
		workerPID:    workerPID,
		client:       NewClient(port),
		logger:       logger,
	}
}

// Start starts the Sync Storage service and registers the worker.
func (r *Runner) Start() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info("Starting Sync Storage service")

	if r.isRunning {
		r.logger.Info("SyncStorage already running")
		return true
	}

	// Check if already running externally
	if r.client.HealthCheck() == nil {
		r.logger.Info("SyncStorage already started externally, connecting")
		r.isRunning = true
		r.isExternal = true
		r.registerWorker()
		return true
	}

	// Download and start sync-storage
	executablePath, err := r.prepareExecutable()
	if err != nil {
		r.logger.Error("Failed to prepare sync-storage executable", "error", err)
		return false
	}

	args := []string{}
	if r.testRunID != "" {
		args = append(args, "--testRunId", r.testRunID)
	}
	if r.port != "" {
		args = append(args, "--port", r.port)
	}
	if r.baseURL != "" {
		args = append(args, "--baseURL", r.baseURL)
	}
	if r.privateToken != "" {
		args = append(args, "--privateToken", r.privateToken)
	}

	r.logger.Info("Starting SyncStorage process", "executable", executablePath, "args", args)

	r.process = exec.Command(executablePath, args...)
	r.process.Dir = filepath.Dir(executablePath)

	// Capture stdout/stderr for logging
	stdout, err := r.process.StdoutPipe()
	if err != nil {
		r.logger.Error("Failed to create stdout pipe", "error", err)
		return false
	}
	r.process.Stderr = r.process.Stdout

	if err := r.process.Start(); err != nil {
		r.logger.Error("Failed to start SyncStorage process", "error", err)
		return false
	}

	// Read output in background
	go r.readOutput(stdout)

	// Wait for startup
	if !r.waitForStartup() {
		r.logger.Error("SyncStorage failed to start within timeout")
		return false
	}

	r.isRunning = true
	r.logger.Info("SyncStorage started successfully", "port", r.port)

	time.Sleep(postStartupDelay)
	r.registerWorker()

	return true
}

// IsMaster returns whether this worker is the master.
func (r *Runner) IsMaster() bool {
	return r.isMaster
}

// IsAlreadyInProgress returns the in-progress flag state.
func (r *Runner) IsAlreadyInProgress() bool {
	return r.isAlreadyInProgress
}

// SetAlreadyInProgress sets the in-progress flag.
func (r *Runner) SetAlreadyInProgress(v bool) {
	r.isAlreadyInProgress = v
}

// IsRunning returns whether sync storage is running.
func (r *Runner) IsRunning() bool {
	return r.isRunning
}

// TestRunID returns the test run ID.
func (r *Runner) TestRunID() string {
	return r.testRunID
}

// SetTestRunID updates the test run ID.
func (r *Runner) SetTestRunID(id string) {
	r.testRunID = id
}

// SendInProgressTestResult sends test result to sync storage if this worker is master.
func (r *Runner) SendInProgressTestResult(projectID, externalID, statusCode, statusType, startedOn string) bool {
	if !r.isMaster {
		r.logger.Debug("Not master worker, skipping send to sync storage")
		return false
	}

	if r.isAlreadyInProgress {
		r.logger.Debug("Test already in progress, skipping duplicate send")
		return false
	}

	model := TestResultCutModel{
		ProjectID:          projectID,
		AutoTestExternalID: externalID,
		StatusCode:         statusCode,
		StatusType:         statusType,
		StartedOn:          startedOn,
	}

	if err := r.client.SendInProgressTestResult(r.testRunID, model); err != nil {
		r.logger.Warn("Failed to send test result to sync storage", "error", err)
		return false
	}

	r.isAlreadyInProgress = true
	r.logger.Debug("Successfully sent test result to sync storage")
	return true
}

// SetWorkerStatus sets the worker status.
func (r *Runner) SetWorkerStatus(status string) {
	if !r.isRunning {
		return
	}

	req := SetWorkerStatusRequest{
		PID:       r.workerPID,
		Status:    status,
		TestRunID: r.testRunID,
	}

	if err := r.client.SetWorkerStatus(req); err != nil {
		r.logger.Error("Failed to set worker status", "error", err)
	} else {
		r.logger.Info("Successfully set worker status", "status", status)
	}
}

func (r *Runner) registerWorker() {
	req := RegisterRequest{
		PID:       r.workerPID,
		TestRunID: r.testRunID,
	}

	resp, err := r.client.Register(req)
	if err != nil {
		r.logger.Error("Failed to register worker", "error", err)
		return
	}

	r.isMaster = resp.IsMaster
	if r.isMaster {
		r.logger.Info("Master worker registered", "pid", r.workerPID)
	} else {
		r.logger.Info("Worker registered", "pid", r.workerPID)
	}
}

func (r *Runner) waitForStartup() bool {
	deadline := time.Now().Add(startupTimeout)
	for time.Now().Before(deadline) {
		if r.client.HealthCheck() == nil {
			return true
		}
		time.Sleep(startupCheckInterval)
	}
	return false
}

func (r *Runner) readOutput(reader io.Reader) {
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			r.logger.Info("SyncStorage", "output", strings.TrimSpace(string(buf[:n])))
		}
		if err != nil {
			break
		}
	}
}

func (r *Runner) prepareExecutable() (string, error) {
	fileName := getExecutableFileName()
	cacheDir := filepath.Join("build", ".caches")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create cache dir: %w", err)
	}

	targetPath := filepath.Join(cacheDir, fileName)

	if _, err := os.Stat(targetPath); err == nil {
		r.logger.Info("Using existing sync-storage executable", "path", targetPath)
		if runtime.GOOS != "windows" {
			os.Chmod(targetPath, 0o755)
		}
		absPath, _ := filepath.Abs(targetPath)
		return absPath, nil
	}

	r.logger.Info("Downloading sync-storage executable")
	downloadURL := fmt.Sprintf("%s%s/%s", syncStorageRepoURL, syncStorageVersion, fileName)

	if err := downloadFile(downloadURL, targetPath); err != nil {
		return "", fmt.Errorf("failed to download sync-storage: %w", err)
	}

	if runtime.GOOS != "windows" {
		os.Chmod(targetPath, 0o755)
	}

	absPath, _ := filepath.Abs(targetPath)
	return absPath, nil
}

func getExecutableFileName() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	var osPart string
	switch osName {
	case "windows":
		osPart = "windows"
	case "darwin":
		osPart = "darwin"
	case "linux":
		osPart = "linux"
	default:
		panic(fmt.Sprintf("unsupported OS: %s", osName))
	}

	var archPart string
	switch arch {
	case "amd64":
		archPart = "amd64"
	case "arm64":
		archPart = "arm64"
	default:
		panic(fmt.Sprintf("unsupported architecture: %s", arch))
	}

	name := fmt.Sprintf("syncstorage-%s-%s_%s", syncStorageVersion, osPart, archPart)
	if osName == "windows" {
		name += ".exe"
	}
	return name
}

func downloadFile(url, targetPath string) error {
	resp, err := (&http.Client{Timeout: 60 * time.Second}).Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	f, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
