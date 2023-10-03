package golang

import (
	"context"
	"github.com/dave/dst/decorator"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/juju/errors"
	"github.com/tfadeyi/errors/internal/logging"
	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"
)

type Parser struct {
	includedDirs        []string
	logger              *logging.Logger
	fset                *token.FileSet
	strict              bool
	applicationPackages map[string]*decorator.Package
}

// Options contains the configuration options available to the Parser
type Options struct {
	Ctx              context.Context
	Logger           *logging.Logger
	InputDirectories []string
	Strict           bool
	Recursive        bool
}

// NewParser client Parser performs all checks at initialization time
func NewParser(opts *Options) *Parser {
	// create default options, these will be overridden
	if opts == nil {
		opts = new(Options)
	}

	logger := opts.Logger
	dirs := opts.InputDirectories
	ctx := opts.Ctx
	recursive := opts.Recursive

	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		log := logging.LoggerFromContext(ctx)
		logger = &log
	}

	applicationPackages := map[string]*decorator.Package{}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			continue
		}
		pkgs, err := decorator.Load(&packages.Config{Dir: dir, Context: ctx,
			Mode: packages.NeedFiles | packages.NeedImports | packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo})
		if err != nil {
			logger.Debug("Failed to load application packages", "dir", dir, "error", err)
			continue
		}

		for _, pkg := range pkgs {
			if _, ok := applicationPackages[pkg.Name]; !ok {
				applicationPackages[pkg.Name] = pkg
			}
		}

		if recursive {
			// walk through the directories and parse the not already found go packages
			err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() {
					pkgs, err = decorator.Load(&packages.Config{Dir: path, Context: ctx,
						Mode: packages.NeedFiles | packages.NeedImports | packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo})
					if err != nil {
						return err
					}

					for _, pkg := range pkgs {
						if _, ok := applicationPackages[pkg.Name]; !ok {
							applicationPackages[pkg.Name] = pkg
						}
					}
				}
				return err
			})
			if err != nil {
				logger.Debug("Failed to load application packages in subdirectories", "dir", dir, "error", err)
				continue
			}
		}
	}

	return &Parser{
		includedDirs:        dirs,
		logger:              logger,
		fset:                token.NewFileSet(),
		strict:              opts.Strict,
		applicationPackages: applicationPackages,
	}
}

// getFile returns the ast go file struct given filename or an io.Reader. If an io.Reader is passed it will take precedence
// over the filename
func getFile(fset *token.FileSet, name string, file io.Reader, mode goparser.Mode) (*ast.File, error) {
	return goparser.ParseFile(fset, name, file, mode)
}

func (p *Parser) warn(err error, keyValues ...interface{}) {
	if p.logger != nil {
		p.logger.Warn(err, keyValues...)
	}
}

func (p *Parser) ParseSource(ctx context.Context, content io.Reader) (*api.Manifest, error) {
	return p.parseSourceContentAndGenerateManifest(ctx, content)
}

func (p *Parser) AnnotateAllErrors(ctx context.Context, wrapErrors, annotateOnlyTodos bool) error {
	return p.annotateAllErrors(ctx, wrapErrors, annotateOnlyTodos)
}

func (p *Parser) ParseAllSources(ctx context.Context) ([]*api.Manifest, error) {
	return p.parseAllSourcesAndGenerateManifest(ctx)
}

func (p *Parser) AnnotateSourceErrors(ctx context.Context, filename string, wrapErrors, annotateOnlyTodos bool) error {
	return p.annotateSourceFileErrors(ctx, filename, wrapErrors, annotateOnlyTodos)
}
