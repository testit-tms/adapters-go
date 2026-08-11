package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTestRunTags_UnmarshalJSON_array(t *testing.T) {
	var tags TestRunTags
	require.NoError(t, json.Unmarshal([]byte(`["smoke","nightly","smoke"]`), &tags))
	require.Equal(t, TestRunTags{"smoke", "nightly"}, tags)
}

func TestTestRunTags_UnmarshalJSON_commaString(t *testing.T) {
	var tags TestRunTags
	require.NoError(t, json.Unmarshal([]byte(`"smoke, nightly ,"`), &tags))
	require.Equal(t, TestRunTags{"smoke", "nightly"}, tags)
}

func TestTestRunTags_UnmarshalText_jsonArrayString(t *testing.T) {
	var tags TestRunTags
	require.NoError(t, tags.UnmarshalText([]byte(`["a","b"]`)))
	require.Equal(t, TestRunTags{"a", "b"}, tags)
}

func TestTestRunTags_empty(t *testing.T) {
	var tags TestRunTags
	require.NoError(t, json.Unmarshal([]byte(`""`), &tags))
	require.Nil(t, tags)
	require.NoError(t, tags.UnmarshalText([]byte("")))
	require.Nil(t, tags)
}

func TestTestRunLinks_UnmarshalJSON(t *testing.T) {
	var links TestRunLinks
	raw := `[
		{"url":"https://ci.example/job/1","title":"CI Job"},
		{"url":"https://ci.example/job/1","title":"dup"},
		{"url":"","title":"skip"},
		{"url":"https://ci.example/job/2","type":"Issue"}
	]`
	require.NoError(t, json.Unmarshal([]byte(raw), &links))
	require.Equal(t, TestRunLinks{
		{Url: "https://ci.example/job/1", Title: "CI Job", Type: "Related"},
		{Url: "https://ci.example/job/2", Type: "Issue"},
	}, links)
}

func TestTestRunLinks_UnmarshalText_invalid(t *testing.T) {
	var links TestRunLinks
	err := links.UnmarshalText([]byte(`not-json`))
	require.Error(t, err)
}

func TestMergeTags(t *testing.T) {
	got := MergeTags([]string{"ui", "smoke"}, []string{"smoke", "nightly"})
	require.Equal(t, []string{"ui", "smoke", "nightly"}, got)
}

func TestHasLinkURL(t *testing.T) {
	require.True(t, HasLinkURL([]string{"https://a"}, "https://a"))
	require.False(t, HasLinkURL([]string{"https://a"}, "https://b"))
}
