package extractors

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"

	"github.com/vorlif/spreak/xspreak/internal/config"
	"github.com/vorlif/spreak/xspreak/internal/util"
)

type Context struct {
	Config *config.Config
	Log    *log.Entry

	OriginalPackages []*packages.Package

	Packages    map[string]*packages.Package
	Inspector   *inspector.Inspector
	CommentMaps map[string]map[string]ast.CommentMap // pkg -> file -> node -> comments

	Definitions Definitions
}

func (c *Context) GetPosition(pos token.Pos) token.Position {
	for _, pkg := range c.Packages {
		if position := pkg.Fset.Position(pos); position.IsValid() {
			return position
		}
	}

	return token.Position{}
}

func (c *Context) GetType(ident *ast.Ident) (*packages.Package, types.Object) {
	for _, pkg := range c.Packages {
		if pkg.Types == nil {
			continue
		}
		if obj, ok := pkg.TypesInfo.Defs[ident]; ok {
			if obj == nil || obj.Type() == nil || obj.Pkg() == nil {
				return nil, nil
			}
			return pkg, obj
		}
		if obj, ok := pkg.TypesInfo.Uses[ident]; ok {
			if obj == nil || obj.Type() == nil || obj.Pkg() == nil {
				return nil, nil
			}
			return pkg, obj
		}
		if obj, ok := pkg.TypesInfo.Implicits[ident]; ok {
			if obj == nil || obj.Type() == nil || obj.Pkg() == nil {
				return nil, nil
			}
			return pkg, obj
		}
	}
	return nil, nil
}

func (c *Context) GetLocalizeTypeToken(expr ast.Expr) TypeToken {
	if expr == nil {
		return TypeNone
	}

	switch v := expr.(type) {
	case *ast.SelectorExpr:
		return c.GetLocalizeTypeToken(v.Sel)
	case *ast.Ident:
		_, vType := c.GetType(v)
		if vType == nil {
			return TypeNone
		}

		if vType.Pkg() == nil || vType.Pkg().Path() != config.SpreakLocalizePackagePath {
			return TypeNone
		}

		tok, ok := stringExtractNames[vType.Name()]
		if !ok {
			return TypeNone
		}

		return tok
	default:
		return TypeNone
	}
}

func (c *Context) SearchIdentAndToken(start ast.Node) (TypeToken, *ast.Ident) {
	switch ident := start.(type) {
	case *ast.Ident:
		if tok := c.GetLocalizeTypeToken(ident); tok != TypeNone {
			return tok, ident
		}

		pkg, obj := c.GetType(ident)
		if pkg == nil {
			break
		}

		if def := c.Definitions.Get(objToKey(obj), ""); def != nil {
			return def.Token, ident
		}
	}

	selector := searchSelector(start)
	if selector == nil {
		return TypeNone, nil
	}

	switch ident := selector.X.(type) {
	case *ast.Ident:
		if tok := c.GetLocalizeTypeToken(ident); tok != TypeNone {
			return tok, ident
		}

		pkg, obj := c.GetType(ident)
		if pkg == nil {
			break
		}

		if def := c.Definitions.Get(objToKey(obj), ""); def != nil {
			return def.Token, ident
		}
		if def := c.Definitions.Get(objToKey(obj), selector.Sel.Name); def != nil {
			return def.Token, ident
		}

		if obj.Type() == nil {
			break
		}

		if pointer, ok := obj.Type().(*types.Pointer); ok {
			if def := c.Definitions.Get(pointer.Elem().String(), selector.Sel.Name); def != nil {
				return def.Token, ident
			}
		}
	}

	if tok := c.GetLocalizeTypeToken(selector.Sel); tok != TypeNone {
		return tok, selector.Sel
	}

	pkg, obj := c.GetType(selector.Sel)
	if pkg == nil {
		return TypeNone, nil
	}

	if def := c.Definitions.Get(objToKey(obj), ""); def != nil {
		return def.Token, selector.Sel
	}
	if def := c.Definitions.Get(objToKey(obj), selector.Sel.Name); def != nil {
		return def.Token, selector.Sel
	}

	if obj.Type() != nil {
		if pointer, ok := obj.Type().(*types.Pointer); ok {
			if def := c.Definitions.Get(pointer.Elem().String(), selector.Sel.Name); def != nil {
				return def.Token, selector.Sel
			}
		}
	}

	return TypeNone, nil
}

func (c *Context) GetComments(pkg *packages.Package, node ast.Node, stack []ast.Node) []*ast.CommentGroup {
	if _, hasPkg := c.CommentMaps[pkg.PkgPath]; !hasPkg {
		return nil
	}

	pos := c.GetPosition(node.Pos())
	if _, hasFile := c.CommentMaps[pkg.PkgPath][pos.Filename]; !hasFile {
		return nil
	}

	var topNode = node
	for i := len(stack) - 1; i >= 0; i-- {
		entry := stack[i]
		entryPos := c.GetPosition(entry.Pos())
		if !entryPos.IsValid() || entryPos.Line < pos.Line {
			break
		}

		topNode = entry
	}

	var comments []*ast.CommentGroup
	ast.Inspect(topNode, func(node ast.Node) bool {
		nodeComments := c.CommentMaps[pkg.PkgPath][pos.Filename][node]
		if len(nodeComments) > 0 {
			comments = append(comments, nodeComments...)
		}
		return true
	})

	return comments
}

func (c *Context) BuildPackages() {
	c.Packages = make(map[string]*packages.Package, len(c.OriginalPackages))
	defer util.TrackTime(time.Now(), "Collect packages")
	for _, originalPackage := range c.OriginalPackages {
		c.collectPackages(originalPackage, 0)
	}

	files := make([]*ast.File, 0, 200)
	for _, pkg := range c.Packages {
		files = append(files, pkg.Syntax...)
	}

	c.Inspector = inspector.New(files)
}

func (c *Context) collectPackages(startPck *packages.Package, depth int) {
	if c.Config.MaxDepth >= 0 && depth >= c.Config.MaxDepth {
		return
	}

	c.Packages[startPck.ID] = startPck
	for _, importedPackage := range startPck.Imports {
		if !c.isPartOfDirectory(importedPackage) {
			continue
		}

		if _, ok := c.Packages[importedPackage.ID]; ok {
			continue
		}
		c.collectPackages(importedPackage, depth+1)
	}
}

func (c *Context) isPartOfDirectory(pkg *packages.Package) bool {
	if config.IsValidSpreakPackage(pkg.PkgPath) {
		return true
	}

	for _, src := range c.OriginalPackages {
		if strings.HasPrefix(pkg.PkgPath, src.PkgPath) {
			return true
		}
	}

	return false
}

func objToKey(obj types.Object) string {
	return fmt.Sprintf("%s.%s", obj.Pkg().Path(), obj.Name())
}

func searchSelector(expr interface{}) *ast.SelectorExpr {
	switch v := expr.(type) {
	case *ast.SelectorExpr:
		return v
	case *ast.Ident:
		if v.Obj == nil {
			break
		}
		return searchSelector(v.Obj.Decl)
	case *ast.ValueSpec:
		return searchSelector(v.Type)
	case *ast.Field:
		return searchSelector(v.Type)
	}
	return nil
}
