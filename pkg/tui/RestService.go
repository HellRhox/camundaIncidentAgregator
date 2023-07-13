package tui

import (
	"camundaIncidentAggregator/pkg/utils/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RestModel struct {
	mode     constants.Mode
	list     list.Model
	spinners spinner.Model
	day      int
	month    int
}

func InitRest(day int, month int) RestModel {
	m := RestModel{mode: constants.RestCalls, day: day, month: month}
	m.resetSpinner("blue")
	if constants.WindowSize.Height != 0 {
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-1)
	}
	m.list = list.New(m.getItems(), list.NewDefaultDelegate(), 8, 8)
	m.list.Title = "Camundas"
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.Keymap.Up,
			constants.Keymap.Down,
			constants.Keymap.Enter,
			constants.Keymap.Back,
			constants.Keymap.Quit,
		}
	}
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

	return constants.DocStyle.Render(m.list.View() + "\n")
}

type listItem struct {
	title       string
	description string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.description }
func (i listItem) FilterValue() string { return i.title }

func (m RestModel) getItems() []list.Item {
	items := make([]list.Item, len(constants.Config.Camundas))
	for i, item := range constants.Config.Camundas {
		items[i] = list.Item(listItem{title: item.URL, description: m.spinners.View()})
	}
	return items
}

func (m RestModel) resetSpinner(color string) {
	m.spinners = spinner.New()
	m.spinners.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	m.spinners.Spinner = spinner.Points
}
