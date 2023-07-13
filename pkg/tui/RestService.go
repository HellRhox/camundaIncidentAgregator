package tui

import (
	"camundaIncidentAggregator/pkg/utils/constants"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type RestModel struct {
}

func InitRest(day int, month int) RestModel {
	m := RestModel{}
	return m
}

func (m RestModel) Init() tea.Cmd {
	return nil
}

func (m RestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Back):
			timeModel, _ := initTimeFrame()
			return timeModel.Update(constants.WindowSize)
		}
	}
	return m, nil
}

func (m RestModel) View() string {
	return ""
}
