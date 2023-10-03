package manifest

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	commonoptions "github.com/tfadeyi/errors/cmd/app/options/common"
	"os"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		Source string
		Strict bool
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
	o.addAppFlags(cmd.PersistentFlags())
	return o
}

// Complete initialises the components needed for the application to function given the options
func (o *Options) Complete() error {
	_, err := os.Stat(o.Source)
	if errors.Is(err, os.ErrNotExist) {
		// TODO write better error
		return errors.New("source file is missing")
	}
	return nil
}

func (o *Options) addAppFlags(fs *pflag.FlagSet) {
	fs.StringVarP(
		&o.Source,
		"file",
		"f",
		"",
		"Application error manifest to parse",
	)
	fs.BoolVar(
		&o.Strict,
		"strict",
		false,
		"Set the parser to strict mode",
	)
}
