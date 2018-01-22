package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/build"
	"github.com/voormedia/kd/pkg/config"
)

var cmdBuild = &cobra.Command{
	Use:   "build [APP...]",
	Short: "Build container images for apps",
	Args:  cobra.ArbitraryArgs,

	Run: func(cmd *cobra.Command, args []string) {
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
