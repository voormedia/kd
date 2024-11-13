package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/build"
	"github.com/voormedia/kd/pkg/config"
)

var buildTag string = ""
var buildCacheTag string = ""

var cmdBuild = &cobra.Command{
	Use:                   "build [app[:tag]]",
	Short:                 "Build container images for an application",
	DisableFlagsInUseLine: true,

	Args:    cobra.RangeArgs(0, 1),
	Aliases: []string{"bld"},

	Long: `Builds a single application. If only one application is configured, the name
can be omitted. Application images will be pushed to the registry and tagged
as "latest" by default. The tag can optionally be specified.`,

	Example: "  kd build my-app\n  kd build my-app:awesome-tag",

	PreRun: func(cmd *cobra.Command, args []string) {
		if conf, err := config.Load(); err == nil {
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

		app, err := conf.ResolveApp(name, buildTag)
		if err != nil {
			log.Fatal(err)
		}

		err = build.Run(log, app, buildCacheTag)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	cmdBuild.Flags().StringVar(&buildTag, "tag", "", "tag to use for the built image")
	cmdBuild.Flags().StringVar(&buildCacheTag, "cache-tag", "", "tag to use for build cache (defaults to git branch)")
	cmdRoot.AddCommand(cmdBuild)
}
