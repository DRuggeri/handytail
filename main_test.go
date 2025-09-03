package main

import (
	"regexp"
	"testing"
)

func TestProcessLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal text",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "single backspace removes previous character",
			input:    "hello\b world",
			expected: "hell world",
		},
		{
			name:     "multiple backspaces",
			input:    "hello\b\b\b world",
			expected: "he world",
		},
		{
			name:     "backspace at beginning does nothing",
			input:    "\bhello",
			expected: "hello",
		},
		{
			name:     "backspace removes entire string",
			input:    "test\b\b\b\b",
			expected: "",
		},
		{
			name:     "carriage return removed",
			input:    "hello\rworld",
			expected: "helloworld",
		},
		{
			name:     "multiple carriage returns removed",
			input:    "hello\r\r\rworld",
			expected: "helloworld",
		},
		{
			name:     "tab character removed (control char)",
			input:    "hello\tworld",
			expected: "helloworld",
		},
		{
			name:     "newline character removed (control char)",
			input:    "hello\nworld",
			expected: "helloworld",
		},
		{
			name:     "form feed removed (control char)",
			input:    "hello\fworld",
			expected: "helloworld",
		},
		{
			name:     "vertical tab removed (control char)",
			input:    "hello\vworld",
			expected: "helloworld",
		},
		{
			name:     "bell character removed (control char)",
			input:    "hello\aworld",
			expected: "helloworld",
		},
		{
			name:     "escape character removed (control char)",
			input:    "hello\x1bworld",
			expected: "helloworld",
		},
		{
			name:     "complex combination",
			input:    "hello\b\rwo\trld\b\b\b\x1b test",
			expected: "hellwo test",
		},
		{
			name:     "unicode characters preserved",
			input:    "hello 世界",
			expected: "hello 世界",
		},
		{
			name:     "backspace with unicode",
			input:    "hello 世\b界",
			expected: "hello 界",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only control characters",
			input:    "\r\n\t\f\v",
			expected: "",
		},
		{
			name:     "only backspaces",
			input:    "\b\b\b",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processLine(tt.input)
			if result != tt.expected {
				t.Errorf("processLine(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRegexSliceString(t *testing.T) {
	var rs RegexSlice

	// Test empty slice
	if rs.String() != "" {
		t.Errorf("empty RegexSlice.String() = %q, want empty string", rs.String())
	}

	// Add some patterns
	rs.Set("test")
	rs.Set("hello.*world")

	result := rs.String()
	expected := "test,hello.*world"
	if result != expected {
		t.Errorf("RegexSlice.String() = %q, want %q", result, expected)
	}
}

func TestRegexSliceSet(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		shouldError bool
	}{
		{
			name:        "valid simple pattern",
			pattern:     "test",
			shouldError: false,
		},
		{
			name:        "valid regex pattern",
			pattern:     "hello.*world",
			shouldError: false,
		},
		{
			name:        "valid complex pattern",
			pattern:     "^(SUCCESS|COMPLETE).*\\d+$",
			shouldError: false,
		},
		{
			name:        "valid case-insensitive pattern",
			pattern:     "(?i)error",
			shouldError: false,
		},
		{
			name:        "invalid pattern - unclosed bracket",
			pattern:     "[abc",
			shouldError: true,
		},
		{
			name:        "invalid pattern - unclosed paren",
			pattern:     "(abc",
			shouldError: true,
		},
		{
			name:        "invalid pattern - bad escape",
			pattern:     "\\",
			shouldError: true,
		},
		{
			name:        "empty pattern is valid",
			pattern:     "",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rs RegexSlice
			err := rs.Set(tt.pattern)

			if tt.shouldError && err == nil {
				t.Errorf("RegexSlice.Set(%q) expected error but got none", tt.pattern)
			}

			if !tt.shouldError && err != nil {
				t.Errorf("RegexSlice.Set(%q) unexpected error: %v", tt.pattern, err)
			}

			if !tt.shouldError && len(rs) != 1 {
				t.Errorf("RegexSlice.Set(%q) expected 1 pattern, got %d", tt.pattern, len(rs))
			}
		})
	}
}

func TestRegexSliceMultiplePatterns(t *testing.T) {
	var rs RegexSlice

	patterns := []string{"test", "hello.*world", "^SUCCESS"}

	for _, pattern := range patterns {
		err := rs.Set(pattern)
		if err != nil {
			t.Fatalf("unexpected error adding pattern %q: %v", pattern, err)
		}
	}

	if len(rs) != len(patterns) {
		t.Errorf("expected %d patterns, got %d", len(patterns), len(rs))
	}

	// Test that patterns work correctly
	testCases := []struct {
		input   string
		matches []bool // which patterns should match
	}{
		{"test", []bool{true, false, false}},
		{"hello beautiful world", []bool{false, true, false}},
		{"SUCCESS: operation complete", []bool{false, false, true}},
		{"no match here", []bool{false, false, false}},
	}

	for _, tc := range testCases {
		for i, pattern := range rs {
			matches := pattern.MatchString(tc.input)
			if matches != tc.matches[i] {
				t.Errorf("pattern %d (%q) matching %q: got %t, want %t",
					i, pattern.String(), tc.input, matches, tc.matches[i])
			}
		}
	}
}

func TestPatternMatching(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		inputs  []string
		matches []bool
	}{
		{
			name:    "literal matching",
			pattern: "SUCCESS",
			inputs:  []string{"SUCCESS", "FAILURE", "SUCCESS: done", "success", ""},
			matches: []bool{true, false, true, false, false},
		},
		{
			name:    "case insensitive",
			pattern: "(?i)success",
			inputs:  []string{"SUCCESS", "success", "Success", "FAILURE", "successful"},
			matches: []bool{true, true, true, false, true},
		},
		{
			name:    "word boundaries",
			pattern: "\\bSUCCESS\\b",
			inputs:  []string{"SUCCESS", "SUCCESSFUL", "SUCCESS!", "MY_SUCCESS", "SUCCESS: done"},
			matches: []bool{true, false, true, false, true},
		},
		{
			name:    "wildcard matching",
			pattern: "BUILD.*SUCCESS",
			inputs:  []string{"BUILD SUCCESS", "BUILD COMPLETE SUCCESS", "BUILD", "SUCCESS", "BUILD FAILED"},
			matches: []bool{true, true, false, false, false},
		},
		{
			name:    "alternatives",
			pattern: "(ERROR|FAILED|TIMEOUT)",
			inputs:  []string{"ERROR", "FAILED", "TIMEOUT", "SUCCESS", "ERROR: connection failed"},
			matches: []bool{true, true, true, false, true},
		},
		{
			name:    "digit matching",
			pattern: "CODE\\s+\\d+",
			inputs:  []string{"CODE 123", "CODE  456", "CODE ABC", "ERROR CODE 789", "CODE"},
			matches: []bool{true, true, false, true, false},
		},
		{
			name:    "start and end anchors",
			pattern: "^START.*END$",
			inputs:  []string{"START something END", "PREFIX START something END", "START something END SUFFIX", "START END"},
			matches: []bool{true, false, false, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regex, err := regexp.Compile(tt.pattern)
			if err != nil {
				t.Fatalf("failed to compile pattern %q: %v", tt.pattern, err)
			}

			for i, input := range tt.inputs {
				result := regex.MatchString(input)
				expected := tt.matches[i]
				if result != expected {
					t.Errorf("pattern %q matching %q: got %t, want %t",
						tt.pattern, input, result, expected)
				}
			}
		})
	}
}

func TestIntegrationProcessLineWithPatterns(t *testing.T) {
	// Test the combination of processLine and pattern matching
	tests := []struct {
		name           string
		input          string
		successPattern string
		failurePattern string
		expectSuccess  bool
		expectFailure  bool
	}{
		{
			name:           "success after processing",
			input:          "BUILD\b SUCCESS",
			successPattern: "BUIL SUCCESS",
			failurePattern: "FAILED",
			expectSuccess:  true,
			expectFailure:  false,
		},
		{
			name:           "failure after removing control chars",
			input:          "ERROR\r\n\t123",
			successPattern: "SUCCESS",
			failurePattern: "ERROR\\d+",
			expectSuccess:  false,
			expectFailure:  true,
		},
		{
			name:           "backspace reveals success",
			input:          "FAILUCESS\b\b\b\b\b\bSS",
			successPattern: "FAISS",
			failurePattern: "FAILURE",
			expectSuccess:  true,
			expectFailure:  false,
		},
		{
			name:           "carriage return hiding failure",
			input:          "SUCCESS\rFAILED",
			successPattern: "SUCCESS",
			failurePattern: "FAILED",
			expectSuccess:  true,
			expectFailure:  true,
		},
		{
			name:           "no match after processing",
			input:          "SOME\b\b\b\bTEST\r\n",
			successPattern: "SUCCESS",
			failurePattern: "FAILED",
			expectSuccess:  false,
			expectFailure:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processedLine := processLine(tt.input)

			successRegex, err := regexp.Compile(tt.successPattern)
			if err != nil {
				t.Fatalf("failed to compile success pattern: %v", err)
			}

			failureRegex, err := regexp.Compile(tt.failurePattern)
			if err != nil {
				t.Fatalf("failed to compile failure pattern: %v", err)
			}

			successMatch := successRegex.MatchString(processedLine)
			failureMatch := failureRegex.MatchString(processedLine)

			if successMatch != tt.expectSuccess {
				t.Errorf("success pattern matching: got %t, want %t (processed: %q)",
					successMatch, tt.expectSuccess, processedLine)
			}

			if failureMatch != tt.expectFailure {
				t.Errorf("failure pattern matching: got %t, want %t (processed: %q)",
					failureMatch, tt.expectFailure, processedLine)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkProcessLine(b *testing.B) {
	testLine := "hello\b\b\rworld\ttest\b\b\nmore text"
	for i := 0; i < b.N; i++ {
		processLine(testLine)
	}
}

func BenchmarkRegexMatching(b *testing.B) {
	regex := regexp.MustCompile("BUILD.*SUCCESS")
	testLine := "BUILD COMPLETE SUCCESS"
	for i := 0; i < b.N; i++ {
		regex.MatchString(testLine)
	}
}

func BenchmarkCompleteProcessing(b *testing.B) {
	regex := regexp.MustCompile("BUILD.*SUCCESS")
	testLine := "BUILD\b\r COMPLETE\t SUCCESS\n"
	for i := 0; i < b.N; i++ {
		processed := processLine(testLine)
		regex.MatchString(processed)
	}
}
