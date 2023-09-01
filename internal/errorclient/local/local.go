package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/tfadeyi/errors/internal/errorclient"
	"github.com/tfadeyi/errors/pkg/api"
	"gopkg.in/yaml.v3"
)

type (
	Client struct {
		errorclient.Options
		Spec *api.Manifest
	}
)

var (
	ErrSpecificationDoesNotExist = errors.New("specification file doesn't exist")
)

const (
	errorDefinitionPath = "errors"
)

func New(opts errorclient.Options) errorclient.Client {
	if opts.ErrorDefinitionPath == "" {
		opts.ErrorDefinitionPath = errorDefinitionPath
	}
	return &Client{
		Options: opts,
		Spec:    nil,
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

	name := strings.TrimSpace(l.Spec.Name)
	baseURL := strings.TrimSpace(l.Spec.BaseUrl)
	summary := strings.TrimSpace(v.Short)

	result := fmt.Sprintf("* %s.", summary)
	if l.ShowErrorURLs {
		url := fmt.Sprintf("%s/%s/%s/%s", baseURL, name, l.ErrorDefinitionPath, code)
		result = fmt.Sprintf("%s Additional information is available at %s", result, url)
	}
	return result, nil
}
