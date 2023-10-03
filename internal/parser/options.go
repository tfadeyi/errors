package parser

import (
	"context"
	"github.com/tfadeyi/errors/internal/logging"
	"github.com/tfadeyi/errors/internal/parser/language"
	"github.com/tfadeyi/errors/internal/parser/language/golang"
)

type (
	// Options is a struct contains all the configurations available for the parser
	Options struct {
		// TargetLanguage is the language targeted by the parser, i.e: go.
		TargetLanguage language.Target

		// IncludedDirs is the array containing all the directories that will be parsed by the parser.
		// SourceFile and SourceContent will override this, if present.
		// Option: func Include(dirs ...string) Option
		IncludedDirs []string

		// Logger is the parser's logger
		// Option: func Logger(logger *logging.Logger) Option
		Logger *logging.Logger

		// Strict set the parser to error hard on errors rather than keep processing
		Strict bool

		// Recursive set the parser to recursively parser the target directories
		Recursive bool
	}
	// Option is a more atomic to configure the different Options rather than passing the entire Options struct.
	Option func(p *Options)
)

// Include configure the parser to parse the given included directories
// SourceFile and SourceContent will override this, if present.
func Include(dirs ...string) Option {
	return func(e *Options) {
		e.IncludedDirs = dirs
	}
}

// Logger configure the parser's logger
func Logger(logger *logging.Logger) Option {
	return func(e *Options) {
		log := logger.WithName("parser")
		e.Logger = &log
	}
}

func Strict(strict bool) Option {
	return func(e *Options) {
		e.Strict = strict
	}
}

func Recursive(recursive bool) Option {
	return func(e *Options) {
		e.Recursive = recursive
	}
}

// Custom configurations

// Go returns the options.Option to run the parser targeting golang source code
func Go(ctx context.Context) Option {
	return func(opts *Options) {
		opts.TargetLanguage = golang.NewParser(&golang.Options{
			Ctx:              ctx,
			Logger:           opts.Logger,
			InputDirectories: opts.IncludedDirs,
			Recursive:        opts.Recursive,
		})
	}
}
