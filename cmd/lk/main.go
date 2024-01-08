package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/ardnew/walk/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func usage(s *walk.Styles) {
	exe, err := os.Executable()
	if err != nil {
		exe = os.Args[0]
	}
	exe = filepath.Base(exe)
        _, _ = fmt.Fprintf(os.Stderr, "\n  "+s.Cursor.Render(" " + exe + " ")+"\n\n  Usage: " + exe + " [flags] [path]\n\n")
        w := tabwriter.NewWriter(os.Stderr, 0, 8, 2, ' ', 0)
        put := func(s string) {
                _, _ = fmt.Fprintln(w, s)
        }
        put("    Arrows, hjkl\tMove cursor")
        put("    Enter\tEnter directory")
        put("    Backspace\tExit directory")
        put("    Space\tToggle preview")
        put("    Esc, q\tExit with cd")
        put("    Ctrl+c\tExit without cd")
        put("    /\tFuzzy search")
        put("    dd\tDelete file or dir")
        put("    y\tYank current directory path to clipboard")
        put("\n  Flags:\n")
        put("    --help\t-h\tdisplay help")
        put("    --version\t-v\tdisplay version")
        put("    --icons\t-i\tdisplay icons")
	put("    --command\t-c\t\"open\" file command line")
	put("         (path replaces first {}, else appended)")
        _ = w.Flush()
        _, _ = fmt.Fprintf(os.Stderr, "\n")
        os.Exit(1)
}

func version(s *walk.Styles) {
        fmt.Printf("\n  %s %s\n\n", s.Cursor.Render(" walk "), walk.Version())
        os.Exit(0)
}

func main() {
	style := walk.DefaultStyle()
	options := []walk.Option{
		walk.Style(style),
		walk.Size(80, 60),
	}

	startPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--help" || os.Args[i] == "-h" {
			usage(style)
		}

		if os.Args[i] == "--version" || os.Args[i] == "-v" {
			version(style)
		}

		if os.Args[i] == "--icons" || os.Args[i] == "-i" {
			options = append(options, walk.Icons())
			continue
		}

		const cmdflag = "--command"
		if strings.HasPrefix(os.Args[i], cmdflag + "=") {
			options = append(options, walk.Command(
				strings.TrimPrefix(os.Args[i], cmdflag + "="),
			))
			continue
		} else if os.Args[i] == cmdflag || os.Args[i] == "-c" {
			i++
			if i < len(os.Args) {
				options = append(options, walk.Command(os.Args[i]))
			}
			continue
		}

		startPath, err = filepath.Abs(os.Args[i])
		if err != nil {
			panic(err)
		}
		options = append(options, walk.Path(startPath))
	}

	output := termenv.NewOutput(os.Stderr)
	lipgloss.SetColorProfile(output.ColorProfile())

	w := walk.New(options...)
	p := tea.NewProgram(w, tea.WithOutput(os.Stderr))

	if _, err := p.Run(); err != nil {
		panic(err)
	}
	w.Exit()
}
