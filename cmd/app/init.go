package app

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	commonoptions "github.com/tfadeyi/errors/cmd/app/options/common"
	initoptions "github.com/tfadeyi/errors/cmd/app/options/init"
	"github.com/tfadeyi/errors/internal/logging"
	"github.com/tfadeyi/errors/internal/parser"
	"github.com/tfadeyi/errors/internal/parser/language"
)

func initCmd(common *commonoptions.Options) *cobra.Command {
	opts := initoptions.New(common)

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes the current project to use error.fyi",
		Long:  ``,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())
			logger = logger.WithName("init")

			if err := opts.Complete(); err != nil {
				return err
			}

			cmd.SetContext(logging.ContextWithLogger(cmd.Context(), logger))
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())

			logger.Info("Initializing project")

			parserOptions := []parser.Option{
				parser.Include(opts.IncludedDirs...),
				parser.Recursive(true),
				parser.Logger(&logger),
			}

			switch opts.Language {
			case language.Go:
				parserOptions = append(parserOptions, parser.Go(cmd.Context()))
			default:
				// do nothing
			}

			logger.Info("Parser configuration",
				"target directories", opts.IncludedDirs,
				"language", opts.Language,
			)

			p := parser.New(
				parserOptions...,
			)

			logger.Info("Annotating error declarations...")
			err := p.AnnotateAllErrors(cmd.Context(), opts.WrapErrors, opts.AnnotateOnlyTodos)
			if err != nil {
				return errors.Annotate(err, "failed to annotate the application")
			}
			logger.Info("Done ✅")

			logger.Info("Project was successfully initialized!")
			logger.Info("Check the target directories", "directories", opts.IncludedDirs)

			return nil
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}
