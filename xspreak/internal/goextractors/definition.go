package goextractors

import (
	"context"
	"go/ast"
	"go/token"
	"time"

	"github.com/vorlif/spreak/xspreak/internal/result"
	"github.com/vorlif/spreak/xspreak/internal/util"

	"github.com/vorlif/spreak/xspreak/internal/config"
	"github.com/vorlif/spreak/xspreak/internal/extract/extractors"
)

type definitionExtractor struct{}

func NewDefinitionExtractor() extractors.Extractor {
	return &definitionExtractor{}
}

func (d *definitionExtractor) Run(ctx context.Context, extractCtx *extractors.Context) ([]result.Issue, error) {
	defer util.TrackTime(time.Now(), "Extract definitions")
	runner := &definitionExtractorRunner{
		ctx:        ctx,
		extractCtx: extractCtx,
	}
	extractCtx.Inspector.Nodes(nil, runner.searchDefinitions)
	return []result.Issue{}, nil
}

func (d definitionExtractor) Name() string {
	return "extract_definitions"
}

type definitionExtractorRunner struct {
	ctx        context.Context
	extractCtx *extractors.Context
}

func (de *definitionExtractorRunner) searchDefinitions(n ast.Node, push bool) bool {
	if !push {
		return true
	}

	switch v := n.(type) {
	case *ast.FuncDecl:
		de.extractFunc(v)
	case *ast.GenDecl:
		switch v.Tok {
		case token.VAR:
			de.extractVar(v)
		case token.TYPE:
			de.extractStruct(v)
		}
	}

	return true
}

// var t localize.Singular.
func (de *definitionExtractorRunner) extractVar(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		selector := searchSelector(valueSpec.Type)
		if selector == nil {
			continue
		}

		tok := de.extractCtx.GetLocalizeTypeToken(selector)
		if tok != extractors.TypeSingular {
			// TODO(fv): log hint
			continue
		}

		for _, name := range valueSpec.Names {
			pkg, obj := de.extractCtx.GetType(name)
			if pkg == nil {
				continue
			}

			def := &extractors.Definition{
				Type:  extractors.VarSingular,
				Token: tok,
				Pck:   pkg,
				Ident: name,
				Path:  obj.Pkg().Path(),
				ID:    objToKey(obj),
				Obj:   obj,
			}

			de.addDefinition(def)
		}

	}

}

/*
type TT struct {
	T localize.Singular
	P localize.Plural
}.
*/
func (de *definitionExtractorRunner) extractStruct(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		pkg, obj := de.extractCtx.GetType(typeSpec.Name)
		if obj == nil {
			continue
		}

		if !config.ShouldScanStruct(pkg.PkgPath) {
			continue
		}

		for i, field := range structType.Fields.List {

			var tok extractors.TypeToken
			switch field.Type.(type) {
			case *ast.Ident:
				if pkg.PkgPath == config.SpreakLocalizePackagePath {
					tok = de.extractCtx.GetLocalizeTypeToken(field.Type)
					break
				}
				if selector := searchSelector(field.Type); selector == nil {
					continue
				} else {
					tok = de.extractCtx.GetLocalizeTypeToken(selector)
				}
			default:
				if selector := searchSelector(field.Type); selector == nil {
					continue
				} else {
					tok = de.extractCtx.GetLocalizeTypeToken(selector)
				}
			}

			if tok == extractors.TypeNone {
				continue
			}

			for ii, fieldName := range field.Names {
				def := &extractors.Definition{
					Type:       extractors.StructField,
					Token:      tok,
					Pck:        pkg,
					Ident:      typeSpec.Name,
					Path:       obj.Pkg().Path(),
					ID:         objToKey(obj),
					Obj:        obj,
					FieldIdent: fieldName,
					FieldName:  fieldName.Name,
					FieldPos:   calculatePosIdx(ii, i),
				}

				de.addDefinition(def)
			}
		}
	}
}

// func translate(msgid localize.Singular, plural localize.Plural)
// func getTranslation() (localize.Singular, localize.Plural).
func (de *definitionExtractorRunner) extractFunc(decl *ast.FuncDecl) {
	pck, obj := de.extractCtx.GetType(decl.Name)
	if pck == nil {
		return
	}

	if decl.Type == nil {
		return
	}

	if decl.Type.Params != nil {
		for i, param := range decl.Type.Params.List {
			if len(param.Names) == 0 {
				continue
			}

			selector := searchSelector(param)
			if selector == nil {
				continue
			}

			tok := de.extractCtx.GetLocalizeTypeToken(selector)
			if tok == extractors.TypeNone || tok == extractors.TypeMessage {
				continue
			}

			for ii, name := range param.Names {
				def := &extractors.Definition{
					Type:       extractors.FunctionParam,
					Token:      tok,
					Pck:        pck,
					Ident:      decl.Name,
					Path:       obj.Pkg().Path(),
					ID:         objToKey(obj),
					Obj:        obj,
					FieldIdent: name,
					FieldName:  name.Name,

					FieldPos: calculatePosIdx(i, ii),
				}
				de.addDefinition(def)
			}
		}
	}
}

func (de *definitionExtractorRunner) addDefinition(d *extractors.Definition) {
	key := d.Key()
	if _, ok := de.extractCtx.Definitions[key]; !ok {
		de.extractCtx.Definitions[key] = make(map[string]*extractors.Definition)
	}

	de.extractCtx.Definitions[key][d.FieldName] = d
}
