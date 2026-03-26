package tms

import (
	"sync"
	"sync/atomic"
	"testing"

	"golang.org/x/exp/slog"
)

var (
	// pendingResults stores test results when importRealtime is false.
	pendingResults   []TestResult
	pendingResultsMu sync.Mutex

	// activeTests tracks the number of currently running tests.
	activeTests int64

	// totalTests tracks total tests that have been started (to avoid premature flush).
	totalTests int64

	// flushOnce ensures flush is only executed once.
	flushOnce sync.Once
)

// Run is the recommended way to use the adapter with importRealtime=false.
// Call it from TestMain to guarantee all buffered results are flushed:
//
//	func TestMain(m *testing.M) {
//	    os.Exit(tms.Run(m))
//	}
func Run(m *testing.M) int {
	code := m.Run()
	if !cfg.ImportRealtime {
		Flush()
	}
	return code
}

// Flush writes all pending test results to TMS and notifies sync-storage.
// Call this explicitly from TestMain (via Run) after m.Run() for guaranteed
// behavior, or rely on automatic flush triggered when the last test completes.
func Flush() {
	flushOnce.Do(doFlush)
}

func doFlush() {
	const op = "tms.doFlush"

	pendingResultsMu.Lock()
	results := make([]TestResult, len(pendingResults))
	copy(results, pendingResults)
	pendingResults = nil
	pendingResultsMu.Unlock()

	if len(results) == 0 {
		logger.Info("No pending test results to flush", slog.String("op", op))
		onBlockCompleted()
		return
	}

	logger.Info("Flushing pending test results",
		slog.Int("count", len(results)),
		slog.String("op", op))

	for i := range results {
		id, err := client.writeTest(results[i])
		if err != nil {
			logger.Error("error writing pending test result",
				"error", err,
				"externalId", results[i].externalId,
				slog.String("op", op))
			continue
		}

		// Update testPhaseObjects with the result ID if available
		if id != "" {
			if tpo, ok := testPhaseObjects[results[i].externalKey]; ok {
				tpo.resultID = id
			}
		}
	}

	logger.Info("Flush completed", slog.Int("count", len(results)), slog.String("op", op))

	// Notify sync-storage that all tests are done
	onBlockCompleted()
}

// addPendingResult adds a test result to the pending buffer.
func addPendingResult(tr TestResult) {
	pendingResultsMu.Lock()
	pendingResults = append(pendingResults, tr)
	pendingResultsMu.Unlock()
}

// trackTestStart increments the active test counter.
func trackTestStart() {
	atomic.AddInt64(&activeTests, 1)
	atomic.AddInt64(&totalTests, 1)
}

// trackTestEnd decrements the active test counter.
func trackTestEnd() {
	atomic.AddInt64(&activeTests, -1)
}
