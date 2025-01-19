package main

import (
	"fmt"
	"github.com/trudso/ginco/pkg/api"
	"github.com/trudso/ginco/internal/serialization"
)

func main() {
	file := api.MetaFile {
		Models: []api.MetaModel{},
	}

	jsonSerializer := serialization.MetaFileJsonSerializer {}

	data, err := jsonSerializer.Serialize( file )
	if err != nil {
		fmt.Println("error")
	}

	fmt.Printf("done: %+v\n", data)
}
