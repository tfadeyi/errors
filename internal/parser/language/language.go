package language

import (
	"context"
	"strings"

	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"
)

type (
	// Target is the source code language the parser should parse.
	Target interface {
		// ParseSource returns manifest(s) struct given a data source, returns error if parsing fails
		ParseSource(ctx context.Context) (map[string]api.Manifest, error)
		// AnnotateErrors annotates with comments and wraps error definitions and declarations
		AnnotateErrors(ctx context.Context, wrapErrors bool) error
	}
)

const (
	Go = "go"
)

// IsSupportedLanguage is returns true if the input is a supported language
func IsSupportedLanguage(input string) bool {
	return strings.ToLower(input) == Go
}
