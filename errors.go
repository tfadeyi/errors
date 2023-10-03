package errors

import (
	"context"
	"fmt"
	"github.com/tfadeyi/errors/internal/errorclient"
	"github.com/tfadeyi/errors/internal/errorclient/local"
	"log"
)

type (
	// WrapperOption is a more atomic to configure the different wrapperOptions rather than passing the entire Options struct.
	WrapperOption func(o *wrapperOptions)

	// wrapperOptions contains all the different configuration values available to the wrapper
	wrapperOptions struct {
		// logger is internal logger for the wrapper, leave nil to avoid logging
		// WrapperOption: func logger(logger *log.Logger) WrapperOption
		logger *log.Logger
		// sourceFilename allows clients to set the location of the file containing the error specification
		// WrapperOption: func ManifestFilename(filename string) WrapperOption
		sourceFilename string
		// source allows clients to pass the contents of the error specification file as a []byte
		// WrapperOption: func Manifest(source []byte) WrapperOption
		source []byte
		// errorDefinitionURLPath is the parent URL path where the errors will be available
		errorDefinitionURLPath string
		// numberOfSuggestionsDisplayed number of suggestions shown when the error is returned
		numberOfSuggestionsDisplayed int
		// showMarkdownErrors enables or disables the error's pretty view being shown at error return time
		showMarkdownErrors bool
		// Silent stops the additional error context from being wrapped into the error
		silent bool
	}

	// Wrapper is the wrapper struct for errors
	Wrapper struct {
		Options *wrapperOptions
		client  errorclient.Client
	}
)

var (
	global = New(ManifestFilename(defaultErrorSpecificationLocation))
)

const (
	defaultErrorSpecificationLocation = "./errors.yaml"
)

func New(opts ...WrapperOption) *Wrapper {
	wrapper := &Wrapper{
		client: nil,
		Options: &wrapperOptions{
			numberOfSuggestionsDisplayed: 1,
			showMarkdownErrors:           true,
			silent:                       false,
		},
	}
	for _, opt := range opts {
		opt(wrapper.Options)
	}

	cl := local.New(errorclient.Options{
		SourceFilename:         wrapper.Options.sourceFilename,
		Source:                 wrapper.Options.source,
		ErrorDefinitionURLPath: wrapper.Options.errorDefinitionURLPath,
		ShowMarkdownErrors:     wrapper.Options.showMarkdownErrors,
		NumberOfSuggestions:    wrapper.Options.numberOfSuggestionsDisplayed,
	})

	wrapper.client = cl

	return wrapper
}

func (w *Wrapper) Manifest(content []byte) {
	w.Options.source = content
	w.updateErrorClient()
}

func (w *Wrapper) ManifestFilename(filepath string) {
	w.Options.sourceFilename = filepath
	w.updateErrorClient()
}

func (w *Wrapper) Logger(logger *log.Logger) {
	w.Options.logger = logger
}

func (w *Wrapper) ErrorParentPath(parentDir string) {
	w.Options.errorDefinitionURLPath = parentDir
	w.updateErrorClient()
}

func (w *Wrapper) Silence(silence bool) {
	w.Options.silent = silence
	w.updateErrorClient()
}

func (w *Wrapper) DisplayedSuggestions(num int) {
	w.Options.numberOfSuggestionsDisplayed = num
	w.updateErrorClient()
}

func (w *Wrapper) MarkdownRender(markdown bool) {
	w.Options.showMarkdownErrors = markdown
	w.updateErrorClient()
}

func (w *Wrapper) updateErrorClient() {
	w.client = local.New(errorclient.Options{
		SourceFilename:         w.Options.sourceFilename,
		Source:                 w.Options.source,
		ErrorDefinitionURLPath: w.Options.errorDefinitionURLPath,
		ShowMarkdownErrors:     w.Options.showMarkdownErrors,
		NumberOfSuggestions:    w.Options.numberOfSuggestionsDisplayed,
	})
}

// ErrorWithContext wraps the incoming error with error defined by the Aloe specification according to the input code.
// if no error is found in the specification, the original error is returned.
func (w *Wrapper) ErrorWithContext(ctx context.Context, err error, code string) error {
	if w == nil || err == nil {
		return err
	}

	if !w.Options.silent {
		newErrMessage, genErr := w.client.GenerateErrorMessageFromCode(ctx, code)
		if genErr != nil {
			w.log(genErr.Error())
			return err
		}
		return fmt.Errorf("[%w]\n%s", err, newErrMessage)
	}

	return err
}

// Error wraps the incoming error with error defined by the application error manifest according to the input code.
// if no error is found in the application error manifest, the original error is returned.
func (w *Wrapper) Error(err error, code string) error {
	return w.ErrorWithContext(context.Background(), err, code)
}

func (w *Wrapper) log(msg string, keyVal ...any) {
	if w.Options.logger != nil {
		w.Options.logger.Printf(msg, keyVal...)
	}
}

// Global Wrapper Functions //

func SetManifest(content []byte) {
	global.Manifest(content)
}

func SetManifestFilename(filepath string) {
	global.ManifestFilename(filepath)
}

func SetLogger(logger *log.Logger) {
	global.Logger(logger)
}

func SetErrorParentPath(parentDir string) {
	global.ErrorParentPath(parentDir)
}

func SetSilence(silence bool) {
	global.Silence(silence)
}

func SetDisplayedSuggestions(num int) {
	global.DisplayedSuggestions(num)
}

func SetMarkdownRender(markdown bool) {
	global.MarkdownRender(markdown)
}

// ErrorWithContext wraps the incoming error with error defined by the application error manifest according to the input code.
// if no error is found in the application error manifest, the original error is returned.
func ErrorWithContext(ctx context.Context, err error, code string) error {
	return global.ErrorWithContext(ctx, err, code)
}

// Error wraps the incoming error with error defined by the application error manifest according to the input code.
// if no error is found in the application error manifest, the original error is returned.
func Error(err error, code string) error {
	return global.Error(err, code)
}

// WrapperOption Functions //

func Manifest(source []byte) WrapperOption {
	return func(o *wrapperOptions) {
		o.source = source
	}
}

func ManifestFilename(filename string) WrapperOption {
	return func(o *wrapperOptions) {
		o.sourceFilename = filename
	}
}

func Logger(logger *log.Logger) WrapperOption {
	return func(o *wrapperOptions) {
		o.logger = logger
	}
}

func ErrorParentPath(parentDir string) WrapperOption {
	return func(o *wrapperOptions) {
		o.errorDefinitionURLPath = parentDir
	}
}

func Silence(silence bool) WrapperOption {
	return func(o *wrapperOptions) {
		o.silent = silence
	}
}

func DisplayedSuggestions(num int) WrapperOption {
	return func(o *wrapperOptions) {
		o.numberOfSuggestionsDisplayed = num
	}
}

func MarkdownRender(markdown bool) WrapperOption {
	return func(o *wrapperOptions) {
		o.showMarkdownErrors = markdown
	}
}
