package tms

import (
	"os"

	"github.com/jtolds/gls"
	"github.com/testit-tms/adapters-go/config"
	"github.com/testit-tms/adapters-go/syncstorage"

	"golang.org/x/exp/slog"
)

var (
	cfg              *config.Config
	client           *tmsClient
	logger           *slog.Logger
	ctxMgr           *gls.ContextManager
	testPhaseObjects map[string]*testPhaseContainer
	syncRunner       *syncstorage.Runner
)

const (
	nodeKey         = "current_step"
	testResultKey   = "current_result_object"
	testInstanceKey = "test_instance"
)

func init() {
	cfg = config.MustLoad()
	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: getLogLevel(cfg.IsDebug)}),
	)
	client = newClient(*cfg)
	if cfg.AdapterMode == "2" {
		callCreateTestRun(client, cfg)
	} else {
		client.updateTestRun()
	}
	ctxMgr = gls.NewContextManager()
	testPhaseObjects = make(map[string]*testPhaseContainer)

	// Initialize Sync Storage
	initSyncStorage()
}

func callCreateTestRun(client *tmsClient, cfg *config.Config) {
	cfg.TestRunId = client.createTestRun()
	client.cfg.TestRunId = cfg.TestRunId
	print("test run id: ", cfg.TestRunId)
}

func getLogLevel(b bool) slog.Level {
	if b {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func initSyncStorage() {
	testRunID := cfg.TestRunId
	if testRunID == "" {
		logger.Warn("No test run ID available, skipping Sync Storage initialization")
		return
	}

	runner := syncstorage.NewRunner(
		testRunID,
		cfg.SyncStoragePort,
		cfg.Url,
		cfg.Token,
		logger,
	)

	if runner.Start() {
		syncRunner = runner
		logger.Info("Sync Storage initialized successfully")
	} else {
		logger.Warn("Failed to start Sync Storage, continuing without it")
	}
}
