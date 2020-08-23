package dialog

import (
	"github.com/jroimartin/gocui"
	"github.com/kulack/pet/config"
	"regexp"
	"strings"
	"os"
)

var (
	views      = []string{}
	layoutStep = 3
	curView    = -1
	idxView    = 0

	//CurrentCommand is the command before assigning to variables
	CurrentCommand string
	//FinalCommand is the command after assigning to variables
	FinalCommand string
)

func insertParams(command string, params map[string]string) string {
	resultCommand := command
	for k, v := range params {
		resultCommand = strings.Replace(resultCommand, k, v, -1)
	}
	return resultCommand
}

// SearchForParams returns variables from a command
//func SearchForLegacyParams(lines []string) map[string]string {
//	// Update from commit fb528be88a82eed4f6c06b4128a7dfac86162552
//	// Old pull request: https://github.com/knqyf263/pet/pull/54/commits
//	// Additionally, added, allow escape parameters by using !<notParameter>
//	re := `[^!]<([^\s=>]+(?:=(?:\\\\|\\>|[^>\\])*)?)>`
//	if len(lines) == 1 {
//		r, _ := regexp.Compile(re)
//
//		params := r.FindAllStringSubmatch(lines[0], -1)
//		if len(params) == 0 {
//			return nil
//		}
//
//		extracted := map[string]string{}
//		for _, p := range params {
//			if p[0][:0] != "<" {
//				// Trim off any leading character that was matched
//				// by the value of [^!] in the regex
//				p[0] = p[0][1:]
//			}
//			splitted := strings.SplitN(p[1], "=", 2)
//			if len(splitted) == 1 {
//				// There is no value specified for the variable, pull the
//				// variable from the environment if it exists.
//				extracted[p[0]] = os.ExpandEnv("$" + splitted[0])
//			} else {
//				extracted[p[0]] = splitted[1]
//			}
//		}
//		return extracted
//	}
//	return nil
//}

func SearchForParams(lines []string) map[string]string {
	// Regular expression syntax is ${VAR=defaultval} where =defaultval is optional
	// Escape the parameter syntax by using \${VAR}
	re := `[^$]\$\{([^\s=\}]+(?:=(?:\\\\|\\\}|[^\}\\])*)?)\}`
	startParam := "{"
	if config.Conf.General.LegacyParams {
		// Update from commit fb528be88a82eed4f6c06b4128a7dfac86162552
		// Old pull request: https://github.com/knqyf263/pet/pull/54/commits
		// Additionally, added, allow escape parameters by using !<notParameter>
		re = `[^!]<([^\s=>]+(?:=(?:\\\\|\\>|[^>\\])*)?)>`
		startParam = "<"
	}
	if len(lines) == 1 {
		r, _ := regexp.Compile(re)

		params := r.FindAllStringSubmatch(lines[0], -1)
		if len(params) == 0 {
			return nil
		}

		extracted := map[string]string{}
		for _, p := range params {
			if p[0][:0] != startParam {
				// Trim off any leading character that was matched
				// by the value of [^x] at the beginning of the regex
				p[0] = p[0][1:]
			}
			splitted := strings.SplitN(p[1], "=", 2)
			if len(splitted) == 1 {
				// There is no value specified for the variable, pull the
				// default variable from the environment if it exists.
				extracted[p[0]] = os.ExpandEnv("$" + splitted[0])
			} else {
				extracted[p[0]] = splitted[1]
			}
		}
		return extracted
	}
	return nil
}

// If the command has no parameters but may have escape sequences
// or shell variable usages, apply those changes and return the command.
func prepareLegacyCommand(command string) string {
	cmd := strings.Replace(command, "!<", "<", -1)
	return cmd
}

func PrepareCommand(command string) string {
	if config.Conf.General.LegacyParams {
		return prepareLegacyCommand(command)
	}
	cmd := strings.Replace(command, "$${", "${", -1)
	return cmd
}

func evaluateParams(g *gocui.Gui, _ *gocui.View) error {
	paramsFilled := map[string]string{}
	for _, v := range views {
		view, _ := g.View(v)
		res := view.Buffer()
		res = strings.Replace(res, "\n", "", -1)
		paramsFilled[v] = strings.TrimSpace(res)
	}
	// TODO: Could be replacing real text like from an XML comment.
	FinalCommand = insertParams(PrepareCommand(CurrentCommand), paramsFilled)
	return gocui.ErrQuit
}
