package tms

import (
	"os"

	"github.com/jtolds/gls"
	"github.com/testit-tms/adapters-go/pkg/tms/config"
	"golang.org/x/exp/slog"
)

var (
	cfg    *config.Config
	client *tmsClient
	logger *slog.Logger
	ctxMgr *gls.ContextManager
)

const (
	nodeKey           = "current_step_container"
	testResultKey     = "test_result_object"
	testInstanceKey   = "test_instance"
)

func init() {
	cfg = config.MustLoad()
	client = New(*cfg)
	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	ctxMgr = gls.NewContextManager()
}
