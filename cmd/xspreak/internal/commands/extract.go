package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/vorlif/spreak/xspreak/internal/goextractors"
	"github.com/vorlif/spreak/xspreak/internal/result"

	"github.com/vorlif/spreak/xspreak/internal/encoder"
	"github.com/vorlif/spreak/xspreak/internal/extract"
	"github.com/vorlif/spreak/xspreak/internal/extract/extractors"
)

func (e *Executor) initExtract() {
	e.extractCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the extraction",
		Run:   e.executeRun,
	}
	e.rootCmd.AddCommand(e.extractCmd)
}

func (e *Executor) executeRun(_ *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), e.cfg.Timeout)
	defer cancel()

	extractedIssues, errE := e.runExtraction(ctx, args)
	if errE != nil {
		e.log.Fatalf("Running error: %s", errE)
	}

	domainIssues := make(map[string][]result.Issue)
	for _, iss := range extractedIssues {
		if _, ok := domainIssues[iss.Domain]; !ok {
			domainIssues[iss.Domain] = []result.Issue{iss}
		} else {
			domainIssues[iss.Domain] = append(domainIssues[iss.Domain], iss)
		}
	}

	if len(extractedIssues) == 0 {
		domainIssues[""] = make([]result.Issue, 0)
		log.Println("No Strings found")
	}

	for domainName, issues := range domainIssues {
		var outputFile string
		if domainName == "" {
			outputFile = filepath.Join(e.cfg.OutputDir, e.cfg.OutputFile)
		} else {
			outputFile = filepath.Join(e.cfg.OutputDir, domainName+".pot")
		}

		outputDir := filepath.Dir(outputFile)
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			log.Printf("Output folder does not exist, trying to create it: %s\n", outputDir)
			if errC := os.MkdirAll(outputDir, os.ModePerm); errC != nil {
				log.Fatalf("Output folder does not exist and could not be created: %s", errC)
			}
		}

		dst, err := os.Create(outputFile)
		if err != nil {
			e.log.WithError(err).Fatal("Output file could not be created")
		}
		defer dst.Close()

		enc := encoder.NewPotEncoder(e.cfg, dst)
		if errEnc := enc.Encode(issues); errEnc != nil {
			e.log.WithError(errEnc).Fatal("Output file could not be written")
		}

		_ = dst.Close()
		log.Printf("File written: %s\n", outputFile)
	}
}

func (e *Executor) runExtraction(ctx context.Context, args []string) ([]result.Issue, error) {
	e.cfg.Args = args

	extractorsToRun := []extractors.Extractor{
		goextractors.NewDefinitionExtractor(),
		goextractors.NewCommentsExtractor(),
		goextractors.NewFuncCallExtractor(),
		goextractors.NewGlobalAssignExtractor(),
		goextractors.NewSliceDefExtractor(),
		goextractors.NewStructDefExtractor(),
		goextractors.NewVariablesExtractor(),
		goextractors.NewErrorExtractor(),
	}

	extractCtx, err := e.contextLoader.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("context loading failed: %w", err)
	}

	runner, err := extract.NewRunner(e.cfg, extractCtx.Packages)
	if err != nil {
		return nil, err
	}

	issues, err := runner.Run(ctx, extractCtx, extractorsToRun)
	if err != nil {
		return nil, err
	}

	return issues, nil
}
