package goextractors

import (
	"go/parser"
	"testing"
)

func TestExtractStringLiteral(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		wantStr   string
		wantFound bool
	}{
		{
			name:      "String extracted",
			code:      `"Extracted string"`,
			wantStr:   `Extracted string`,
			wantFound: true,
		},
		{
			name:      "Even addition is merged",
			code:      `"Extracted " + "string"`,
			wantStr:   `Extracted string`,
			wantFound: true,
		},
		{
			name:      "Odd addition is merged",
			code:      `"Extracted " + "string" + " is combined"`,
			wantStr:   `Extracted string is combined`,
			wantFound: true,
		},
		{
			name:      "Backquotes are removed",
			code:      "`Extracted string`",
			wantStr:   "Extracted string",
			wantFound: true,
		},
		{
			name:      "Multiline text with backquotes are formatted correctly",
			code:      "`This is an multiline\nstring`",
			wantStr:   "This is an multiline\nstring",
			wantFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.code)
			if err != nil {
				t.Errorf("Expression %s could not be parsed: %v", tt.code, expr)
			}
			extractedStr, found := ExtractStringLiteral(expr)
			if extractedStr != tt.wantStr {
				t.Errorf("ExtractStringLiteral() string = %v, want %v", extractedStr, tt.wantStr)
			}
			if (found == nil) == tt.wantFound {
				t.Errorf("ExtractStringLiteral() got1 = %v, want %v", found, (found == nil))
			}
		})
	}

	t.Run("Nil is ignored", func(t *testing.T) {
		extractedStr, found := ExtractStringLiteral(nil)
		if extractedStr != "" {
			t.Errorf("ExtractStringLiteral() string = %v, want %v", extractedStr, "")
		}
		if found != nil {
			t.Errorf("ExtractStringLiteral() got1 = %v, want %v", found, false)
		}
	})
}
