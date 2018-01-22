package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/deploy"
)

var cmdDeploy = &cobra.Command{
	Use:   "deploy <TARGET>",
	Short: "Configure and deploy apps to k8s cluster",
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		apps, err := config.ResolveAppNames([]string{})
		if err != nil {
			log.Fatal(err)
		}

		tgt, err := config.ResolveTargetName(args[0])
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
