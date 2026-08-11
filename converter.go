package tms

import (
	"fmt"
	"strconv"
	"strings"

	tmsclient "github.com/testit-tms/adapters-go/v2/adaptersapi"
	"github.com/testit-tms/adapters-go/v2/config"
	"github.com/testit-tms/adapters-go/v2/htmlutils"
)

// TODO: validate that hasInfo always true is correct
const defaultHasInfo = true
const defaultLinkType = tmsclient.LINKTYPE_RELATED

func testToAutotestModel(test TestResult, projectId string) tmsclient.AutoTestCreateApiModel {
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

	if len(test.tags) != 0 {
		req.SetTags(test.tags)
	}

	if len(test.links) != 0 {
		links := make([]tmsclient.LinkCreateApiModel, 0, len(test.links))
		for _, link := range test.links {
			linkType := defaultLinkType
			if link.LinkType != "" {
				parsedLinkType, err := tmsclient.NewLinkTypeFromValue(string(link.LinkType))
				if err != nil {
					logger.Error("error converting link type", "error", err)
				} else {
					linkType = *parsedLinkType
				}
			}

			l := tmsclient.NewLinkCreateApiModel(link.Url, linkType)
			l.SetTitle(link.Title)
			l.SetDescription(link.Description)

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

func stepToAutoTestStepModel(s []StepResult) []tmsclient.AutoTestStepApiModel {
	steps := make([]tmsclient.AutoTestStepApiModel, 0, len(s))
	for _, step := range s {
		model := tmsclient.NewAutoTestStepApiModel(step.Name)
		model.SetDescription(step.Description)

		if len(step.ChildrenSteps) != 0 {
			model.SetSteps(stepToAutoTestStepModel(step.ChildrenSteps))
		}

		steps = append(steps, *model)
	}

	// Apply HTML escaping to the steps slice
	htmlutils.EscapeHtmlInObjectSlice(steps)

	return steps
}

func testToUpdateAutotestModel(test TestResult, autotest tmsclient.AutoTestApiResult) tmsclient.AutoTestUpdateApiModel {
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

	if len(test.tags) != 0 {
		req.SetTags(test.tags)
	}

	if test.title != "" {
		req.SetTitle(test.title)
	}

	if len(test.links) != 0 {
		links := make([]tmsclient.LinkUpdateApiModel, 0, len(test.links))
		for _, link := range test.links {
			linkType := defaultLinkType
			if link.LinkType != "" {
				parsedLinkType, err := tmsclient.NewLinkTypeFromValue(string(link.LinkType))
				if err != nil {
					logger.Error("error converting link type", "error", err)
				} else {
					linkType = *parsedLinkType
				}
			}

			l := tmsclient.NewLinkUpdateApiModel(link.Url, linkType)
			l.SetTitle(link.Title)
			l.SetDescription(link.Description)

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

// passed failed skipped from framework
func mapType(status string) tmsclient.TestStatusType {
	status = strings.ToLower(status)
	switch status {
	case "passed":
		return tmsclient.TESTSTATUSTYPE_SUCCEEDED
	case "failed":
		return tmsclient.TESTSTATUSTYPE_FAILED
	case "skipped":
		return tmsclient.TESTSTATUSTYPE_INCOMPLETE
	case "inprogress":
		return tmsclient.TESTSTATUSTYPE_IN_PROGRESS
	}
	return tmsclient.TESTSTATUSTYPE_INCOMPLETE
}

func testToResultModel(test TestResult, confID string) ([]tmsclient.AutoTestResultsForTestRunModel, error) {
	req := tmsclient.NewAutoTestResultsForTestRunModel(confID, test.externalId)
	req.SetStatusType(mapType(test.status))
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
			linkType := defaultLinkType
			if link.LinkType != "" {
				parsedLinkType, err := tmsclient.NewLinkTypeFromValue(string(link.LinkType))
				if err != nil {
					logger.Error("error converting link type", "error", err)
				} else {
					linkType = *parsedLinkType
				}
			}

			l := tmsclient.NewLinkPostModel(link.Url, linkType, defaultHasInfo)
			l.SetTitle(link.Title)
			l.SetDescription(link.Description)
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

func stepToAttachmentPutModelAutoTestStepResultsModel(s []StepResult) ([]tmsclient.AttachmentPutModelAutoTestStepResultsModel, error) {
	steps := make([]tmsclient.AttachmentPutModelAutoTestStepResultsModel, 0, len(s))
	for _, step := range s {
		model := tmsclient.NewAttachmentPutModelAutoTestStepResultsModel()
		model.SetTitle(step.Name)
		model.SetDescription(step.Description)
		outcome, err := tmsclient.NewAvailableTestResultOutcomeFromValue(step.Status)
		if err != nil {
			return nil, err
		}
		model.SetOutcome(*outcome)
		model.SetStartedOn(step.StartedOn)
		model.SetCompletedOn(step.CompletedOn)
		model.SetDuration(step.Duration)

		if len(step.Attachments) != 0 {
			attachs := make([]tmsclient.AttachmentPutModel, 0, len(step.Attachments))
			for _, attach := range step.Attachments {
				a := tmsclient.NewAttachmentPutModel(attach)
				attachs = append(attachs, *a)
			}
			model.SetAttachments(attachs)
		}

		if len(step.ChildrenSteps) != 0 {
			cs, err := stepToAttachmentPutModelAutoTestStepResultsModel(step.ChildrenSteps)
			if err != nil {
				return nil, err
			}
			model.SetStepResults(cs)
		}

		if len(step.Parameters) != 0 {
			params := make(map[string]string, len(step.Parameters))
			for k, v := range step.Parameters {
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

func testToUpdateResultModel(model *tmsclient.TestResultResponse, test TestResult) (tmsclient.TestResultUpdateRequest, error) {
	tearDownsAttachments, err := stepToAttachmentPutModelAutoTestStepResultsModel(test.teardowns)
	if err != nil {
		return tmsclient.TestResultUpdateRequest{}, err
	}

	tearDowns, err := mapAttachmentsToStepResults(tearDownsAttachments)
	if err != nil {
		return tmsclient.TestResultUpdateRequest{}, fmt.Errorf("error mapping tearDowns: %w", err)
	}

	setupsAttachments, err := stepToAttachmentPutModelAutoTestStepResultsModel(test.setups)
	if err != nil {
		return tmsclient.TestResultUpdateRequest{}, err
	}

	setups, err := mapAttachmentsToStepResults(setupsAttachments)
	if err != nil {
		return tmsclient.TestResultUpdateRequest{}, fmt.Errorf("error mapping setups: %w", err)
	}

	req := tmsclient.NewTestResultUpdateRequest()
	req.SetTeardownResults(tearDowns)
	req.SetSetupResults(setups)
	req.SetDuration(model.GetDurationInMs())
	req.SetLinks(mapLinkApiResultsToCreateLinkApiModels(model.GetLinks()))
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

	req.SetStatusType(mapType(test.status))

	// Apply HTML escaping to the update request
	htmlutils.EscapeHtmlInObject(req)

	return *req, nil
}

func mapLinkApiResultsToCreateLinkApiModels(links []tmsclient.LinkApiResult) []tmsclient.CreateLinkApiModel {
	result := make([]tmsclient.CreateLinkApiModel, 0, len(links))
	for _, link := range links {
		m := tmsclient.NewCreateLinkApiModel(link.Url, link.Type)
		if link.Title.IsSet() {
			m.SetTitle(link.GetTitle())
		}
		if link.Description.IsSet() {
			m.SetDescription(link.GetDescription())
		}
		result = append(result, *m)
	}
	return result
}

func buildUpdateEmptyTestRunApiModel(testRun *tmsclient.TestRunApiResult) *tmsclient.UpdateEmptyTestRunApiModel {
	model := tmsclient.NewUpdateEmptyTestRunApiModel(testRun.Id, testRun.Name)
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
		updateLink := tmsclient.NewUpdateLinkApiModel(link.Url, link.Type)
		if link.Id.IsSet() {
			updateLink.SetId(link.GetId())
		}
		if link.Title.IsSet() {
			updateLink.SetTitle(link.GetTitle())
		}
		if link.Description.IsSet() {
			updateLink.SetDescription(link.GetDescription())
		}

		updateLinks[i] = *updateLink
	}

	return updateLinks
}

func configLinksToCreateLinkApiModels(links []config.TestRunLink) []tmsclient.CreateLinkApiModel {
	result := make([]tmsclient.CreateLinkApiModel, 0, len(links))
	for _, link := range links {
		createLink, ok := configLinkToCreateLinkApiModel(link)
		if !ok {
			continue
		}
		result = append(result, createLink)
	}
	return result
}

func configLinkToCreateLinkApiModel(link config.TestRunLink) (tmsclient.CreateLinkApiModel, bool) {
	linkType, err := tmsclient.NewLinkTypeFromValue(link.Type)
	if err != nil {
		linkType = tmsclient.LINKTYPE_RELATED.Ptr()
	}
	if link.Url == "" || linkType == nil {
		return tmsclient.CreateLinkApiModel{}, false
	}
	m := tmsclient.NewCreateLinkApiModel(link.Url, *linkType)
	if link.Title != "" {
		m.SetTitle(link.Title)
	}
	if link.Description != "" {
		m.SetDescription(link.Description)
	}
	return *m, true
}

func configLinkToUpdateLinkApiModel(link config.TestRunLink) (tmsclient.UpdateLinkApiModel, bool) {
	linkType, err := tmsclient.NewLinkTypeFromValue(link.Type)
	if err != nil {
		linkType = tmsclient.LINKTYPE_RELATED.Ptr()
	}
	if link.Url == "" || linkType == nil {
		return tmsclient.UpdateLinkApiModel{}, false
	}
	m := tmsclient.NewUpdateLinkApiModel(link.Url, *linkType)
	if link.Title != "" {
		m.SetTitle(link.Title)
	}
	if link.Description != "" {
		m.SetDescription(link.Description)
	}
	return *m, true
}
