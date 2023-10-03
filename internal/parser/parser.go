package parser

import (
	"context"
	"github.com/juju/errors"
	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"
	"io"
)

type (
	// Parser parses source files containing the sloth definitions
	Parser struct {
		// Opts contains the different options available to the parser.
		// These are applied by the parser constructor during initialization
		Opts *Options
	}
)

var (
	ErrNoTargetLanguage = errors.New("no target source language was set")
)

// New creates a new instance of the parser. See options.Option for more info on the available configuration.
func New(opts ...Option) *Parser {
	defaultOpts := new(Options)
	for _, opt := range opts {
		opt(defaultOpts)
	}
	for _, opt := range opts {
		opt(defaultOpts)
	}

	return &Parser{defaultOpts}
}

// ParseSource parses the data source for the target annotations using the given parser configurations and returns a parsed specification.
func (p *Parser) ParseSource(ctx context.Context, content io.Reader) (*api.Manifest, error) {
	if p.Opts.TargetLanguage == nil {
		return nil, ErrNoTargetLanguage
	}
	return p.Opts.TargetLanguage.ParseSource(ctx, content)
}

// ParseAllSources parses the data source for the target annotations using the given parser configurations and returns a parsed specification.
func (p *Parser) ParseAllSources(ctx context.Context) ([]*api.Manifest, error) {
	if p.Opts.TargetLanguage == nil {
		return nil, ErrNoTargetLanguage
	}
	return p.Opts.TargetLanguage.ParseAllSources(ctx)
}

// AnnotateErrors parses the data source for error declarations using the given parser configurations and returns them.
func (p *Parser) AnnotateAllErrors(ctx context.Context, wrapErrors, annotateOnlyTodos bool) error {
	if p.Opts.TargetLanguage == nil {
		return ErrNoTargetLanguage
	}
	return p.Opts.TargetLanguage.AnnotateAllErrors(ctx, wrapErrors, annotateOnlyTodos)
}

// AnnotateErrors parses the data source for error declarations using the given parser configurations and returns them.
func (p *Parser) AnnotateSourceErrors(ctx context.Context, filename string, wrapErrors, annotateOnlyTodos bool) error {
	if p.Opts.TargetLanguage == nil {
		return ErrNoTargetLanguage
	}
	return p.Opts.TargetLanguage.AnnotateSourceErrors(ctx, filename, wrapErrors, annotateOnlyTodos)
}
