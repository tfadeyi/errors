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
		// ErrorDefinitionURLPath is the parent URL path where the errors will be available
		ErrorDefinitionURLPath string
		// showErrorURL enables and disables the errors' URL being shown when the error is returned
		showErrorURL bool
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
	defaultErrorSpecificationLocation = "./default.yaml"
)

func New(opts ...WrapperOption) *Wrapper {
	wrapper := &Wrapper{
		client:  nil,
		Options: &wrapperOptions{},
	}
	for _, opt := range opts {
		opt(wrapper.Options)
	}

	cl := local.New(errorclient.Options{
		SourceFilename: wrapper.Options.sourceFilename,
		Source:         wrapper.Options.source,
	})

	wrapper.client = cl

	return wrapper
}

func (w *Wrapper) SetManifest(content []byte) {
	w.Options.source = content
	w.client = local.New(errorclient.Options{
		SourceFilename: w.Options.sourceFilename,
		Source:         w.Options.source,
	})
}

func (w *Wrapper) SetManifestFilename(filepath string) {
	w.Options.sourceFilename = filepath
	w.client = local.New(errorclient.Options{
		SourceFilename: w.Options.sourceFilename,
		Source:         w.Options.source,
	})
}

func (w *Wrapper) SetLogger(logger *log.Logger) {
	w.Options.logger = logger
}

func (w *Wrapper) SetErrorParentPath(parentDir string) {
	w.Options.ErrorDefinitionURLPath = parentDir
}

func (w *Wrapper) ShowErrorURL(show bool) {
	w.Options.showErrorURL = show
}

// ErrorWithContext wraps the incoming error with error defined by the Aloe specification according to the input code.
// if no error is found in the specification, the original error is returned.
func (w *Wrapper) ErrorWithContext(ctx context.Context, err error, code string) error {
	if w == nil || err == nil {
		return err
	}

	newErrMessage, genErr := w.client.GenerateErrorMessageFromCode(ctx, code)
	if genErr != nil {
		w.log(genErr.Error())
		return err
	}

	return fmt.Errorf("%s: [%w]", newErrMessage, err)
}

// Error wraps the incoming error with error defined by the application error manifest according to the input code.
// if no error is found in the application error manifest, the original error is returned.
func (w *Wrapper) Error(err error, code string) error {
	if w == nil || err == nil {
		return err
	}
	newErrMessage, errFromGenerator := w.client.GenerateErrorMessageFromCode(context.Background(), code)
	if errFromGenerator != nil {
		w.log(errFromGenerator.Error())
		return err
	}

	return fmt.Errorf("[%w]\n%s", err, newErrMessage)
}

func (w *Wrapper) log(msg string, keyVal ...any) {
	if w.Options.logger != nil {
		w.Options.logger.Printf(msg, keyVal...)
	}
}

// Global Wrapper Functions //

func SetManifest(content []byte) {
	global.SetManifest(content)
}

func SetManifestFilename(filepath string) {
	global.SetManifestFilename(filepath)
}

func SetLogger(logger *log.Logger) {
	global.SetLogger(logger)
}

func SetErrorParentPath(parentDir string) {
	global.SetErrorParentPath(parentDir)
}

func ShowErrorURL(show bool) {
	global.ShowErrorURL(show)
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
		o.ErrorDefinitionURLPath = parentDir
	}
}

func DisableErrorURL() WrapperOption {
	return func(o *wrapperOptions) {
		o.showErrorURL = false
	}
}
