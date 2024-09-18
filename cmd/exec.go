package cmd

import (
	"al.essio.dev/pkg/shellescape"
	"fmt"
	"os"
	"strings"

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
		options = append(options, fmt.Sprintf("--query %s", shellescape.Quote(flag.Query)))
	}

	commands, err := filter(options, flag.FilterTag)
	if err != nil {
		return err
	}
	command := strings.Join(commands, "; ")

	// Show final command before executing it
	fmt.Printf("> %s\n", command)

	if flag.History || config.Conf.General.History {
		var histfile = "/tmp/pet.histfile"
		var escaped = strings.Replace(command, "'", "'\\''", -1)
		var hist = fmt.Sprintf("HISTFILE=%s history -s '%s'; history -w %s",
			histfile, escaped, histfile)
		if err := run(hist, os.Stdin, os.Stdout); err != nil {
			fmt.Println("Failed to update history: ", err)
		}
	}
	return run(command, os.Stdin, os.Stdout)
}

func init() {
	RootCmd.AddCommand(execCmd)
	execCmd.Flags().BoolVarP(&config.Flag.Color, "color", "", false,
		`Enable colorized output (only fzf)`)
	execCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	execCmd.Flags().StringVarP(&config.Flag.FilterTag, "tag", "t", "",
		`Filter tag`)
	execCmd.Flags().BoolVarP(&config.Flag.History, "history", "H", false,
		`Write History to /tmp/pet.histfile`)
}
