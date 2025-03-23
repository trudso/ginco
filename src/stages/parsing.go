package stages

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

type TokenType int

const (
	TT_COMMENT TokenType = iota
	TT_SCOPE
	TT_NUMBER
	TT_IDENTIFIER
	TT_SYMBOL
	TT_RUNE
	TT_INVALID
	TT_EOF
)

var skippables = []string{" ", "\t", "\r", "\n"}

type Token struct {
	Type     TokenType
	Position int
	Value    string
}

func isEOF(content string, startIdx int) bool {
	return firstValidTokenIndex(content, startIdx) == -1
}

func popToken(content string, startIdx int) (Token, int, error) {
	startIdx = firstValidTokenIndex(content, startIdx)
	if startIdx == -1 {
		return Token{
			Type:     TT_EOF,
			Position: startIdx,
			Value:    "",
		}, -1, nil
	}

	r := rune(content[startIdx])
	if unicode.IsDigit(r) {
		return popNumber(content, startIdx)
	}
	if unicode.IsLetter(r) {
		return popIdentifier(content, startIdx)
	}

	switch r {
	case '#':
		return popComment(content, startIdx)
	case '{':
		return popScope(content, startIdx)
	default:
		return Token{
			Type:     TT_SYMBOL,
			Position: startIdx,
			Value:    string(r),
		}, startIdx + 1, nil
	}
}

// popSymbol returns the first valid rune
// that is not a digit/letter/scope
func popSymbol(content string, startIdx int) (Token, int, error) {
	realStartIdx := firstValidTokenIndex(content, startIdx)

	if realStartIdx == -1 {
		return Token{}, startIdx, formatParsingError("Found EOF and not a symbol", content, startIdx)
	}

	r := rune(content[realStartIdx])
	if unicode.IsDigit(r) || unicode.IsLetter(r) || r == '{' || r == '}' {
		return Token{}, startIdx, formatParsingError("No symbol found", content, startIdx)
	}

	return Token{
		Type:     TT_SYMBOL,
		Position: realStartIdx,
		Value:    string(r),
	}, realStartIdx + 1, nil
}

// popSingleRune returns the first valid rune
// that is not a scope
func popSingleRune(content string, startIdx int) (Token, int, error) {
	realStartIdx := firstValidTokenIndex(content, startIdx)

	if realStartIdx == -1 {
		return Token{}, startIdx, formatParsingError("Found EOF and not a symbol", content, startIdx)
	}

	r := rune(content[realStartIdx])
	if r == '{' || r == '}' {
		return Token{}, startIdx, formatParsingError("No rune found", content, startIdx)
	}

	return Token{
		Type:     TT_RUNE,
		Position: realStartIdx,
		Value:    string(r),
	}, realStartIdx + 1, nil
}

func popExpectedToken(content string, startIdx int, tokenType TokenType, value string) (Token, int, error) {
	token, nextIdx, err := popToken(content, startIdx)
	if err != nil {
		return token, nextIdx, err
	}

	if token.Type != tokenType {
		return token, nextIdx, formatParsingError(fmt.Sprintf("Expected token type %v, but found %v", tokenType, token.Type), content, startIdx)
	}

	if token.Value != value {
		return token, nextIdx, formatParsingError(fmt.Sprintf("Expected value %q, but found %q", value, token.Value), content, startIdx)
	}

	return token, nextIdx, err
}

func popIdentifier(content string, startIdx int) (Token, int, error) {
	realStartIdx := firstValidTokenIndex(content, startIdx)
	if realStartIdx == -1 {
		return Token{}, startIdx, formatParsingError("No identifier found", content, startIdx)
	}
	endIdx := realStartIdx
	for i := realStartIdx; i < len(content); i++ {
		r := rune(content[i])
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			break
		}

		endIdx = i + 1
	}

	if realStartIdx == endIdx {
		return Token{}, realStartIdx, formatParsingError("No identifier found", content, realStartIdx)
	}

	return Token{
		Type:     TT_IDENTIFIER,
		Position: realStartIdx,
		Value:    content[realStartIdx:endIdx],
	}, endIdx, nil
}

func popNumber(content string, startIdx int) (Token, int, error) {
	realStartIdx := firstValidTokenIndex(content, startIdx)
	if realStartIdx == -1 {
		return Token{}, startIdx, formatParsingError("No number found", content, startIdx)
	}
	endIdx := realStartIdx
	for i := realStartIdx; i < len(content); i++ {
		r := rune(content[i])
		if !unicode.IsDigit(r) {
			break
		}

		endIdx = i + 1
	}

	if realStartIdx == endIdx {
		return Token{}, startIdx, formatParsingError("No number found", content, startIdx)
	}

	return Token{
		Type:     TT_NUMBER,
		Position: realStartIdx,
		Value:    content[realStartIdx:endIdx],
	}, endIdx, nil
}

func popScope(content string, startIdx int) (Token, int, error) {
	realStartIdx := firstValidTokenIndex(content, startIdx)
	if realStartIdx == -1 {
		return Token{
			TT_SCOPE,
			startIdx,
			"",
		}, startIdx, formatParsingError("No scope found", content, startIdx)
	}

	beginScopeCount, endScopeCount := 0, 0
	for i := realStartIdx; i < len(content); i++ {
		switch rune(content[i]) {
		case '{':
			beginScopeCount++
		case '}':
			endScopeCount++
		}

		if beginScopeCount == 0 {
			return Token{
				TT_SCOPE,
				realStartIdx,
				"",
			}, -1, formatParsingError("No scope found", content, startIdx)
		}

		if beginScopeCount == endScopeCount && endScopeCount > 0 {
			return Token{
				Type:     TT_SCOPE,
				Position: startIdx + 1,
				Value:    getSubString(content, realStartIdx+1, i),
			}, i + 1, nil
		}
	}

	if beginScopeCount == 0 {
		return Token{
			TT_SCOPE,
			realStartIdx,
			"",
		}, -1, formatParsingError("No scope found", content, startIdx)
	}

	return Token{
		TT_SCOPE,
		realStartIdx,
		"",
	}, -1, formatParsingError("Unable to find matching scope identifiers", content, startIdx)
}

func popComment(content string, startIdx int) (Token, int, error) {
	realStartIdx := firstValidTokenIndex(content, startIdx)
	if realStartIdx == -1 {
		return Token{}, startIdx, formatParsingError("No comment found", content, startIdx)
	}

	if content[realStartIdx] != '#' {
		return Token{}, startIdx, formatParsingError("No comment found", content, startIdx)
	}

	endIdx := nextIndexOf(content, "\n", realStartIdx)
	if endIdx == -1 {
		endIdx = len(content)
	}
	return Token{
		Type:     TT_COMMENT,
		Position: realStartIdx + 1,
		Value:    getSubString(content, realStartIdx+1, endIdx),
	}, endIdx + 1, nil
}

func firstValidTokenIndex(content string, startPos int) int {
	for i := startPos; i < len(content); i++ {
		if !slices.Contains(skippables, string(content[i])) {
			return i
		}
	}

	return -1
}

func nextIndexOf(content string, substr string, startIdx int) int {
	rest := content[startIdx:]
	restIdx := strings.Index(rest, substr)
	if restIdx == -1 {
		return restIdx
	}

	return startIdx + restIdx
}

func getSubString(content string, startIdx, endIdx int) string {
	if endIdx == -1 {
		return content[startIdx:]
	}
	return content[startIdx:endIdx]
}

func formatParsingError(errorMessage, content string, atIdx int) error {
	// trim idx
	idx := max(0, min(max(atIdx, 0), len(content)-1))

	contentToIdx := content[:idx]
	line := strings.Count(contentToIdx, "\n") + 1
	lineStartIdx := strings.LastIndex(contentToIdx, "\n")

	endIdx := idx + min(20, len(content)-idx)
	return fmt.Errorf("[%d:%d] ...%s: %s", line, idx-lineStartIdx, content[idx:endIdx], errorMessage)
}
