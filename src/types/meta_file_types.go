package types

// Data structure
type MetaFile struct {
	Models []MetaModel
}

type MetaModel struct {
	Name   string
	Fields []MetaModelField
}

type MetaModelField struct {
	Name string
	Type MetaType
}

type MetaType struct {
	Package string
	Name    string
}

