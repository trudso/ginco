package stages

import (
	"errors"
	"fmt"
	"io"

	"github.com/trudso/ginco/types"
)

// tokens
const (
	COMPOSITION = "="
	AGGREGATION = "-"
	COLLECTION  = "*"
	NON_NULL    = "1"
	NULLABLE    = "?"

	IMPORT       = "import"
	PACKAGE      = "package"
	TRAIT_SYMBOL = "@"
	MODEL        = "model"
	MODEL_FIELDS = "fields"
)

var ERR_IMPORT_NOT_SUPPORTED = errors.New("import keyword not supported yet")
var ERR_UNEXPECTED_TOKEN = errors.New("Unexpected token")
var ERR_TRAIT_NOT_DEFINED = errors.New("Trait not defined")
var ERR_MODEL_NAME_NOT_DEFINED = errors.New("Model name not defined")

type GincoMetaFileParser struct{}

func (self GincoMetaFileParser) Parse(reader io.Reader) (types.MetaFile, error) {
	return types.MetaFile{}, nil
}

/*
	package roleplaying {
		@changeset
		model Character {
			fields {
				@noChangeset
				=1 id uuid
				=? name string
				=? age number
				-1 type CharacterType
				=* skills Skill
	    	}
	    }
	}
*/
func parsePackage(content string, idx int) (types.MetaPackage, int, error) {
	pkg := types.MetaPackage {}
	_, nextIdx, err := popExpectedToken(content, idx, TT_IDENTIFIER, PACKAGE) 
	if err != nil {
		return pkg, idx, err
	}

	// name
	nameToken, nextIdx, err := popIdentifier(content, nextIdx)
	if err != nil {
		return pkg, nextIdx, err
	}

	pkg.Name = nameToken.Value

	// content
	scope, nextIdx, err := popScope(content, nextIdx)
	if err != nil {
		return pkg, nextIdx, err
	}

	scopeIdx := 0
	var trait *types.MetaTrait = nil
	for {
		token, _, err := popToken( scope.Value, scopeIdx )
		if err != nil {
			return pkg, scopeIdx, err
		}

		if token.Type == TT_EOF {
			return pkg, nextIdx, nil
		}

		if token.Type == TT_SYMBOL && token.Value == TRAIT_SYMBOL {
			//var newTrait types.MetaTrait
			newTrait, nextIdx, err := parseTrait(scope.Value, scopeIdx)
			trait = &newTrait
			if err != nil {
				return pkg, scopeIdx, err
			}
			scopeIdx = nextIdx
		}

		if token.Type == TT_IDENTIFIER && token.Value == MODEL {
			model, nextIdx, err := parseModel(scope.Value, scopeIdx)
			if err != nil {
				return pkg, scopeIdx, err
			}
			if trait != nil {
				model.Traits = append( model.Traits, *trait)
				trait = nil
			}

			pkg.Models = append( pkg.Models, model )
			scopeIdx = nextIdx
		}
	}
}

func parseModel(content string, idx int) (types.MetaModel, int, error) {
	model := types.MetaModel{}
	token, nextIdx, err := popIdentifier(content, idx)
	if err != nil {
		return model, idx, err
	}

	if token.Value != MODEL {
		return model, idx, formatParsingError(fmt.Sprintf("Expected model keyword, but found %s", token.Value), content, idx)
	}

	token, nextIdx, err = popIdentifier(content, nextIdx)
	model.Name = token.Value

	token, nextIdx, err = popToken(content, nextIdx)
	if token.Type == TT_SCOPE {
		scopeContent := token.Value
		token, _, err := popToken(scopeContent, 0)
		if err != nil {
			return model, nextIdx, err
		}

		if token.Type == TT_IDENTIFIER && token.Value == MODEL_FIELDS {
			fields, _, err := parseModelFields(scopeContent, 0)
			if err != nil {
				return model, nextIdx, err
			}

			model.Fields = fields
		}
	}

	return model, nextIdx, nil
}

func parseModelFields(content string, idx int) ([]types.MetaModelField, int, error) {
	_, nextIdx, err := popExpectedToken(content, idx, TT_IDENTIFIER, MODEL_FIELDS)
	if err != nil {
		return nil, nextIdx, err
	}

	scope, nextIdx, err := popScope(content, nextIdx)
	if err != nil {
		return nil, nextIdx, err
	}

	fields := []types.MetaModelField{}
	scopeIdx := 0
	for !isEOF(scope.Value, scopeIdx) {
		var field types.MetaModelField
		field, scopeIdx, err = parseModelField(scope.Value, scopeIdx)
		if err != nil {
			return fields, scopeIdx, err
		}

		fields = append(fields, field)
	}

	return fields, scopeIdx, nil
}

func parseTrait(content string, idx int) (types.MetaTrait, int, error) {
	trait := types.MetaTrait{}
	symbolToken, nextIdx, err := popSymbol(content, idx)
	if err != nil {
		return trait, nextIdx, err
	}

	if symbolToken.Value != TRAIT_SYMBOL {
		return trait, idx, formatParsingError("No trait found", content, symbolToken.Position)
	}

	identifier, nextIdx, err := popIdentifier(content, nextIdx)
	if err != nil {
		return trait, nextIdx, err
	}

	trait.Name = identifier.Value
	return trait, nextIdx, nil
}

func parseModelField(content string, idx int) (types.MetaModelField, int, error) {
	field := types.MetaModelField{}
	curIdx := idx

	for {
		symbolToken, nextIdx, err := popSymbol(content, curIdx)
		if err != nil {
			return field, curIdx, err
		}

		if symbolToken.Value == TRAIT_SYMBOL {
			trait, nextIdx, err := parseTrait(content, curIdx)
			if err != nil {
				return field, nextIdx, err
			}

			field.Traits = append(field.Traits, trait)
			curIdx = nextIdx
		} else if isOwnershipSymbol(symbolToken.Value) {
			switch symbolToken.Value {
			case COMPOSITION:
				field.Ownership = types.Composition
			case AGGREGATION:
				field.Ownership = types.Aggregation
			}

			multiplicy, nextIdx, err := popSingleRune(content, nextIdx)
			if err != nil {
				return field, nextIdx, err
			}
			switch multiplicy.Value {
			case NULLABLE:
				field.Cardinality = types.ZeroOrOne
			case NON_NULL:
				field.Cardinality = types.One
			case COLLECTION:
				field.Cardinality = types.Collection
			}

			name, nextIdx, err := popIdentifier(content, nextIdx)
			if err != nil {
				return field, nextIdx, err
			}

			field.Name = name.Value

			metaType, nextIdx, err := parseMetaType(content, nextIdx)
			if err != nil {
				return field, nextIdx, err
			}

			field.Type = metaType
			return field, nextIdx, err
		} else {
			return field, curIdx, formatParsingError(fmt.Sprintf("Unexpected symbol %s", symbolToken.Value), content, curIdx)
		}
	}
}

func parseMetaType(content string, idx int) (types.MetaType, int, error) {
	metaType := types.MetaType{}

	token, nextIdx, err := popIdentifier(content, idx)
	if err != nil {
		return metaType, idx, err
	}

	metaType.Name = token.Value
	return metaType, nextIdx, nil
}

func isOwnershipSymbol(symbol string) bool {
	return symbol == COMPOSITION || symbol == AGGREGATION
}
