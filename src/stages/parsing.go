package stages

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

const (
	TRAIT       = iota + 1 // @
	SCOPE_BEGIN            // {
	SCOPE_END              // }
	COMMENT                // #
	SYMBOL                 // any non ANSI and non digit element
	IDENTIFIER             // ANSI
	NUMBER                 // all digits
	EOF                    // end of file
)

var ERR_NO_SCOPE_FOUND = errors.New("No scope found")
var ERR_NO_MATCHING_SCOPE_BRACKETS = errors.New("unable to find matching scope brackets")
var ERR_NOT_A_NUMBER = errors.New("Not a number")
var ERR_EXPECTED_ELEMENT_NOT_FOUND = errors.New("None of the expected elements were found")

func popWord(content string) (string, string) {
	return pop(content, " \n")
}

func popLine(content string) (string, string) {
	return pop(content, "\n")
}

func pop(content string, delimiters string) (string, string) {
	trimedContent := strings.Trim(content, delimiters)
	index := strings.IndexAny(trimedContent, delimiters)
	if index == -1 {
		return trimedContent, ""
	}

	return content[:index], content[index+1:]
}

func popScope(content string) (string, string, error) {
	trimmedContent := trim(content)
	startScopeCount := 0
	endScopeCount := 0
	scope := ""

	for index, rune := range trimmedContent {
		scope += string(rune)
		if rune == '{' {
			startScopeCount++
		}

		if rune == '}' {
			endScopeCount++
		}

		if startScopeCount == 0 {
			return "", "", ERR_NO_SCOPE_FOUND
		}

		if startScopeCount == endScopeCount && startScopeCount > 0 {
			return trim(scope[1 : len(scope)-1]), trim(trimmedContent[index+1:]), nil
		}
	}

	if startScopeCount == 0 {
		return "", "", ERR_NO_SCOPE_FOUND
	}

	return "", "", ERR_NO_MATCHING_SCOPE_BRACKETS
}

func popExpectedToken(content string, expectations []string) (string, string, error) {
	trimmedContent := trim(content)
	for _, expectation := range expectations {
		if strings.HasPrefix(trimmedContent, expectation) {
			return expectation, trim(trimmedContent[len(expectation):]), nil
		}
	}

	_, foundToken, _ := peekToken(trimmedContent)
	return "", trimmedContent, errors.Join(ERR_EXPECTED_ELEMENT_NOT_FOUND, fmt.Errorf("Expected: [%q], but found %q", expectations, foundToken))
}

func peekToken(content string) (int, string, error) {
	trimmedContent := trim(content)
	if len(trimmedContent) == 0 {
		return EOF, "", nil
	}

	var symbol = trimmedContent[0]
	switch symbol {
	case '{':
		return SCOPE_BEGIN, string(symbol), nil
	case '}':
		return SCOPE_END, string(symbol), nil
	case '@':
		return TRAIT, string(symbol), nil
	case '#':
		return COMMENT, string(symbol), nil
	}

	if unicode.IsLetter(rune(symbol)) {
		return IDENTIFIER, peekWord(trimmedContent), nil
	}

	if unicode.IsDigit(rune(symbol)) {
		result, err := peekNumber(trimmedContent)
		return NUMBER, strconv.Itoa(result), err
	}

	result := string(symbol)
	nextTokenType, nextToken, err := peekToken(trimmedContent[1:])

	if err != nil {
		return SYMBOL, result, err
	}

	if nextTokenType != SYMBOL {
		return SYMBOL, result, nil
	}

	result += string(nextToken)
	return SYMBOL, result, nil
}

func peekWord(content string) string {
	trimedContent := trim(content)
	index := strings.IndexAny(trimedContent, " \t\n")
	if index == -1 {
		return trimedContent
	}

	return content[:index]
}

func peekNumber(content string) (int, error) {
	trimmedContent := trim(content)
	result := ""
	for _, r := range trimmedContent {
		if !unicode.IsDigit(r) {
			break
		}

		result += string(r)
	}

	num, err := strconv.Atoi(result)
	if err != nil {
		return 0, ERR_NOT_A_NUMBER
	}

	return num, nil
}

func trim(content string) string {
	return strings.Trim(content, " \t\n")
}
