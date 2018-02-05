package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/build"
	"github.com/voormedia/kd/pkg/config"
)

var cmdBuild = &cobra.Command{
	Use:   "build [app[:tag] ...]",
	Short: "Build container images for all or some applications",
	Args:  cobra.ArbitraryArgs,
	// Aliases: []string{"bld"},

	Long: `Builds either all applications, or a single application. Application images
will be pushed to the registry and tagged as 'latest' by default. The tag can
optionally be specified per application.`,

	Example: "  kd build my-app my-other-app",

	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.ValidArgs = config.AppNames()
	},

	Run: func(_ *cobra.Command, args []string) {
		apps, err := config.ResolveAppNames(args)
		if err != nil {
			log.Fatal(err)
		}

		err = build.Run(verbose, log, apps)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdBuild)
}
