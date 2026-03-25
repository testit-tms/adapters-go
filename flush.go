package tms

import (
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/exp/slog"
)

const (
	// flushDebounceDelay is the time to wait after the last test finishes
	// before automatically flushing pending results.
	flushDebounceDelay = 2 * time.Second
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

	// debounceTimer is the current debounce timer (can be reset).
	debounceTimer *time.Timer
	debounceMu    sync.Mutex
)

// Flush writes all pending test results to TMS and notifies sync-storage.
// Call this explicitly from TestMain after m.Run() for guaranteed behavior,
// or rely on automatic flush triggered 2 seconds after the last test completes.
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

// trackTestEnd decrements the active test counter and schedules
// a debounced flush if no tests are running.
func trackTestEnd() {
	current := atomic.AddInt64(&activeTests, -1)
	if current == 0 && atomic.LoadInt64(&totalTests) > 0 && !cfg.ImportRealtime {
		scheduleFlush()
	}
}

// scheduleFlush starts or resets the debounce timer for auto-flush.
func scheduleFlush() {
	debounceMu.Lock()
	defer debounceMu.Unlock()

	if debounceTimer != nil {
		debounceTimer.Stop()
	}

	debounceTimer = time.AfterFunc(flushDebounceDelay, func() {
		// Double-check no new tests started during the delay
		if atomic.LoadInt64(&activeTests) == 0 {
			Flush()
		}
	})
}
