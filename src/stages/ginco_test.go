package stages

import (
	"testing"

	"github.com/trudso/ginco/types"
)

func TestParseTrait(t *testing.T) {
	testCases := []struct {
		content               string
		expectedTraitName     string
		expectedRest          string
		expectedErrorContains []string
	}{
		{"", "", "", []string{"None of the expected elements were found"}},
		{"@", "", "", []string{"Trait not defined"}},
		{"@Trait1", "Trait1", "", nil},
		{"@Trait\nmodel Something {", "Trait", "model Something {", nil},
		{"Not a trait", "", "Not a trait", []string{"None of the expected elements were found"}},
	}

	for _, tc := range testCases {
		trait, rest, err := parseTrait(tc.content)
		assertErrorContains(t, err, tc.expectedErrorContains)
		if trait.Name != tc.expectedTraitName || rest != tc.expectedRest {
			t.Errorf("parseTrait(%q) = %q, %q, %q; expected %q, %q", tc.content,
				trait, rest, err,
				tc.expectedTraitName, tc.expectedRest)
		}
	}
}

func TestParseModel(t *testing.T) {
	testCases := []struct {
		content       string
		expectedModel types.MetaModel
		expectedRest  string
		expectedError error
	}{
		{` model Character {
				fields {
					=1 id uuid
	    		}
			}`, types.MetaModel{
			Name:   "Character",
			Fields: []types.MetaModelField{},
		},
			"", nil,
		},
	}

	for _, tc := range testCases {
		model, rest, err := parseModel(tc.content)
		if model.Name != tc.expectedModel.Name || rest != tc.expectedRest || err != tc.expectedError {
			t.Errorf("parseModel(%q) = %+v, %q, %q; expected %+v, %q, %q", tc.content,
				model, rest, err,
				tc.expectedModel, tc.expectedRest, tc.expectedError)
		}
	}
}

func TestParseModelField(t *testing.T) {
	testCases := []struct {
		content       string
		expectedField types.MetaModelField
		expectedRest  string
		expectedError error
	}{
		{"=? F1 string", types.MetaModelField{
			Name:        "F1",
			Cardinality: "?",
			Type:        types.MetaType{Package: "", Name: "string"},
		}, "", nil},
	}

	for _, tc := range testCases {
		field, rest, err := parseModelField(tc.content)
		if err != tc.expectedError || rest != tc.expectedRest ||
			field.Name != tc.expectedField.Name ||
			field.Cardinality != tc.expectedField.Cardinality ||
			field.Type.Package != tc.expectedField.Type.Package ||
			field.Type.Name != tc.expectedField.Type.Name {
			t.Errorf("parseModel(%q) = %+v, %q, %q; expected %+v, %q, %q", tc.content,
				field, rest, err,
				tc.expectedField, tc.expectedRest, tc.expectedError)
		}
	}
}
