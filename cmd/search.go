package cmd

import (
	"fmt"
	"strings"


	"github.com/kulack/pet/config"
	"github.com/kulack/pet/dialog"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var delimiter string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search snippets",
	Long:  `Search snippets interactively (default filtering tool: peco)`,
	RunE:  search,
}

func search(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag

	var options []string
	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query \"%s\"", flag.Query))
		if ! flag.NoSingleMatch {
			// Convert the single match flag into the search tools required flag
			options = append(options, config.Conf.General.SingleMatch)
		}
	}
	if flag.Exact {
		options = append(options, "--exact")
	}
	commands, err := filter(options)
	if err != nil {
		return err
	}

	command := strings.Join(commands, flag.Delimiter);
	fmt.Print(dialog.PrepareCommand(command))

	if terminal.IsTerminal(1) {
		fmt.Print("\n")
	}
	return nil
}

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVarP(&config.Flag.Color, "color", "", false,
		`Enable colorized output (only fzf)`)
	searchCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	searchCmd.Flags().BoolVarP(&config.Flag.NoSingleMatch, "nosingle", "n", false,
		`If a query (from -q|--query) matches only a single command, do not use it immediately, instead stay on the search page`)
	searchCmd.Flags().BoolVarP(&config.Flag.Exact, "exact", "e", false,
		`A query (from -q|--query) matches exactly not fuzzy (depending on the configured selection command: fzf by default)`)
	searchCmd.Flags().StringVarP(&config.Flag.Delimiter, "delimiter", "d", "; ",
		`Use delim as the command delimiter character`)
}
