package tui

import (
	camunda "camundaIncidentAggregator/pkg/utils"
	"camundaIncidentAggregator/pkg/utils/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type RestModel struct {
	mode               constants.Mode
	list               list.Model
	spinners           spinner.Model
	day                int
	month              int
	responseSuccessful bool
	callstate          []int
	incidentCount      []int
	resolvedCount      []int
	autoRetires        int
	callStarted        bool
}

type responseMsg struct {
	success bool
}

func InitRest(day int, month int) RestModel {
	m := RestModel{mode: constants.RestCalls, day: day, month: month, callstate: make([]int, len(constants.Config.Camundas)), incidentCount: make([]int, len(constants.Config.Camundas)), resolvedCount: make([]int, len(constants.Config.Camundas))}
	m.spinners = m.resetSpinner("69")
	m.list = list.New(m.getItems(), list.NewDefaultDelegate(), 8, 8)
	m.list.Title = "Camundas"
	m.Update(m.list.StartSpinner())
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

func (m *RestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// defaulting to those so it get´s shown in the goLand execution view
		if msg.Width == 0 && msg.Height == 0 {
			msg.Width += 244
			msg.Height += 20
		}
		constants.WindowSize = msg
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	// managing key inputs
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, constants.Keymap.Back):
			timeModel, _ := initTimeFrame()
			return timeModel.Update(constants.WindowSize)
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinners, cmd = m.spinners.Update(msg)
		m.list.SetItems(m.getItems())
		cmds = append(cmds, cmd)
		newListModel, cmd := m.list.Update(msg)
		m.list = newListModel
		cmds = append(cmds, cmd)
		if !m.responseSuccessful && !m.callStarted {
			cmds = append(cmds, m.getCounts)
			m.callStarted = true
		}
	case responseMsg:
		m.responseSuccessful = msg.success
		if !m.responseSuccessful {
			m.callStarted = false
		}
		return m, nil
	}
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd, m.spinners.Tick)

	return m, tea.Batch(cmds...)
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

func (m *RestModel) getItems() []list.Item {
	items := make([]list.Item, len(constants.Config.Camundas))
	for i, item := range constants.Config.Camundas {
		var description string
		if !m.responseSuccessful {
			description = m.spinners.View()
		} else if m.callstate[i] == 200 {
			description = "Incidents Totel:" + strconv.Itoa(m.incidentCount[i]+m.resolvedCount[i]) + " Open:" + strconv.Itoa(m.incidentCount[i]) + " Resolved:" + strconv.Itoa(m.resolvedCount[i])
		} else if m.callstate[i] != 200 {
			description = "HTTP ERROR:" + strconv.Itoa(m.callstate[i])
		}
		items[i] = list.Item(listItem{title: item.URL, description: description})
	}
	return items
}

func (m RestModel) resetSpinner(color string) spinner.Model {
	m.spinners = spinner.New()
	m.spinners.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	m.spinners.Spinner = spinner.Points
	return m.spinners
}

func (m *RestModel) getCounts() tea.Msg {
	/* for showing purpose only
	for i := range constants.Config.Camundas {

		m.callstate[i] = randomCallstat()
		m.responses[i] = camunda.ListCountResponse{Count: 50}
	}
	time.Sleep(20)
	return responseMsg{success: rand.Intn(2) == 1}
	*/
	log.Debug("METHOD FOR REST CALLS STARTED")
	var restClient camunda.CamundaRest
	startDay := strings.Replace(time.Now().AddDate(0, -m.month, -m.day).Format("2006-01-02T15:04:05.000-0700"), "+", "-", 1)
	endDay := strings.Replace(time.Now().Format("2006-01-02T15:04:05.000-0700"), "+", "-", 1)
	success := false
	for i, entry := range constants.Config.Camundas {
		restClient = restClient.CreatClient(entry.URL, entry.User, entry.Password)
		err, currentIncidentResponse := restClient.GetListOfIncidentsCount(startDay, endDay)
		historicErr, historyIncidentsResponse := restClient.GetListOfHistoricIncidentsCount(startDay, endDay)
		if err != nil {
			log.With(err).Error("ERROR RETRIEVING CURRENT INCIDENT COUNT")
			if m.autoRetires >= 3 {
				m.autoRetires++
				log.Fatal("TO MANY AUTO-RETRIES RETRIEVING ACTIVE INCIDENTS")
			}
			return responseMsg{success: false}
		} else if historicErr != nil {
			log.With(historicErr).Error("ERROR RETRIEVING HISTORIC INCIDENT COUNT")
		} else {
			success = true
			m.callstate[i] = currentIncidentResponse.StatusCode
			m.incidentCount[i] = currentIncidentResponse.Count
			m.resolvedCount[i] = historyIncidentsResponse.Count
		}
		log.Debug("REST-CALLS SUCCESSFUL ")
	}

	return responseMsg{success: success}
}

func randomCallstat() int {
	callstateBool := rand.Intn(2) == 1
	if callstateBool {
		return 200
	} else {
		return 404
	}
}
