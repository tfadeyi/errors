package grammar

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"
)

type (
	// Grammar is the participle grammar use to parse the Sloth comment groups in source files
	Grammar struct {
		// Stmts is a list of Sloth grammar Statements
		Stmts []*Statement `@@*`
	}
	// Statement is any comment starting with @sloth keyword
	Statement struct {
		Scope Scope  `@@`
		Value string `Whitespace* @(String (Whitespace|EOL)*)+`
	}
	// Scope defines the statement scope, similar to a code function
	Scope struct {
		// Type is the specification struct a statement refers to
		Type  string `(Fyi @((".error" (".suggestion")? ))?)`
		Value string `Whitespace* @("name"|"id"|"title"|"description"|"base_url"|"version"|"short"|"long"|"code"|"error_code"|"severity")`
	}
)

// GetType returns the type of the statement scope
func (k Scope) GetType() string {
	return k.Type
}

func parseAndAssignStructFields(attr string, value string, fields []reflect.StructField, pValue reflect.Value) error {
	for _, field := range fields {
		tag, ok := field.Tag.Lookup("yaml")
		if !ok {
			return nil
		}
		key := strings.Split(tag, ",")[0]
		if attr == key {
			// set field value
			v := pValue.FieldByName(field.Name)
			if v.IsValid() {
				if v.CanSet() {
					switch v.Kind() {
					case reflect.Pointer:
						if field.Name == "Severity" {
							severity := api.ErrorSeverity(value)
							v.Set(reflect.ValueOf(&severity))
						} else {
							v.Set(reflect.ValueOf(&value))
						}
					case reflect.Bool:
						b, err := strconv.ParseBool(value)
						if err != nil {
							return err
						}
						v.SetBool(b)
					case reflect.Float64:
						f, err := strconv.ParseFloat(value, 64)
						if err != nil {
							return err
						}
						v.SetFloat(f)
					case reflect.Map:
						// label or annotation
						m := strings.Split(value, " ")
						v.SetMapIndex(reflect.ValueOf(m[0]), reflect.ValueOf(m[1]))
					default:
						v.Set(reflect.ValueOf(value))
					}
				}
			}
		}
	}
	return nil
}

func (g Grammar) parse() (*api.Manifest, error) {
	var spec = &api.Manifest{
		BaseUrl:           "",
		Description:       nil,
		ErrorsDefinitions: api.ErrorDefinitions{},
		Name:              "",
		Title:             nil,
		Version:           "",
	}

	lowSeverity := api.ErrorSeverityLow
	var foundErr = &api.Error{
		Code:        "",
		Long:        nil,
		Meta:        &api.ErrorMeta{Loc: nil},
		Short:       "",
		Title:       "",
		Severity:    &lowSeverity,
		Suggestions: api.Suggestions{},
	}
	var foundSolution = &api.Suggestion{
		Id:     "",
		Short:  "",
		DocRef: nil,
	}

	for _, attr := range g.Stmts {
		switch attr.Scope.GetType() {
		case ".error.suggestion":
			fields := reflect.VisibleFields(reflect.TypeOf(*foundSolution))
			pValue := reflect.ValueOf(foundSolution).Elem()
			if err := parseAndAssignStructFields(strings.ToLower(attr.Scope.Value), strings.TrimSpace(attr.Value), fields, pValue); err != nil {
				continue
			}
			foundSolution.ErrorCode = foundErr.Code
			foundSolution.Id = fmt.Sprintf("%d", len(foundErr.Suggestions)+1)
			foundErr.Suggestions[foundSolution.Id] = *foundSolution
		case ".error":
			fields := reflect.VisibleFields(reflect.TypeOf(*foundErr))
			pValue := reflect.ValueOf(foundErr).Elem()
			if err := parseAndAssignStructFields(strings.ToLower(attr.Scope.Value), strings.TrimSpace(attr.Value), fields, pValue); err != nil {
				continue
			}
			spec.ErrorsDefinitions[foundErr.Code] = *foundErr
		default:
			fields := reflect.VisibleFields(reflect.TypeOf(*spec))
			pValue := reflect.ValueOf(spec).Elem()
			if err := parseAndAssignStructFields(strings.ToLower(attr.Scope.Value), strings.TrimSpace(attr.Value), fields, pValue); err != nil {
				continue
			}
		}
	}

	return spec, nil
}

func createGrammar(filename, source string, options ...participle.ParseOption) (*Grammar, error) {
	ast, err := participle.Build[Grammar](
		participle.Lexer(lexerDefinition),
	)
	if err != nil {
		return nil, err
	}

	return ast.ParseString(filename, source, options...)
}

// Eval evaluates the source input against the grammar and returns an instance of *sloth.spec
func Eval(source string, options ...participle.ParseOption) (*api.Manifest, error) {
	grammar, err := createGrammar("", source, options...)
	if err != nil {
		return nil, err
	}

	return grammar.parse()
}
