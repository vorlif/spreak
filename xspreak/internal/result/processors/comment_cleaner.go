package processors

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/vorlif/spreak/xspreak/internal/config"
	"github.com/vorlif/spreak/xspreak/internal/result"
)

const flagPrefix = "xspreak:"

var (
	reRange = regexp.MustCompile(`^range:\s+\d+\.\.\d+\s*$`)
)

type commentCleaner struct {
	allowPrefixes []string
}

var _ Processor = (*commentCleaner)(nil)

func NewCommentCleaner(cfg *config.Config) Processor {
	c := &commentCleaner{
		allowPrefixes: make([]string, 0, len(cfg.CommentPrefixes)),
	}

	for _, prefix := range cfg.CommentPrefixes {
		prefix = strings.TrimSpace(prefix)
		if prefix != "" {
			c.allowPrefixes = append(c.allowPrefixes, prefix)
		}
	}

	return c
}

func (s commentCleaner) Process(inIssues []result.Issue) ([]result.Issue, error) {
	outIssues := make([]result.Issue, 0, len(inIssues)/10)

	for _, iss := range inIssues {
		// remove duplicates and extract text
		commentGroups := make(map[*ast.CommentGroup][]string)
		for _, commentGroup := range iss.CommentGroups {
			commentGroups[commentGroup] = strings.Split(commentGroup.Text(), "\n")
		}

		// filter text
		for _, lines := range commentGroups {
			isTranslatorComment := false
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if s.hasTranslatorPrefix(line) {
					isTranslatorComment = true
				} else if strings.HasPrefix(line, flagPrefix) {
					iss.Flags = append(iss.Flags, parseFlags(line)...)
					isTranslatorComment = false
					continue
				} else if len(line) == 0 {
					isTranslatorComment = false
					continue
				}

				if isTranslatorComment {
					iss.Comment = append(iss.Comment, line)
				}
			}
		}

		// remove groups
		iss.CommentGroups = nil
		outIssues = append(outIssues, iss)
	}

	return outIssues, nil
}

func (s commentCleaner) hasTranslatorPrefix(line string) bool {
	for _, prefix := range s.allowPrefixes {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}

	return false
}

func (s commentCleaner) Name() string {
	return "comment_cleaner"
}

func parseFlags(line string) []string {
	possibleFlags := strings.Split(strings.TrimPrefix(line, flagPrefix), ",")
	flags := make([]string, 0, len(possibleFlags))
	for _, flag := range possibleFlags {
		flag = strings.TrimSpace(flag)

		if strings.HasPrefix(flag, "range:") {
			if !reRange.MatchString(flag) {
				log.WithField("input", flag).Warn("Invalid range flag")
				continue
			}

			rangeFlag := fmt.Sprintf("range: %s", strings.TrimSpace(strings.TrimPrefix(flag, "range:")))
			flags = append(flags, rangeFlag)
		}

		if flag == "ignore" {
			flags = append(flags, flag)
		}
	}

	return flags
}
