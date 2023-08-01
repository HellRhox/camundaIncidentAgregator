package tui

import (
	camunda "camundaIncidentAggregator/pkg/utils"
	"camundaIncidentAggregator/pkg/utils/constants"
	"camundaIncidentAggregator/pkg/utils/timeFormat"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"strconv"
)

type detailedSingleView struct {
	mode               constants.Mode
	paging             paginator.Model
	list               list.Model
	spinner            spinner.Model
	indexOfElement     int
	day                int
	month              int
	responseSuccessful bool
	callstate          []int
	currentDetails     map[string][]camunda.ListResponseEntre
	historicDetails    map[string][]camunda.ListResponseEntre
	keys               []string
	autoRetires        int
	callStarted        bool
}

func initDetailView(day int, month int, index int) detailedSingleView {
	m := detailedSingleView{}
	var totalPage int
	m.day = day
	m.month = month
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Hamburger
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 1
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(totalPage)
	m.paging = p
	m.list = list.New(nil, list.NewDefaultDelegate(), 8, 8)
	m.list.SetShowHelp(false)
	m.Update(m.list.StartSpinner())
	m.indexOfElement = index
	m.responseSuccessful = false

	m.mode = constants.Results
	return m
}

type responseDetailMsg struct {
	success bool
}

func (m *detailedSingleView) Init() tea.Cmd {
	return nil
}

func (m *detailedSingleView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// defaulting to those, so it get´s shown in the goLand execution view
		if msg.Width == 0 && msg.Height == 0 {
			msg.Width += 244
			msg.Height += 20
		}
		constants.WindowSize = msg
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
		cmds = append(cmds, m.spinner.Tick, m.getDetails)
	// managing key inputs
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Back):
			restService := InitRest(m.day, m.month)
			return restService.Update(constants.WindowSize)
		case key.Matches(msg, m.paging.KeyMap.PrevPage):
			cmds = append(cmds, m.updateList()...)
			break
		case key.Matches(msg, m.paging.KeyMap.NextPage):
			cmds = append(cmds, m.updateList()...)
			break
		}
	case responseDetailMsg:
		m.responseSuccessful = msg.success
		if !m.responseSuccessful {
			m.callStarted = false
		}
		m.paging.TotalPages = getTotalPagesFromMaximum(len(m.currentDetails), len(m.historicDetails))
		m.updateList()
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.paging, cmd = m.paging.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *detailedSingleView) View() string {
	var output string
	if m.responseSuccessful {
		listString := m.list.View()
		helpString := lipgloss.NewStyle().Margin(0, 2).AlignVertical(lipgloss.Bottom).Align(lipgloss.Left).Render(helpView())
		pagingString := lipgloss.NewStyle().AlignVertical(lipgloss.Bottom).PaddingLeft(constants.WindowSize.Width - len(helpString)).Align(lipgloss.Right).Render(m.paging.View())
		output += lipgloss.JoinVertical(lipgloss.Left, listString, lipgloss.PlaceVertical(constants.WindowSize.Height-50, lipgloss.Bottom, helpString+pagingString))

	} else {
		output += lipgloss.Place(constants.WindowSize.Width, constants.WindowSize.Height, lipgloss.Center, lipgloss.Center, "Loading Details\n"+m.spinner.View())
	}
	return output
}

func (m *detailedSingleView) getItems(index int) ([]list.Item, []tea.Cmd) {
	items := make([]list.Item, 3)
	var cmds []tea.Cmd
	total := len(m.currentDetails[m.keys[index]]) + len(m.historicDetails[m.keys[index]])
	if m.responseSuccessful {
		item := detailedViewListItem{title: "Total", description: strconv.Itoa(total), progress: progress.New()}
		item.progress.ShowPercentage = false
		cmds = append(cmds, item.progress.SetPercent(1))
		items[0] = list.Item(item)
		item = detailedViewListItem{title: "Opend", description: strconv.Itoa(len(m.currentDetails[m.keys[index]])), progress: progress.New()}
		item.progress.ShowPercentage = false
		currentPercent := (1 / float64(total)) * float64(len(m.currentDetails[m.keys[index]]))
		cmds = append(cmds, item.progress.SetPercent(currentPercent))
		items[1] = list.Item(item)
		item = detailedViewListItem{title: "Resolved", description: strconv.Itoa(len(m.historicDetails[m.keys[index]])), progress: progress.New()}
		item.progress.ShowPercentage = false
		historicPercent := (1 / float64(total)) * float64(len(m.historicDetails[m.keys[index]]))
		cmds = append(cmds, item.progress.SetPercent(historicPercent))
		items[2] = list.Item(item)
	}
	return items, cmds
}

func (m *detailedSingleView) getDetailsCurrent() camunda.ListResponse {
	var restClient camunda.CamundaRest
	restClient = restClient.CreatClient(constants.Config.Camundas[m.indexOfElement].URL, constants.Config.Camundas[m.indexOfElement].User, constants.Config.Camundas[m.indexOfElement].Password)
	log.Debug("Preparing rest call for current incidents")
	err, currentDetails := restClient.GetListOfIncidents(timeFormat.GetTimeFormatForRest(m.month, m.day))
	if err != nil {
		log.With(err).Fatal("COULD NOT GET CURRENT DETAILS FOR INCIDENTS")
	}
	return currentDetails
}

func (m *detailedSingleView) getDetails() tea.Msg {
	m.callStarted = true
	m.currentDetails = aggregateData(m.getDetailsCurrent())
	m.historicDetails = aggregateData(m.getDetailsHistoric())
	m.getKeys()
	return responseDetailMsg{success: true}
}

func (m *detailedSingleView) getDetailsHistoric() camunda.ListResponse {
	var restClient camunda.CamundaRest
	restClient = restClient.CreatClient(constants.Config.Camundas[m.indexOfElement].URL, constants.Config.Camundas[m.indexOfElement].User, constants.Config.Camundas[m.indexOfElement].Password)
	log.Debug("Preparing rest call for historic incidents")
	err, historicDetails := restClient.GetListOfHistoricIncidents(timeFormat.GetTimeFormatForRest(m.month, m.day))
	if err != nil {
		log.With(err).Fatal("COULD NOT GET HISTORIC DETAILS FOR INCIDENTS")
	}
	return historicDetails
}

func (m *detailedSingleView) getKeys() {
	for mapKey, _ := range m.currentDetails {
		m.keys = append(m.keys, mapKey)
	}
	for mapKey, _ := range m.historicDetails {
		if contains(m.keys, mapKey) {
			continue
		} else {
			m.keys = append(m.keys, mapKey)
		}
	}

}

func (m *detailedSingleView) updateList() []tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	if m.responseSuccessful {
		index := m.paging.Page
		m.list.Title = m.keys[index]
		items, returnedCmds := m.getItems(index)
		cmds = append(cmds, returnedCmds...)
		m.list.SetItems(items)
		m.list, cmd = m.list.Update(constants.WindowSize)
		cmds = append(cmds, cmd)
	}
	return cmds
}

func aggregateData(entities camunda.ListResponse) map[string][]camunda.ListResponseEntre {
	returnMap := make(map[string][]camunda.ListResponseEntre)
	log.Debug("Generating sorted map with list")
	for _, entry := range entities {
		if _, ok := returnMap[entry.ProcessDefinitionId]; ok {
			returnMap[entry.ProcessDefinitionId] = append(returnMap[entry.ProcessDefinitionId], entry)
		} else {
			returnMap[entry.ProcessDefinitionId] = append(make([]camunda.ListResponseEntre, 1), entry)
		}
	}
	return returnMap
}

func helpView() string {
	var keyBindings string
	for _, keyBinding := range getKeybindingsDetail() {
		keyBindings += keyBinding.Help().Desc + ":" + keyBinding.Help().Key + " "
	}
	return constants.HelpStyle(keyBindings)
}

func getKeybindingsDetail() []key.Binding {
	constants.Keymap.Left.SetHelp("⇦", "prev. page")
	constants.Keymap.Right.SetHelp("⇨", "next page")
	return []key.Binding{
		constants.Keymap.Up,
		constants.Keymap.Down,
		constants.Keymap.Left,
		constants.Keymap.Right,
		constants.Keymap.Enter,
		constants.Keymap.Back,
		constants.Keymap.Quit,
	}
}
func getTotalPagesFromMaximum(firstVal int, secondVal int) int {
	if firstVal >= secondVal {
		return firstVal
	} else {
		return secondVal
	}
}

func contains(slice []string, searched string) bool {
	for _, value := range slice {
		if value == searched {
			return true
		}
	}
	return false
}

type detailedViewListItem struct {
	title       string
	description string
	progress    progress.Model
}

func (i detailedViewListItem) Title() string { return i.title }
func (i detailedViewListItem) Description() string {
	return i.progress.ViewAs(i.progress.Percent()) + " " + i.description
}
func (i detailedViewListItem) FilterValue() string { return i.title }
