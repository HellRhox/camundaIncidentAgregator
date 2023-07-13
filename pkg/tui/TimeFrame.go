package tui

import (
	"camundaIncidentAggregator/pkg/utils/constants"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"strings"
)

var (
	choices = []string{"Day", "Week", "Month"}
)

type changeToRestMsg struct {
	day   int
	month int
}

type TimeModel struct {
	mode   constants.Mode
	cursor int
}

func (m TimeModel) Init() tea.Cmd {
	return nil
}

func initTimeFrame() (tea.Model, tea.Cmd) {
	m := TimeModel{mode: constants.TimeInput, cursor: 0}

	return m, func() tea.Msg {
		return constants.ErrMsg{}
	}
}

func (m TimeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Enter):
			var muleMsg changeToRestMsg
			if m.cursor == 0 {
				muleMsg = changeToRestMsg{day: 1}
			} else if m.cursor == 1 {
				muleMsg = changeToRestMsg{day: 7}
			} else if m.cursor == 2 {
				muleMsg = changeToRestMsg{month: 1}
			}
			log.Debug("User choose:" + choices[m.cursor])
			return m.Update(muleMsg)

		case key.Matches(msg, constants.Keymap.Down):
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}

		case key.Matches(msg, constants.Keymap.Up):
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		case key.Matches(msg, constants.Keymap.Back):
			//Implement easter egg

		default:
			return m, nil
		}
	case changeToRestMsg:
		restService := InitRest(msg.day, msg.month)
		return restService.Update(constants.WindowSize)

	}
	return m, nil
}

func (m TimeModel) View() string {
	s := strings.Builder{}
	s.WriteString("What timeframe do u want to query?\n\n")
	output := ""
	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	output += s.String()

	return lipgloss.JoinVertical(lipgloss.Left, "\n", output, m.helpView())
}

func (m TimeModel) helpView() string {
	var keyBindings string
	for _, keyBinding := range getKeybindings() {
		keyBindings += keyBinding.Help().Desc + ":" + keyBinding.Help().Key + " "
	}
	return constants.HelpStyle(keyBindings)
}

func getKeybindings() []key.Binding {
	return []key.Binding{
		constants.Keymap.Up,
		constants.Keymap.Down,
		constants.Keymap.Enter,
		constants.Keymap.Quit,
	}
}
