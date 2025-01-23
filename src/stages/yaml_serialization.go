package stages

import (
	"github.com/trudso/ginco/types"
	"io"
	"gopkg.in/yaml.v3"
)

type YamlMetaFile struct {}

func (self YamlMetaFile) Parse(reader io.Reader) (types.MetaFile, error) {
	result := types.MetaFile {}
	decoder := yaml.NewDecoder( reader )
	err := decoder.Decode( &result );
	return result, err
}
