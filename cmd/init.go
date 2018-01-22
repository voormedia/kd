package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/scaffold"
)

var cmdInit = &cobra.Command{
	Use:   "init",
	Short: "Generate initial configuration files",

	Run: func(cmd *cobra.Command, args []string) {
		if err := scaffold.Run(log); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdInit)
}
