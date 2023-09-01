package app

import (
	"github.com/spf13/cobra"

	commonoptions "github.com/tfadeyi/errors/cmd/app/options/common"
	versionoptions "github.com/tfadeyi/errors/cmd/app/options/version"
	"github.com/tfadeyi/errors/internal/logging"
	"github.com/tfadeyi/errors/internal/version"
)

func versionCmd(common *commonoptions.Options) *cobra.Command {
	opts := versionoptions.New(common)
	cmd := &cobra.Command{
		Use:           "version",
		Short:         "Returns the binary build information.",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := opts.Validate()
			if err != nil {
				return err
			}
			err = opts.Complete()
			return err
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			log := logging.LoggerFromContext(ctx)
			log = log.WithName("version")
			log.Debug(version.BuildInfo())
			log.Info(version.Info())
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}
