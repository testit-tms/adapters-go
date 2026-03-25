package tms

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/jtolds/gls"
	"github.com/testit-tms/adapters-go/models"
)

type TestMetadata struct {
	ClassName   string
	Title       string
	NameSpace   string
	Description string
	DisplayName string
	Parameters  map[string]interface{}
	Links       []models.Link
	Labels      []string
	Tags        []string
	ExternalId  string
	WorkItemIds []string
}

func Test(t *testing.T, m TestMetadata, f func()) {
	tr := newTestResult(m, t)
	testPhaseObjects := getCurrentTestPhaseObject(t)
	testPhaseObjects.test = tr

	// Track active tests for auto-flush
	trackTestStart()

	// Notify sync-storage that test execution started
	onRunningStarted()

	defer func() {
		panicObject := recover()
		tr.completedOn = time.Now()
		tr.duration = tr.completedOn.UnixMilli() - tr.startedOn.UnixMilli()
		if panicObject != nil {
			t.Fail()
			tr.message = fmt.Sprintf("%+v", panicObject)
			tr.trace = string(debug.Stack())
		}

		if testPhaseObjects.before != nil {
			tr.addBefore(testPhaseObjects.before.convertToStepResult())
		}

		if tr.status == "" {
			if testPhaseObjects.before != nil && testPhaseObjects.before.status == models.Failed {
				tr.status = models.Failed
				tr.message = testPhaseObjects.before.message
				tr.trace = testPhaseObjects.before.trace
			} else {
				tr.status = models.GetTestStatus(t)
			}
		}

		id := tr.write()
		if id != "" {
			testPhaseObjects.resultID = id
		}

		// In realtime mode, notify sync-storage per test.
		// In non-realtime mode, onBlockCompleted is called once during Flush().
		if cfg.ImportRealtime {
			onBlockCompleted()
		}

		// Track test end — may trigger debounced auto-flush
		trackTestEnd()
	}()

	if testPhaseObjects.before != nil && testPhaseObjects.before.status == models.Failed {
		return
	}

	ctxMgr.SetValues(gls.Values{
		testResultKey:   tr,
		nodeKey:         tr,
		testInstanceKey: t,
	}, f)
}

func newTestResult(m TestMetadata, t *testing.T) *TestResult {
	TestResult := &TestResult{
		startedOn:   time.Now(),
		displayName: m.DisplayName,
		className:   m.ClassName,
		nameSpace:   m.NameSpace,
		description: m.Description,
		title:       m.Title,
		links:       m.Links,
		labels:      m.Labels,
		tags:        m.Tags,
		externalId:  m.ExternalId,
		workItemIds: m.WorkItemIds,
		parameters:  m.Parameters,
	}

	if TestResult.displayName == "" {
		TestResult.displayName = "Test"
	}

	if TestResult.title == "" {
		TestResult.title = TestResult.displayName
	}

	if TestResult.className == "" {
		TestResult.className = t.Name()
	}

	if TestResult.nameSpace == "" {
		programCounters := make([]uintptr, 10)
		callersCount := runtime.Callers(0, programCounters)

		var testFile string
		for i := 0; i < callersCount; i++ {
			_, testFile, _, _ = runtime.Caller(i)
			if strings.Contains(testFile, "_test.go") {
				break
			}
		}
		ts := strings.Split(testFile, "/")
		TestResult.nameSpace = strings.TrimSuffix(ts[len(ts)-1], ".go")
	}

	if TestResult.externalId == "" {
		hash := md5.Sum([]byte(fmt.Sprintf("%s%s%s", TestResult.displayName, TestResult.className, TestResult.nameSpace)))
		TestResult.externalId = hex.EncodeToString(hash[:])
	}

	TestResult.externalKey = t.Name()

	return TestResult
}
