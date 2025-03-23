package types

type Cardinality int

const (
	ZeroOrOne Cardinality = iota
	One
	Collection
)

type Ownership int

const (
	Composition Ownership = iota
	Aggregation
)

// Data structure
type MetaFile struct {
	Packages []MetaPackage
}

type MetaPackage struct {
	Name   string
	Models []MetaModel
}

type MetaTrait struct {
	Name string
}

type MetaModel struct {
	Name   string
	Traits []MetaTrait
	Fields []MetaModelField
}

type MetaModelField struct {
	Name string
	Type MetaType
	// Kind        string ?
	Cardinality Cardinality
	Ownership   Ownership
	Nullable    bool
	Traits      []MetaTrait
}

type MetaType struct {
	Package string
	Name    string
}

type MetaEnum struct {
	Name     string
	Traits   []MetaTrait
	Literals []string
}
