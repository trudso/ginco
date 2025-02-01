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
)

var ERR_IMPORT_NOT_SUPPORTED = errors.New("import keyword not supported yet")
var ERR_UNEXPECTED_TOKEN = errors.New("Unexpected token")

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
			1= id uuid
			?= name string
			?= age number
			1- type CharacterType
			*= skills Skill
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

	scope := trim(content)
	_, scope, err := popExpectedToken(scope, []string{PACKAGE})
	if err != nil {
		return pkg, scope, err
	}

	scope, restResult, err := popScope(scope)

	// traits and models
	for {
		tokenType, token, err := peekToken(scope)
		if err != nil {
			return pkg, restResult, err
		}

		if tokenType == SCOPE_END {
			break
		}

		switch token {
		case TRAIT_SYMBOL:
			_, scope, err := parseTrait(scope)
			if err != nil {
				return pkg, scope, err
			}
		case MODEL:
			_, scope, err := parseModel(scope)
			if err != nil {
				return pkg, scope, err
			}
		}
	}

	return pkg, restResult, nil
}

func parseTrait(content string) (types.MetaTrait, string, error) {
	trait := types.MetaTrait{}
	return trait, content, fmt.Errorf("Traits not supported yet")
}

func parseModel(content string) (types.MetaModel, string, error) {
	model := types.MetaModel{}
	return model, content, fmt.Errorf("Models not supported yet")
}

// --- token helper functions --- //
