package parser

import (
	"context"
	"github.com/juju/errors"
	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"

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
	ErrNoTargetLanguage = errors.New("no target source language was set")
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

// ParseSource parses the data source for the target annotations using the given parser configurations and returns a parsed specification.
func (p *Parser) ParseSource(ctx context.Context) (map[string]api.Manifest, error) {
	if p.Opts.TargetLanguage == nil {
		return nil, ErrNoTargetLanguage
	}
	return p.Opts.TargetLanguage.ParseSource(ctx)
}

// AnnotateErrors parses the data source for error declarations using the given parser configurations and returns them.
func (p *Parser) AnnotateErrors(ctx context.Context, wrapErrors bool) error {
	if p.Opts.TargetLanguage == nil {
		return ErrNoTargetLanguage
	}
	return p.Opts.TargetLanguage.AnnotateErrors(ctx, wrapErrors)
}
