package tms

import (
	"context"
	"fmt"
	"strings"

	"github.com/testit-tms/adapters-go/pkg/tms/config"
	tmsclient "github.com/testit-tms/api-client-golang"
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

func (c *tmsClient) writeTest(test testResult) error {
	const op = "tmsClient.writeTest"
	logger := logger.With("op", op)

	ctx := context.WithValue(context.Background(), tmsclient.ContextAPIKeys, map[string]tmsclient.APIKey{
		"Bearer or PrivateToken": {
			Key:    c.cfg.Token,
			Prefix: "PrivateToken",
		},
	})

	nulBool := new(bool)
	*nulBool = false
	logger.Debug("searching for test", "externalId", test.externalId)
	resp, _, err := c.client.AutoTestsApi.ApiV2AutoTestsSearchPost(ctx).
		ApiV2AutoTestsSearchPostRequest(tmsclient.ApiV2AutoTestsSearchPostRequest{
			Filter: &tmsclient.AutotestsSelectModelFilter{
				ExternalIds: []string{test.externalId},
				ProjectIds:  []string{c.cfg.ProjectId},
				IsDeleted:   *tmsclient.NewNullableBool(nulBool),
			},
		}).Execute()

	if err != nil {
		logger.Error("failed to search for test", "error", err)
		return err
	}

	var autotestID string
	if len(resp) == 0 {
		req := testToAutotestModel(test, c.cfg.ProjectId)
		logger.Debug("create new autotest", "request", req)
		na, _, err := c.client.AutoTestsApi.CreateAutoTest(ctx).
			CreateAutoTestRequest(req).
			Execute()

		if err != nil {
			logger.Error("failed to create new autotest", "error", err)
			return err
		}

		autotestID = *na.Id
	} else {
		req := testToUpdateAutotestModel(test, resp[0])
		logger.Debug("update existing autotest", "request", req)
		_, err = c.client.AutoTestsApi.UpdateAutoTest(ctx).
			UpdateAutoTestRequest(req).
			Execute()

		if err != nil {
			logger.Error("failed to update existing autotest", "error", err)
			return err
		}

		autotestID = *resp[0].Id
	}

	if len(test.workItemIds) != 0 {
		for _, v := range test.workItemIds {
			logger.Debug("link autotest to workitem", "workItemId", v, "autotestId", autotestID)
			_, err = c.client.AutoTestsApi.LinkAutoTestToWorkItem(ctx, autotestID).
				LinkAutoTestToWorkItemRequest(tmsclient.LinkAutoTestToWorkItemRequest{
					Id: v,
				}).
				Execute()
		}
		if err != nil {
			logger.Error("failed to link autotest to workitem", "error", err)
		}
	}

	req, err := testToResultModel(test, c.cfg.ConfigurationId)
	if err != nil {
		logger.Error("failed to convert test to result model", "error", err)
		return err
	}
	logger.Debug("upload result to test run", "request", req)
	_, _, err = c.client.TestRunsApi.SetAutoTestResultsForTestRun(ctx, c.cfg.TestRunId).
		AutoTestResultsForTestRunModel(req).
		Execute()

	if err != nil {
		logger.Error("failed to upload result to test run", "error", err)
		return err
	}

	return nil
}
