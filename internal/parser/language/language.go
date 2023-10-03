package language

import (
	"context"
	"io"
	"strings"

	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"
)

type (
	// Target is the source code language the parser should parse.
	Target interface {
		// ParseSource parses the data source for @fyi annotations and returns an error manifest if successfully parsed.
		// Returns error if the parsing process fails
		ParseSource(ctx context.Context, content io.Reader) (*api.Manifest, error)
		// ParseAllSources parses the included directories for @fyi annotations and returns an error manifest if successfully parsed.
		// Returns error if the parsing process fails
		ParseAllSources(ctx context.Context) ([]*api.Manifest, error)
		// AnnotateAllErrors annotates all errors in the included directories with fyi comments and error wrapping
		AnnotateAllErrors(ctx context.Context, wrapErrors, annotateOnlyTodos bool) error
		// AnnotateSourceErrors annotates all errors in the content with fyi comments and error wrapping
		AnnotateSourceErrors(ctx context.Context, filename string, wrapErrors, annotateOnlyTodos bool) error
	}
)

const (
	// Go language parser const
	Go = "go"
)

// IsSupportedLanguage is returns true if the input is a supported language
func IsSupportedLanguage(input string) bool {
	return strings.ToLower(input) == Go
}
