package cmd

import (
	"bytes"
	"os"

	"github.com/spf13/cobra"
)

var cmdCompletion = &cobra.Command{
	Use:                   "completion zsh",
	Short:                 "Output shell completion code for the specified shell",
	DisableFlagsInUseLine: true,

	ValidArgs: []string{"zsh"},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}

		if err := cobra.OnlyValidArgs(cmd, args); err != nil {
			return err
		}

		return nil
	},

	Long: `Output shell completion code for zsh. The shell code must be evaluated
to provide interactive completion of kd commands. This can be done by sourcing
it in .zprofile.`,

	Example: "  echo 'autoload -U compinit && compinit && source <(kd completion zsh)' >> ~/.zprofile",

	Run: func(cmd *cobra.Command, args []string) {
		var buf bytes.Buffer

		if err := cmd.Parent().GenZshCompletion(&buf); err != nil {
			log.Fatal(err)
		}

		os.Stdout.Write(buf.Bytes())
	},
}

func init() {
	cmdRoot.AddCommand(cmdCompletion)
}
