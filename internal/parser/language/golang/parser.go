package golang

import (
	"context"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/juju/errors"
	"github.com/tfadeyi/errors/internal/logging"
	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"
)

type Parser struct {
	// specs contains references to all the service specifications that have been parsed
	specs map[string]api.Manifest
	// current references the current service specification being parsed
	current *api.Manifest
	// sourceFile is the path to the target file to be parsed, i.e: -f file.go
	sourceFile string
	// sourceContent is the reader to the content to be parsed
	sourceContent io.ReadCloser
	includedDirs  []string
	logger        *logging.Logger
	fset          *token.FileSet
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
		specs:         map[string]api.Manifest{},
		current:       nil,
		sourceFile:    sourceFile,
		sourceContent: sourceContent,
		includedDirs:  dirs,
		logger:        logger,
		fset:          token.NewFileSet(),
	}
}

// getAllGoPackages fetches all the available golang packages in the target directory and subdirectories and returns
// a map of *ast.Package(s) keyed on the package name
func getAllGoPackages(fset *token.FileSet, dir string, mode goparser.Mode) (map[string]*ast.Package, error) {
	pkgs, err := goparser.ParseDir(fset, dir, nil, mode)
	if err != nil {
		return map[string]*ast.Package{}, err
	}

	// walk through the directories and parse the not already found go packages
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			foundPkgs, err := goparser.ParseDir(fset, path, nil, mode)
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
func getFile(fset *token.FileSet, name string, file io.ReadCloser, mode goparser.Mode) (*ast.File, error) {
	if file != nil {
		defer func(file io.ReadCloser) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)
	}
	return goparser.ParseFile(fset, name, file, mode)
}

func (p *Parser) warn(err error, keyValues ...interface{}) {
	if p.logger != nil {
		p.logger.Warn(err, keyValues...)
	}
}

// Parse will parse the source code for fyi annotations.
// In case of error during parsing, Parse returns an empty sloth.Spec
func (p *Parser) ParseSource(ctx context.Context) (map[string]api.Manifest, error) {
	return p.parseManifest(ctx)
}

// AnnotateErrors will parse the source code and annotate all error definitions with fyi comments
func (p *Parser) AnnotateErrors(ctx context.Context, wrapErrors bool) error {
	return p.annotateErrorDefinitions(ctx, wrapErrors)
}

func (p *Parser) stats() {
	for _, manifest := range p.specs {
		p.logger.Info("Found", "application", manifest.Name, "errors", len(manifest.ErrorsDefinitions))
	}
}
