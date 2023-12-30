package main

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func main() {
	startPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	showIcons := false
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--help" || os.Args[1] == "-h" {
			usage()
		}

		if os.Args[i] == "--version" || os.Args[1] == "-v" {
			version()
		}

		if os.Args[i] == "--icons" {
			showIcons = true
			continue
		}

		startPath, err = filepath.Abs(os.Args[1])
		if err != nil {
			panic(err)
		}
	}

	output := termenv.NewOutput(os.Stderr)
	lipgloss.SetColorProfile(output.ColorProfile())

	m := &model{
		path:      startPath,
		width:     80,
		height:    60,
		positions: make(map[string]position),
		showIcons: showIcons,
	}
	m.list()

	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	os.Exit(m.exitCode)
}
