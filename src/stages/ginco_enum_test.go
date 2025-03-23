package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEnum(t *testing.T) {
	testCases := []struct {
		content             string
		expectedEnumName    string
		expectedLiterals    []string
		expectedNextIds     int
		expectedErrorValues []string
	}{
		{`enum CharacterType {
				literals {
					warrior
					rogue
					wizard
				}
			}`, "CharacterType", []string{"warrior", "rogue", "wizard"}, 82, nil},
	}

	for _, tc := range testCases {
		enum, nextIdx, err := parseEnum(tc.content, 0)

		assertErrorContains(t, err, tc.expectedErrorValues)

		assert.Equal(t, tc.expectedEnumName, enum.Name)
		assert.Equal(t, len(tc.expectedLiterals), len(enum.Literals))
		assert.Equal(t, tc.expectedNextIds, nextIdx)

		for _, literal := range enum.Literals {
			assert.Contains(t, tc.expectedLiterals, literal)
		}
	}

}
