package yaml

import (
	"bytes"
	"context"
	"github.com/juju/errors"
	"io"

	"github.com/tfadeyi/errors/internal/generate/helpers"
	"github.com/tfadeyi/errors/internal/logging"
	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"
	"gopkg.in/yaml.v3"
)

type Generator struct {
	logger *logging.Logger
	writer io.Writer
	output string
	header string
}

// Options contains the configuration options available to the Generator
type Options struct {
	Logger *logging.Logger
	Writer io.Writer
	Output string
	Header string
}

func New(opts *Options) *Generator {
	// create default options, these will be overridden
	if opts == nil {
		opts = new(Options)
	}
	return &Generator{
		logger: opts.Logger,
		writer: opts.Writer,
		output: opts.Output,
		header: opts.Header,
	}
}

func (g *Generator) GenerateManifests(ctx context.Context, specs map[string]api.Manifest) error {
	return writeYAML(ctx, g.writer, specs, g.output != "", g.output, g.header)
}

func writeYAML(ctx context.Context, writer io.Writer, specs map[string]api.Manifest, toFile bool, output, header string) error {
	for _, spec := range specs {
		// handle signals with context
		select {
		case <-ctx.Done():
			return errors.New("termination signal was received, terminating process...")
		default:
		}

		var files = make(map[string][]byte)

		body, err := yaml.Marshal(spec)
		if err != nil {
			return err
		}

		file := output
		files[file] = bytes.Join([][]byte{[]byte("---"), []byte(header), body}, []byte("\n"))
		if err := helpers.Clean(file); err != nil {
			return err
		}

		if toFile {
			if err := helpers.WriteToFile(files); err != nil {
				return err
			}
			continue
		}

		if err := helpers.Write(writer, files); err != nil {
			return err
		}
	}

	return nil
}
