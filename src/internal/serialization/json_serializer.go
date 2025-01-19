package serialization

import "github.com/trudso/ginco/pkg/api"

type MetaFileJsonSerializer struct {}

func (self MetaFileJsonSerializer) Serialize(file api.MetaFile) ([]byte, error)  {
	return nil, nil
}

