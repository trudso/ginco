package stages

import (
	"io"

	"github.com/trudso/ginco/types"
)

type GincoMetaFileParser struct{}

func (self GincoMetaFileParser) Parse(reader io.Reader) (types.MetaFile, error) {
	return types.MetaFile{}, nil
}
