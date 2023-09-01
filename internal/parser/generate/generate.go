package generate

import (
	"context"
	"strings"
)

type (
	// Target is an abstraction for content generators, i.e: yaml and markdown.
	Target interface {
		// Generate transcribes the content from the specs to the writer
		Generate(ctx context.Context, specs map[string]any) error
	}
)

// IsSupportedOutputFormat checks if the given output format is a supported one
func IsSupportedOutputFormat(format string) bool {
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case Yaml, Markdown:
		return true
	}
	return false
}

const (
	Yaml     = "yaml"
	Markdown = "markdown"
)
