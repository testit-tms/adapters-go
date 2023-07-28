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
)

type TestMetadata struct {
	ClassName   string
	Title       string
	NameSpace   string
	Description string
	DisplayName string
	Parameters  map[string]interface{}
	Links       []Link
	Labels      []string
	ExternalId  string
	WorkItemIds []string
}

func Test(t *testing.T, m TestMetadata, f func()) {
	tr := newTestResult(m, t)

	defer func() {
		panicObject := recover()
		// getCurrentTestPhaseObject(t).Test = tr
		tr.completedOn = time.Now()
		tr.duration = tr.completedOn.UnixMilli() - tr.startedOn.UnixMilli()
		if panicObject != nil {
			t.Fail()
			tr.message = fmt.Sprintf("%+v", panicObject)
			tr.trace = string(debug.Stack())
		}

		if tr.status == "" {
			tr.status = getTestStatus(t)
		}

		tr.write()
		// setups := getCurrentTestPhaseObject(t).Befores
		// for _, setup := range setups {
		// 	setup.Children = append(setup.Children, r.UUID)
		// 	err := setup.writeResultsFile()
		// 	if err != nil {
		// 		log.Println("Failed to write content of result to json file", err)
		// 	}
		// }

		//if panicObject != nil {
		//	panic(panicObject)
		//}
	}()

	ctxMgr.SetValues(gls.Values{
		testResultKey:   tr,
		nodeKey:         tr,
		testInstanceKey: t,
	}, f)
}

const (
	passed  = "Passed"
	failed  = "Failed"
	skipped = "Skipped"
)

func getTestStatus(t *testing.T) string {
	if t.Failed() {
		return failed
	} else if t.Skipped() {
		return skipped
	} else {
		return passed
	}
}

func newTestResult(m TestMetadata, t *testing.T) *testResult {
	testResult := &testResult{
		startedOn:   time.Now(),
		displayName: m.DisplayName,
		className:   m.ClassName,
		nameSpace:   m.NameSpace,
		description: m.Description,
		title:       m.Title,
		links:       m.Links,
		labels:      m.Labels,
		externalId:  m.ExternalId,
		workItemIds: m.WorkItemIds,
		parameters: m.Parameters,
	}

	if testResult.displayName == "" {
		testResult.displayName = "Test"
	}

	if testResult.title == "" {
		testResult.title = testResult.displayName
	}

	if testResult.className == "" {
		testResult.className = t.Name()
	}

	if testResult.nameSpace == "" {
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
		testResult.nameSpace = strings.TrimSuffix(ts[len(ts)-1], ".go")
	}

	if testResult.externalId == "" {
		hash := md5.Sum([]byte(fmt.Sprintf("%s%s%s", testResult.displayName, testResult.className, testResult.nameSpace)))
		testResult.externalId = hex.EncodeToString(hash[:])
	}

	return testResult
}
