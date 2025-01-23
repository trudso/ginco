package stages

import (
	"io"

	"github.com/trudso/ginco/types"
)

// data structures
type CodeGeneratorResult struct {
	Path    string
	Content string
}

// stages
type MetaFileParser interface {
	Parse(reader io.Reader) (types.MetaFile, error)
}

type ModelTransformer interface {
	Transform(model types.MetaModel) ([]types.MetaModel, error)
}

type ModelCodeGenerator interface {
	Generate(model types.MetaModel) ([]CodeGeneratorResult, error)
}
