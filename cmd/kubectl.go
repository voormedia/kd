package cmd

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/config"
	"github.com/voormedia/kd/pkg/kubectl"
)

var cmdKubectl = &cobra.Command{
	Use:                   "kubectl <target> [commands ...]",
	Short:                 "Invoke kubectl with project context and namespace",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Aliases:               []string{"ctl"},

	Long: `Invokes kubectl with the project context and namespace defined by the given
target. This ensures you always send commands to the correct cluster, with the
correct credentials and namespace.

This is meant to be used as a replacement for invoking kubectl directly.`,

	Example: "  kd ctl production get pods -o wide",

	DisableFlagParsing: true,

	PreRun: func(cmd *cobra.Command, args []string) {
		if conf, err := config.Load(); err == nil {
			cmd.ValidArgs = conf.TargetNames()
		}
	},

	Run: func(_ *cobra.Command, args []string) {
		err := kubectl.Run(log, args...)
		if err != nil {
			log.Fatal(err)
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					os.Exit(status.ExitStatus())
				}
			}
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdKubectl)
}
