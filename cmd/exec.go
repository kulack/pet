package cmd

import (
	"fmt"
	"os"
	"strings"
	"github.com/fatih/color"
	"github.com/kulack/pet/config"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run the selected commands",
	Long:  `Run the selected commands directly`,
	RunE:  execute,
}

func execute(cmd *cobra.Command, args []string) (err error) {
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

	runOptions := []int{}
	if ! flag.Silent {
		runOptions = append(runOptions, config.RunOptionEcho)
	}
	commands, err := filter(options)
	if err != nil {
		return err
	}
	command := strings.Join(commands, "; ")
	if config.Flag.Debug {
		fmt.Printf("Command: %s\n", command)
	}
	if config.Flag.Command {
		fmt.Printf("%s: %s\n", color.YellowString("Command"), command)
	}
	if flag.History || config.Conf.General.History {
		var histfile = "/tmp/pet.histfile"
		var escaped = strings.Replace(command, "'", "'\\''", -1)
		var hist = fmt.Sprintf("HISTFILE=%s history -s '%s'; history -w %s",
			histfile, escaped, histfile)
		var histOptions = []int{};
		if config.Flag.Debug {
			histOptions = runOptions
		}
		if err := run(hist, os.Stdin, os.Stdout, histOptions); err != nil {
			fmt.Println("Failed to update history: ", err)
		}
	}
	return run(command, os.Stdin, os.Stdout, runOptions)
}

func init() {
	RootCmd.AddCommand(execCmd)
	execCmd.Flags().BoolVarP(&config.Flag.Color, "color", "", false,
		`Enable colorized output (only fzf)`)
	execCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	execCmd.Flags().BoolVarP(&config.Flag.Command, "command", "c", false,
		`Show the command with the plain text before executing`)
	execCmd.Flags().BoolVarP(&config.Flag.NoSingleMatch, "nosingle", "n", false,
		`If a query (from -q|--query) matches only a single command, do not run it immediately, instead stay on the search page`)
	execCmd.Flags().BoolVarP(&config.Flag.Exact, "exact", "e", false,
		`A query (from -q|--query) matches exactly not fuzzy (depending on the configured selection command: fzf by default)`)
	execCmd.Flags().BoolVarP(&config.Flag.Silent, "silent", "s", false,
		`Silent, do not echo the command before running it`)
	execCmd.Flags().BoolVarP(&config.Flag.History, "history", "H", false,
		`Add executed command to history. Requires PROMPT_COMMAND='history -r /tmp/pet.histfile; echo "" > /tmp/pet.histfile' in your shell environment to actually import pet history if it exists`)
}
