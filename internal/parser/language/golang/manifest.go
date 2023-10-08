package golang

import (
	"context"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"io"
	"k8s.io/utils/pointer"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
	"github.com/microcosm-cc/bluemonday"
	"github.com/tfadeyi/errors/internal/parser/grammar"
	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"
)

func (p *Parser) parseManifestComments(fset *token.FileSet, manifest api.Manifest, comments ...*ast.CommentGroup) (api.Manifest, error) {
	for _, comment := range comments {
		if !strings.HasPrefix(strings.TrimSpace(comment.Text()), "@fyi") {
			continue
		}

		commentString := strings.TrimSpace(comment.Text())

		p.logger.Debug("Parsing", "comment", commentString)
		// partialManifest contains the partially parsed sloth Specification for a given comment group
		// this means the parsed spec will only contain data for the fields that are present in the comments, making the spec only partially accurate
		partialManifest, err := grammar.Eval(commentString)
		if err != nil {
			if p.strict {
				return api.Manifest{}, err
			}
			p.logger.Debug("Could not parse comments", "error", err)
			continue
		}

		if manifest.Name == "" {
			manifest.Name = bluemonday.UGCPolicy().Sanitize(partialManifest.Name)
		}
		if manifest.Version == "" {
			manifest.Version = bluemonday.UGCPolicy().Sanitize(partialManifest.Version)
		}
		if manifest.BaseUrl == "" {
			manifest.BaseUrl = bluemonday.UGCPolicy().Sanitize(partialManifest.BaseUrl)
		}
		if manifest.Repository == "" {
			manifest.Repository = bluemonday.UGCPolicy().Sanitize(partialManifest.Repository)
		}
		if manifest.Description == nil && partialManifest.Description != nil {
			manifest.Description = pointer.String(bluemonday.UGCPolicy().Sanitize(*partialManifest.Description))
		}
		if manifest.Title == nil && partialManifest.Title != nil {
			manifest.Title = pointer.String(bluemonday.UGCPolicy().Sanitize(*partialManifest.Title))
		}

		for key, definition := range partialManifest.ErrorsDefinitions {
			definition.Meta = &api.ErrorMeta{Loc: &api.ErrorMetaLoc{
				Filename: getErrorLocationPath(manifest, fset.File(comment.Pos()).Name()),
				Line:     fset.Position(comment.Pos()).Line,
			}}
			// sanitize all fields
			definition.Short = bluemonday.UGCPolicy().Sanitize(definition.Short)
			definition.Code = bluemonday.UGCPolicy().Sanitize(definition.Code)
			definition.Title = bluemonday.UGCPolicy().Sanitize(definition.Title)
			// sanitize all suggestions
			for _, suggestion := range definition.Suggestions {
				suggestion.Short = bluemonday.UGCPolicy().Sanitize(suggestion.Short)
			}

			manifest.ErrorsDefinitions[key] = definition
		}
	}

	return manifest, nil
}

// parseSourceContentAndGenerateManifest will parse the source code for @fyi annotations and return an error manifest
func (p *Parser) parseSourceContentAndGenerateManifest(ctx context.Context, source io.Reader) (*api.Manifest, error) {
	if source == nil {
		return nil, errors.New("invalid source content")
	}

	var manifest = api.Manifest{
		BaseUrl:           "",
		Description:       nil,
		ErrorsDefinitions: api.ErrorDefinitions{},
		Name:              "",
		Title:             nil,
		Version:           "",
	}
	// collect all fyi annotations from the file and add them to the manifest struct

	file, err := getFile(p.fset, "", source, goparser.ParseComments)
	if err != nil {
		// error hard as we can't extract more data for the spec
		return nil, err
	}

	p.logger.Debug("Parsing source code", "file", file.Name)
	manifest, err = p.parseManifestComments(p.fset, manifest, file.Comments...)
	if err != nil {
		return nil, err
	}

	p.logger.Debug("Parsed source code", "file", file.Name)
	return &manifest, nil
}

// parseAllSourcesAndGenerateManifest will parse the source code for @fyi annotations.
func (p *Parser) parseAllSourcesAndGenerateManifest(ctx context.Context) ([]*api.Manifest, error) {
	var manifest = api.Manifest{
		BaseUrl:           "",
		Description:       nil,
		ErrorsDefinitions: api.ErrorDefinitions{},
		Name:              "",
		Title:             nil,
		Version:           "",
	}

	// collect all @fyi annotations from packages and add them to the manifest struct
	for _, pkg := range p.applicationPackages {
		// Prioritise parsing the main.go if present in the package
		for _, file := range pkg.Package.Syntax {
			// handle signals with context
			select {
			case <-ctx.Done():
				return nil, errors.New("termination signal was received, terminating process...")
			default:
			}

			filename := pkg.Fset.File(file.Pos()).Name()
			if strings.Contains(filename, "main.go") {
				p.logger.Debug("Parsing source code", "package", pkg.Name, "file", filename)

				var err error
				manifest, err = p.parseManifestComments(pkg.Fset, manifest, file.Comments...)
				if err != nil {
					return nil, err
				}

				p.logger.Debug("Parsed source code", "package", pkg.Name, "file", filename)
				break
			}
		}

		// parse the rest of the files, skipping main.go
		for _, file := range pkg.Package.Syntax {
			// handle signals with context
			select {
			case <-ctx.Done():
				return nil, errors.New("termination signal was received, terminating process...")
			default:
			}

			filename := pkg.Fset.File(file.Pos()).Name()
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

			var err error
			manifest, err = p.parseManifestComments(pkg.Fset, manifest, file.Comments...)
			if err != nil {
				return nil, err
			}

			p.logger.Debug("Parsed source code", "package", pkg.Name, "file", filename)
		}
	}

	return []*api.Manifest{&manifest}, nil
}

func getErrorLocationPath(manifest api.Manifest, path string) string {
	repo := manifest.Repository
	_, after, ok := strings.Cut(path, repo)
	if !ok {
		return path
	}
	return filepath.Join(repo, after)
}
