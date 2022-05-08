package commands

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/vorlif/spreak/xspreak/internal/config"
	"github.com/vorlif/spreak/xspreak/internal/extract"
)

type Executor struct {
	rootCmd    *cobra.Command
	extractCmd *cobra.Command

	cfg *config.Config
	log *logrus.Entry

	contextLoader *extract.ContextLoader
}

func NewExecutor() *Executor {
	e := &Executor{
		cfg: config.NewDefault(),
		log: logrus.WithField("service", "executor"),
	}

	e.rootCmd = &cobra.Command{
		Use:     "xspreak",
		Version: VersionName,
		Args:    cobra.ArbitraryArgs,
		Short:   "String extraction for spreak.",
		Long:    `Simple tool to extract strings and create POT files for application translations.`,
		Run:     e.executeRun,
	}

	fs := pflag.NewFlagSet("config flag set", pflag.ContinueOnError)
	fs.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}
	fs.Usage = func() {}
	fs.SortFlags = false
	e.rootCmd.PersistentFlags().SortFlags = false
	initRootFlags(e.rootCmd.PersistentFlags(), config.NewDefault())
	initRootFlags(fs, e.cfg)

	if err := fs.Parse(os.Args); err != nil {
		if err != pflag.ErrHelp {
			logrus.WithError(err).Fatal("Args could not be parsed")
		}
	}

	if err := e.cfg.Prepare(); err != nil {
		logrus.Fatalf("Configuration could not be processed %v", err)
	}

	if e.cfg.IsVerbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	e.log.Debug("Starting execution...")

	e.initExtract()

	e.contextLoader = extract.NewContextLoader(e.cfg)

	return e
}

func (e *Executor) Execute() error {
	return e.rootCmd.Execute()
}

func initRootFlags(fs *pflag.FlagSet, cfg *config.Config) {
	def := config.NewDefault()

	fs.BoolVarP(&cfg.IsVerbose, "verbose", "V", def.IsVerbose, "author name for copyright attribution")
	fs.StringVarP(&cfg.SourceDir, "directory", "D", def.SourceDir, "set DIRECTORY to list for input files search")
	fs.StringVarP(&cfg.OutputDir, "output-dir", "p", def.OutputDir, "output files will be placed in directory DIR.")
	fs.StringVarP(&cfg.OutputFile, "output", "o", def.OutputFile, "write output to specified file")
	fs.StringSliceVarP(&cfg.CommentPrefixes, "add-comments", "c", def.CommentPrefixes, "place comment blocks starting with TAG and preceding keyword lines in output file")
	fs.BoolVarP(&cfg.ExtractErrors, "extract-errors", "e", def.ExtractErrors, "strings from errors.New(STRING) are extracted")
	fs.StringVar(&cfg.ErrorContext, "errors-context", def.ErrorContext, "context which is automatically assigned to extracted errors")

	fs.StringVarP(&cfg.DefaultDomain, "default-domain", "d", def.DefaultDomain, "use NAME.po for output (instead of messages.po)")
	fs.BoolVar(&cfg.WriteNoLocation, "no-location", def.WriteNoLocation, "do not write '#: filename:line' lines")

	fs.IntVarP(&cfg.WrapWidth, "width", "w", def.WrapWidth, "set output page width")
	fs.BoolVar(&cfg.DontWrap, "no-wrap", def.DontWrap, "do not break long message lines, longer than the output page width, into several lines")

	fs.BoolVar(&cfg.OmitHeader, "omit-header", def.OmitHeader, "don't write header with 'msgid \"\"' entry")
	fs.StringVar(&cfg.CopyrightHolder, "copyright-holder", def.CopyrightHolder, "set copyright holder in output")
	fs.StringVar(&cfg.PackageName, "package-name", def.PackageName, "set package name in output")
	fs.StringVar(&cfg.BugsAddress, "msgid-bugs-address", def.BugsAddress, "set report address for msgid bugs")

	fs.DurationVar(&cfg.Timeout, "timeout", def.Timeout, "Timeout for total work")
}
