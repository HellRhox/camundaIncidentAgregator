package constants

import (
	configuration "camundaIncidentAggregator/pkg/config"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// P the current tea program
	P *tea.Program
	// WindowSize store the size of the terminal window
	WindowSize tea.WindowSizeMsg

	Config *configuration.Config
)

type ErrMsg struct {
	Error error
}

type Mode int

// constants for tracking the mode of the programm
const (
	TimeInput Mode = iota
	RestCalls
	Results
)

// AppPath for opening application in browser
var AppPath = "/camunda/app/cockpit/"

// DocStyle main style
var DocStyle = lipgloss.NewStyle().Margin(0, 2)

// HelpStyle styling for help context menu
var HelpStyle = lipgloss.NewStyle().Align(lipgloss.Bottom, lipgloss.Left).Foreground(lipgloss.Color("241")).Render

type keymap struct {
	Enter      key.Binding
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Back       key.Binding
	Quit       key.Binding
	OpenAsLink key.Binding
	Export     key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "select"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "j"),
		key.WithHelp("⇪/j", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "k"),
		key.WithHelp("⇩/k", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("⇦", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("⇨", "right"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc/⌫", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	OpenAsLink: key.NewBinding(
		key.WithKeys("ctrl+o"),
		key.WithHelp("ctrl+o", "open as link"),
	),

	Export: key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "export as .csv-file"),
	),
}
