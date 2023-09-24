package annotate

import (
	"os"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	fyi "github.com/tfadeyi/errors"
	commonoptions "github.com/tfadeyi/errors/cmd/app/options/common"
	"github.com/tfadeyi/errors/internal/parser/language"
)

type (
	// Options is the list of options/flag available to the application,
	// plus the clients needed by the application to function.
	Options struct {
		IncludedDirs []string
		Source       string
		Language     string
		WrapErrors   bool
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
	// check the target file is present
	_, err := os.Stat(o.Source)
	if errors.Is(err, os.ErrNotExist) {
		// @fyi.error code annotate_error_43
		// @fyi.error title Input File Does Not Exist
		// @fyi.error short The tool tried fetching the target file but could not find it.
		return fyi.Error(err, "annotate_error_43")
	}

	// check the language is part of the supported group
	if !language.IsSupportedLanguage(o.Language) {
		// @fyi.error code annotate_error_52
		// @fyi.error title Unsupported Language
		// @fyi.error short The selected language is not supported by the tool's parser, available: go.
		return fyi.Error(errors.Errorf("unsupported language: %s", o.Language), "annotate_error_52")
	}
	return nil
}

func (o *Options) addAppFlags(fs *pflag.FlagSet) {
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
	fs.BoolVar(
		&o.WrapErrors,
		"wrap",
		false,
		"Wrap errors with error.fyi error wrapper library",
	)
}
