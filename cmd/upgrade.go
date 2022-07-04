package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/upgrade"
)

var cmdUpgrade = &cobra.Command{
	Use:                   "upgrade",
	Short:                 "Upgrade configuration to the latest version",
	DisableFlagsInUseLine: true,

	Run: func(_ *cobra.Command, args []string) {
		if err := upgrade.Run(log); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdUpgrade)
}
