package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/trudso/ginco/stages"
)

func main() {
	reader := strings.NewReader(`
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
	}`)

	parser := stages.GincoMetaFileParser{}
	file, err := parser.Parse(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unexpected error: %+v\n", err)
	}

	fmt.Printf("Ginco out: %+v\n", file)
}
