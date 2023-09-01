package parser

import (
	"context"
	"github.com/juju/errors"

	"github.com/tfadeyi/errors/internal/parser/options"
)

type (
	// Parser parses source files containing the sloth definitions
	Parser struct {
		// Opts contains the different options available to the parser.
		// These are applied by the parser constructor during initialization
		Opts *options.Options
	}
)

var (
	ErrNoContentGenerator = errors.New("no target content generator was set")
	ErrNoTargetLanguage   = errors.New("no target source language was set")
)

// New creates a new instance of the parser. See options.Option for more info on the available configuration.
func New(opts ...options.Option) *Parser {
	defaultOpts := new(options.Options)
	for _, opt := range opts {
		opt(defaultOpts)
	}
	for _, opt := range opts {
		opt(defaultOpts)
	}

	return &Parser{defaultOpts}
}

// Parse parses the data source for the target annotations using the given parser configurations and returns a parsed specification.
func (p *Parser) Parse(ctx context.Context) (map[string]any, error) {
	if p.Opts.TargetLanguage == nil {
		return nil, ErrNoTargetLanguage
	}
	return p.Opts.TargetLanguage.Parse(ctx)
}

func (p *Parser) Generate(ctx context.Context, specs map[string]any) error {
	if p.Opts.TargetGenerator == nil {
		return ErrNoContentGenerator
	}
	return p.Opts.TargetGenerator.Generate(ctx, specs)
}
