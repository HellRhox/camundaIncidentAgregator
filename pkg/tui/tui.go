package tui

import (
	"camundaIncidentAggregator/pkg/utils/constants"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func StartTea() error {
	m, _ := initTimeFrame()
	constants.P = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := constants.P.Run(); err != nil {
		log.Fatal(err.Error())
	}
	return nil
}
