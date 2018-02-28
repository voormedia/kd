package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

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
			if verbose {
				tw := tabwriter.NewWriter(os.Stdout, 10, 4, 3, ' ', 0)
				fmt.Fprintf(tw, "NAME\tALIASES\tCONTEXT\tNAMESPACE\n")
				for _, tgt := range conf.Targets {
					fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", tgt.Name, strings.Join(tgt.Alias, ", "), tgt.Context, tgt.Namespace)
				}
				tw.Flush()
			} else {
				fmt.Println(strings.Join(conf.TargetNames(), "\n"))
			}
		case "apps":
			if verbose {
				tw := tabwriter.NewWriter(os.Stdout, 10, 4, 2, ' ', 0)
				fmt.Fprintf(tw, "NAME\tPATH\tROOT\n")
				for _, app := range conf.Apps {
					fmt.Fprintf(tw, "%s\t%s\t%s\n", app.Name, app.Path, app.Root)
				}
				tw.Flush()
			} else {
				fmt.Println(strings.Join(conf.AppNames(), "\n"))
			}
		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdList)
}
