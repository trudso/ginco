package main

import (
	"fmt"
	"strings"

	"github.com/trudso/ginco/stages"
)

func main() {
	reader := strings.NewReader(`Models:
	Test:
		Fields:
			IntField:
				Type:
					Name: int
	`)

	yamlMetaFile := stages.YamlMetaFile{}
	metafile, err := yamlMetaFile.Parse(reader)
	if err != nil {
		fmt.Printf("Error while parsing yaml file: %+v\n", err)
	}

	fmt.Printf("Ginco out: %+v", metafile)
}
