package tms

import (
	"fmt"
	"strconv"

	"github.com/testit-tms/adapters-go/htmlutils"
	tmsclient "github.com/testit-tms/api-client-golang/v3"
)

// TODO: validate that hasInfo always true is correct
const defaultHasInfo = true

func testToAutotestModel(test testResult, projectId string) tmsclient.AutoTestCreateApiModel {
	req := tmsclient.NewAutoTestCreateApiModel(projectId, test.externalId, test.displayName)
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
		labels := make([]tmsclient.LabelApiModel, 0, len(test.labels))
		for _, label := range test.labels {
			labels = append(labels, tmsclient.LabelApiModel{
				Name: label,
			})
		}
		req.SetLabels(labels)
	}

	if len(test.links) != 0 {
		links := make([]tmsclient.LinkCreateApiModel, 0, len(test.links))
		for _, link := range test.links {

			l := tmsclient.NewLinkCreateApiModel(link.Url, defaultHasInfo)
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

	req.SetExternalKey(test.externalKey)

	// Apply HTML escaping to the model
	htmlutils.EscapeHtmlInObject(req)

	return *req
}

func stepToAutoTestStepModel(s []stepresult) []tmsclient.AutoTestStepApiModel {
	steps := make([]tmsclient.AutoTestStepApiModel, 0, len(s))
	for _, step := range s {
		model := tmsclient.NewAutoTestStepApiModel(step.name)
		model.SetDescription(step.description)

		if len(step.childrenSteps) != 0 {
			model.SetSteps(stepToAutoTestStepModel(step.childrenSteps))
		}

		steps = append(steps, *model)
	}

	// Apply HTML escaping to the steps slice
	htmlutils.EscapeHtmlInObjectSlice(steps)

	return steps
}

func testToUpdateAutotestModel(test testResult, autotest tmsclient.AutoTestApiResult) tmsclient.AutoTestUpdateApiModel {
	req := tmsclient.NewAutoTestUpdateApiModel(autotest.ProjectId, test.externalId, test.displayName)

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
		labels := make([]tmsclient.LabelApiModel, 0, len(test.labels))
		for _, label := range test.labels {
			labels = append(labels, tmsclient.LabelApiModel{
				Name: label,
			})
		}
		req.SetLabels(labels)
	}

	if test.title != "" {
		req.SetTitle(test.title)
	}

	if len(test.links) != 0 {
		links := make([]tmsclient.LinkUpdateApiModel, 0, len(test.links))
		for _, link := range test.links {
			l := tmsclient.NewLinkUpdateApiModel(link.Url, defaultHasInfo)
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

	// Apply HTML escaping to the model
	htmlutils.EscapeHtmlInObject(req)

	return *req
}

func testToResultModel(test testResult, confID string) ([]tmsclient.AutoTestResultsForTestRunModel, error) {
	req := tmsclient.NewAutoTestResultsForTestRunModel(confID, test.externalId)
	req.SetStatusCode(test.status)
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

	if len(test.resultLinks) != 0 {
		links := make([]tmsclient.LinkPostModel, 0, len(test.resultLinks))
		for _, link := range test.resultLinks {
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
		req.SetLinks(links)
	}

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

	// Apply HTML escaping to the request
	htmlutils.EscapeHtmlInObject(req)

	return []tmsclient.AutoTestResultsForTestRunModel{*req}, nil
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

	// Apply HTML escaping to the steps slice
	htmlutils.EscapeHtmlInObjectSlice(steps)

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

	// Apply HTML escaping to the search request
	htmlutils.EscapeHtmlInObject(req)

	return *req
}

func mapAttachmentsToStepResults(attachments []tmsclient.AttachmentPutModelAutoTestStepResultsModel) ([]tmsclient.AutoTestStepResultUpdateRequest, error) {
	results := make([]tmsclient.AutoTestStepResultUpdateRequest, len(attachments))
	for i, attachment := range attachments {
		result := tmsclient.NewAutoTestStepResultUpdateRequest()
		result.SetTitle(attachment.GetTitle())
		result.SetDescription(attachment.GetDescription())

		outcome, err := tmsclient.NewAvailableTestResultOutcomeFromValue(string(attachment.GetOutcome()))
		if err != nil {
			return nil, err
		}
		result.SetOutcome(*outcome)
		result.SetStartedOn(attachment.GetStartedOn())
		result.SetCompletedOn(attachment.GetCompletedOn())
		result.SetDuration(attachment.GetDuration())

		// Mapping nested attachments at the step level is not supported in this model.
		// Attachments should be linked to the test result as a whole.

		if attachment.HasStepResults() {
			nestedResults, err := mapAttachmentsToStepResults(attachment.GetStepResults())
			if err != nil {
				return nil, err
			}
			result.SetStepResults(nestedResults)
		}

		result.SetParameters(attachment.GetParameters())

		results[i] = *result
	}
	return results, nil
}

func testToUpdateResultModel(model *tmsclient.TestResultResponse, test testResult) (tmsclient.TestResultUpdateV2Request, error) {
	tearDownsAttachments, err := stepToAttachmentPutModelAutoTestStepResultsModel(test.teardowns)
	if err != nil {
		return tmsclient.TestResultUpdateV2Request{}, err
	}

	tearDowns, err := mapAttachmentsToStepResults(tearDownsAttachments)
	if err != nil {
		return tmsclient.TestResultUpdateV2Request{}, fmt.Errorf("error mapping tearDowns: %w", err)
	}

	setupsAttachments, err := stepToAttachmentPutModelAutoTestStepResultsModel(test.setups)
	if err != nil {
		return tmsclient.TestResultUpdateV2Request{}, err
	}

	setups, err := mapAttachmentsToStepResults(setupsAttachments)
	if err != nil {
		return tmsclient.TestResultUpdateV2Request{}, fmt.Errorf("error mapping setups: %w", err)
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

	req.SetStatusCode(test.status)

	// Apply HTML escaping to the update request
	htmlutils.EscapeHtmlInObject(req)

	return *req, nil
}

func buildUpdateEmptyTestRunApiModel(testRun *tmsclient.TestRunV2ApiResult) *tmsclient.UpdateEmptyTestRunApiModel {
	model := tmsclient.NewUpdateEmptyTestRunApiModel(testRun.Id, testRun.Name)
	model.Description = testRun.Description
	model.LaunchSource = testRun.LaunchSource
	model.Attachments = buildAssignAttachmentApiModel(testRun.Attachments)
	model.Links = buildUpdateLinkApiModel(testRun.Links)

	return model
}

func buildAssignAttachmentApiModel(attachments []tmsclient.AttachmentApiResult) []tmsclient.AssignAttachmentApiModel {
	updateAttachments := make([]tmsclient.AssignAttachmentApiModel, len(attachments))
	for i, attachment := range attachments {
		updateAttachment := tmsclient.NewAssignAttachmentApiModel(attachment.Id)
		updateAttachments[i] = *updateAttachment
	}

	return updateAttachments
}

func buildUpdateLinkApiModel(links []tmsclient.LinkApiResult) []tmsclient.UpdateLinkApiModel {
	updateLinks := make([]tmsclient.UpdateLinkApiModel, len(links))
	for i, link := range links {
		updateLink := tmsclient.NewUpdateLinkApiModel(link.Url, link.HasInfo)
		updateLink.Id = link.Id
		updateLink.Title = link.Title
		updateLink.Description = link.Description

		updateLinks[i] = *updateLink
	}

	return updateLinks
}
