package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/voormedia/kd/pkg/config"
)

var output outputType = "table"

var cmdList = &cobra.Command{
	Use:                   "list (targets | apps)",
	Short:                 "List configured targets or applications",
	DisableFlagsInUseLine: false,

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
			if output == "name" {
				fmt.Println(strings.Join(conf.TargetNames(), "\n"))
			} else {
				tw := tabwriter.NewWriter(os.Stdout, 10, 4, 3, ' ', 0)
				fmt.Fprintf(tw, "NAME\tALIASES\tCONTEXT\tNAMESPACE\n")
				for _, tgt := range conf.Targets {
					fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", tgt.Name, strings.Join(tgt.Alias, ", "), tgt.Context, tgt.Namespace)
				}
				tw.Flush()
			}
		case "apps":
			if output == "name" {
				fmt.Println(strings.Join(conf.AppNames(), "\n"))
			} else {
				tw := tabwriter.NewWriter(os.Stdout, 10, 4, 2, ' ', 0)
				fmt.Fprintf(tw, "NAME\tPATH\tROOT\n")
				for _, app := range conf.Apps {
					fmt.Fprintf(tw, "%s\t%s\t%s\n", app.Name, app.Path, app.Root)
				}
				tw.Flush()
			}
		}
	},
}

func init() {
	cmdList.Flags().VarP(&output, "output", "o", `output format, either "table" or "name"`)
	cmdRoot.AddCommand(cmdList)
}

type outputType string

const (
	outputTable outputType = "table"
	outputName  outputType = "name"
)

func (e *outputType) String() string {
	return `"` + string(*e) + `"`
}

func (e *outputType) Set(v string) error {
	switch v {
	case "table", "name":
		*e = outputType(v)
		return nil
	default:
		return errors.New(`must be either "table" or "name"`)
	}
}

func (e *outputType) Type() string {
	return "type"
}
