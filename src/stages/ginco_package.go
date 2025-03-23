package stages

import "github.com/trudso/ginco/types"

/*
	package roleplaying {
		@changeset
		model A {
		}
		model B {
		}
		model C {
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
