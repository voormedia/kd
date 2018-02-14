package cmd

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/config"
)

var cmdKubectl = &cobra.Command{
	Use:     "kubectl <target> [commands ...]",
	Short:   "Invoke kubectl with project context and namespace",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"ctl"},

	Long: `Invokes kubectl with the project context and namespace defined by the given
target. This ensures you always send commands to the correct cluster, with the
correct credentials and namespace.

This is meant to be used as a replacement for invoking kubectl directly.`,

	Example: "  kd kubectl production get pods -o wide",

	DisableFlagParsing: true,

	PreRun: func(cmd *cobra.Command, args []string) {
		if conf, err := config.Load(); err == nil {
			cmd.ValidArgs = conf.TargetNames()
		}
	},

	Run: func(_ *cobra.Command, args []string) {
		conf, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}

		tgt, err := conf.ResolveTarget(args[0])
		if err != nil {
			log.Fatal(err)
		}

		args = append([]string{
			"--context", tgt.Context,
			"--namespace", tgt.Namespace,
		}, args[1:]...)

		cmd := exec.Command("kubectl", args...)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			log.Error(err)
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
