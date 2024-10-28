package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/util"
)

var version = "master"
var log = util.NewLogger("kd")

var cmdRoot = &cobra.Command{
	Use:   "kd",
	Short: "Build and deploy apps to k8s cluster",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			log.SetLevel(util.Debug)
		}
	},
}

func init() {
	cmdRoot.Version = version
	cmdRoot.PersistentFlags().Bool("verbose", false, "verbose output")
}

func setCommand(cmdRoot *cobra.Command, command string) {
	cmdRoot.SetArgs(append([]string{command}, os.Args[1:]...))
}

func Execute() {
	executableName := filepath.Base(os.Args[0])

	switch executableName {
	case "kbuild":
		setCommand(cmdRoot, "build")
	case "kctl":
		setCommand(cmdRoot, "ctl")
	case "kdeploy":
		setCommand(cmdRoot, "deploy")
	}

	cmdRoot.Execute()
}
