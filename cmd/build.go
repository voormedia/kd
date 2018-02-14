package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/build"
	"github.com/voormedia/kd/pkg/config"
)

var cmdBuild = &cobra.Command{
	Use:     "build [app[:tag]]",
	Short:   "Build container images for an application",
	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"bld"},

	Long: `Builds a single application. If only one application is configured, the name
can be omitted.Application images will be pushed to the registry and tagged
as 'latest' by default. The tag can optionally be specified per application.`,

	Example: "  kd build my-app",

	PreRun: func(cmd *cobra.Command, args []string) {
		if conf, err := config.Load(); err != nil {
			cmd.ValidArgs = conf.AppNames()
		}
	},

	Run: func(_ *cobra.Command, args []string) {
		conf, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		name := ""
		if len(args) > 0 {
			name = args[0]
		}

		app, err := conf.ResolveApp(name)
		if err != nil {
			log.Fatal(err)
		}

		err = build.Run(verbose, log, app)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdBuild)
}
