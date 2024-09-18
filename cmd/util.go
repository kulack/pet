package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/kulack/pet/config"
	"github.com/kulack/pet/dialog"
	"github.com/kulack/pet/snippet"
)

func editFile(command, file string, startingLine int) error {
	// Note that this works for most unix editors (nano, vi, vim, etc)
	command += " " + file
	return run(command, os.Stdin, os.Stdout)
}

func filter(options []string, tag string) (commands []string, err error) {
	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return commands, fmt.Errorf("load snippet failed: %v", err)
	}

	if 0 < len(tag) {
		var filteredSnippets snippet.Snippets
		for _, snippet := range snippets.Snippets {
			for _, t := range snippet.Tag {
				if tag == t {
					filteredSnippets.Snippets = append(filteredSnippets.Snippets, snippet)
				}
			}
		}
		snippets = filteredSnippets
	}

	snippetTexts := map[string]snippet.SnippetInfo{}
	var text string
	for _, s := range snippets.Snippets {
		command := s.Command
		if strings.ContainsAny(command, "\n") {
			command = strings.Replace(command, "\n", "\\n", -1)
		}

		tags := ""
		for _, tag := range s.Tag {
			tags += fmt.Sprintf("#%s ", tag)
		}

		format := "[$description]: $command $tags"
		if config.Conf.General.Format != "" {
			format = config.Conf.General.Format
		}

		t := strings.Replace(format, "$command", command, 1)
		t = strings.Replace(t, "$description", s.Description, 1)
		t = strings.Replace(t, "$tags", tags, 1)

		snippetTexts[t] = s
		if config.Flag.Color || config.Conf.General.Color {
			t = strings.Replace(format, "$command", command, 1)
			t = strings.Replace(t, "$description", color.HiRedString(s.Description), 1)
			t = strings.Replace(t, "$tags", color.HiCyanString(tags), 1)
		}
		text += t + "\n"
	}

	var buf bytes.Buffer
	selectCmd := fmt.Sprintf("%s %s",
		config.Conf.General.SelectCmd, strings.Join(options, " "))
	err = run(selectCmd, strings.NewReader(text), &buf)
	if err != nil {
		return nil, nil
	}

	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")
	var params [][2]string

	// If only one line is selected, search for params in the command
	if len(lines) == 1 {
		snippetInfo := snippetTexts[lines[0]]
		params = dialog.SearchForParams(snippetInfo.Command)
	} else {
		params = nil
	}

	if params != nil {
		snippetInfo := snippetTexts[lines[0]]
		dialog.CurrentCommand = snippetInfo.Command
		dialog.GenerateParamsLayout(params, dialog.CurrentCommand)
		res := []string{dialog.FinalCommand}
		return res, nil
	}
	for _, line := range lines {
		snippetInfo := snippetTexts[line]
		commands = append(commands, fmt.Sprint(snippetInfo.Command))
	}
	return commands, nil
}

func selectFile(options []string, tag string) (snippetFile string, err error) {
	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return snippetFile, fmt.Errorf("load snippet failed: %v", err)
	}

	if 0 < len(tag) {
		var filteredSnippets snippet.Snippets
		for _, snippet := range snippets.Snippets {
			for _, t := range snippet.Tag {
				if tag == t {
					filteredSnippets.Snippets = append(filteredSnippets.Snippets, snippet)
				}
			}
		}
		snippets = filteredSnippets
	}

	snippetTexts := map[string]snippet.SnippetInfo{}
	var text string
	for _, s := range snippets.Snippets {
		command := s.Command
		if strings.ContainsAny(command, "\n") {
			command = strings.Replace(command, "\n", "\\n", -1)
		}
		t := fmt.Sprintf("[%s]: %s", s.Description, command)

		tags := ""
		for _, tag := range s.Tag {
			tags += fmt.Sprintf(" #%s", tag)
		}
		t += tags

		snippetTexts[t] = s
		text += t + "\n"
	}

	var buf bytes.Buffer
	selectCmd := fmt.Sprintf("%s %s",
		config.Conf.General.SelectCmd, strings.Join(options, " "))
	err = run(selectCmd, strings.NewReader(text), &buf)
	if err != nil {
		return snippetFile, nil
	}

	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")

	for _, line := range lines {
		snippetInfo := snippetTexts[line]
		snippetFile = fmt.Sprint(snippetInfo.Filename)
	}
	return snippetFile, nil
}

// CountLines returns the number of lines in a certain buffer
func CountLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
