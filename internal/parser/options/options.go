package options

import (
	"io"

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

		// SourceFile is the file the parser will parse. Shouldn't be used together with SourceContent
		// Option: func SourceFile(file string) Option
		SourceFile string

		// SourceContent is the io.Reader the parser will parse. Shouldn't be used together with SourceFile
		// Option: func SourceContent(content io.ReadCloser) Option
		SourceContent io.ReadCloser

		Output string
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

// SourceFile configure the parser to parse a specific file
// Shouldn't be used together with SourceContent
func SourceFile(file string) Option {
	return func(e *Options) {
		e.SourceFile = file
	}
}

// SourceContent configure the parser to parse a specific io.Reader
// Shouldn't be used together with SourceFile
func SourceContent(content io.ReadCloser) Option {
	return func(e *Options) {
		e.SourceContent = content
	}
}

func Output(dir string) Option {
	return func(e *Options) {
		e.Output = dir
	}
}

// Custom configurations

// Go returns the options.Option to run the parser targeting golang source code
func Go() Option {
	return func(opts *Options) {
		opts.TargetLanguage = golang.NewParser(&golang.Options{
			Logger:           opts.Logger,
			SourceFile:       opts.SourceFile,
			SourceContent:    opts.SourceContent,
			InputDirectories: opts.IncludedDirs,
		})
	}
}
