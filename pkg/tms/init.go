package tms

import (
	"os"

	"github.com/jtolds/gls"
	"github.com/testit-tms/adapters-go/pkg/tms/config"
	"golang.org/x/exp/slog"
)

var (
	cfg              *config.Config
	client           *tmsClient
	logger           *slog.Logger
	ctxMgr           *gls.ContextManager
	testPhaseObjects map[string]*testPhaseContainer
)

const (
	nodeKey         = "current_step"
	testResultKey   = "current_result_object"
	testInstanceKey = "test_instance"
)

func init() {
	cfg = config.MustLoad()
	client = newClient(*cfg)
	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: getLogLevel(cfg.IsDebug)}),
	)
	ctxMgr = gls.NewContextManager()
	testPhaseObjects = make(map[string]*testPhaseContainer)
}

func getLogLevel(b bool) slog.Level {
	if b {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}
