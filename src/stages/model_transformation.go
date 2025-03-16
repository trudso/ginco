package stages

import "github.com/trudso/ginco/types"

type ModelTransformer interface {
	Transform(model types.MetaModel) ([]types.MetaModel, error)
}

func TransformModel(model types.MetaModel, transformers []ModelTransformer) ([]types.MetaModel, error) {
	return []types.MetaModel{model}, nil
}
