package manifest

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	commonoptions "github.com/tfadeyi/errors/cmd/app/options/common"
	"github.com/tfadeyi/errors/internal/parser/language"
	"os"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		IncludedDirs      []string
		Source            string
		Language          string
		WrapErrors        bool
		AnnotateOnlyTodos bool
		Recursive         bool
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
	fs.BoolVar(
		&o.WrapErrors,
		"wrap",
		false,
		"Wrap errors with error.fyi error wrapper library",
	)
	fs.BoolVar(
		&o.AnnotateOnlyTodos,
		"todo",
		false,
		"Annotates only the errors with a TODO comment above",
	)
}
