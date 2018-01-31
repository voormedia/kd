package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/deploy"
)

var cmdDeploy = &cobra.Command{
	Use:   "deploy [APP] <TARGET>",
	Short: "Configure and deploy apps to k8s cluster",
	Args:  cobra.RangeArgs(1, 2),

	Run: func(cmd *cobra.Command, args []string) {
		names := []string{}
		if len(args) > 1 {
			names = args[0:1]
		}

		apps, err := config.ResolveAppNames(names)
		if err != nil {
			log.Fatal(err)
		}

		tgt, err := config.ResolveTargetName(args[len(args)-1])
		if err != nil {
			log.Fatal(err)
		}

		err = deploy.Run(verbose, log, apps, tgt)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdDeploy)
}
