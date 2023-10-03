package app

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	annotateoptions "github.com/tfadeyi/errors/cmd/app/options/annotate"
	commonoptions "github.com/tfadeyi/errors/cmd/app/options/common"
	"github.com/tfadeyi/errors/internal/logging"
	"github.com/tfadeyi/errors/internal/parser"
	"github.com/tfadeyi/errors/internal/parser/language"
)

func annotateCmd(common *commonoptions.Options) *cobra.Command {
	opts := annotateoptions.New(common)

	cmd := &cobra.Command{
		Use:          "annotate",
		Short:        "Annotates the target file with error.fyi comments",
		Long:         ``,
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())
			logger = logger.WithName("annotate")

			if err := opts.Complete(); err != nil {
				return err
			}

			cmd.SetContext(logging.ContextWithLogger(cmd.Context(), logger))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())

			logger.Info("Annotating target file", "file", opts.Source)

			parserOptions := []parser.Option{
				parser.Logger(&logger),
			}

			switch opts.Language {
			case language.Go:
				parserOptions = append(parserOptions, parser.Go(cmd.Context()))
			default:
				// do nothing
			}

			logger.Info("Parser configuration",
				"file", opts.Source,
				"language", opts.Language,
			)

			p := parser.New(
				parserOptions...,
			)

			logger.Info("Annotating file...")
			err := p.AnnotateSourceErrors(cmd.Context(), opts.Source, opts.WrapErrors, opts.AnnotateOnlyTodos)
			if err != nil {
				return errors.Annotate(err, "failed to annotate the application")
			}
			logger.Info("Done ✅")

			logger.Info("@fyi annotations were successfully added!")
			logger.Info("Check the updated file", "file", opts.Source)

			return nil
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}
