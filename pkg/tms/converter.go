package tms

import (
	tmsclient "github.com/testit-tms/api-client-golang"
)

func testToAutotestModel(test testResult, projectId string) tmsclient.CreateAutoTestRequest {
	req := tmsclient.NewCreateAutoTestRequest(test.externalId, projectId, test.displayName)
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

	if len(test.links) != 0 {
		links := make([]tmsclient.LinkPostModel, 0, len(test.links))
		for _, link := range test.links {
			l := tmsclient.NewLinkPostModel(link.Url)
			l.SetTitle(link.Title)
			l.SetDescription(link.Description)

			if link.LinkType != "" {
				linkType, err := tmsclient.NewLinkTypeFromValue(string(link.LinkType))
				if err != nil {
					logger.Error("Error converting link type", "error", err)
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

	return *req
}

func stepToAutoTestStepModel(s []step) []tmsclient.AutoTestStepModel {
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

func testToUpdateAutotestModel(test testResult, autotest tmsclient.AutoTestModel) tmsclient.UpdateAutoTestRequest {
	req := tmsclient.NewUpdateAutoTestRequest(test.externalId, autotest.ProjectId, test.displayName)

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
			l := tmsclient.NewLinkPutModel(link.Url)
			l.SetTitle(link.Title)
			l.SetDescription(link.Description)

			if link.LinkType != "" {
				linkType, err := tmsclient.NewLinkTypeFromValue(string(link.LinkType))
				if err != nil {
					logger.Error("Error converting link type", "error", err)
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

	req.SetIsFlaky(*autotest.IsFlaky.Get())
	req.SetId(*autotest.Id)

	return *req
}

func testToResultModel(test testResult, confID string) ([]tmsclient.AutoTestResultsForTestRunModel, error) {
	outcome, err := tmsclient.NewAvailableTestResultOutcomeFromValue(test.status)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		req.SetStepResults(steps)
	}

	if len(test.resultLinks) != 0 {
		links := make([]tmsclient.LinkPostModel, 0, len(test.resultLinks))
		for _, link := range test.resultLinks {
			l := tmsclient.NewLinkPostModel(link.Url)
			l.SetTitle(link.Title)
			l.SetDescription(link.Description)
			if link.LinkType != "" {
				linkType, err := tmsclient.NewLinkTypeFromValue(string(link.LinkType))
				if err != nil {
					logger.Error("Error converting link type", "error", err)
				} else {
					l.SetType(*linkType)
				}
			}
			links = append(links, *l)
		}
		req.SetLinks(links)
	}

	return []tmsclient.AutoTestResultsForTestRunModel{*req}, nil
}

func stepToAttachmentPutModelAutoTestStepResultsModel(s []step) ([]tmsclient.AttachmentPutModelAutoTestStepResultsModel, error) {
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

		if len(step.childrenSteps) != 0 {
			cs, err := stepToAttachmentPutModelAutoTestStepResultsModel(step.childrenSteps)
			if err != nil {
				return nil, err
			}
			model.SetStepResults(cs)
		}

		steps = append(steps, *model)
	}

	return steps, nil
}
