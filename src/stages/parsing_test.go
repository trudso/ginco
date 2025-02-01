package stages

import (
	"testing"
)

func TestPopWord(t *testing.T) {
	testCases := []struct {
		content      string
		expectedWord string
		expectedRest string
	}{
		{"", "", ""},
		{"test", "test", ""},
		{" hello", "hello", ""},
		{"hello world, dude", "hello", "world, dude"},
		{"line\nbreak", "line", "break"},
	}

	for _, tc := range testCases {
		word, rest := popWord(tc.content)
		if (word != tc.expectedWord) || (rest != tc.expectedRest) {
			t.Errorf("popWord(%q) = %q, %q; expected %q, %q",
				tc.content, word, rest, tc.expectedWord, tc.expectedRest)
		}
	}
}

func TestPopScope(t *testing.T) {
	testCases := []struct {
		content       string
		expectedScope string
		expectedRest  string
		expectedError error
	}{
		{"{{{}}}", "{{}}", "", nil},
		{"{}", "", "", nil},
		{"{{{}", "", "", ERR_NO_MATCHING_SCOPE_BRACKETS},
		{"{ this is the scope } and this is the rest", "this is the scope", "and this is the rest", nil},
		{"", "", "", ERR_NO_SCOPE_FOUND},
		{" { inner thing { with some stuff }} and some other stuff", "inner thing { with some stuff }", "and some other stuff", nil},
		{" { inner thing { with some { stuff }} and some other stuff", "", "", ERR_NO_MATCHING_SCOPE_BRACKETS},
		{"something first { some inner scope stuff }", "", "", ERR_NO_SCOPE_FOUND},
	}

	for _, tc := range testCases {
		scope, rest, err := popScope(tc.content)

		if scope != tc.expectedScope || rest != tc.expectedRest || err != tc.expectedError {
			t.Errorf("popScope(%q) = %q, %q, %q; expected %q, %q, %q",
				tc.content, scope, rest, err, tc.expectedScope, tc.expectedRest, tc.expectedError)
		}
	}
}

func TestPopExpectedToken(t *testing.T) {
	testCases := []struct {
		content       string
		expectations  []string
		expectedToken string
		expectedRest  string
		expectedError error
	}{
		{"", []string{"@", "package", "model"}, "", "", ERR_EXPECTED_ELEMENT_NOT_FOUND},
		{"some content", []string{"@", "package", "model"}, "", "some content", ERR_EXPECTED_ELEMENT_NOT_FOUND},
		{"model TestModel {", []string{"model", "package"}, "model", "TestModel {", nil},
		{" fields {", []string{"fields"}, "fields", "{", nil},
		{"	=1 Name string", []string{"@", "-", "="}, "=", "1 Name string", nil},
		{"	1 Name string", []string{"0", "1", "*"}, "1", "Name string", nil},
	}

	for _, tc := range testCases {
		token, rest, err := popExpectedToken(tc.content, tc.expectations)

		if token != tc.expectedToken || rest != tc.expectedRest || err != tc.expectedError {
			t.Errorf("popScope(%q, %q) = %q, %q, %q; expected %q, %q, %q",
				tc.content, tc.expectations,
				token, rest, err,
				tc.expectedToken, tc.expectedRest, tc.expectedError)
		}
	}
}

func TestPeekWord(t *testing.T) {
	testCases := []struct {
		content      string
		expectedWord string
	}{
		{"one", "one"},
		{"", ""},
		{"one two", "one"},
		{"one two three", "one"},
	}

	for _, tc := range testCases {
		word := peekWord(tc.content)
		if word != tc.expectedWord {
			t.Errorf("peekWord(%q) = %q; expected %q", tc.content, word, tc.expectedWord)
		}
	}
}

func TestPeekNumber(t *testing.T) {
	testCases := []struct {
		content        string
		expectedNumber int
		expectedError  error
	}{
		{"123", 123, nil},
		{"123ABC", 123, nil},
		{"A123BC", 0, ERR_NOT_A_NUMBER},
	}

	for _, tc := range testCases {
		number, err := peekNumber(tc.content)
		if number != tc.expectedNumber || err != tc.expectedError {
			t.Errorf("peekNumber(%q) = %q, %q; expected %q, %q", tc.content, number, err, tc.expectedNumber, tc.expectedError)
		}
	}
}

func TestPeekToken(t *testing.T) {
	testCases := []struct {
		content           string
		expectedTokenType int
		expectedToken     string
		expectedError     error
	}{
		{"", EOF, "", nil},
		{" \n\n ", EOF, "", nil},
		{"package test", IDENTIFIER, "package", nil},
		{" package test", IDENTIFIER, "package", nil},
		{"@SomeTrait", TRAIT, "@", nil},
		{"{ some scope }", SCOPE_BEGIN, "{", nil},
		{" }\n model test {", SCOPE_END, "}", nil},
		{"# some commment", COMMENT, "#", nil},
		{"123", NUMBER, "123", nil},
		{"*123", SYMBOL, "*", nil},
		{"$!&%AB123", SYMBOL, "$!&%", nil},
	}

	for _, tc := range testCases {
		tokenType, token, err := peekToken(tc.content)
		if tokenType != tc.expectedTokenType || token != tc.expectedToken || err != tc.expectedError {
			t.Errorf("peekToken(%q) = %q, %q, %q; expected %q, %q, %q", tc.content, tokenType, token, err,
				tc.expectedTokenType, tc.expectedToken, tc.expectedError)
		}
	}
}
