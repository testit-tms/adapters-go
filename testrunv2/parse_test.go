package testrunv2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse_keepsMergeFieldsIgnoresExtras(t *testing.T) {
	raw := []byte(`{
		"id": "b74aa525-7ffe-4269-9afc-86284b946bcc",
		"name": "TestRun_2026-09-11T09:24:35",
		"description": "213123123",
		"launchSource": "Test IT",
		"projectId": "28e826ff-69e4-4752-92f3-9e22f275a9e4",
		"testResults": [{"id": "ignored"}],
		"tags": ["smoke"],
		"links": [{
			"id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"url": "https://existing.example",
			"title": "Old",
			"type": "Related"
		}],
		"attachments": [{
			"id": "11111111-2222-3333-4444-555555555555",
			"fileId": "file-1",
			"type": "text/plain",
			"size": 1,
			"name": "note.txt",
			"createdDate": "2026-09-11T06:24:36.04Z"
		}]
	}`)

	got, err := Parse(raw)
	require.NoError(t, err)
	require.Equal(t, "b74aa525-7ffe-4269-9afc-86284b946bcc", got.Id)
	require.Equal(t, "213123123", got.Description)
	require.Equal(t, "Test IT", got.LaunchSource)
	require.Equal(t, []string{"smoke"}, got.Tags)
	require.Equal(t, []Link{{
		Id:    "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
		Title: "Old",
		Url:   "https://existing.example",
		Type:  "Related",
	}}, got.Links)
	require.Equal(t, []Attachment{{Id: "11111111-2222-3333-4444-555555555555"}}, got.Attachments)
}
