package po

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner_ScanComment(t *testing.T) {
	t.Run("header comment", func(t *testing.T) {
		comment := `# This file is distributed under the same license as the Django package.
#
# Translators:
# F Wolff <friedel@translate.org.za>, 2019-2020
# Stephen Cox <stephencoxmail@gmail.com>, 2011-2012   
# unklphil <villiers.strauss@gmail.com>, 2014,2019`

		s := newScanner(bytes.NewReader([]byte(comment)))
		assert.NotNil(t, s)

		tests := []string{
			"# This file is distributed under the same license as the Django package.",
			"#",
			"# Translators:",
			"# F Wolff <friedel@translate.org.za>, 2019-2020",
			"# Stephen Cox <stephencoxmail@gmail.com>, 2011-2012",
			"# unklphil <villiers.strauss@gmail.com>, 2014,2019",
		}

		for _, tt := range tests {
			tok, lit := s.scan()
			assert.Equal(t, commentTranslator, tok)
			assert.Equal(t, tt, lit)
		}
	})
}
