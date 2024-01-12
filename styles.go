package walk

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Warning, Preview, Cursor, Bar, Search, Danger lipgloss.Style
}

func NewStyles() *Styles { return new(Styles).Default() }

func (s *Styles) Default() *Styles {
	if s == nil {
		s = new(Styles)
	}
	s.Warning = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).PaddingLeft(1).PaddingRight(1)
	s.Preview = lipgloss.NewStyle().PaddingLeft(2)
	s.Cursor = lipgloss.NewStyle().Background(lipgloss.Color("#825DF2")).Foreground(lipgloss.Color("#FFFFFF"))
	s.Bar = lipgloss.NewStyle().Background(lipgloss.Color("#5C5C5C")).Foreground(lipgloss.Color("#FFFFFF"))
	s.Search = lipgloss.NewStyle().Background(lipgloss.Color("#499F1C")).Foreground(lipgloss.Color("#FFFFFF"))
	s.Danger = lipgloss.NewStyle().Background(lipgloss.Color("#FF0000")).Foreground(lipgloss.Color("#FFFFFF"))
	return s
}
