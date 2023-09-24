package manifest

import (
	"context"
	"io"

	"github.com/juju/errors"
	api "github.com/tfadeyi/errors/pkg/api/v0.1.0"
	"gopkg.in/yaml.v3"
)

type (
	// Target is an abstraction for content generators, i.e: yaml and markdown.
	Target interface {
		// GenerateManifests transcribes the content from the specs to the writer
		GenerateManifests(ctx context.Context, manifests map[string]api.Manifest) error
	}

	Generator struct {
		Opts *Options
	}
)

var (
	ErrNoContentGenerator = errors.New("no target content generator was set")
)

// New creates a new instance of the Generator. See options.Option for more info on the available configuration.
func New(opts ...Option) *Generator {
	defaultOpts := new(Options)
	for _, opt := range opts {
		opt(defaultOpts)
	}
	for _, opt := range opts {
		opt(defaultOpts)
	}

	return &Generator{defaultOpts}
}

func (g *Generator) GenerateManifests(ctx context.Context, manifests map[string]api.Manifest) error {
	if g.Opts.TargetManifestGenerator == nil {
		return ErrNoContentGenerator
	}
	return g.Opts.TargetManifestGenerator.GenerateManifests(ctx, manifests)
}

func ValidateFromReader(r io.Reader) (*api.Manifest, error) {
	var manifest api.Manifest
	err := yaml.NewDecoder(r).Decode(&manifest)
	return &manifest, err
}
