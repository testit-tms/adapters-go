package tms

import (
	"fmt"
	"strconv"

	tmsclient "github.com/testit-tms/api-client-golang/v3"
)

// hasInfo is about integrations with other systems,
// by default it is false
const defaultHasInfo = false

func testToAutotestModel(test testResult, projectId string) tmsclient.AutoTestPostModel {
	req := tmsclient.NewAutoTestPostModel(test.externalId, projectId, test.displayName)
	req.SetTitle(test.title)

	if test.description != "" {
		req.SetDescription(test.description)
	}

	if test.className != "" {
		req.SetClassname(test.className)
	}

	if test.nameSpace != "" {
		req.SetNamespace(test.nameSpace)
	}

	if len(test.labels) != 0 {
		labels := make([]tmsclient.LabelPostModel, 0, len(test.labels))
		for _, label := range test.labels {
			labels = append(labels, tmsclient.LabelPostModel{
				Name: label,
			})
		}
		req.SetLabels(labels)
	}

	links_post_models := resultLinksToLinkPostModels(test.links)
	req.SetLinks(links_post_models)

	if len(test.steps) != 0 {
		req.SetSteps(stepToAutoTestStepModel(test.steps))
	}

	if len(test.setups) != 0 {
		req.SetSetup(stepToAutoTestStepModel(test.setups))
	}

	req.SetExternalKey(test.externalKey)

	return *req
}

func stepToAutoTestStepModel(s []stepresult) []tmsclient.AutoTestStepModel {
	steps := make([]tmsclient.AutoTestStepModel, 0, len(s))
	for _, step := range s {
		model := tmsclient.NewAutoTestStepModel(step.name)
		model.SetDescription(step.description)

		if len(step.childrenSteps) != 0 {
			model.SetSteps(stepToAutoTestStepModel(step.childrenSteps))
		}

		steps = append(steps, *model)
	}

	return steps
}

func testToUpdateAutotestModel(test testResult, autotest tmsclient.AutoTestApiResult) tmsclient.AutoTestPutModel {
	req := tmsclient.NewAutoTestPutModel(test.externalId, autotest.ProjectId, test.displayName)

	if test.description != "" {
		req.SetDescription(test.description)
	}

	if test.className != "" {
		req.SetClassname(test.className)
	}

	if test.nameSpace != "" {
		req.SetNamespace(test.nameSpace)
	}

	if len(test.labels) != 0 {
		labels := make([]tmsclient.LabelPostModel, 0, len(test.labels))
		for _, label := range test.labels {
			labels = append(labels, tmsclient.LabelPostModel{
				Name: label,
			})
		}
		req.SetLabels(labels)
	}

	if test.title != "" {
		req.SetTitle(test.title)
	}

	if len(test.links) != 0 {
		links := make([]tmsclient.LinkPutModel, 0, len(test.links))
		for _, link := range test.links {
			l := tmsclient.NewLinkPutModel(link.Url, defaultHasInfo)
			l.SetTitle(link.Title)
			l.SetDescription(link.Description)

			if link.LinkType != "" {
				linkType, err := tmsclient.NewLinkTypeFromValue(string(link.LinkType))
				if err != nil {
					logger.Error("error converting link type", "error", err)
				} else {
					l.SetType(*linkType)
				}
			}

			links = append(links, *l)
		}
		req.SetLinks(links)
	}

	if len(test.steps) != 0 {
		req.SetSteps(stepToAutoTestStepModel(test.steps))
	}

	if len(test.setups) != 0 {
		req.SetSetup(stepToAutoTestStepModel(test.setups))
	}

	if len(test.teardowns) != 0 {
		req.SetTeardown(stepToAutoTestStepModel(test.teardowns))
	}

	req.SetExternalKey(test.externalKey)
	req.SetIsFlaky(autotest.IsFlaky)
	req.SetId(autotest.Id)

	return *req
}

func testToResultModel(test testResult, confID string) ([]tmsclient.AutoTestResultsForTestRunModel, error) {
	outcome, err := tmsclient.NewAvailableTestResultOutcomeFromValue(test.status)
	if err != nil {
		return nil, fmt.Errorf("error converting test status to outcome: %w", err)
	}
	req := tmsclient.NewAutoTestResultsForTestRunModel(confID, test.externalId, *outcome)
	req.SetDuration(test.duration)
	req.SetMessage(test.message)
	req.SetTraces(test.trace)
	req.SetStartedOn(test.startedOn)
	req.SetCompletedOn(test.completedOn)

	if len(test.steps) != 0 {
		steps, err := stepToAttachmentPutModelAutoTestStepResultsModel(test.steps)
		if err != nil {
			return nil, fmt.Errorf("error converting steps to attachment model: %w", err)
		}
		req.SetStepResults(steps)
	}

	if len(test.setups) != 0 {
		steps, err := stepToAttachmentPutModelAutoTestStepResultsModel(test.setups)
		if err != nil {
			return nil, fmt.Errorf("error converting setups to attachment model: %w", err)
		}
		req.SetSetupResults(steps)
	}

	links := resultLinksToLinkPostModels(test.resultLinks)
	req.SetLinks(links)

	if len(test.attachments) != 0 {
		attachs := make([]tmsclient.AttachmentPutModel, 0, len(test.attachments))
		for _, attach := range test.attachments {
			a := tmsclient.NewAttachmentPutModel(attach)
			attachs = append(attachs, *a)
		}
		req.SetAttachments(attachs)
	}

	if len(test.parameters) != 0 {
		params := make(map[string]string, len(test.parameters))
		for k, v := range test.parameters {
			params[k] = parseValueParameter(v)
		}
		req.SetParameters(params)
	}

	return []tmsclient.AutoTestResultsForTestRunModel{*req}, nil
}

func resultLinksToLinkPostModels(source []Link) []tmsclient.LinkPostModel {
	if len(source) == 0 {
		return []tmsclient.LinkPostModel{}
	}

	links := make([]tmsclient.LinkPostModel, 0, len(source))
	for _, link := range source {
		l := tmsclient.NewLinkPostModel(link.Url, defaultHasInfo)
		l.SetTitle(link.Title)
		l.SetDescription(link.Description)
		if link.LinkType != "" {
			linkType, err := tmsclient.NewLinkTypeFromValue(string(link.LinkType))
			if err != nil {
				logger.Error("error converting link type", "error", err)
			} else {
				l.SetType(*linkType)
			}
		}
		links = append(links, *l)
	}

	return links
}

func stepToAttachmentPutModelAutoTestStepResultsModel(s []stepresult) ([]tmsclient.AttachmentPutModelAutoTestStepResultsModel, error) {
	steps := make([]tmsclient.AttachmentPutModelAutoTestStepResultsModel, 0, len(s))
	for _, step := range s {
		model := tmsclient.NewAttachmentPutModelAutoTestStepResultsModel()
		model.SetTitle(step.name)
		model.SetDescription(step.description)
		outcome, err := tmsclient.NewAvailableTestResultOutcomeFromValue(step.status)
		if err != nil {
			return nil, err
		}
		model.SetOutcome(*outcome)
		model.SetStartedOn(step.startedOn)
		model.SetCompletedOn(step.completedOn)
		model.SetDuration(step.duration)

		if len(step.attachments) != 0 {
			attachs := make([]tmsclient.AttachmentPutModel, 0, len(step.attachments))
			for _, attach := range step.attachments {
				a := tmsclient.NewAttachmentPutModel(attach)
				attachs = append(attachs, *a)
			}
			model.SetAttachments(attachs)
		}

		if len(step.childrenSteps) != 0 {
			cs, err := stepToAttachmentPutModelAutoTestStepResultsModel(step.childrenSteps)
			if err != nil {
				return nil, err
			}
			model.SetStepResults(cs)
		}

		if len(step.parameters) != 0 {
			params := make(map[string]string, len(step.parameters))
			for k, v := range step.parameters {
				params[k] = parseValueParameter(v)
			}
			model.SetParameters(params)
		}

		steps = append(steps, *model)
	}

	return steps, nil
}

func parseValueParameter(value interface{}) string {

	switch value.(type) {
	case []byte:
		return string(value.([]byte))
	case uintptr:
		return strconv.Itoa(int(value.(uintptr)))
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case complex64:
		return fmt.Sprintf("%f i%f", real(value.(complex64)), imag(value.(complex64)))
	case complex128:
		return fmt.Sprintf("%f i%f", real(value.(complex128)), imag(value.(complex128)))
	case uint:
		return strconv.FormatUint(uint64(value.(uint)), 10)
	case uint8:
		return strconv.FormatUint(uint64(value.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(value.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(value.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(value.(uint64), 10)
	case int:
		return strconv.FormatInt(int64(value.(int)), 10)
	case int8:
		return strconv.FormatInt(int64(value.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(value.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(value.(int32)), 10)
	case int64:
		return strconv.FormatInt(value.(int64), 10)
	case bool:
		return strconv.FormatBool(value.(bool))
	case string:
		return value.(string)
	default:
		return fmt.Sprintf("%+v", value)
	}
}

func getSearchRequest(externalID, projectID string) tmsclient.AutoTestSearchApiModel {
	f := tmsclient.NewAutoTestFilterApiModel()
	f.SetExternalIds([]string{externalID})
	f.SetProjectIds([]string{projectID})
	f.SetIsDeleted(false)

	req := tmsclient.NewAutoTestSearchApiModel()
	req.SetFilter(*f)

	return *req
}

func testToUpdateResultModel(model *tmsclient.TestResultResponse, test testResult) (tmsclient.TestResultUpdateV2Request, error) {
	tearDowns, err := stepToAttachmentPutModelAutoTestStepResultsModel(test.teardowns)
	if err != nil {
		return tmsclient.TestResultUpdateV2Request{}, err
	}

	setups, err := stepToAttachmentPutModelAutoTestStepResultsModel(test.setups)
	if err != nil {
		return tmsclient.TestResultUpdateV2Request{}, err
	}

	req := tmsclient.NewTestResultUpdateV2Request()
	req.SetTeardownResults(tearDowns)
	req.SetSetupResults(setups)

	req.SetDurationInMs(model.GetDurationInMs())
	req.SetLinks(model.GetLinks())
	req.SetStepResults(model.GetStepResults())
	req.SetFailureClassIds(model.GetFailureClassIds())
	req.SetComment(model.GetComment())

	if len(model.Attachments) != 0 {
		attachs := make([]tmsclient.AttachmentUpdateRequest, 0, len(model.Attachments))
		for _, attach := range model.Attachments {
			a := tmsclient.NewAttachmentUpdateRequest(attach.Id)
			attachs = append(attachs, *a)
		}

		req.SetAttachments(attachs)
	}

	req.SetOutcome(tmsclient.TestResultOutcome(test.status))

	return *req, nil
}
