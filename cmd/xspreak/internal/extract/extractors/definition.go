package extractors

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type DefinitionType int

const (
	VarSingular DefinitionType = iota
	Array
	FunctionReturn
	FunctionParam
	StructField
)

type TypeToken string

const (
	TypeNone     TypeToken = ""
	TypeSingular TypeToken = "Singular"
	TypePlural   TypeToken = "Plural"
	TypeDomain   TypeToken = "Domain"
	TypeContext  TypeToken = "Context"
	TypeMessage  TypeToken = "Message"
)

type Definitions map[string]map[string]*Definition // path.name -> field || "" -> Definition

func (defs Definitions) Get(key, fieldName string) *Definition {
	if _, ok := defs[key]; !ok {
		return nil
	}

	if _, ok := defs[key][fieldName]; !ok {
		return nil
	}

	return defs[key][fieldName]
}

func (defs Definitions) GetFields(key string) map[string]*Definition {
	if _, ok := defs[key]; !ok {
		return nil
	}

	if len(defs[key]) == 0 {
		return nil
	}

	return defs[key]
}

type Definition struct {
	Type  DefinitionType
	Token TypeToken
	Pck   *packages.Package
	Ident *ast.Ident
	Path  string // github.com/name/repo/package/pack
	ID    string // github.com/name/repo/package/pack.StructName
	Obj   types.Object

	// only for functions and structs
	FieldIdent *ast.Ident
	FieldName  string

	// functions only
	FieldPos int
}

func (d *Definition) Key() string {
	return d.ID
}

var stringExtractNames = map[string]TypeToken{
	"MsgID":    TypeSingular,
	"Singular": TypeSingular,
	"Plural":   TypePlural,
	"Domain":   TypeDomain,
	"Context":  TypeContext,
}
