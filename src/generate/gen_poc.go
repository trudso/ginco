package main

import (
	"fmt"
	"log"
	"os"

	"github.com/trudso/ginco/stages"
	"github.com/trudso/ginco/types"
)

func main() {
	generator := stages.PocModelCodeGenerator{}
	results, _ := generator.Generate(types.MetaModel{})
	for _, result := range results {
		file, err := os.Create(result.Path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		fmt.Fprintf(file, result.Content)

		fmt.Println("Done")
	}
}
