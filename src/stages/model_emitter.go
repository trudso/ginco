package stages

import "github.com/trudso/ginco/types"

type ModelEmitterResult struct {
	Path    string
	Content string
}

type ModelEmitter interface {
	Generate(model types.MetaModel) ([]ModelEmitterResult, error)
}

func EmitModel(model types.MetaModel, emitters []ModelEmitter) ([]ModelEmitterResult, error) {
	return []ModelEmitterResult{}, nil
}
