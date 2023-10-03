package docs

import (
	"os"
	"strings"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	commonoptions "github.com/tfadeyi/errors/cmd/app/options/common"
	"github.com/tfadeyi/errors/internal/generate"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		Format          string
		OutputDirectory string
		Source          string
		ErrorTemplate   string
		InfoTemplate    string
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
	if selectedFormat != generate.Markdown {
		return errors.Errorf("the output format given %q is not valid", o.Format)
	}

	// Check if output is a directory and error if the format chosen is YAML
	if file, err := os.Stat(o.OutputDirectory); !errors.Is(err, os.ErrNotExist) {
		if !file.IsDir() {
			return errors.Errorf("output %q is not a directory", o.OutputDirectory)
		}
	}

	return nil
}

func (o *Options) addAppFlags(fs *pflag.FlagSet) {
	fs.StringVar(
		&o.Format,
		"format",
		generate.Markdown,
		"Output format",
	)
	fs.StringVarP(
		&o.OutputDirectory,
		"output",
		"o",
		"./errors",
		"Target output file or directory to store the generated output",
	)
	fs.StringVar(
		&o.Source,
		"manifest",
		"",
		"Path to application error manifest file",
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
