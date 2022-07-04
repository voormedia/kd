package cmd

import (
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/util"
)

var version = "2.0.0"
var log = util.NewLogger("kd")

var cmdRoot = &cobra.Command{
	Use:   "kd",
	Short: "Build and deploy apps to k8s cluster",
}

func init() {
	cmdRoot.Version = version
}

func Execute() {
	cmdRoot.Execute()
}
