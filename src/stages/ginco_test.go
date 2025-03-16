package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trudso/ginco/types"
)

func TestParseTrait(t *testing.T) {
	testCases := []struct {
		content             string
		expectedTraitName   string
		expectedNextIdx     int
		expectedErrorValues []string
	}{
		{"", "", 0, []string{"Found EOF"}},
		{"@", "", 1, []string{"No identifier found"}},
		{"@Trait1", "Trait1", 7, nil},
		{" @Trait1", "Trait1", 8, nil},
		{"@Trait\nmodel Something {", "Trait", 6, nil},
		{"Not a trait", "", 0, []string{"No symbol found"}},
	}

	for _, tc := range testCases {
		trait, nextIdx, err := parseTrait(tc.content, 0)
		assertErrorContains(t, err, tc.expectedErrorValues)
		assert.Equal(t, tc.expectedTraitName, trait.Name)
		assert.Equal(t, tc.expectedNextIdx, nextIdx)
	}
}

func TestParseModelField(t *testing.T) {
	testCases := []struct {
		content             string
		expectedFieldName   string
		expectedTypeName		string
		expectedOwnership   types.Ownership
		expectedCardinality types.Cardinality
		expectedNextIdx     int
		expectedTraitNames []string
		expectedErrorValues []string
	}{
		{"", "", "", types.Composition, types.ZeroOrOne, 0, nil, []string{"Found EOF"}},
		{"=1 name string", "name", "string", types.Composition, types.One, 14, nil, nil},
		{"  =1    name     string", "name", "string", types.Composition, types.One, 23, nil, nil},
		{"=1 id uuid", "id", "uuid", types.Composition, types.One, 10, nil, nil},
		{"-1 type CharacterType", "type", "CharacterType", types.Aggregation, types.One, 21, nil, nil},
		{"=* skills Skill", "skills", "Skill", types.Composition, types.Collection, 15, nil, nil},
		{"@FieldTrait1\n=* skills Skill", "skills", "Skill", types.Composition, types.Collection, 28, []string {"FieldTrait1"}, nil},
		{"@FieldTrait1\n@FieldTrait2\n@FieldTrait3\n=* skills Skill", "skills", "Skill", types.Composition, types.Collection, 54, []string {"FieldTrait1", "FieldTrait2", "FieldTrait3"}, nil}	,
	}

	for _, tc := range testCases {
		field, nextIdx, err := parseModelField(tc.content, 0)
		assertErrorContains(t, err, tc.expectedErrorValues)
		assert.Equal(t, tc.expectedFieldName, field.Name)
		assert.Equal(t, tc.expectedOwnership, field.Ownership)
		assert.Equal(t, tc.expectedCardinality, field.Cardinality)
		assert.Equal(t, tc.expectedNextIdx, nextIdx)
		assert.Equal(t, len(field.Traits), len(tc.expectedTraitNames))

		for _, trait := range field.Traits {
			assert.Contains(t, tc.expectedTraitNames, trait.Name)
		}
	}
}

func TestParseSimpleModel(t *testing.T) {
	inputTest := `model Character {
		fields {
				@noChangeset
				=1 id uuid
				=? name string
				=? age number
				-1 type CharacterType
				=* skills Skill
	   	}
	  }`

	model, nextIdx, err := parseModel( inputTest, 0)
	assert.NoError(t, err)
	assert.NotNil(t, model )
	assert.Equal(t, len(inputTest), nextIdx)
	assert.Equal(t, 5, len(model.Fields))
}

func TestParseSimplePackage(t *testing.T) {
	inputTest := `package roleplaying {
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
	}`

	pkg, nextIdx, err := parsePackage( inputTest, 0)
	assert.NoError(t, err)
	assert.NotNil(t, pkg )
	assert.Equal(t, 1, len(pkg.Models))
	assert.Equal(t, len(inputTest), nextIdx) 
}
