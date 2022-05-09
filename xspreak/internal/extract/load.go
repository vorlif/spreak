package extract

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"

	"github.com/vorlif/spreak/xspreak/internal/config"
	"github.com/vorlif/spreak/xspreak/internal/extract/extractors"
	"github.com/vorlif/spreak/xspreak/internal/util"
)

var loadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedSyntax |
	packages.NeedTypes |
	packages.NeedTypesInfo |
	packages.NeedImports |
	packages.NeedDeps

type ContextLoader struct {
	config *config.Config
	log    *logrus.Entry
}

func NewContextLoader(cfg *config.Config) *ContextLoader {
	return &ContextLoader{
		config: cfg,
		log:    logrus.WithField("service", "ContextLoader"),
	}
}

func (cl *ContextLoader) buildArgs() []string {
	args := cl.config.Args
	if len(args) == 0 {
		return []string{"./..."}
	}

	var retArgs []string
	for _, arg := range args {
		if strings.HasPrefix(arg, ".") || filepath.IsAbs(arg) {
			retArgs = append(retArgs, arg)
		} else if strings.ContainsRune(arg, filepath.Separator) {
			retArgs = append(retArgs, fmt.Sprintf(".%c%s", filepath.Separator, arg), arg)
		} else {
			retArgs = append(retArgs, arg)
		}
	}

	return retArgs
}

func (cl *ContextLoader) Load(ctx context.Context) (*extractors.Context, error) {
	pkgConf := &packages.Config{
		Context: ctx,
		Mode:    loadMode,
		Dir:     cl.config.SourceDir,
		Logf:    logrus.WithField("service", "package-loader").Debugf,
		Tests:   false,
	}

	originalPkgs, err := cl.loadPackages(ctx, pkgConf)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}

	if len(originalPkgs) == 0 {
		return nil, errors.New("no go files to analyze")
	}

	pkgs := make(map[string]*packages.Package, len(originalPkgs))

	ret := &extractors.Context{
		OriginalPackages: originalPkgs,
		Packages:         pkgs,
		Config:           cl.config,
		Log:              cl.log,
		Definitions:      make(extractors.Definitions, 200),
		CommentMaps:      make(map[string]map[string]ast.CommentMap),
	}

	ret.BuildPackages()

	return ret, nil
}

func (cl *ContextLoader) loadPackages(ctx context.Context, pkgCfg *packages.Config) ([]*packages.Package, error) {
	args := cl.buildArgs()
	cl.log.Debugf("Built loader args are %s", args)
	pkgs, err := loadPackagesFromDir(pkgCfg, args)
	if err != nil {
		return nil, fmt.Errorf("%w failed to load with go/packages", err)
	}

	if ctx.Err() != nil {
		return nil, fmt.Errorf("%w timed out to load packages", ctx.Err())
	}

	cl.debugPrintLoadedPackages(pkgs)
	return pkgs, nil
}

func (cl *ContextLoader) debugPrintLoadedPackages(pkgs []*packages.Package) {
	cl.log.Debugf("loaded %d pkgs", len(pkgs))
	for i, pkg := range pkgs {
		var syntaxFiles []string
		for _, sf := range pkg.Syntax {
			syntaxFiles = append(syntaxFiles, pkg.Fset.Position(sf.Pos()).Filename)
		}
		cl.log.Debugf("Loaded pkg #%d: ID=%s GoFiles=%d Syntax=%d",
			i, pkg.ID, len(pkg.GoFiles), len(syntaxFiles))
	}
}

func loadPackagesFromDir(pkgCfg *packages.Config, args []string) ([]*packages.Package, error) {
	defer util.TrackTime(time.Now(), "Loading source packages")
	pkgs, err := packages.Load(pkgCfg, args...)
	if err != nil {
		return nil, err
	}

	if packages.PrintErrors(pkgs) > 0 {
		logrus.Warn("There are files with errors, the extraction may fail")
	}

	return pkgs, nil
}
