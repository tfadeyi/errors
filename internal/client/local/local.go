package local

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"errors"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/tfadeyi/go-aloe/internal/client"
	"github.com/tfadeyi/go-aloe/pkg/api"
	yaml "gopkg.in/yaml.v3"
)

type (
	Client struct {
		SpecFilename string
		SpecSource   []byte
		Spec         *api.Application
		Logger       *log.Logger
	}
)

var (
	ErrSpecNotExist      = errors.New("specification file doesn't exist")
	ErrFetchingSpec      = errors.New("couldn't fetch the specification file data")
	ErrUnsupportedFormat = errors.New("the specification is in an invalid format")
	ErrFailedParsingSpec = errors.New("could not parse the aloe specification file")
)

// New returns a new instance of the client.Client, it performs checks on the inputs and errors if the check fails
// filename is path to aloe specification file
// source is the in-memory aloe specification
func New(filename string, source []byte) (client.Client, error) {
	if filename != "" {
		_, err := os.Stat(filename)
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrSpecNotExist
		}

		source, err = os.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	}

	spec, err := decodeSpec(filename, source)
	if err != nil {
		return nil, err
	}

	if spec == nil || (spec.Name == "") {
		return nil, ErrFailedParsingSpec
	}

	return &Client{SpecFilename: filename, SpecSource: source, Spec: spec}, nil
}

func decodeSpec(filename string, buf []byte) (*api.Application, error) {
	var spec api.Application
	var err error

	if filename == "" {
		if err = toml.Unmarshal(buf, &spec); err != nil {
			if err = yaml.Unmarshal(buf, &spec); err != nil {
				return nil, fmt.Errorf("%s: [%w]", ErrFailedParsingSpec, err)
			}
		}
		return &spec, nil
	}

	switch filepath.Ext(filename) {
	case ".yaml":
		err = yaml.Unmarshal(buf, &spec)
	case ".json":
		err = json.Unmarshal(buf, &spec)
	case ".toml":
		err = toml.Unmarshal(buf, &spec)
	default:
		err = ErrUnsupportedFormat
	}
	return &spec, err
}

func NewLazy(filename string, source []byte) client.Client {
	return &Client{
		SpecFilename: filename,
		SpecSource:   source,
		Spec:         nil,
	}
}

func (l *Client) GenerateErrorMessageFromCode(ctx context.Context, code string) (string, error) {
	code = strings.TrimSpace(code)
	if l.Spec == nil {
		var err error
		if l.SpecFilename != "" {
			_, err = os.Stat(l.SpecFilename)
			if errors.Is(err, os.ErrNotExist) {
				return "", ErrSpecNotExist
			}
			l.SpecSource, err = os.ReadFile(l.SpecFilename)
			if err != nil {
				return "", err
			}
		}

		if l.Spec, err = decodeSpec(l.SpecFilename, l.SpecSource); err != nil {
			return "", err
		}
	}

	// build url
	v, ok := l.Spec.ErrorsDefinitions[code]
	if !ok {
		return "", errors.New("no error was not found in the error specification file")
	}

	// TODO validate inputs
	name := strings.TrimSpace(l.Spec.Name)
	baseURL := strings.TrimSpace(l.Spec.BaseUrl)
	summary := strings.TrimSpace(v.Summary)

	url := fmt.Sprintf("%s/%s/errors_definitions/%s", baseURL, name, code)
	result := fmt.Sprintf("%s\nfor additional info check %s", summary, url)
	return result, nil
}
