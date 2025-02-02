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

func parseFile(content string) (types.MetaFile, error) {
	file := types.MetaFile{
		Packages: []types.MetaPackage{},
	}
	rest := content
	for {
		tokenType, token, err := peekToken(rest)
		if err != nil {
			return file, err
		}

		if tokenType == EOF {
			return file, nil
		}

		switch token {
		case IMPORT:
			return file, ERR_IMPORT_NOT_SUPPORTED
		case PACKAGE:
			pkg, newRest, err := parsePackage(rest)
			if err != nil {
				return file, err
			}
			file.Packages = append(file.Packages, pkg)
			rest = newRest
		default:
			return file, fmt.Errorf("Unexpected token %q", token)
		}
	}
}

func parsePackage(content string) (types.MetaPackage, string, error) {
	pkg := types.MetaPackage{
		Models: []types.MetaModel{},
	}

	trimmedContent := trim(content)
	_, rest, err := popExpectedToken(trimmedContent, []string{PACKAGE})
	if err != nil {
		return pkg, rest, err
	}

	pkgScope, rest, err := popScope(rest)

	// traits and models
	for {
		tokenType, token, err := peekToken(pkgScope)
		if err != nil {
			return pkg, rest, err
		}

		if tokenType == SCOPE_END || tokenType == EOF {
			break
		}

		switch token {
		case TRAIT_SYMBOL:
			_, pkgScope, err = parseTrait(pkgScope)
			if err != nil {
				return pkg, pkgScope, err
			}
		case MODEL:
			_, pkgScope, err = parseModel(pkgScope)
			if err != nil {
				return pkg, pkgScope, err
			}
		}
	}

	return pkg, rest, nil
}

func parseTrait(content string) (types.MetaTrait, string, error) {
	trait := types.MetaTrait{}
	_, rest, err := popExpectedToken(content, []string{TRAIT_SYMBOL})
	if err != nil {
		return trait, rest, err
	}

	traitName, rest := popWord(rest)
	if len(traitName) == 0 {
		return trait, rest, ERR_TRAIT_NOT_DEFINED
	}

	trait.Name = traitName
	return trait, rest, nil
}

func parseModel(content string) (types.MetaModel, string, error) {
	model := types.MetaModel{}
	_, rest, err := popExpectedToken(content, []string{"model"})
	if err != nil {
		return model, rest, err
	}

	// get name
	name, rest := popWord(rest)
	if len(name) == 0 {
		return model, rest, ERR_MODEL_NAME_NOT_DEFINED
	}
	model.Name = name

	// get scope
	scope, resultRest, err := popScope(rest)
	if err != nil {
		return model, rest, err
	}

	scopeRest := scope

	for {
		tokenType, token, err := peekToken(scopeRest)
		if err != nil {
			return model, resultRest, err
		}

		if tokenType == EOF || tokenType == SCOPE_END {
			break
		}

		switch token {
		case MODEL_FIELDS:
			fields, modelRest, err := parseModelFields(scopeRest)
			if err != nil {
				return model, resultRest, err
			}
			model.Fields = append(model.Fields, fields...)
			scopeRest = modelRest
		default:
			return model, resultRest, ERR_UNEXPECTED_TOKEN
		}
	}

	return model, resultRest, nil
}

/*
	fields {
		@noChangeset
		=1 id uuid
		=? name string
		=? age number
		-1 type CharacterType
		=* skills Skill
	}
*/
func parseModelFields(content string) ([]types.MetaModelField, string, error) {
	fields := []types.MetaModelField{}
	_, rest, err := popExpectedToken(content, []string{MODEL_FIELDS})
	if err != nil {
		return fields, rest, err
	}
	scope, resultRest, err := popScope(rest)
	if err != nil {
		return fields, resultRest, err
	}

	//
	scopeRest := scope
	fieldTraits := []types.MetaTrait{}
	for {
		tokenType, _, err := peekToken(scopeRest)
		if err != nil {
			return fields, rest, err
		}

		if tokenType == EOF || tokenType == SCOPE_END {
			break
		}

		switch tokenType {
		case TRAIT:
			trait, newRest, err := parseTrait(scopeRest)
			if err != nil {
				return fields, resultRest, err
			}
			scopeRest = newRest
			fieldTraits = append(fieldTraits, trait)
		default:
			field, newRest, err := parseModelField(scopeRest)
			if err != nil {
				return fields, resultRest, err
			}
			field.Traits = append(field.Traits, fieldTraits...)
			fields = append(fields, field)
			fieldTraits = []types.MetaTrait{} // reset traits
			scopeRest = newRest
		}
	}

	return fields, resultRest, nil
}

func parseModelField(content string) (types.MetaModelField, string, error) {
	field := types.MetaModelField{}

	// kind
	kind, rest, err := popExpectedToken(content, []string{COMPOSITION, AGGREGATION})
	if err != nil {
		return field, rest, err
	}
	field.Kind = kind

	// cardinality
	cardinality, rest, err := popExpectedToken(rest, []string{COLLECTION, NULLABLE, NON_NULL})
	if err != nil {
		return field, rest, err
	}
	field.Cardinality = cardinality

	// field name
	name, rest := popWord(rest)
	field.Name = name

	// type
	typeName, rest := popLine(rest)
	field.Type.Name = typeName

	return field, rest, nil
}
