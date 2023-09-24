package golang

import (
	"context"
	"go/ast"
	goparser "go/parser"
	"os"
	"strings"

	"github.com/juju/errors"
	"github.com/microcosm-cc/bluemonday"
	"github.com/tfadeyi/errors/internal/parser/grammar"
	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"
)

func (p *Parser) parseManifestComments(filename string, comments ...*ast.CommentGroup) error {
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

	p.logger.Debug("Current application being parsed", "application", p.current.Name)
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
			if p.current != nil && (p.current.Name == partialServiceSpec.Name ||
				p.current.Name == "") {
				p.specs[partialServiceSpec.Name] = *p.current
			}

			spec, ok := p.specs[partialServiceSpec.Name]
			if !ok {
				p.specs[partialServiceSpec.Name] = *partialServiceSpec
				p.current = partialServiceSpec
			} else {
				p.current = &spec
			}
		}

		if p.current.Name == "" {
			p.current.Name = partialServiceSpec.Name
		}
		if p.current.Version == "" {
			p.current.Version = partialServiceSpec.Version
		}
		if p.current.BaseUrl == "" {
			p.current.BaseUrl = partialServiceSpec.BaseUrl
		}
		if p.current.Description == nil {
			p.current.Description = partialServiceSpec.Description
		}
		if p.current.Title == nil {
			p.current.Title = partialServiceSpec.Title
		}

		for key, definition := range partialServiceSpec.ErrorsDefinitions {
			definition.Meta = &api.ErrorMeta{Loc: &api.ErrorMetaLoc{
				Filename: p.fset.File(comment.Pos()).Name(),
				Line:     p.fset.Position(comment.Pos()).Line,
			}}
			p.current.ErrorsDefinitions[key] = definition
		}

	}
	return nil
}

// parseManifest will parse the source code for @fyi annotations.
func (p *Parser) parseManifest(ctx context.Context) (map[string]api.Manifest, error) {
	// collect all fyi annotations from the file and add them to the spec struct
	if (p.sourceFile != "" && p.sourceFile != "-") || p.sourceContent != nil {
		file, err := getFile(p.fset, p.sourceFile, p.sourceContent, goparser.ParseComments)
		if err != nil {
			// error hard as we can't extract more data for the spec
			return nil, err
		}

		p.logger.Debug("Parsing source code", "file", file.Name)
		if err := p.parseManifestComments("", file.Comments...); err != nil {
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
		foundPkgs, err := getAllGoPackages(p.fset, dir, goparser.ParseComments)
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

	// collect all @fyi annotations from packages and add them to the spec struct
	for _, pkg := range applicationPackages {
		// Prioritise parsing the main.go if present in the package
		for filename, file := range pkg.Files {
			if strings.Contains(filename, "main.go") {
				p.logger.Debug("Parsing source code", "package", pkg.Name, "file", filename)

				if err := p.parseManifestComments(filename, file.Comments...); err != nil {
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

			if err := p.parseManifestComments(filename, file.Comments...); err != nil {
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
