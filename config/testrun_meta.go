package config

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TestRunLink is a link attached to a test run (not to an autotest/result).
type TestRunLink struct {
	Url         string `json:"url"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

// TestRunTags accepts JSON array, JSON-array-as-string, or comma-separated values.
type TestRunTags []string

func (t *TestRunTags) UnmarshalJSON(data []byte) error {
	data = bytesTrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		*t = nil
		return nil
	}
	if data[0] == '[' {
		var items []string
		if err := json.Unmarshal(data, &items); err != nil {
			return fmt.Errorf("testRunTags: %w", err)
		}
		*t = normalizeTags(items)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("testRunTags: %w", err)
	}
	return t.UnmarshalText([]byte(s))
}

func (t *TestRunTags) UnmarshalText(text []byte) error {
	*t = parseTags(string(text))
	return nil
}

// SetValue is used by cleanenv for environment variables.
func (t *TestRunTags) SetValue(s string) error {
	return t.UnmarshalText([]byte(s))
}

// TestRunLinks accepts a JSON array of link objects (or the same JSON as a string for env).
type TestRunLinks []TestRunLink

func (l *TestRunLinks) UnmarshalJSON(data []byte) error {
	data = bytesTrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		*l = nil
		return nil
	}
	if data[0] == '[' {
		var items []TestRunLink
		if err := json.Unmarshal(data, &items); err != nil {
			return fmt.Errorf("testRunLinks: %w", err)
		}
		*l = normalizeLinks(items)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("testRunLinks: %w", err)
	}
	return l.UnmarshalText([]byte(s))
}

func (l *TestRunLinks) UnmarshalText(text []byte) error {
	s := strings.TrimSpace(string(text))
	if s == "" {
		*l = nil
		return nil
	}
	var items []TestRunLink
	if err := json.Unmarshal([]byte(s), &items); err != nil {
		return fmt.Errorf("testRunLinks: expected JSON array: %w", err)
	}
	*l = normalizeLinks(items)
	return nil
}

// SetValue is used by cleanenv for environment variables.
func (l *TestRunLinks) SetValue(s string) error {
	return l.UnmarshalText([]byte(s))
}

func parseTags(raw string) TestRunTags {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	if strings.HasPrefix(raw, "[") {
		var items []string
		if err := json.Unmarshal([]byte(raw), &items); err == nil {
			return normalizeTags(items)
		}
	}
	parts := strings.Split(raw, ",")
	return normalizeTags(parts)
}

func normalizeTags(items []string) TestRunTags {
	out := make(TestRunTags, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		tag := strings.TrimSpace(item)
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		out = append(out, tag)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func normalizeLinks(items []TestRunLink) TestRunLinks {
	out := make(TestRunLinks, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		url := strings.TrimSpace(item.Url)
		if url == "" {
			continue
		}
		if _, ok := seen[url]; ok {
			continue
		}
		seen[url] = struct{}{}
		linkType := strings.TrimSpace(item.Type)
		if linkType == "" {
			linkType = "Related"
		}
		out = append(out, TestRunLink{
			Url:         url,
			Title:       strings.TrimSpace(item.Title),
			Description: strings.TrimSpace(item.Description),
			Type:        linkType,
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// MergeTags keeps existing tags and appends new ones without duplicates.
func MergeTags(existing, added []string) []string {
	out := make([]string, 0, len(existing)+len(added))
	seen := make(map[string]struct{}, len(existing)+len(added))
	for _, tag := range existing {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		out = append(out, tag)
	}
	for _, tag := range added {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		out = append(out, tag)
	}
	return out
}

// MergeLinkURLs returns true if url is already present in existingURLs.
func HasLinkURL(existingURLs []string, url string) bool {
	for _, existing := range existingURLs {
		if existing == url {
			return true
		}
	}
	return false
}

func bytesTrimSpace(data []byte) []byte {
	return []byte(strings.TrimSpace(string(data)))
}
