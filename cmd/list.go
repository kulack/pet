package cmd

import (
	"fmt"
	"github.com/kulack/pet/dialog"
	"strings"

	"github.com/fatih/color"
	"github.com/kulack/pet/config"
	"github.com/kulack/pet/snippet"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	column = 40
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all snippets",
	Long:  `Show all snippets`,
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return err
	}

	col := config.Conf.General.Column
	if col == 0 {
		col = column
	}

	padTo := 0
	width := 10000000
	if terminal.IsTerminal(1) {
		padTo = col
		width = 100
		if w, _, err := terminal.GetSize(1); err == nil {
			width = w
		}
	} else {
		padTo = 0
	}

	for _, snippet := range snippets.Snippets {
		command := snippet.Command
		if ! config.Flag.Raw {
			command = dialog.PrepareCommand(command)
		}
		if config.Flag.OneLine {
			description := runewidth.FillRight(runewidth.Truncate(snippet.Description, col, "..."), padTo)
			command := runewidth.Truncate(snippet.Command, width-4-col, "...")
			// make sure multiline command printed as oneline
			command = strings.Replace(command, "\n", "\\n", -1)
			fmt.Fprintf(color.Output, "%s : %s\n",
				color.GreenString(description), color.YellowString(command))
		} else {
			fmt.Fprintf(color.Output, "%12s %s\n",
				color.GreenString("Description:"), snippet.Description)
			if strings.Contains(snippet.Command, "\n") {
				lines := strings.Split(snippet.Command, "\n")
				firstLine, restLines := lines[0], lines[1:]
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), firstLine)
				for _, line := range restLines {
					fmt.Fprintf(color.Output, "%12s %s\n",
						" ", line)
				}
			} else {
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), snippet.Command)
			}
			if snippet.Tag != nil {
				tag := strings.Join(snippet.Tag, " ")
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.CyanString("        Tag:"), tag)
			}
			if snippet.Output != "" {
				output := strings.Replace(snippet.Output, "\n", "\n             ", -1)
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.RedString("     Output:"), output)
			}
			fmt.Println(strings.Repeat("-", 30))
		}
	}
	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&config.Flag.OneLine, "oneline", "", false,
		`Display snippets in one line`)
	listCmd.Flags().BoolVarP(&config.Flag.Raw, "raw", "", false,
		`Display the raw source of the commands without escapes removed`)
}
