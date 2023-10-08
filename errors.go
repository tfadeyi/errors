package errors

import (
	"context"
	"fmt"
	"github.com/tfadeyi/errors/pkg/errorclient"
	"github.com/tfadeyi/errors/pkg/errorclient/local"
	"log"
)

type (
	// HandlerOption is a more atomic to configure the different HandlerOptions rather than passing the entire ErrClientOptions struct.
	HandlerOption func(o *HandlerOptions)

	// HandlerOptions contains all the different configuration values available to the wrapper
	HandlerOptions struct {
		// Logger is internal Logger for the wrapper, leave nil to avoid logging
		Logger *log.Logger

		// Silent globally sets the error handler to stop adding the additional error context in the input error of Error() and ErrorWithContext()
		Silent bool

		errorclient.ErrClientOptions
	}

	// Handler is the wrapper struct for errors
	Handler struct {
		Options *HandlerOptions
	}
)

var (
	global = New()
)

const (
	defaultErrorSpecificationLocation = "./errors.yaml"
)

/*
New creates and configures a new instance of the error handler.

example:

errHandler := New(WithManifest("content"),WithNoLogger(),WithRenderMarkdown(true))

errHandler.Error(err, "error code")
*/
func New(opts ...HandlerOption) *Handler {
	newOpts := &HandlerOptions{
		log.Default(),
		false,
		errorclient.ErrClientOptions{
			SourceFilename:         defaultErrorSpecificationLocation,
			Source:                 nil,
			ErrorDefinitionURLPath: "",
			DisplayMarkdownErrors:  true,
			NumberOfSuggestions:    1,
			OverrideErrorURL:       "",
			DisplayShortSummary:    true,
			DisplayErrorURL:        true,
		},
	}

	for _, opt := range opts {
		opt(newOpts)
	}
	for _, opt := range opts {
		opt(newOpts)
	}

	return &Handler{
		Options: newOpts,
	}
}

// SetManifest sets the source manifest used by the Handler
// note: it is not required if the SetManifestFilename is set
func (w *Handler) SetManifest(content []byte) {
	w.Options.Source = content
}

// SetManifestFilename sets the file path of the manifest used by the Handler.
// note: it is not required if SetManifest is set
func (w *Handler) SetManifestFilename(filepath string) {
	w.Options.SourceFilename = filepath
}

// SetLogger sets the logger used by the Handler internally.
// note: pass nil if no logging is wanted
func (w *Handler) SetLogger(logger *log.Logger) {
	w.Options.Logger = logger
}

// SetErrorParentPath
func (w *Handler) SetErrorParentPath(parentDir string) {
	w.Options.ErrorDefinitionURLPath = parentDir
}

// SetSilence
func (w *Handler) SetSilence(silence bool) {
	w.Options.Silent = silence
}

// SetDisplayedSuggestions
func (w *Handler) SetDisplayedSuggestions(num int) {
	w.Options.NumberOfSuggestions = num
}

// SetMarkdownRender
func (w *Handler) SetMarkdownRender(markdown bool) {
	w.Options.DisplayMarkdownErrors = markdown
}

// SetOverrideErrorURL
func (w *Handler) SetOverrideErrorURL(url string) {
	w.Options.OverrideErrorURL = url
}

// SetShowShortSummary
func (w *Handler) SetShowShortSummary(flag bool) {
	w.Options.DisplayShortSummary = flag
}

// SetDisplayErrorURL
func (w *Handler) SetDisplayErrorURL(flag bool) {
	w.Options.DisplayErrorURL = flag
}

// ErrorWithContext wraps the incoming error with error defined by the Aloe specification according to the input code.
// if no error is found in the specification, the original error is returned.
func (w *Handler) ErrorWithContext(ctx context.Context, err error, code string, opts ...HandlerOption) error {
	if w == nil || err == nil {
		return err
	}

	for _, opt := range opts {
		opt(w.Options)
	}
	client := local.New(w.Options.ErrClientOptions)

	if !w.Options.Silent {
		newErrMessage, genErr := client.GenerateErrorMessageFromCode(ctx, code)
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
func (w *Handler) Error(err error, code string, opts ...HandlerOption) error {
	return w.ErrorWithContext(context.Background(), err, code, opts...)
}

func (w *Handler) log(msg string, keyVal ...any) {
	if w.Options.Logger != nil {
		w.Options.Logger.Printf(msg, keyVal...)
	}
}

// Global Handler Functions //

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

func SetSilence(silence bool) {
	global.SetSilence(silence)
}

func SetDisplayedSuggestions(num int) {
	global.SetDisplayedSuggestions(num)
}

func SetMarkdownRender(markdown bool) {
	global.SetMarkdownRender(markdown)
}

// SetOverrideErrorURL
func SetOverrideErrorURL(url string) {
	global.Options.OverrideErrorURL = url
}

// SetShowShortSummary
func SetShowShortSummary(flag bool) {
	global.Options.DisplayShortSummary = flag
}

// SetDisplayErrorURL
func SetDisplayErrorURL(flag bool) {
	global.Options.DisplayErrorURL = flag
}

// ErrorWithContext wraps the incoming error with error defined by the application error manifest according to the input code.
// if no error is found in the application error manifest, the original error is returned.
func ErrorWithContext(ctx context.Context, err error, code string, opts ...HandlerOption) error {
	return global.ErrorWithContext(ctx, err, code, opts...)
}

// Error wraps the incoming error with error defined by the application error manifest according to the input code.
// if no error is found in the application error manifest, the original error is returned.
func Error(err error, code string, opts ...HandlerOption) error {
	return global.Error(err, code, opts...)
}

// HandlerOption Functions //

func WithManifest(source []byte) HandlerOption {
	return func(o *HandlerOptions) {
		o.Source = source
	}
}

func WithManifestFilename(filename string) HandlerOption {
	return func(o *HandlerOptions) {
		o.SourceFilename = filename
	}
}

func WithLogger(logger *log.Logger) HandlerOption {
	return func(o *HandlerOptions) {
		o.Logger = logger
	}
}

func WithNoLogger() HandlerOption {
	return func(o *HandlerOptions) {
		o.Logger = nil
	}
}

func WithErrorParentPath(parentDir string) HandlerOption {
	return func(o *HandlerOptions) {
		o.ErrorDefinitionURLPath = parentDir
	}
}

func WithSilence(silence bool) HandlerOption {
	return func(o *HandlerOptions) {
		o.Silent = silence
	}
}

func WithNumberOfSuggestions(num int) HandlerOption {
	return func(o *HandlerOptions) {
		o.NumberOfSuggestions = num
	}
}

func WithRenderMarkdown(markdown bool) HandlerOption {
	return func(o *HandlerOptions) {
		o.DisplayMarkdownErrors = markdown
	}
}

func WithOverrideErrorURL(url string) HandlerOption {
	return func(o *HandlerOptions) {
		o.OverrideErrorURL = url
	}
}

func WithShowShortSummary(flag bool) HandlerOption {
	return func(o *HandlerOptions) {
		o.DisplayShortSummary = flag
	}
}

func WithDisplayErrorURL(flag bool) HandlerOption {
	return func(o *HandlerOptions) {
		o.DisplayErrorURL = flag
	}
}
