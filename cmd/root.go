package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/util"
)

var version = "master"
var log = util.NewLogger("kd")

var verbose bool

var cmdRoot = &cobra.Command{
	Use:   "kd",
	Short: "Build and deploy apps to k8s cluster",
	BashCompletionFunction: customCompletion,
}

func init() {
	cmdRoot.Version = version
	cmdRoot.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func Execute() {
	cmdRoot.Execute()
}
