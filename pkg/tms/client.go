package tms

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/testit-tms/adapters-go/pkg/tms/config"
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
func (c *tmsClient) writeTest(test testResult) (string, error) {
	const op = "tmsClient.writeTest"
	logger := logger.With("op", op)

	ctx := context.WithValue(context.Background(), tmsclient.ContextAPIKeys, map[string]tmsclient.APIKey{
		"Bearer or PrivateToken": {
			Key:    c.cfg.Token,
			Prefix: "PrivateToken",
		},
	})

	logger.Debug("searching for test", "externalId", test.externalId, slog.String("op", op))
	sr := getSearchRequest(test.externalId, c.cfg.ProjectId)
	resp, r, err := c.client.AutoTestsAPI.ApiV2AutoTestsSearchPost(ctx).
		AutoTestSearchApiModel(sr).
		Execute()
	if err != nil {
		logger.Error("failed to search for test", "error", err, slog.String("response", respToString(r.Body)), slog.String("op", op))
		return "", fmt.Errorf("%s: failed to search for test: %w", op, err)
	}

	var autotestID string
	if len(resp) == 0 {
		cr := testToAutotestModel(test, c.cfg.ProjectId)
		if c.cfg.AutomaticCreationTestCases {
			cr.SetShouldCreateWorkItem(c.cfg.AutomaticCreationTestCases)
		}

		logger.Debug("create new autotest", "request", cr)
		na, _, err := c.client.AutoTestsAPI.CreateAutoTest(ctx).
			AutoTestPostModel(cr).
			Execute()

		if err != nil {
			logger.Error("failed to create new autotest", "error", err, slog.String("response", respToString(r.Body)), slog.String("op", op))
			return "", fmt.Errorf("%s: failed to create new autotest: %w", op, err)
		}

		autotestID = na.Id
	} else {
		ur := testToUpdateAutotestModel(test, resp[0])
		logger.Debug("update existing autotest", "request", ur)
		r, err = c.client.AutoTestsAPI.UpdateAutoTest(ctx).
			AutoTestPutModel(ur).
			Execute()

		if err != nil {
			logger.Error("failed to update existing autotest", "error", err, slog.String("response", respToString(r.Body)), slog.String("op", op))
			return "", fmt.Errorf("%s: failed to update existing autotest: %w", op, err)
		}

		autotestID = resp[0].Id
	}

	if len(test.workItemIds) != 0 {
		var linkedWorkItems []tmsclient.WorkItemIdentifierModel
		linkedWorkItems, r, err = c.client.AutoTestsAPI.GetWorkItemsLinkedToAutoTest(ctx, autotestID).
			Execute()

		if err != nil {
			logger.Error("failed to get linked workitems to autotest", "error", err, slog.String("response", respToString(r.Body)), slog.String("op", op))
		}

		for _, v := range linkedWorkItems {
			var linkedWorkItemId string = strconv.Itoa(int(v.GetGlobalId()))
			var index int = getIndex(test.workItemIds, linkedWorkItemId)

			if index != -1 {
				test.workItemIds = remove(test.workItemIds, index)

				continue
			}

			if c.cfg.AutomaticUpdationLinksToTestCases {
				for i := 0; i < maxTries; i++ {
					r, err = c.client.AutoTestsAPI.DeleteAutoTestLinkFromWorkItem(ctx, autotestID).
						WorkItemId(linkedWorkItemId).
						Execute()
					if err != nil {
						logger.Error("failed to unlink autotest from workitem", "error", err, slog.String("response", respToString(r.Body)), slog.String("op", op))
						time.Sleep(waitingTime * time.Millisecond)
					} else {
						break
					}
				}
			}
		}

		for _, v := range test.workItemIds {
			logger.Debug("link autotest to workitem", "workItemId", v, "autotestId", autotestID)
			for i := 0; i < maxTries; i++ {
				r, err = c.client.AutoTestsAPI.LinkAutoTestToWorkItem(ctx, autotestID).
					WorkItemIdModel(tmsclient.WorkItemIdModel{
						Id: v,
					}).Execute()
				if err != nil {
					logger.Error("failed to link autotest to workitem", "error", err, slog.String("response", respToString(r.Body)), slog.String("op", op))
					time.Sleep(waitingTime * time.Millisecond)
				} else {
					break
				}
			}
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
		if r != nil && r.Body != nil {
			logger.Error("failed to upload result to test run", "error", err, slog.String("response", respToString(r.Body)), slog.String("op", op))
		} else {
			logger.Error("failed to upload result to test run", "error", err, slog.String("op", op))
		}

		return "", fmt.Errorf("%s: failed to upload result to test run: %w", op, err)
	}

	return ids[0], nil
}

func (c *tmsClient) writeAttachments(paths ...string) []string {
	const op = "tmsClient.writeAttachment"
	logger := logger.With("op", op)

	ctx := context.WithValue(context.Background(), tmsclient.ContextAPIKeys, map[string]tmsclient.APIKey{
		"Bearer or PrivateToken": {
			Key:    c.cfg.Token,
			Prefix: "PrivateToken",
		},
	})

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
			logger.Error("failed to upload attachment", "error", err, slog.String("response", respToString(r.Body)), slog.String("op", op))
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

func (c *tmsClient) updateTest(test testResult) error {
	const op = "tmsClient.updateTest"
	logger := logger.With("op", op)

	ctx := context.WithValue(context.Background(), tmsclient.ContextAPIKeys, map[string]tmsclient.APIKey{
		"Bearer or PrivateToken": {
			Key:    c.cfg.Token,
			Prefix: "PrivateToken",
		},
	})

	logger.Debug("searching for test", "externalId", test.externalId)
	sr := getSearchRequest(test.externalId, c.cfg.ProjectId)
	resp, r, err := c.client.AutoTestsAPI.ApiV2AutoTestsSearchPost(ctx).
		AutoTestSearchApiModel(sr).
		Execute()

	if err != nil {
		logger.Error("failed to search for test", "error", err, slog.String("response", respToString(r.Body)), slog.String("op", op))
		return fmt.Errorf("%s: failed to search for test: %w", op, err)
	}

	ur := testToUpdateAutotestModel(test, resp[0])

	r, err = c.client.AutoTestsAPI.UpdateAutoTest(ctx).
		AutoTestPutModel(ur).
		Execute()

	if err != nil {
		logger.Error("failed to update existing autotest", "error", err, slog.String("response", respToString(r.Body)), slog.String("op", op))
		return fmt.Errorf("%s: failed to update existing autotest: %w", op, err)
	}

	return nil
}

func (c *tmsClient) updateTestResult(resultId string, test testResult) error {
	const op = "tmsClient.updateTestResult"
	logger := logger.With("op", op)

	ctx := context.WithValue(context.Background(), tmsclient.ContextAPIKeys, map[string]tmsclient.APIKey{
		"Bearer or PrivateToken": {
			Key:    c.cfg.Token,
			Prefix: "PrivateToken",
		},
	})

	m, r, err := c.client.TestResultsAPI.ApiV2TestResultsIdGet(ctx, resultId).Execute()
	if err != nil {
		logger.Error("failed to get test result", "error", err, slog.String("response", respToString(r.Body)), slog.String("op", op))
		return fmt.Errorf("%s: failed to get test result: %w", op, err)
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
		logger.Error("failed to update test result", "error", err, slog.String("response", respToString(r.Body)), slog.String("op", op))
		return fmt.Errorf("%s: failed to update test result: %w", op, err)
	}

	return nil
}
