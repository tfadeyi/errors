package spec

import (
	"os"
	"strings"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	errhandler "github.com/tfadeyi/errors"
	commonoptions "github.com/tfadeyi/errors/cmd/app/options/common"
	"github.com/tfadeyi/errors/internal/parser/generate"
	"github.com/tfadeyi/errors/internal/parser/language"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		Format                 string
		IncludedDirs           []string
		OutputFileAndDirectory string
		Source                 string
		Language               string
		ErrorTemplate          string
		InfoTemplate           string
		*commonoptions.Options
	}
)

// New creates a new instance of the application's options
func New(common *commonoptions.Options) *Options {
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
	if !generate.IsSupportedOutputFormat(selectedFormat) {
		// @fyi.error code invalid_output_format
		// @fyi.error title invalid_output_format
		// @fyi.error short the output format passed to --format was invalid, valid: yaml, markdown
		return errhandler.Error(errors.Errorf("the output format given %q is not valid", o.Format), "invalid_output_format")
	}

	// Check if output is a directory and error if the format chosen is YAML
	if file, err := os.Stat(o.OutputFileAndDirectory); !errors.Is(err, os.ErrNotExist) {
		if file.IsDir() && o.Format == generate.Yaml {
			// Here we add more specific info about the error they may encounter in this configuration

			// @fyi.error code invalid_yaml_output_file
			// @fyi.error title invalid_yaml_output_file
			// @fyi.error short the output file passed to the CLI is a directory not a file, please point a file
			return errhandler.Error(errors.Errorf("output %q should be a file not a directory", o.OutputFileAndDirectory), "invalid_yaml_output_file")
		}
	}

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
		&o.OutputFileAndDirectory,
		"output",
		"o",
		"",
		"Target output file or directory to store the generated output",
	)
	fs.StringVarP(
		&o.Source,
		"file",
		"f",
		"",
		"Source code file to parse",
	)
	fs.StringVarP(
		&o.Language,
		"language",
		"l",
		language.Go,
		"Target source code language",
	)
	fs.StringVar(
		&o.ErrorTemplate,
		"error-template",
		"",
		"Custom application error go-template filepath (markdown)",
	)
	fs.StringVar(
		&o.InfoTemplate,
		"info-template",
		"",
		"Custom application information go-template filepath (markdown)",
	)
}
