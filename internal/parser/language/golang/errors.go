package golang

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/guess"
	"github.com/juju/errors"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	ErrorFyiGoClientLibrary = "github.com/tfadeyi/errors"
)

// annotateErrors will parse the source code and annotate all error definitions with fyi comments
func (p *Parser) annotateAllErrors(ctx context.Context, wrapErrors, annotateOnlyTodos bool) error {
	for _, pkg := range p.applicationPackages {
		// parse the rest of the files, skipping main.go
		for _, file := range pkg.Package.Syntax {
			filename := pkg.Fset.File(file.Pos()).Name()
			p.logger.Debug("Parsing source code", "package", pkg.Name, "file", filename)
			// handle signals with context
			select {
			case <-ctx.Done():
				return errors.New("termination signal was received, terminating process...")
			default:
			}
			var err error

			out, err := annotateFile(file, pkg, wrapErrors, annotateOnlyTodos)
			if err != nil {
				p.warn(err)
				continue
			}
			writer, err := os.OpenFile(filename, os.O_RDWR, 0644)
			if err != nil {
				p.warn(err)
				continue
			}

			if _, err := io.Copy(writer, out); err != nil {
				p.warn(err)
				continue
			}
			p.logger.Debug("Parsed source code", "package", pkg.Name, "file", filename)
		}
	}
	return nil
}

// annotateSourceFileErrors will parse the source code and annotate all error definitions with fyi comments
func (p *Parser) annotateSourceFileErrors(ctx context.Context, sourceFile string, wrapErrors, annotateOnlyTodos bool) error {
	if sourceFile == "" {
		// TODO more user friendly error
		return errors.New("invalid source file")
	}

	p.logger.Debug("Parsing source code", "file", sourceFile)
	pkgs, err := decorator.Load(&packages.Config{Dir: filepath.Dir(sourceFile), Context: ctx,
		Mode: packages.NeedFiles | packages.NeedImports | packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo})
	if err != nil {
		p.logger.Debug("Failed to load package information", "package", sourceFile, "error", err)
		// TODO more user friendly error
		return err
	}

	pkg := pkgs[0]
	var astFile *ast.File
	for _, f := range pkg.Package.Syntax {
		fpath := pkg.Fset.File(f.Pos()).Name()
		if filepath.Base(fpath) == filepath.Base(sourceFile) {
			astFile = f
		}
	}

	out, err := annotateFile(astFile, pkg, wrapErrors, annotateOnlyTodos)
	if err != nil {
		p.logger.Debug("Failed to annotate file", "file", sourceFile, "error", err)
		// TODO more user friendly error
		return err
	}

	writer, err := os.OpenFile(sourceFile, os.O_RDWR, 0644)
	if err != nil {
		p.logger.Debug("Failed to open to file", "file", sourceFile, "error", err)
		// TODO more user friendly error
		return err
	}

	if _, err := io.Copy(writer, out); err != nil {
		p.logger.Debug("Failed to write to file", "file", sourceFile, "error", err)
		// TODO more user friendly error
		return err
	}

	p.logger.Debug("Parsed source code", "file", sourceFile)
	return nil
}

// annotateFile annotates with @fyi annotation, errors and main functions if they are present in the ast file.
func annotateFile(astFile *ast.File, pkg *decorator.Package, shouldWrap, annotateOnlyTodos bool) (io.Reader, error) {
	dec := pkg.Decorator
	decorateFile, err := dec.DecorateFile(astFile)
	if err != nil {
		return nil, err
	}

	ast.Inspect(astFile, func(node ast.Node) bool {
		switch node.(type) {
		case *ast.FuncDecl:
			v := node.(*ast.FuncDecl)
			funcName := v.Name.String()
			if funcName == "main" && pkg.Name == "main" {
				dec.Dst.Nodes[v].Decorations().Before = dst.NewLine
				if !strings.Contains(strings.Join(dec.Dst.Nodes[node].Decorations().Start.All(), ","), "@fyi") {
					dec.Dst.Nodes[v].Decorations().Start.Append(
						"// @fyi name <CHANGE ME>",
						"// @fyi title <CHANGE ME>",
						"// @fyi base_url <CHANGE ME>",
						"// @fyi version v0.1.0",
						"// @fyi description <CHANGE ME>",
					)
				}
			}
		case *ast.ReturnStmt:
			ast.Inspect(node, func(returnDecl ast.Node) bool {
				if err := annotateReturnStmt(astFile.Name.String(), node, returnDecl, pkg, shouldWrap, annotateOnlyTodos); err != nil {
					return true
				}
				return false
			})
		case *ast.GenDecl:
			ast.Inspect(node, func(genDecl ast.Node) bool {
				if err := annotateGenericDecl(astFile.Name.String(), node, genDecl, pkg, shouldWrap, annotateOnlyTodos); err != nil {
					return true
				}
				return false
			})
		default:
			return true
		}

		return true
	})

	output := bytes.NewBuffer([]byte{})
	// creating a new restorer instance plus creating an alias for the error.fyi client
	res := decorator.NewRestorerWithImports(astFile.Name.String(), guess.New())
	fr := res.FileRestorer()
	fr.Alias[ErrorFyiGoClientLibrary] = "fyi"
	if err := fr.Fprint(output, decorateFile); err != nil {
		return nil, err
	}

	return output, nil
}

// isError returns true if ast.Node is an error
func isError(v ast.Expr, info *types.Info) bool {
	if n, ok := info.TypeOf(v).(*types.Named); ok {
		o := n.Obj()
		return o != nil && o.Pkg() == nil && o.Name() == "error"
	}
	return false
}

// annotateReturnStmt annotates the errors on a return statement node
func annotateReturnStmt(filename string, parentNode ast.Node, errorNode ast.Node, pkg *decorator.Package,
	shouldWrap, annotateOnlyTodos bool) error {
	v, ok := errorNode.(ast.Expr)
	if !ok {
		return errors.New("node is not return statement")
	}

	info := pkg.TypesInfo
	fset := pkg.Fset
	dec := pkg.Decorator

	if isError(v, info) {
		errName := fmt.Sprintf("%s_error_%d", filename, fset.Position(v.Pos()).Line)
		// check that no @fyi is already present in the comments before updating the comments
		comments := strings.Join(dec.Dst.Nodes[parentNode].Decorations().Start.All(), ",")
		annotate := func() {
			// Wrap the error with fyi.Error
			if shouldWrap {
				switch v.(type) {
				case *ast.CallExpr:
					// store the current error definition, so it can be added as an argument
					currentErrorFunc := dst.Clone(dec.Dst.Nodes[v]).(*dst.CallExpr)
					//dst.Print(dec.Dst.Nodes[v])
					dec.Dst.Nodes[v].(*dst.CallExpr).Fun.(*dst.Ident).Name = "Error"
					dec.Dst.Nodes[v].(*dst.CallExpr).Fun.(*dst.Ident).Path = ErrorFyiGoClientLibrary
					dec.Dst.Nodes[v].(*dst.CallExpr).Args = []dst.Expr{currentErrorFunc, &dst.BasicLit{Value: fmt.Sprintf("%q", errName), Kind: token.STRING}}
				case *ast.Ident:
					// store the current error definition, so it can be added as an argument
					currentError := dst.Clone(dec.Dst.Nodes[v]).(*dst.Ident)
					dec.Dst.Nodes[v].(*dst.Ident).Name = fmt.Sprintf("Error(%s, %q)", currentError.Name, errName)
					dec.Dst.Nodes[v].(*dst.Ident).Path = ErrorFyiGoClientLibrary
				}
			}

			// regardless of the node (*ast.CallExpr or *ast.Ident) we should annotate it
			dec.Dst.Nodes[parentNode].Decorations().Before = dst.NewLine
			dec.Dst.Nodes[parentNode].Decorations().Start.Append(
				fmt.Sprintf("// @fyi.error code %s", errName),
				"// @fyi.error title <CHANGE ME TO A NICE TITLE>",
				"// @fyi.error short <TL;DR ON THE CAUSE OF THE ERROR>",
				"// @fyi.error severity severe|medium|low",
				"// @fyi.error.suggestion short <TL;DR ON HOW CAN THE ERROR BE FIXED?>",
			)
		}

		if !strings.Contains(comments, "@fyi") {
			if annotateOnlyTodos {
				if strings.Contains(comments, "TODO") {
					annotate()
				}
			} else {
				annotate()
			}
		}

	}
	return nil
}

// annotateGenericDecl annotates the errors on a genericDecl node
func annotateGenericDecl(filename string, parentNode ast.Node, errorNode ast.Node, pkg *decorator.Package,
	shouldWrap, annotateOnlyTodos bool) error {
	v, ok := errorNode.(ast.Expr)
	if !ok {
		return errors.New("node is not generic declaration")
	}

	info := pkg.TypesInfo
	fset := pkg.Fset
	dec := pkg.Decorator

	// check if the node is an error if so, add comment annotations and wrap error
	if isError(v, info) {
		errName := fmt.Sprintf("%s_error_%d", filename, fset.Position(v.Pos()).Line)
		// In case of generic declarations we always expect the node to be ast.CallExpr so no other types need to be checked
		if _, ok := v.(*ast.CallExpr); !ok {
			return errors.New("invalid error node")
		}

		comments := strings.Join(dec.Dst.Nodes[parentNode].Decorations().Start.All(), ",")
		annotate := func() {
			dec.Dst.Nodes[parentNode].Decorations().Before = dst.NewLine
			// Wrap the error with fyi.Error
			if shouldWrap {
				// store the current error definition, so it can be added as an argument
				currentErrorFunc := dst.Clone(dec.Dst.Nodes[v]).(*dst.CallExpr)
				dec.Dst.Nodes[v].(*dst.CallExpr).Fun.(*dst.Ident).Name = "Error"
				dec.Dst.Nodes[v].(*dst.CallExpr).Fun.(*dst.Ident).Path = ErrorFyiGoClientLibrary
				dec.Dst.Nodes[v].(*dst.CallExpr).Args = []dst.Expr{currentErrorFunc, &dst.BasicLit{Value: fmt.Sprintf("%q", errName), Kind: token.STRING}}
			}

			dec.Dst.Nodes[parentNode].Decorations().Start.Append(
				fmt.Sprintf("// @fyi.error code %s", errName),
				"// @fyi.error title <CHANGE ME TO A NICE TITLE>",
				"// @fyi.error short <TL;DR ON THE CAUSE OF THE ERROR>",
				"// @fyi.error severity severe|medium|low ",
				"// @fyi.error.suggestion short <TL;DR ON HOW CAN THE ERROR BE FIXED?>",
			)
		}

		if !strings.Contains(comments, "@fyi") {
			if annotateOnlyTodos {
				if strings.Contains(comments, "TODO") {
					annotate()
				}
			} else {
				annotate()
			}
		}
	}
	return nil
}