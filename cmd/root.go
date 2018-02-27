package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/util"
)

var log = util.NewLogger("kd")

var cmdRoot = &cobra.Command{
	Use:     "kd",
	Short:   "Build and deploy apps to k8s cluster",
	Version: "1.2.2",
}

var verbose bool

func init() {
	cmdRoot.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func Execute() {
	cmdRoot.Execute()
}
