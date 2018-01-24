package cmd

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/config"
)

var cmdKubectl = &cobra.Command{
	Use:   "kubectl <TARGET>",
	Short: "Invoke kubectl with project context and namespace",
	Args:  cobra.MinimumNArgs(1),

	Run: func(_ *cobra.Command, args []string) {
		tgt, err := config.ResolveTargetName(args[0])
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
