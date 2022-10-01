package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
	Help   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select}, // first column
		{k.Help, k.Quit},         // second column
	}
}

var defaultKeyMap = keyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("⬆/k", "move cursor up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("⬇/j", "move cursor down"),
	),
	Select: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("⏎/spacebar", "de-/select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("ctrl + c/q", "quit program"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

type ScannerSelectionModel struct {
	availableScanners []string
	cursor            int
	selectedScanners  map[int]struct{}
	help              help.Model
	textStyle         lipgloss.Style
	cursorStyle       lipgloss.Style
	selectionStyle    lipgloss.Style
	qutting           bool
}

func NewModel() ScannerSelectionModel {
	return ScannerSelectionModel{
		availableScanners: []string{"Holy", "Moly"},
		selectedScanners:  make(map[int]struct{}),
		textStyle:         lipgloss.NewStyle().Foreground(lipgloss.Color("#5CBA45")),
		cursorStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#4784B3")),
		selectionStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("#1B5C29")),
		help:              help.New(),
	}
}

func (m ScannerSelectionModel) Init() tea.Cmd {
	return nil
}

func (m ScannerSelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, defaultKeyMap.Quit):
			m.qutting = true
			return m, tea.Quit
		case key.Matches(msg, defaultKeyMap.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, defaultKeyMap.Down):
			if m.cursor < len(m.availableScanners)-1 {
				m.cursor++
			}
		case key.Matches(msg, defaultKeyMap.Select):
			_, ok := m.selectedScanners[m.cursor]
			if ok {
				delete(m.selectedScanners, m.cursor)
			} else {
				m.selectedScanners[m.cursor] = struct{}{}
			}
		case key.Matches(msg, defaultKeyMap.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return m, nil
}

func (m ScannerSelectionModel) View() (s string) {
	if m.qutting {
		return m.textStyle.Render("It was grooty see ya!\n")
	}
	s = "Which scanners should be enabled?\n\n"
	for i, choice := range m.availableScanners {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selectedScanners[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", m.cursorStyle.Render(cursor), m.selectionStyle.Render(checked), choice)
	}
	helpView := m.help.View(defaultKeyMap)
	height := 8 - strings.Count(s, "\n") - strings.Count(helpView, "\n")
	return "\n" + m.textStyle.Render(s) + strings.Repeat("\n", height) + helpView
}
