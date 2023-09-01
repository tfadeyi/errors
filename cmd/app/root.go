/*
Copyright © 2023 Oluwole Fadeyi (oluwolefadeyi@gmail.com)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package app

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	commonoptions "github.com/tfadeyi/errors/cmd/app/options/common"
	"github.com/tfadeyi/errors/internal/logging"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

func cmd(opts *commonoptions.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "errctl",
		Short: "Errctl.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logger := logging.LoggerFromContext(cmd.Context())
			logger = logger.WithName("root")

			if err := opts.Complete(); err != nil {
				logger.Error(err, "flag argument error")
				return err
			}
			if opts.LogLevel != "" {
				logger = logger.SetLevel(opts.LogLevel)
			}
			cmd.SetContext(logging.ContextWithLogger(cmd.Context(), logger))
			return nil
		},
	}
	opts = opts.Prepare(cmd)
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	opts := commonoptions.New()
	rootCmd = cmd(opts)
	rootCmd.AddCommand(specGenerateCmd(opts))
	rootCmd.AddCommand(specValidateCmd(opts))
	rootCmd.AddCommand(versionCmd(opts))
}
