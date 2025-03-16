package stages

import "github.com/trudso/ginco/types"

type PocModelCodeGenerator struct{}

func (self PocModelCodeGenerator) Generate(model types.MetaModel) ([]ModelEmitterResult, error) {
	return []ModelEmitterResult{
		{
			Path: "PocTest.go",
			Content: `package main

import "fmt"

func PocTestFunc() {
	fmt.Println("POC Test")
}
`,
		},
	}, nil
}
