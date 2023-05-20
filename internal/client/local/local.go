package local

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"errors"
	api "github.com/tfadeyi/aloe-bindings/src/go/src"
	"github.com/tfadeyi/go-aloe/internal/client"
	"gopkg.in/yaml.v3"
)

type (
	Client struct {
		SpecSource string
		Spec       *api.Spec
		Logger     *log.Logger
	}
)

var (
	ErrSpecNotExist      = errors.New("specification file doesn't exist")
	ErrFetchingSpec      = errors.New("couldn't fetch the specification file data")
	ErrUnsupportedFormat = errors.New("the specification is in an invalid format")
)

func New(source string) (client.Client, error) {
	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		return nil, ErrSpecNotExist
	}
	buf, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}
	var spec api.Spec
	err = yaml.Unmarshal(buf, &spec)
	if err != nil {
		return nil, err
	}
	return &Client{SpecSource: source, Spec: &spec}, nil
}

func NewLazy(source string) client.Client {
	return &Client{
		SpecSource: source,
		Spec:       nil,
	}
}

func (l *Client) GenerateErrorMessageFromCode(ctx context.Context, code string) (string, error) {
	if l.Spec == nil {
		if _, err := os.Stat(l.SpecSource); errors.Is(err, os.ErrNotExist) {
			return "", ErrSpecNotExist
		}
		buf, err := os.ReadFile(l.SpecSource)
		if err != nil {
			return "", err
		}

		switch filepath.Ext(l.SpecSource) {
		case ".yaml":
			err = yaml.Unmarshal(buf, &l.Spec)
		case ".json":
			err = json.Unmarshal(buf, &l.Spec)
		default:

		}
		if err != nil {
			return "", err
		}
	}

	// build url
	v, ok := l.Spec.ErrorsDefinitions[code]
	if !ok {
		return "", errors.New("no error was not found in the error specification file")
	}
	url := fmt.Sprintf("%s/%s/%s", l.Spec.BaseURL, l.Spec.Name, code)
	result := fmt.Sprintf("%s\nfor additional info check %s.", v.Summary, url)
	return result, nil
}
