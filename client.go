package tms

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/testit-tms/adapters-go/client_helpers"
	"github.com/testit-tms/adapters-go/config"
	"github.com/testit-tms/adapters-go/htmlutils"
	tmsclient "github.com/testit-tms/api-client-golang/v3"
	"golang.org/x/exp/slog"
)

const (
	maxTries    = 10
	waitingTime = 100
)

type tmsClient struct {
	cfg    config.Config
	client *tmsclient.APIClient
}

func newClient(cfg config.Config) *tmsClient {
	var scheme string

	if strings.Contains(cfg.Url, "https") {
		scheme = "https"
	} else {
		scheme = "http"
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.CertValidation},
	}
	hc := &http.Client{Transport: tr}
	configuration := tmsclient.NewConfiguration()
	configuration.Host = strings.TrimPrefix(strings.TrimSuffix(cfg.Url, "/"), fmt.Sprintf("%s://", scheme))
	configuration.Scheme = scheme
	configuration.HTTPClient = hc
	return &tmsClient{
		client: tmsclient.NewAPIClient(configuration),
		cfg:    cfg,
	}
}

// TODO: Refactoring is needed
func (c *tmsClient) writeTest(test TestResult) (string, error) {
	const op = "tmsClient.writeTest"
	logger := logger.With("op", op)

	ctx := client_helpers.AuthContext(c.cfg.Token)

	logger.Debug("searching for test", "externalId", test.externalId, slog.String("op", op))
	sr := getSearchRequest(test.externalId, c.cfg.ProjectId)
	resp, r, err := c.client.AutoTestsAPI.ApiV2AutoTestsSearchPost(ctx).
		AutoTestSearchApiModel(sr).
		Execute()
	if err != nil {
		return "", client_helpers.LogAndWrapAPIError(logger, op, "failed to search for test", err, r)
	}

	var autotestID string
	if len(resp) == 0 {
		cr := testToAutotestModel(test, c.cfg.ProjectId)
		if c.cfg.AutomaticCreationTestCases {
			cr.SetShouldCreateWorkItem(c.cfg.AutomaticCreationTestCases)
		}

		logger.Debug("create new autotest", "request", cr)
		na, createResp, err := c.client.AutoTestsAPI.CreateAutoTest(ctx).
			AutoTestCreateApiModel(cr).
			Execute()

		if err != nil {
			return "", client_helpers.LogAndWrapAPIError(logger, op, "failed to create new autotest", err, createResp)
		}

		autotestID = na.Id
	} else {
		ur := testToUpdateAutotestModel(test, resp[0])
		logger.Debug("update existing autotest", "request", ur)
		r, err = c.client.AutoTestsAPI.UpdateAutoTest(ctx).
			AutoTestUpdateApiModel(ur).
			Execute()

		if err != nil {
			return "", client_helpers.LogAndWrapAPIError(logger, op, "failed to update existing autotest", err, r)
		}

		autotestID = resp[0].Id
	}

	if len(test.workItemIds) != 0 {
		var linkedWorkItems []tmsclient.AutoTestWorkItemIdentifierApiResult
		linkedWorkItems, r, err = c.client.AutoTestsAPI.GetWorkItemsLinkedToAutoTest(ctx, autotestID).
			Execute()

		if err != nil {
			_ = client_helpers.LogAndWrapAPIError(logger, op, "failed to get linked workitems to autotest", err, r)
		}

		for _, v := range linkedWorkItems {
			var linkedWorkItemId string = strconv.Itoa(int(v.GetGlobalId()))
			var index int = getIndex(test.workItemIds, linkedWorkItemId)

			if index != -1 {
				test.workItemIds = remove(test.workItemIds, index)

				continue
			}

			if c.cfg.AutomaticUpdationLinksToTestCases {
				_ = client_helpers.Retry(maxTries, waitingTime*time.Millisecond, func() error {
					r, err = c.client.AutoTestsAPI.DeleteAutoTestLinkFromWorkItem(ctx, autotestID).
						WorkItemId(linkedWorkItemId).
						Execute()
					if err != nil {
						_ = client_helpers.LogAndWrapAPIError(logger, op, "failed to unlink autotest from workitem", err, r)
					}
					return err
				})
			}
		}

		for _, v := range test.workItemIds {
			logger.Debug("link autotest to workitem", "workItemId", v, "autotestId", autotestID)
			_ = client_helpers.Retry(maxTries, waitingTime*time.Millisecond, func() error {
				r, err = c.client.AutoTestsAPI.LinkAutoTestToWorkItem(ctx, autotestID).
					WorkItemIdApiModel(tmsclient.WorkItemIdApiModel{
						Id: v,
					}).Execute()
				if err != nil {
					_ = client_helpers.LogAndWrapAPIError(logger, op, "failed to link autotest to workitem", err, r)
				}
				return err
			})
		}
	}

	rr, err := testToResultModel(test, c.cfg.ConfigurationId)
	if err != nil {
		logger.Error("failed to convert test to result model", "error", err, slog.String("op", op))
		return "", fmt.Errorf("%s: failed to convert test to result model: %w", op, err)
	}
	logger.Debug("upload result to test run", "request", rr)
	ids, r, err := c.client.TestRunsAPI.SetAutoTestResultsForTestRun(ctx, c.cfg.TestRunId).
		AutoTestResultsForTestRunModel(rr).
		Execute()

	if err != nil {
		return "", client_helpers.LogAndWrapAPIError(logger, op, "failed to upload result to test run", err, r)
	}

	if len(ids) == 0 {
		return "", fmt.Errorf("%s: failed to upload result to test run: empty result id list", op)
	}
	return ids[0], nil
}

// return test run id
func (c *tmsClient) createTestRun() string {
	const op = "tmsClient.createTestRun"
	logger := logger.With("op", op)

	ctx := client_helpers.AuthContext(c.cfg.Token)

	model := tmsclient.NewCreateEmptyTestRunApiModel(c.cfg.ProjectId)
	model.SetName(c.cfg.TestRunName)

	// Apply HTML escaping to the model
	htmlutils.EscapeHtmlInObject(model)

	testRun, r, err := c.client.TestRunsAPI.CreateEmpty(ctx).
		CreateEmptyTestRunApiModel(*model).
		Execute()

	if err != nil {
		_ = client_helpers.LogAndWrapAPIError(logger, op, "failed to create test run", err, r)
		return ""
	}

	return testRun.Id
}

// return test run
func (c *tmsClient) getTestRun() *tmsclient.TestRunV2ApiResult {
	const op = "tmsClient.getTestRun"
	logger := logger.With("op", op)

	ctx := client_helpers.AuthContext(c.cfg.Token)

	testRun, r, err := c.client.TestRunsAPI.GetTestRunById(ctx, c.cfg.TestRunId).
		Execute()

	if err != nil {
		_ = client_helpers.LogAndWrapAPIError(logger, op, "failed to get test run", err, r)
		return nil
	}

	return testRun
}

func (c *tmsClient) updateTestRun() {
	const op = "tmsClient.updateTestRun"
	logger := logger.With("op", op)

	if c.cfg.TestRunName == "" {
		return
	}

	ctx := client_helpers.AuthContext(c.cfg.Token)

	testRun := c.getTestRun()

	if testRun == nil || testRun.Name == c.cfg.TestRunName {
		return
	}

	testRun.Name = c.cfg.TestRunName
	model := buildUpdateEmptyTestRunApiModel(testRun)

	// Apply HTML escaping to the model
	htmlutils.EscapeHtmlInObject(model)

	r, err := c.client.TestRunsAPI.UpdateEmpty(ctx).
		UpdateEmptyTestRunApiModel(*model).
		Execute()

	if err != nil {
		_ = client_helpers.LogAndWrapAPIError(logger, op, "failed to update test run", err, r)
		return
	}
}

func (c *tmsClient) writeAttachments(paths ...string) []string {
	const op = "tmsClient.writeAttachment"
	logger := logger.With("op", op)

	ctx := client_helpers.AuthContext(c.cfg.Token)

	attachmanetsIds := make([]string, 0, len(paths))
	for _, p := range paths {
		logger.Debug("uploading attachment", "path", p, slog.String("op", op))

		f, err := os.Open(p)
		if err != nil {
			logger.Error("failed to open file", "error", err)
			continue
		}
		resp, r, err := c.client.AttachmentsAPI.ApiV2AttachmentsPost(ctx).
			File(f).
			Execute()

		if err != nil {
			_ = client_helpers.LogAndWrapAPIError(logger, op, "failed to upload attachment", err, r)
			continue
		}

		logger.Debug("attachment uploaded", "id", resp.Id, "path", p, slog.String("op", op))

		attachmanetsIds = append(attachmanetsIds, resp.Id)
	}

	return attachmanetsIds
}

func respToString(r io.ReadCloser) string {
	respBytes, err := io.ReadAll(r)
	if err != nil {
		logger.Error("failed to read response body", "error", err)
		return ""
	}
	return string(respBytes)
}

func getIndex(list []string, item string) int {
	for i := 0; i < len(list); i++ {
		if item == list[i] {
			return i
		}
	}
	return -1
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func (c *tmsClient) updateTest(test TestResult) error {
	const op = "tmsClient.updateTest"
	logger := logger.With("op", op)

	ctx := client_helpers.AuthContext(c.cfg.Token)

	logger.Debug("searching for test", "externalId", test.externalId)
	sr := getSearchRequest(test.externalId, c.cfg.ProjectId)
	resp, r, err := c.client.AutoTestsAPI.ApiV2AutoTestsSearchPost(ctx).
		AutoTestSearchApiModel(sr).
		Execute()

	if err != nil {
		return client_helpers.LogAndWrapAPIError(logger, op, "failed to search for test", err, r)
	}

	ur := testToUpdateAutotestModel(test, resp[0])

	r, err = c.client.AutoTestsAPI.UpdateAutoTest(ctx).
		AutoTestUpdateApiModel(ur).
		Execute()

	if err != nil {
		return client_helpers.LogAndWrapAPIError(logger, op, "failed to update existing autotest", err, r)
	}

	return nil
}

func (c *tmsClient) updateTestResult(resultId string, test TestResult) error {
	const op = "tmsClient.updateTestResult"
	logger := logger.With(
		"op", op,
		"resultId", resultId,
		"externalId", test.externalId,
	)

	ctx := client_helpers.AuthContext(c.cfg.Token)

	logger.Debug("getting test result", "resultId", resultId, slog.String("op", op))
	m, r, err := c.client.TestResultsAPI.ApiV2TestResultsIdGet(ctx, resultId).Execute()
	if err != nil {
		return client_helpers.LogAndWrapAPIError(logger, op, "failed to get test result", err, r)
	}

	ur, err := testToUpdateResultModel(m, test)
	if err != nil {
		logger.Error("failed to convert test to result model", "error", err, slog.String("op", op))
		return fmt.Errorf("%s: failed to convert test to result model: %w", op, err)
	}

	logger.Debug("update test result", "request", ur, slog.String("op", op))

	r, err = c.client.TestResultsAPI.ApiV2TestResultsIdPut(ctx, resultId).
		TestResultUpdateV2Request(ur).
		Execute()

	if err != nil {
		return client_helpers.LogAndWrapAPIError(logger, op, "failed to update test result", err, r)
	}

	return nil
}
