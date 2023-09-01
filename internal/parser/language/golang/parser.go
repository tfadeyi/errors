package golang

import (
	"context"
	"github.com/microcosm-cc/bluemonday"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
	"github.com/tfadeyi/errors/internal/logging"
	"github.com/tfadeyi/errors/internal/parser/grammar"
	"github.com/tfadeyi/errors/pkg/api"
)

type Parser struct {
	// specs contains references to all the service specifications that have been parsed
	specs map[string]any
	// current references the current service specification being parsed
	current any
	// sourceFile is the path to the target file to be parsed, i.e: -f file.go
	sourceFile string
	// sourceContent is the reader to the content to be parsed
	sourceContent io.ReadCloser
	includedDirs  []string
	logger        *logging.Logger
}

// Options contains the configuration options available to the Parser
type Options struct {
	Logger *logging.Logger
	// SourceFile is the path to the target file to be parsed, i.e: -f file.go
	SourceFile string
	// SourceContent is the reader to the content to be parsed
	SourceContent    io.ReadCloser
	InputDirectories []string
}

// NewParser client Parser performs all checks at initialization time
func NewParser(opts *Options) *Parser {
	// create default options, these will be overridden
	if opts == nil {
		opts = new(Options)
	}

	logger := opts.Logger
	dirs := opts.InputDirectories
	sourceFile := opts.SourceFile
	sourceContent := opts.SourceContent

	return &Parser{
		specs:         map[string]any{},
		current:       nil,
		sourceFile:    sourceFile,
		sourceContent: sourceContent,
		includedDirs:  dirs,
		logger:        logger,
	}
}

// getAllGoPackages fetches all the available golang packages in the target directory and subdirectories
func getAllGoPackages(dir string) (map[string]*ast.Package, error) {
	fset := token.NewFileSet()
	pkgs, err := goparser.ParseDir(fset, dir, nil, goparser.ParseComments)
	if err != nil {
		return map[string]*ast.Package{}, err
	}

	// walk through the directories and parse the not already found go packages
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			foundPkgs, err := goparser.ParseDir(fset, path, nil, goparser.ParseComments)
			if err != nil {
				return err
			}
			for pkgName, pkg := range foundPkgs {
				if _, ok := pkgs[pkgName]; !ok {
					pkgs[pkgName] = pkg
				}
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	if len(pkgs) == 0 {
		return nil, errors.Errorf("no go packages were found in the target directory and subdirectories: %s", dir)
	}

	return pkgs, nil
}

// getFile returns the ast go file struct given filename or an io.Reader. If an io.Reader is passed it will take precedence
// over the filename
func getFile(name string, file io.ReadCloser) (*ast.File, error) {
	fset := token.NewFileSet()
	if file != nil {
		defer func(file io.ReadCloser) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)
	}
	return goparser.ParseFile(fset, name, file, goparser.ParseComments)
}

func (p *Parser) parseErrorAnnotations(filename string, comments ...*ast.CommentGroup) error {
	if p.current == nil {
		p.current = &api.Manifest{
			BaseUrl:           "",
			Description:       nil,
			ErrorsDefinitions: api.ErrorDefinitions{},
			Name:              "",
			Title:             nil,
			Version:           "",
		}
	}

	p.logger.Debug("Current application being parsed", "application", p.current.(*api.Manifest).Name)
	for _, comment := range comments {
		if !strings.HasPrefix(strings.TrimSpace(comment.Text()), "@fyi") {
			continue
		}

		commentString := bluemonday.UGCPolicy().Sanitize(strings.TrimSpace(comment.Text()))

		p.logger.Debug("Parsing", "comment", commentString)
		// partialServiceSpec contains the partially parsed sloth Specification for a given comment group
		// this means the parsed spec will only contain data for the fields that are present in the comments, making the spec only partially accurate
		partialServiceSpec, err := grammar.Eval(commentString)
		if err != nil {
			p.warn(err)
			continue
		}

		// if the comment group contains a reference to the service name
		// check if service was parsed before else add it the collection of specs.
		// Set the found service spec as the current service spec.
		if partialServiceSpec.Name != "" {
			if p.current != nil && (p.current.(*api.Manifest).Name == partialServiceSpec.Name ||
				p.current.(*api.Manifest).Name == "") {
				p.specs[partialServiceSpec.Name] = p.current
			}

			spec, ok := p.specs[partialServiceSpec.Name]
			if !ok {
				p.specs[partialServiceSpec.Name] = partialServiceSpec
				p.current = partialServiceSpec
			} else {
				p.current = spec.(*api.Manifest)
			}
		}

		if p.current.(*api.Manifest).Name == "" {
			p.current.(*api.Manifest).Name = partialServiceSpec.Name
		}
		if p.current.(*api.Manifest).Version == "" {
			p.current.(*api.Manifest).Version = partialServiceSpec.Version
		}
		if p.current.(*api.Manifest).BaseUrl == "" {
			p.current.(*api.Manifest).BaseUrl = partialServiceSpec.BaseUrl
		}
		if p.current.(*api.Manifest).Description == nil {
			p.current.(*api.Manifest).Description = partialServiceSpec.Description
		}
		if p.current.(*api.Manifest).Title == nil {
			p.current.(*api.Manifest).Title = partialServiceSpec.Title
		}

		for key, definition := range partialServiceSpec.ErrorsDefinitions {
			definition.Meta = &api.ErrorMeta{Loc: &api.ErrorMetaLoc{
				Path: filename,
			}}
			p.current.(*api.Manifest).ErrorsDefinitions[key] = definition
		}
	}
	return nil
}

func (p *Parser) warn(err error, keyValues ...interface{}) {
	if p.logger != nil {
		p.logger.Warn(err, keyValues...)
	}
}

// Parse will parse the source code for sloth annotations.
// In case of error during parsing, Parse returns an empty sloth.Spec
func (p *Parser) Parse(ctx context.Context) (map[string]any, error) {
	// collect all sloth annotations from the file and add them to the spec struct
	if p.sourceFile != "" || p.sourceContent != nil {
		file, err := getFile(p.sourceFile, p.sourceContent)
		if err != nil {
			// error hard as we can't extract more data for the spec
			return nil, err
		}

		p.logger.Debug("Parsing source code", "file", file.Name)
		if err := p.parseErrorAnnotations("", file.Comments...); err != nil {
			return nil, err
		}

		p.logger.Debug("Parsed source code", "file", file.Name)
		return p.specs, nil
	}

	applicationPackages := map[string]*ast.Package{}
	for _, dir := range p.includedDirs {
		// handle signals with context
		select {
		case <-ctx.Done():
			return nil, errors.New("termination signal was received, terminating process...")
		default:
		}

		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			// skip if dir doesn't exists
			p.warn(err)
			continue
		}
		foundPkgs, err := getAllGoPackages(dir)
		if err != nil {
			p.warn(err)
			continue
		}

		for pkgName, pkg := range foundPkgs {
			if _, ok := applicationPackages[pkgName]; !ok {
				applicationPackages[pkgName] = pkg
			}
		}
	}

	// collect all sloth annotations from packages and add them to the spec struct
	for _, pkg := range applicationPackages {
		// Prioritise parsing the main.go if present in the package
		for filename, file := range pkg.Files {
			if strings.Contains(filename, "main.go") {
				p.logger.Debug("Parsing source code", "package", pkg.Name, "file", filename)

				if err := p.parseErrorAnnotations(filename, file.Comments...); err != nil {
					p.warn(err)
					break
				}

				p.logger.Debug("Parsed source code", "package", pkg.Name, "file", filename)
				break
			}
		}

		// parse the rest of the files, skipping main.go
		for filename, file := range pkg.Files {
			if strings.Contains(filename, "main.go") {
				continue
			}
			p.logger.Debug("Parsing source code", "package", pkg.Name, "file", filename)
			// handle signals with context
			select {
			case <-ctx.Done():
				return nil, errors.New("termination signal was received, terminating process...")
			default:
			}

			if err := p.parseErrorAnnotations(filename, file.Comments...); err != nil {
				p.warn(err)
				continue
			}

			p.logger.Debug("Parsed source code", "package", pkg.Name, "file", filename)
		}
	}

	// print statistics
	p.stats()

	return p.specs, nil
}

func (p *Parser) stats() {
	for _, spec := range p.specs {
		s := spec.(*api.Manifest)
		p.logger.Info("Found", "application", s.Name, "errors", len(s.ErrorsDefinitions))
	}
}
