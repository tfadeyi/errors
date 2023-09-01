package yaml

import (
	"bytes"
	"context"
	"github.com/tfadeyi/errors/internal/logging"
	"github.com/tfadeyi/errors/internal/parser/generate/helpers"
	"gopkg.in/yaml.v3"
	"io"
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

func (g *Generator) Generate(ctx context.Context, specs map[string]any) error {
	return writeYAMLSpecifications(g.writer, specs, g.output != "", g.output, g.header)
}

func writeYAMLSpecifications(writer io.Writer, specs map[string]any, toFile bool, output, header string) error {
	for _, spec := range specs {
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
