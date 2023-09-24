package docs

import (
	"context"

	"github.com/juju/errors"
	"github.com/tfadeyi/errors/internal/generate/docs/markdown"
	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"
)

type (
	// Target is an abstraction for content generators
	Target interface {
		// GenerateDocumentation
		GenerateDocumentation(ctx context.Context, manifest *api.Manifest) error
	}

	Generator struct {
		Opts *Options
	}
)

var (
	ErrNoContentGenerator = errors.New("no target content generator was set. i.e: Markdown")
)

// New creates a new instance of the Generator. See options.Option for more info on the available configuration.
func New(opts ...Option) *Generator {
	defaultOpts := new(Options)
	for _, opt := range opts {
		opt(defaultOpts)
	}
	for _, opt := range opts {
		opt(defaultOpts)
	}

	return &Generator{defaultOpts}
}

func (g *Generator) GenerateDocumentation(ctx context.Context, manifest *api.Manifest) error {
	if g.Opts.TargetDocGenerator == nil {
		return ErrNoContentGenerator
	}
	return g.Opts.TargetDocGenerator.GenerateDocumentation(ctx, manifest)
}

// Markdown returns the options.Option to run the parser generator for Markdown
func Markdown() Option {
	return func(opts *Options) {
		opts.TargetDocGenerator = markdown.New(&markdown.Options{
			Logger:        opts.Logger,
			Writer:        opts.Writer,
			Output:        opts.Output,
			InfoTmplFile:  opts.CustomInfoTemplateFilepath,
			ErrorTmplFile: opts.CustomErrorTemplateFilepath,
		})
	}
}
