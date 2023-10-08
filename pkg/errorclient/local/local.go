package local

import (
	"context"
	"errors"
	"fmt"
	"github.com/tfadeyi/errors/pkg/errorclient"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"
	"gopkg.in/yaml.v3"
)

type (
	Client struct {
		*errorclient.ErrClientOptions
		Spec *api.Manifest
	}
)

var (
	ErrSpecificationDoesNotExist = errors.New("manifest file doesn't exist")
)

const (
	errorDefinitionPath = "errors"
)

func New(opts errorclient.ErrClientOptions) errorclient.ErrClient {
	if opts.ErrorDefinitionURLPath == "" {
		opts.ErrorDefinitionURLPath = errorDefinitionPath
	}
	return &Client{
		ErrClientOptions: &opts,
		Spec:             nil,
	}
}

func decodeSpec(buf []byte) (*api.Manifest, error) {
	var spec api.Manifest
	var err = yaml.Unmarshal(buf, &spec)
	return &spec, err
}

func (l *Client) GenerateErrorMessageFromCode(ctx context.Context, code string) (string, error) {
	select {
	case <-ctx.Done():
		return "", errors.New("termination signal was received, terminating process")
	default:
	}

	code = strings.TrimSpace(code)
	if l.Spec == nil {
		var err error
		if l.SourceFilename != "" && l.Source == nil {
			_, err = os.Stat(l.SourceFilename)
			if errors.Is(err, os.ErrNotExist) {
				return "", ErrSpecificationDoesNotExist
			}
			l.Source, err = os.ReadFile(l.SourceFilename)
			if err != nil {
				return "", err
			}
		}

		if l.Spec, err = decodeSpec(l.Source); err != nil {
			return "", err
		}
	}

	// build error urls
	v, ok := l.Spec.ErrorsDefinitions[code]
	if !ok {
		return "", errors.New("no error was not found in the error specification file")
	}

	if l.DisplayMarkdownErrors {
		return l.printMarkdownError(v)
	}

	return l.printTextError(v)
}

func (l *Client) printTextError(er api.Error) (string, error) {
	var content strings.Builder
	var url = fmt.Sprintf("%s/%s/%s/%s", l.Spec.BaseUrl, l.Spec.Name, l.ErrorDefinitionURLPath, er.Code)

	if er.Title != "" {
		if _, err := content.WriteString(fmt.Sprintf("%s\n", er.Title)); err != nil {
			return "", err
		}
	}

	if er.Short != "" && l.DisplayShortSummary {
		if _, err := content.WriteString("What caused the error\n"); err != nil {
			return "", err
		}
		if _, err := content.WriteString(strings.TrimSpace(er.Short)); err != nil {
			return "", err
		}
		if _, err := content.WriteString("\n"); err != nil {
			return "", err
		}
	}

	if l.DisplayErrorURL {
		if l.OverrideErrorURL != "" {
			url = l.OverrideErrorURL
		}
		if _, err := content.WriteString(fmt.Sprintf("\nAdditional information is available at: %s\n", url)); err != nil {
			return "", err
		}
	}

	if len(er.Suggestions) > 0 && l.NumberOfSuggestions > 0 {
		if _, err := content.WriteString("Quick Solutions\n"); err != nil {
			return "", err
		}
		count := 0
		for _, suggestion := range er.Suggestions {
			if count == l.NumberOfSuggestions {
				break
			}
			if suggestion.Short == "" {
				continue
			}
			if _, err := content.WriteString(fmt.Sprintf("* Suggestion: %s\n", suggestion.Short)); err != nil {
				return "", err
			}
		}
	}

	return content.String(), nil
}

func (l *Client) printMarkdownError(er api.Error) (string, error) {
	var content strings.Builder
	var url = fmt.Sprintf("%s/%s/%s/%s", l.Spec.BaseUrl, l.Spec.Name, l.ErrorDefinitionURLPath, er.Code)

	if er.Title != "" {
		if _, err := content.WriteString(fmt.Sprintf("# %s\n", strings.ToTitle(er.Title))); err != nil {
			return "", err
		}
	}

	if er.Short != "" && l.DisplayShortSummary {
		if _, err := content.WriteString("## What caused the error\n"); err != nil {
			return "", err
		}
		if _, err := content.WriteString(strings.TrimSpace(er.Short)); err != nil {
			return "", err
		}
		if _, err := content.WriteString("\n"); err != nil {
			return "", err
		}
	}

	if l.DisplayErrorURL {
		if l.OverrideErrorURL != "" {
			url = l.OverrideErrorURL
		}
		if _, err := content.WriteString(fmt.Sprintf("\n> Additional information is available at: %s\n", url)); err != nil {
			return "", err
		}
	}

	if len(er.Suggestions) > 0 && l.NumberOfSuggestions > 0 {
		if _, err := content.WriteString("## Quick Solutions\n"); err != nil {
			return "", err
		}
		count := 0
		for _, suggestion := range er.Suggestions {
			if count == l.NumberOfSuggestions {
				break
			}
			if suggestion.Short == "" {
				continue
			}
			if _, err := content.WriteString(fmt.Sprintf("* **Suggestion**: %s\n", suggestion.Short)); err != nil {
				return "", err
			}
		}
	}
	return glamour.Render(content.String(), "dark")
}
