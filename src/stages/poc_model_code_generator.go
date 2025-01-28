package stages

import "github.com/trudso/ginco/types"

type PocModelCodeGenerator struct {}

func (self PocModelCodeGenerator) Generate( model types.MetaModel) ([]CodeGeneratorResult, error) {
	return []CodeGeneratorResult{
		 {
			Path: "PocTest.go",
			Content: `package poc

import "fmt"

func PocTestFunc() {
	fmt.Println("POC Test")
}
`,
		},
	}, nil
}
