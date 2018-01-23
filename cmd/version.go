package cmd

import (
	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Show version",

	Run: func(cmd *cobra.Command, args []string) {
		log.Log("version", cmdRoot.Version)
	},
}

func init() {
	cmdRoot.AddCommand(cmdVersion)
}
