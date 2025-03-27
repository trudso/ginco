package stages

import (
	"slices"

	"github.com/trudso/ginco/types"
)

const (
	ENUM = "enum"
	LITERALS = "literals"
)

/*
	enum CharacterType {
		literals {
			warrior		
			rogue
			wizard
		}
	}
*/
func parseEnum(content string, idx int) (types.MetaEnum, int, error) {
	enum := types.MetaEnum{}
	token, nextIdx, err := popExpectedToken( content, idx, TT_IDENTIFIER, ENUM )
	if err != nil {
		return enum, idx, err
	}

	token, nextIdx, err = popIdentifier(content, nextIdx)
	if err != nil {
		return enum, nextIdx, err
	}

	enum.Name = token.Value
	token, nextIdx, err = popToken(content, nextIdx)
	if token.Type == TT_SCOPE {
		scopeContent := token.Value
		token, _, err := popToken(scopeContent, 0)
		if err != nil {
			return enum, nextIdx, err
		}

		if token.Type == TT_IDENTIFIER && token.Value == LITERALS {
			literals, _, err := parseEnumLiterals(scopeContent, 0)
			if err != nil {
				return enum, nextIdx, err
			}

			enum.Literals = literals
		}
	}

	return enum, nextIdx, nil
}

func parseEnumLiterals(content string, idx int) ([]string, int, error) {
	literals := []string{}

	_, nextIdx, err := popExpectedToken(content, idx, TT_IDENTIFIER, LITERALS)
	if err != nil {
		return nil, nextIdx, err
	}

	scope, nextIdx, err := popScope(content, nextIdx)
	if err != nil {
		return nil, nextIdx, err
	}

	scopeIdx := 0
	for !isEOF(scope.Value, scopeIdx) {
		var literalToken Token
		literalToken, scopeIdx, err = popIdentifier(scope.Value, scopeIdx)
		if err != nil {
			return literals, scopeIdx, err
		}

		if slices.Contains( literals, literalToken.Value ) {
			return literals, nextIdx, formatParsingError("duplicate literal found", content, literalToken.Position)
		}

		literals = append(literals, literalToken.Value)
	}

	return literals, nextIdx, nil
}
