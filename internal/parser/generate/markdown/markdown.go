package markdown

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/tfadeyi/errors/internal/parser/generate/helpers"
	"io"
	"path/filepath"
	"text/template"

	"github.com/juju/errors"
	"github.com/tfadeyi/errors/internal/logging"
	"github.com/tfadeyi/errors/pkg/api"
)

//go:embed templates/error.md.tmpl
var errorDefinitionMarkdownTmpl string

//go:embed templates/info.md.tmpl
var applicationInfoMarkdownTmpl string

type Generator struct {
	logger                      *logging.Logger
	output                      string
	writer                      io.Writer
	infoTmplFile, errorTmplFile string
}

// Options contains the configuration options available to the Generator
type Options struct {
	Logger                      *logging.Logger
	Writer                      io.Writer
	Output                      string
	InfoTmplFile, ErrorTmplFile string
}

func New(opts *Options) *Generator {
	// create default options, these will be overridden
	if opts == nil {
		opts = new(Options)
	}

	return &Generator{
		logger:        opts.Logger,
		output:        opts.Output,
		writer:        opts.Writer,
		infoTmplFile:  opts.InfoTmplFile,
		errorTmplFile: opts.ErrorTmplFile,
	}
}

func (g *Generator) Generate(ctx context.Context, specs map[string]any) error {
	return writeMarkdownSpecifications(g.writer, specs, g.output != "", g.output, g.infoTmplFile, g.errorTmplFile)
}

func writeMarkdownSpecifications(writer io.Writer, specs map[string]any, toFile bool, outputDirectory string, infoTmpl, errorTmpl string) error {
	for _, spec := range specs {
		foundSpec, ok := spec.(*api.Manifest)
		if !ok {
			return errors.New("found invalid application errors manifest")
		}

		var files map[string][]byte
		var err error

		if infoTmpl != "" && errorTmpl != "" {
			files, err = generateMarkdownWithCustomTemplates(foundSpec, outputDirectory, infoTmpl, errorTmpl)
			if err != nil {
				return err
			}
		} else {
			files, err = generateMarkdown(foundSpec, outputDirectory)
			if err != nil {
				return err
			}
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

func generateMarkdown(spec *api.Manifest, outputDir string) (map[string][]byte, error) {
	files := make(map[string][]byte)
	root := filepath.Join(outputDir, "index.md")
	// parse application general information
	tmpl, err := template.New(spec.Name).Parse(applicationInfoMarkdownTmpl)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, spec)
	if err != nil {
		return nil, err
	}
	if _, ok := files[root]; !ok {
		files[root] = bluemonday.UGCPolicy().SanitizeBytes(buf.Bytes())
	}

	for code, def := range spec.ErrorsDefinitions {
		tmpl, err := template.New(code).Parse(errorDefinitionMarkdownTmpl)
		if err != nil {
			return nil, err
		}
		buf := bytes.NewBuffer([]byte{})
		err = tmpl.Execute(buf, def)
		if err != nil {
			return nil, err
		}
		path := filepath.Join(outputDir, "errors", fmt.Sprintf("%s.md", code))
		if _, ok := files[path]; !ok {
			files[path] = bluemonday.UGCPolicy().SanitizeBytes(buf.Bytes())
		}
	}
	return files, nil
}

func generateMarkdownWithCustomTemplates(spec *api.Manifest, outputDir string, infoTmplFile, errorTmplFile string) (map[string][]byte, error) {
	files := make(map[string][]byte)
	root := filepath.Join(outputDir, "index.md")
	// parse application general information
	tmpl, err := template.New(spec.Name).ParseFiles(infoTmplFile)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, spec)
	if err != nil {
		return nil, err
	}
	if _, ok := files[root]; !ok {
		files[root] = bluemonday.UGCPolicy().SanitizeBytes(buf.Bytes())
	}

	for code, def := range spec.ErrorsDefinitions {
		tmpl, err := template.New(code).ParseFiles(errorTmplFile)
		if err != nil {
			return nil, err
		}
		buf := bytes.NewBuffer([]byte{})
		err = tmpl.Execute(buf, def)
		if err != nil {
			return nil, err
		}
		path := filepath.Join(outputDir, "errors", fmt.Sprintf("%s.md", code))
		if _, ok := files[path]; !ok {
			files[path] = bluemonday.UGCPolicy().SanitizeBytes(buf.Bytes())
		}
	}
	return files, nil
}
