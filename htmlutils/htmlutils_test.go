package htmlutils

import (
	"os"
	"testing"
)

func TestEscapeHtmlTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "String without HTML tags",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			name:     "String with HTML tags",
			input:    "Hello <b>world</b>",
			expected: "Hello &lt;b&gt;world&lt;/b&gt;",
		},
		{
			name:     "String with self-closing tag",
			input:    "Hello <br/> world",
			expected: "Hello \\<br/&gt; world",
		},
		{
			name:     "String with already escaped characters",
			input:    "Hello \\<b\\>world\\</b\\>",
			expected: "Hello \\&lt;b\\&gt;world\\&lt;/b\\&gt;",
		},
		{
			name:     "String with script tag",
			input:    "<script>alert('xss')</script>",
			expected: "&lt;script&gt;alert('xss')&lt;/script&gt;",
		},
		{
			name:     "String with multiple tags",
			input:    "<div><p>Hello</p></div>",
			expected: "&lt;div&gt;&lt;p&gt;Hello&lt;/p&gt;&lt;/div&gt;",
		},
		{
			name:     "String with attributes",
			input:    "<a href='http://example.com'>Link</a>",
			expected: "&lt;a href='http://example.com'&gt;Link&lt;/a&gt;",
		},
		{
			name:     "String with only < or > (no complete tags)",
			input:    "5 < 10 and 10 > 5",
			expected: "5 < 10 and 10 > 5", // No HTML tags detected
		},
		{
			name:     "String with incomplete tag",
			input:    "Hello < world",
			expected: "Hello < world", // No complete HTML tags
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeHtmlTags(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeHtmlTags(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestEscapeHtmlInObject(t *testing.T) {
	type TestStruct struct {
		Name        string
		Description string
		Count       int
		IsActive    bool
	}

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
		},
		{
			name: "Struct with HTML in string fields",
			input: &TestStruct{
				Name:        "Test <b>Name</b>",
				Description: "<script>alert('xss')</script>",
				Count:       5,
				IsActive:    true,
			},
			expected: &TestStruct{
				Name:        "Test &lt;b&gt;Name&lt;/b&gt;",
				Description: "&lt;script&gt;alert('xss')&lt;/script&gt;",
				Count:       5,
				IsActive:    true,
			},
		},
		{
			name: "Struct without HTML",
			input: &TestStruct{
				Name:        "Test Name",
				Description: "Simple description",
				Count:       5,
				IsActive:    true,
			},
			expected: &TestStruct{
				Name:        "Test Name",
				Description: "Simple description",
				Count:       5,
				IsActive:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeHtmlInObject(tt.input)

			if tt.input == nil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
				return
			}

			// Compare struct fields
			if inputStruct, ok := tt.input.(*TestStruct); ok {
				expectedStruct := tt.expected.(*TestStruct)

				if inputStruct.Name != expectedStruct.Name {
					t.Errorf("Name: got %q, want %q", inputStruct.Name, expectedStruct.Name)
				}
				if inputStruct.Description != expectedStruct.Description {
					t.Errorf("Description: got %q, want %q", inputStruct.Description, expectedStruct.Description)
				}
				if inputStruct.Count != expectedStruct.Count {
					t.Errorf("Count: got %d, want %d", inputStruct.Count, expectedStruct.Count)
				}
				if inputStruct.IsActive != expectedStruct.IsActive {
					t.Errorf("IsActive: got %t, want %t", inputStruct.IsActive, expectedStruct.IsActive)
				}
			}
		})
	}
}

func TestEscapeHtmlInObjectSlice(t *testing.T) {
	type TestStruct struct {
		Name string
		Desc string
	}

	tests := []struct {
		name     string
		input    interface{}
		expected []TestStruct
	}{
		{
			name:     "Nil slice",
			input:    nil,
			expected: nil,
		},
		{
			name: "Slice with HTML in strings",
			input: []TestStruct{
				{Name: "Test <b>1</b>", Desc: "<script>alert(1)</script>"},
				{Name: "Test <i>2</i>", Desc: "<div>content</div>"},
			},
			expected: []TestStruct{
				{Name: "Test \\<b&gt;1\\</b&gt;", Desc: "\\<script&gt;alert(1)\\</script&gt;"},
				{Name: "Test \\<i&gt;2\\</i&gt;", Desc: "\\<div&gt;content\\</div&gt;"},
			},
		},
		{
			name: "Slice without HTML",
			input: []TestStruct{
				{Name: "Test 1", Desc: "Description 1"},
				{Name: "Test 2", Desc: "Description 2"},
			},
			expected: []TestStruct{
				{Name: "Test 1", Desc: "Description 1"},
				{Name: "Test 2", Desc: "Description 2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeHtmlInObjectSlice(tt.input)

			if tt.input == nil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
				return
			}

			// Convert result back to slice
			if inputSlice, ok := tt.input.([]TestStruct); ok {
				for i, item := range inputSlice {
					if item.Name != tt.expected[i].Name {
						t.Errorf("Item %d Name: got %q, want %q", i, item.Name, tt.expected[i].Name)
					}
					if item.Desc != tt.expected[i].Desc {
						t.Errorf("Item %d Desc: got %q, want %q", i, item.Desc, tt.expected[i].Desc)
					}
				}
			}
		})
	}
}

func TestEscapeHtmlWithMap(t *testing.T) {
	input := map[string]string{
		"title":       "Test <b>Title</b>",
		"description": "<script>alert('xss')</script>",
		"plain":       "Plain text",
	}

	expected := map[string]string{
		"title":       "Test &lt;b&gt;Title&lt;/b&gt;",
		"description": "&lt;script&gt;alert('xss')&lt;/script&gt;",
		"plain":       "Plain text",
	}

	EscapeHtmlInObject(input)

	for key, expectedValue := range expected {
		if input[key] != expectedValue {
			t.Errorf("Map[%s]: got %q, want %q", key, input[key], expectedValue)
		}
	}
}

func TestEscapeHtmlDisabledByEnvVar(t *testing.T) {
	// Set environment variable to disable escaping
	originalValue := os.Getenv(NoEscapeHtmlEnvVar)
	os.Setenv(NoEscapeHtmlEnvVar, "true")
	defer os.Setenv(NoEscapeHtmlEnvVar, originalValue)

	input := "Hello <b>world</b>"
	result := EscapeHtmlTags(input)

	// Should return original string without escaping
	if result != input {
		t.Errorf("Expected escaping to be disabled, got %q, want %q", result, input)
	}

	// Test with object
	type TestStruct struct {
		Name string
	}

	obj := &TestStruct{Name: "Test <b>Name</b>"}
	EscapeHtmlInObject(obj)

	if obj.Name != "Test <b>Name</b>" {
		t.Errorf("Expected escaping to be disabled for object, got %q", obj.Name)
	}
}

func TestNestedStructs(t *testing.T) {
	type NestedStruct struct {
		Title string
	}

	type MainStruct struct {
		Name   string
		Nested NestedStruct
	}

	input := &MainStruct{
		Name:   "Main <b>Name</b>",
		Nested: NestedStruct{Title: "Nested <i>Title</i>"},
	}

	EscapeHtmlInObject(input)

	expectedMainName := "Main &lt;b&gt;Name&lt;/b&gt;"
	expectedNestedTitle := "Nested &lt;i&gt;Title&lt;/i&gt;"

	if input.Name != expectedMainName {
		t.Errorf("Main Name: got %q, want %q", input.Name, expectedMainName)
	}

	if input.Nested.Title != expectedNestedTitle {
		t.Errorf("Nested Title: got %q, want %q", input.Nested.Title, expectedNestedTitle)
	}
}

func TestDoubleEscapingPrevention(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Already escaped characters should not be double escaped",
			input:    "Hello \\<b\\>world\\</b\\>",
			expected: "Hello \\<b\\&gt;world\\</b\\&gt;",
		},
		{
			name:     "Mix of escaped and unescaped",
			input:    "Hello \\<b>world</b>",
			expected: "Hello \\<b\\&gt;world\\</b\\&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeHtmlTags(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeHtmlTags(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkEscapeHtmlTags(b *testing.B) {
	input := "Hello <b>world</b> with <script>alert('test')</script>"

	for i := 0; i < b.N; i++ {
		EscapeHtmlTags(input)
	}
}

func BenchmarkEscapeHtmlInObject(b *testing.B) {
	type TestStruct struct {
		Name        string
		Description string
		Content     string
	}

	input := &TestStruct{
		Name:        "Test <b>Name</b>",
		Description: "<script>alert('xss')</script>",
		Content:     "Some <div>content</div> here",
	}

	for i := 0; i < b.N; i++ {
		// Reset the struct for each iteration
		input.Name = "Test <b>Name</b>"
		input.Description = "<script>alert('xss')</script>"
		input.Content = "Some <div>content</div> here"

		EscapeHtmlInObject(input)
	}
}
