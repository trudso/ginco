package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopIdentifier(t *testing.T) {
	testCases := []struct {
		input              string
		expectedIdentifier string
		expectedErrors     []string
	}{
		{"width=800", "width", nil },
		{"", "", []string {"No identifier found at 0:"}},
		{"!800", "", []string {"No identifier found at 0:"}},
		{"a=900", "a", nil },
	}

	for _, tc := range testCases {
		token, _, err := popIdentifier(tc.input, 0)
		assertErrorContains(t, err, tc.expectedErrors)
		assert.Equal(t, tc.expectedIdentifier, token.Value)
	}
}

func TestPopNumber(t *testing.T) {
	testCases := []struct {
		input              string
		expectedNumber string
		expectedErrors     []string
	}{
		{"", "", []string {"No number found at 0:"}},
		{"!800", "", []string {"No number found at 0:"}},
		{"900", "900", nil },
	}

	for _, tc := range testCases {
		token, _, err := popNumber(tc.input, 0)
		assertErrorContains(t, err, tc.expectedErrors)
		assert.Equal(t, tc.expectedNumber, token.Value)
	}
}

func assertErrorContains(t *testing.T, err error, expectedElements []string) {
	assert.False(t, err != nil && len(expectedElements) == 0, "expected no error but got: %q", err)
	assert.False(t, err == nil && len(expectedElements) > 0, "expected error to contain %q but got no error", expectedElements)

	for _, e := range expectedElements {
		assert.ErrorContains(t, err, e )
	}
}

func TestPopScope(t *testing.T) {
	testCases := []struct {
		content             string
		expectedValue       string
		expectedPosition    int
		expectedNextIdx     int
		expectedErrorValues []string
	}{
		{"{{{}}}", "{{}}", 1, 6, nil},
		{"{}", "", 1, 2, nil},
		{"{{{}", "", 0, -1, []string{"Unable to find matching scope"}},
		{"{ this is the scope } and this is the rest", " this is the scope ", 1, 21, nil},
		{"", "", 0, 0, []string{"No scope found"}},
		{" { inner thing { with some stuff }} and some other stuff", " inner thing { with some stuff }", 1, 35, nil},
		{" { inner thing { with some { stuff }} and some other stuff", "", 1, -1, []string{"Unable to find matching scope"}},
		{"something first { some inner scope stuff }", "", 0, -1, []string{"No scope found"}},
	}

	for _, tc := range testCases {
		token, nextIdx, err := popScope(tc.content, 0)
		assertErrorContains(t, err, tc.expectedErrorValues)

		assert.Equal( t, TT_SCOPE, token.Type )
		assert.Equal( t, tc.expectedValue, token.Value )
		assert.Equal( t, tc.expectedPosition, token.Position ) 
		assert.Equal( t, tc.expectedNextIdx, nextIdx )
	}
}

func TestPopComment(t *testing.T) {
	testCases := []struct {
		content             string
		expectedValue       string
		expectedPosition    int
		expectedNextIdx     int
		expectedErrorValues []string
	}{
		{"", "", 0, 0, []string{"No comment found"}},
		{"test", "", 0, 0, []string{"No comment found"}},
		{"123", "", 0, 0, []string{"No comment found"}},
		{"#some comment", "some comment", 1, 14, nil},
		{"# some comment", " some comment", 1, 15, nil},
		{" # some comment\n123 Test", " some comment", 2, 16, nil},
	}

	for _, tc := range testCases {
		token, nextIdx, err := popComment(tc.content, 0)
		assertErrorContains(t, err, tc.expectedErrorValues)

		assert.Equal( t, TT_COMMENT, token.Type )
		assert.Equal( t, tc.expectedValue, token.Value )
		assert.Equal( t, tc.expectedPosition, token.Position ) 
		assert.Equal( t, tc.expectedNextIdx, nextIdx )
	}
}
