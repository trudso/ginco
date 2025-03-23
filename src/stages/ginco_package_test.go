package stages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


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

func TestParseMultiModelPackage(t *testing.T) {
	inputTest := `package people {
		@changeset
		model Person {
			fields {
				=1 firstName string
				=1 lastName string
				=1 birthDate date
				=? deathDate date
				=1 residence Address
				=? billingAddress Address
				-? partner Person
			}
		}

		@changeset
		model Address {
			fields {
				=1 streetName string
				=1 streetNumber string
				=1 city string
				=1 zipCode string
			}
		}
	}`

	pkg, nextIdx, err := parsePackage( inputTest, 0)
	assert.NoError(t, err)
	assert.NotNil(t, pkg )
	assert.Equal(t, 2, len(pkg.Models))
	assert.Equal(t, len(inputTest), nextIdx) 
}
