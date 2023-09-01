package app

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/tfadeyi/errors/internal/parser"
	"github.com/tfadeyi/errors/internal/parser/generate"
	"github.com/tfadeyi/errors/internal/parser/language"
	"github.com/tfadeyi/errors/internal/parser/options"
	"io"

	fyi "github.com/tfadeyi/errors"
	commonoptions "github.com/tfadeyi/errors/cmd/app/options/common"
	specoptions "github.com/tfadeyi/errors/cmd/app/options/spec"
	"github.com/tfadeyi/errors/internal/logging"
)

func specGenerateCmd(common *commonoptions.Options) *cobra.Command {
	opts := specoptions.New(common)
	var inputReader io.ReadCloser

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates the targeted application error manifest from a given source code",
		Long:  ``,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())
			logger = logger.WithName("generate")

			if err := opts.Complete(); err != nil {
				return err
			}

			if opts.Source == "-" {
				inputReader = io.NopCloser(cmd.InOrStdin())
			}

			cmd.SetContext(logging.ContextWithLogger(cmd.Context(), logger))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())

			logger.Info("Parsing source code for @fyi error definitions ⚙️",
				"directories", opts.IncludedDirs,
			)

			parserOptions := []options.Option{
				options.Include(opts.IncludedDirs...),
				options.Logger(&logger),
				options.Output(opts.OutputFileAndDirectory),
				options.SourceFile(opts.Source),
				options.SourceContent(inputReader),
				options.Watermark(`# Code generated by errctl: https://github.com/tfadeyi/errors.
# DO NOT EDIT.`),
				options.CustomManifestInfoTemplate(opts.InfoTemplate),
				options.CustomManifestErrorTemplate(opts.ErrorTemplate),
			}

			switch opts.Language {
			case language.Go:
				parserOptions = append(parserOptions, options.Go())
			default:
				// do nothing
			}

			switch opts.Format {
			case generate.Yaml:
				parserOptions = append(parserOptions, options.YAML(cmd.OutOrStdout()))
			case generate.Markdown:
				parserOptions = append(parserOptions, options.Markdown(cmd.OutOrStdout()))
			}

			// @fyi.error code clean_artefacts_error
			// @fyi.error title Error Removing Previous Artefacts
			// @fyi.error short The tool has failed to delete the artefacts from the previous execution.
			// @fyi.error long The tool has failed to delete the artefacts from the previous execution. Try manually deleting them before running the tool again.

			p := parser.New(
				parserOptions...,
			)

			apps, err := p.Parse(cmd.Context())
			if err != nil {
				return errors.Annotate(err, "failed to parse the application(s) error manifests")
			}

			logger.Info("Source code was successfully parsed ✅")

			if err = p.Generate(cmd.Context(), apps); err != nil {
				return errors.Annotate(err, "failed to printout the application(s) error manifests")
			}

			logger.Info("Application(s) error manifest were successfully generated ✅")

			return nil
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}

func specValidateCmd(common *commonoptions.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "validate",
		Short:         "Validates a given aloe specification",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// @fyi.error code validate_not_implemented
			// @fyi.error title validate_not_implemented
			// @fyi.error short spec validate command has not been implemented yet
			// @fyi.error long specification validate command has not been implemented yet, will be implemented shortly
			return fyi.Error(errors.New("not implemented"), "validate_not_implemented")
		},
	}
	return cmd
}
