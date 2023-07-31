package tms

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/testit-tms/adapters-go/pkg/tms/config"
	tmsclient "github.com/testit-tms/api-client-golang"
	"golang.org/x/exp/slog"
)

type tmsClient struct {
	cfg    config.Config
	client *tmsclient.APIClient
}

func New(cfg config.Config) *tmsClient {
	var scheme string

	if strings.Contains(cfg.Url, "https") {
		scheme = "https"
	} else {
		scheme = "http"
	}

	configuration := tmsclient.NewConfiguration()
	configuration.Host = strings.TrimPrefix(strings.TrimSuffix(cfg.Url, "/"), fmt.Sprintf("%s://", scheme))
	configuration.Scheme = scheme
	return &tmsClient{
		client: tmsclient.NewAPIClient(configuration),
		cfg:    cfg,
	}
}

func (c *tmsClient) writeTest(test testResult) (string, error) {
	const op = "tmsClient.writeTest"
	logger := logger.With("op", op)

	ctx := context.WithValue(context.Background(), tmsclient.ContextAPIKeys, map[string]tmsclient.APIKey{
		"Bearer or PrivateToken": {
			Key:    c.cfg.Token,
			Prefix: "PrivateToken",
		},
	})

	logger.Debug("searching for test", "externalId", test.externalId)
	sr := getSearchRequest(test.externalId, c.cfg.ProjectId)
	resp, r, err := c.client.AutoTestsApi.ApiV2AutoTestsSearchPost(ctx).
		ApiV2AutoTestsSearchPostRequest(sr).
		Execute()

	if err != nil {
		logger.Error("failed to search for test", "error", err, slog.String("response", respToString(r.Body)))
		return "", err
	}

	var autotestID string
	if len(resp) == 0 {
		cr := testToAutotestModel(test, c.cfg.ProjectId)
		logger.Debug("create new autotest", "request", cr)
		na, _, err := c.client.AutoTestsApi.CreateAutoTest(ctx).
			CreateAutoTestRequest(cr).
			Execute()

		if err != nil {
			logger.Error("failed to create new autotest", "error", err, slog.String("response", respToString(r.Body)))
			return "", err
		}

		autotestID = *na.Id
	} else {
		ur := testToUpdateAutotestModel(test, resp[0])
		logger.Debug("update existing autotest", "request", ur)
		r, err = c.client.AutoTestsApi.UpdateAutoTest(ctx).
			UpdateAutoTestRequest(ur).
			Execute()

		if err != nil {
			logger.Error("failed to update existing autotest", "error", err, slog.String("response", respToString(r.Body)))
			return "", err
		}

		autotestID = *resp[0].Id
	}

	if len(test.workItemIds) != 0 {
		for _, v := range test.workItemIds {
			logger.Debug("link autotest to workitem", "workItemId", v, "autotestId", autotestID)
			r, err = c.client.AutoTestsApi.LinkAutoTestToWorkItem(ctx, autotestID).
				LinkAutoTestToWorkItemRequest(tmsclient.LinkAutoTestToWorkItemRequest{
					Id: v,
				}).
				Execute()
		}
		if err != nil {
			logger.Error("failed to link autotest to workitem", "error", err, slog.String("response", respToString(r.Body)))
		}
	}

	rr, err := testToResultModel(test, c.cfg.ConfigurationId)
	if err != nil {
		logger.Error("failed to convert test to result model", "error", err)
		return "", err
	}
	logger.Debug("upload result to test run", "request", rr)
	ids, r, err := c.client.TestRunsApi.SetAutoTestResultsForTestRun(ctx, c.cfg.TestRunId).
		AutoTestResultsForTestRunModel(rr).
		Execute()

	if err != nil {
		logger.Error("failed to upload result to test run", "error", err, slog.String("response", respToString(r.Body)))
		return "", err
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
		logger.Debug("uploading attachment", "path", p)

		f, err := os.Open(p)
		if err != nil {
			logger.Error("failed to open file", "error", err)
			continue
		}
		resp, r, err := c.client.AttachmentsApi.ApiV2AttachmentsPost(ctx).
			File(f).
			Execute()

		if err != nil {
			logger.Error("failed to upload attachment", "error", err, slog.String("response", respToString(r.Body)))
			continue
		}

		logger.Debug("attachment uploaded", "id", resp.Id, "path", p)

		attachmanetsIds = append(attachmanetsIds, resp.Id)
	}

	return attachmanetsIds
}

func respToString(r io.ReadCloser) string {
	respBytes, err := ioutil.ReadAll(r)
	if err != nil {
		logger.Error("failed to read response body", "error", err)
		return ""
	}
	return string(respBytes)
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
	resp, r, err := c.client.AutoTestsApi.ApiV2AutoTestsSearchPost(ctx).
		ApiV2AutoTestsSearchPostRequest(sr).
		Execute()

	if err != nil {
		logger.Error("failed to search for test", "error", err, slog.String("response", respToString(r.Body)))
		return err
	}

	ur := testToUpdateAutotestModel(test, resp[0])

	r, err = c.client.AutoTestsApi.UpdateAutoTest(ctx).
		UpdateAutoTestRequest(ur).
		Execute()

	if err != nil {
		logger.Error("failed to update existing autotest", "error", err, slog.String("response", respToString(r.Body)))
		return err
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
	
	ur, err := testToUpdateResultModel(test)
	if err != nil {
		logger.Error("failed to convert test to result model", "error", err)
		return err
	}

	logger.Debug("update test result", "request", ur)

	r, err := c.client.TestResultsApi.ApiV2TestResultsIdPut(ctx, resultId).
		ApiV2TestResultsIdPutRequest(ur).
		Execute()

	if err != nil {
		logger.Error("failed to update test result", "error", err, slog.String("response", respToString(r.Body)))
		return err
	}

	return nil
}
