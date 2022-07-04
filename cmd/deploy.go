package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/deploy"
)

var cmdDeploy = &cobra.Command{
	Use:                   "deploy [app[:tag]] <target>",
	Short:                 "Configure and deploy an application to a cluster",
	DisableFlagsInUseLine: true,

	Args:    cobra.RangeArgs(1, 2),
	Aliases: []string{"dep"},

	Long: `Deploys a single application to the given target. If only one application
is configured, the name can be omitted. By default the application image with
the 'latest' tag in the registry will be deployed. The tag of the image to
deploy can optionally be specified.

Any image that was successfully deployed will be tagged with the name of the
target to which it was deployed.`,

	Example: "  kd deploy my-app production",

	PreRun: func(cmd *cobra.Command, args []string) {
		if conf, err := config.Load(); err == nil {
			cmd.ValidArgs = conf.TargetNames()
		}
	},

	Run: func(_ *cobra.Command, args []string) {
		conf, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		name := ""
		if len(args) > 1 {
			name = args[0]
		}

		tgt, err := conf.ResolveTarget(args[len(args)-1])
		if err != nil {
			log.Fatal(err)
		}

		app, err := conf.ResolveApp(name)
		if err != nil {
			log.Fatal(err)
		}

		err = deploy.Run(log, app, tgt)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdDeploy)
}
