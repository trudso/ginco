package types

// Data structure
type MetaFile struct {
	Packages []MetaPackage
}

type MetaPackage struct {
	Models []MetaModel
}

type MetaTrait struct {
	Name string
}

type MetaModel struct {
	Name   string
	Traits []string
	Fields []MetaModelField
}

type MetaModelField struct {
	Name        string
	Type        MetaType
	Kind        string
	Cardinality string
	Nullable    bool
	Traits      []MetaTrait
}

type MetaType struct {
	Package string
	Name    string
}
