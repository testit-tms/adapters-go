package testrunv2

import (
	"encoding/json"
	"fmt"
)

// Snapshot is the subset of GET /api/v2/testRuns/{id} fields needed for metadata merge.
type Snapshot struct {
	Id           string
	Name         string
	Description  string
	LaunchSource string
	Tags         []string
	Links        []Link
	Attachments  []Attachment
}

type Link struct {
	Id          string
	Title       string
	Url         string
	Description string
	Type        string
}

type Attachment struct {
	Id string
}

// Parse decodes a v2 test-run payload, ignoring unrelated fields (testResults, etc.).
func Parse(data []byte) (*Snapshot, error) {
	var raw struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		LaunchSource string `json:"launchSource"`
		Tags         []string `json:"tags"`
		Links        []struct {
			Id          *string `json:"id"`
			Title       *string `json:"title"`
			Url         string  `json:"url"`
			Description *string `json:"description"`
			Type        string  `json:"type"`
		} `json:"links"`
		Attachments []struct {
			Id string `json:"id"`
		} `json:"attachments"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	if raw.Id == "" {
		return nil, fmt.Errorf("test run id is empty")
	}

	links := make([]Link, 0, len(raw.Links))
	for _, l := range raw.Links {
		if l.Url == "" {
			continue
		}
		link := Link{Url: l.Url, Type: l.Type}
		if l.Id != nil {
			link.Id = *l.Id
		}
		if l.Title != nil {
			link.Title = *l.Title
		}
		if l.Description != nil {
			link.Description = *l.Description
		}
		links = append(links, link)
	}

	attachments := make([]Attachment, 0, len(raw.Attachments))
	for _, a := range raw.Attachments {
		if a.Id == "" {
			continue
		}
		attachments = append(attachments, Attachment{Id: a.Id})
	}

	return &Snapshot{
		Id:           raw.Id,
		Name:         raw.Name,
		Description:  raw.Description,
		LaunchSource: raw.LaunchSource,
		Tags:         raw.Tags,
		Links:        links,
		Attachments:  attachments,
	}, nil
}
