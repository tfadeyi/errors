//go:generate rm -f ./pkg/api/api.go
//go:generate gojsonschema -p api ./schema/schema.json -o ./pkg/api/api.go

package goaloe

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tfadeyi/go-aloe/internal/client"
	"github.com/tfadeyi/go-aloe/internal/client/local"
)

type (
	errorHandler struct {
		client client.Client
		logger *log.Logger
	}

	// Options for the aloe error handler. Use it to configure the aloe error handler
	Options struct {
		// SourceFilename is the location of the file containing the aloe specification for the target service
		SourceFilename string
		// Source is the in-memory aloe specification for the target service
		Source []byte
		// Logger is the logger used by the error handler to report errors that might occur during execution.
		// Set as nil, if no logging is wanted.
		Logger *log.Logger
	}
)

const (
	DefaultAloeFilename = "default.aloe"
)

// newInstance generates a new instance of the aloe error handler
func newInstance(opts Options) (*errorHandler, error) {
	cl, err := local.New(opts.SourceFilename, opts.Source)
	if err != nil {
		return nil, err
	}

	return &errorHandler{
		logger: opts.Logger,
		client: cl,
	}, nil
}

func lazyInstance(opts Options) *errorHandler {
	cl := local.NewLazy(opts.SourceFilename, opts.Source)

	return &errorHandler{
		logger: opts.Logger,
		client: cl,
	}
}

// DefaultOrDie returns an instance of the Aloe error handler, using the default configuration, or panics.
// It expects the aloe specification to be in the current working directory, as default.aloe.(yaml|json|toml) .
func DefaultOrDie() *errorHandler {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("go-aloe was unable to get the current working directory: [%w]", err))
	}

	specs := []string{
		fmt.Sprintf("%s/%s.toml", dir, DefaultAloeFilename),
		fmt.Sprintf("%s/%s.yaml", dir, DefaultAloeFilename),
		fmt.Sprintf("%s/%s.json", dir, DefaultAloeFilename)}
	var instance *errorHandler

	for _, src := range specs {
		instance, err = newInstance(Options{
			SourceFilename: src,
			Logger:         log.Default(),
		})
		if err == nil {
			break
		}
	}

	if err != nil {
		panic(fmt.Errorf("go-aloe default handler failed to initialise: [%w]", err))
	}

	return instance
}

// Default returns an instance of the Aloe error handler, using the default configuration, or errors.
// It expects the aloe specification to be in the current working directory, as default.aloe.(yaml|json|toml) .
func Default() (*errorHandler, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("go-aloe was unable to get current working directory: [%w]", err)
	}

	specs := []string{
		fmt.Sprintf("%s/%s.toml", dir, DefaultAloeFilename),
		fmt.Sprintf("%s/%s.yaml", dir, DefaultAloeFilename),
		fmt.Sprintf("%s/%s.json", dir, DefaultAloeFilename)}
	var instance *errorHandler

	for _, src := range specs {
		instance, err = newInstance(Options{
			SourceFilename: src,
			Logger:         log.Default(),
		})
	}

	if err != nil {
		return nil, fmt.Errorf("go-aloe default handler failed to initialise: [%w]", err)
	}

	return instance, nil
}

// WithOptions returns an instance of the Aloe error handler, with the given configuration.
// Use this to remove logging or specify a different location for the Aloe specification file.
func WithOptions(opts Options) *errorHandler {
	instance := lazyInstance(opts)
	return instance
}

// ErrorWithContext wraps the incoming error with error defined by the Aloe specification according to the input code.
// if no error is found in the specification, the original error is returned.
func (instance *errorHandler) ErrorWithContext(ctx context.Context, err error, code string) error {
	if err == nil {
		return err
	}

	newErrMessage, genErr := instance.client.GenerateErrorMessageFromCode(ctx, code)
	if genErr != nil {
		instance.log(genErr.Error())
		return err
	}

	return fmt.Errorf("%s: [%w]", newErrMessage, err)
}

// Error wraps the incoming error with error defined by the Aloe specification according to the input code.
// if no error is found in the specification, the original error is returned.
func (instance *errorHandler) Error(err error, code string) error {
	if err == nil {
		return err
	}

	newErrMessage, genErr := instance.client.GenerateErrorMessageFromCode(context.Background(), code)
	if genErr != nil {
		instance.log(genErr.Error())
		return err
	}

	return fmt.Errorf("%s: [%w]", newErrMessage, err)
}

func (instance *errorHandler) log(msg string, keyVal ...any) {
	if instance.logger != nil {
		instance.logger.Printf(msg, keyVal...)
	}
}
