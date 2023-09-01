package language

import (
	"context"
)

type (
	// Target is the source code language the parser should parse.
	Target interface {
		// Parse returns specification(s) struct given a data source, returns error if parsing fails
		Parse(ctx context.Context) (map[string]any, error)
	}
)

const (
	Go = "go"
)
