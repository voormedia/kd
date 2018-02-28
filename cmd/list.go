package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/config"
)

var cmdList = &cobra.Command{
	Use:   "list (targets | apps)",
	Short: "List configured targets or applications",
	DisableFlagsInUseLine: true,

	ValidArgs: []string{"targets", "apps"},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		if err := cobra.OnlyValidArgs(cmd, args); err != nil {
			return err
		}

		return nil
	},

	Long: `List all available targets or applications based on the configuration that
was found in the current directory.`,

	Example: "  kd list targets",

	Run: func(_ *cobra.Command, args []string) {
		conf, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		switch args[0] {
		case "targets":
			fmt.Println(strings.Join(conf.TargetNames(), "\n"))
		case "apps":
			fmt.Println(strings.Join(conf.AppNames(), "\n"))
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdList)
}
