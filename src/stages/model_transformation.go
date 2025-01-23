package stages

import "github.com/trudso/ginco/types"

func TransformModel(model types.MetaModel, transformers []ModelTransformer) ([]types.MetaModel, error) {
	return []types.MetaModel{model}, nil
}
