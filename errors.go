package goaloe

import (
	"context"
	"fmt"
	"log"

	"github.com/tfadeyi/go-aloe/internal/client"
	"github.com/tfadeyi/go-aloe/internal/client/local"
)

type (
	Aloe struct {
		source string
		client client.Client
		logger *log.Logger
	}

	Options struct {
		Source string
		Logger *log.Logger
	}
)

func newInstance(opts Options) *Aloe {
	cl := local.NewLazy(opts.Source)

	return &Aloe{
		source: opts.Source,
		logger: opts.Logger,
		client: cl,
	}
}

func Default() *Aloe {
	instance := newInstance(Options{
		Source: "default.aloe.yaml",
		Logger: log.Default(),
	})
	return instance
}

func WithOptions(opts Options) *Aloe {
	instance := newInstance(opts)
	return instance
}

func (instance *Aloe) ErrorWithContext(ctx context.Context, err error, code string) error {
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

func (instance *Aloe) Error(err error, code string) error {
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

func (instance *Aloe) log(msg string, keyVal ...any) {
	if instance.logger != nil {
		instance.logger.Printf(msg, keyVal)
	}
}
