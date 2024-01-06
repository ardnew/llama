package main

import (
	"os"
	"path/filepath"

	"github.com/ardnew/walk/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func main() {
	options := []walk.Option{}

	startPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--help" || os.Args[1] == "-h" {
			walk.Usage()
		}

		if os.Args[i] == "--version" || os.Args[1] == "-v" {
			walk.Version()
		}

		if os.Args[i] == "--icons" {
			options = append(options, walk.Icons())
			continue
		}

		startPath, err = filepath.Abs(os.Args[1])
		if err != nil {
			panic(err)
		}
	}

	options = append(options,
		walk.Path(startPath),
		walk.Size(80, 60),
	)

	output := termenv.NewOutput(os.Stderr)
	lipgloss.SetColorProfile(output.ColorProfile())

	w := walk.New(options...)
	p := tea.NewProgram(w, tea.WithOutput(os.Stderr))

	if _, err := p.Run(); err != nil {
		panic(err)
	}
	w.Exit()
}
