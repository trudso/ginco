package api

// data api
type MetaFile struct {
	Models []MetaModel
}

type MetaModel struct {
	Name string;
}

// stages
type ModelTransformer interface {
	Transform( model MetaModel ) ([]MetaModel, error)
}

// serialization
type MetaFileDeserializer interface {
	Deserialize( data []byte ) (MetaFile, error)
}

type MetaFileSerializer interface {
	Serialize(file MetaFile) ([]byte, error)
}
