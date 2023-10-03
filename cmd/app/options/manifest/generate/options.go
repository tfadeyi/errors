package generate

import (
	"github.com/tfadeyi/errors/internal/generate/helpers"
	"os"
	"strings"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	errhandler "github.com/tfadeyi/errors"
	"github.com/tfadeyi/errors/cmd/app/options/manifest"
	"github.com/tfadeyi/errors/internal/generate"
	"github.com/tfadeyi/errors/internal/parser/language"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		Format       string
		IncludedDirs []string
		Output       string
		Language     string
		Recursive    bool
		*manifest.Options
	}
)

// New creates a new instance of the application's options
func New(common *manifest.Options) *Options {
	opts := new(Options)
	opts.Options = common
	return opts
}

// Prepare assigns the applications flag/options to the cobra cli
func (o *Options) Prepare(cmd *cobra.Command) *Options {
	o.addAppFlags(cmd.Flags())
	return o
}

// Complete initialises the components needed for the application to function given the options
func (o *Options) Complete() error {
	selectedFormat := strings.ToLower(strings.TrimSpace(o.Format))
	if !helpers.IsSupportedOutputFormat(selectedFormat) {
		// @fyi.error code invalid_output_format
		// @fyi.error title invalid_output_format
		// @fyi.error short the output format passed to --format was invalid, valid: yaml, markdown
		return errhandler.Error(errors.Errorf("the output format given %q is not valid", o.Format), "invalid_output_format")
	}

	// Check if output is a directory and error if the format chosen is YAML
	if file, err := os.Stat(o.Output); !errors.Is(err, os.ErrNotExist) {
		if file.IsDir() {
			// @fyi.error code invalid_yaml_output_file
			// @fyi.error title invalid_yaml_output_file
			// @fyi.error short the output file passed to the CLI is a directory not a file, please point a file
			return errhandler.Error(errors.Errorf("output %q should be a file not a directory", o.Output), "invalid_yaml_output_file")
		}
	}

	// TODO check language

	return nil
}

func getWorkingDirOrDie() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

func (o *Options) addAppFlags(fs *pflag.FlagSet) {
	fs.StringSliceVarP(
		&o.IncludedDirs,
		"include",
		"d",
		[]string{getWorkingDirOrDie()},
		"Comma separated list of directories to be parses by the tool",
	)
	fs.StringVar(
		&o.Format,
		"format",
		generate.Yaml,
		"Output format (yaml,markdown)",
	)
	fs.StringVarP(
		&o.Output,
		"output",
		"o",
		"",
		"Target output file or directory to store the generated output",
	)
	fs.BoolVarP(
		&o.Recursive,
		"recursive",
		"r",
		true,
		"Recursively parses the target directories",
	)
	fs.StringVarP(
		&o.Language,
		"language",
		"l",
		language.Go,
		"Target source code language",
	)
}
