package files

import "testing"

func TestTitleToFilename(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{"SimpleTitle", "SimpleTitle"},
		{"A Sample: Invalid/File*Name?", "A Sample  Invalid File Name"},
		{"Hello World", "Hello World"},
		{"<Invalid*Name>", "Invalid Name"},
		{"Special:Char*", "Special Char"},
		{"   Leading and trailing spaces   ", "Leading and trailing spaces"},
		{"Mixed<Chars*>In/Title", "Mixed Chars  In Title"},
		{"", ""},                            // Empty input
		{"Control\x00Char", "Control Char"}, // Control characters
	}

	for _, tt := range tests {
		result := ToValidFilename(tt.title)
		if result != tt.expected {
			t.Errorf("For title %q, expected filename %q but got %q", tt.title, tt.expected, result)
		}
	}
}
