package main

import (
	configuration "camundaIncidentAggregator/pkg/config"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	APPNAME = "CamundaIncidentAggregator"
)

var (
	config configuration.Config
)

type model struct {
	status         []int
	text           []string
	err            []error
	totalRequests  int
	currentRequest int
	output         string
	spinner        spinner.Model
}

type statusMsg Response

type endMsg string

type errMsg struct{ error }

func (e errMsg) Error() string { return e.error.Error() }

func main() {
	var configurationError error
	config, configurationError = configuration.LoadConfig()
	if configurationError != nil {
		log.With(configurationError).Fatal("Error loading configuration file")
	}
	dirError := os.MkdirAll(config.LogPath, 0777)
	if dirError != nil {
		log.With(dirError).Error("Error creating logging directory")
	}
	logFile, fileError := os.OpenFile(config.LogPath+APPNAME+(time.Now()).Format("02-01-2006")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if fileError != nil {
		log.With(fileError).Error("Error writing to log file. (program will work but not as intended)")
	}
	log.SetOutput(logFile)
	m := model{}
	m.totalRequests = len(config.Camundas)
	log.SetLevel(log.ParseLevel(config.LogLevel))
	log.Debug(config)
	m.resetSpinner("69")
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.checkServer)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.resetSpinner("71")
			return m, tea.Quit
		default:
			return m, nil
		}

	case statusMsg:
		m.currentRequest += 1
		m.status = append(m.status, msg.StatusMsg)
		m.text = append(m.text, msg.objectBody.Fact)
		m.spinner.Update(msg)
		return m, tea.Batch(m.checkServer, m.spinner.Tick)

	case endMsg:
		m.resetSpinner("70")
		return m, nil

	case errMsg:
		m.err = append(m.err, msg)
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	default:
		m.spinner.Update(msg)
		return m, nil
	}
}

func (m model) View() string {
	for index := 0; index < m.currentRequest; index++ {
		m.output += fmt.Sprintf("Checking %s %s", config.Camundas[index].URL, m.spinner.View())
		if m.status != nil || m.text != nil || m.err != nil {
			if m.err != nil && m.err[index] != nil {
				m.output += fmt.Sprintf("something went wrong: %s", m.err[index])
			} else if m.status[index] != 0 {
				m.output += fmt.Sprintf("%d %s \n%s \n", m.status[index], http.StatusText(m.status[index]), m.text[index])
			}
		}
	}
	return m.output + "\n"
}

func (m model) checkServer() tea.Msg {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	log.Debug("doing api call")
	if m.currentRequest == m.totalRequests {
		return endMsg("")
	}
	res, err := c.Get(config.Camundas[m.currentRequest].URL)
	if err != nil {
		log.Error(err.Error())
		return errMsg{err}
	}
	defer res.Body.Close()
	log.Debug("reading body")
	body, err := io.ReadAll(res.Body)
	var catfacts CatFact
	json.Unmarshal(body, &catfacts)
	var response Response
	response.rawBody = string(body)
	response.StatusMsg = res.StatusCode
	response.objectBody = catfacts
	return statusMsg(response)

}

func (m *model) resetSpinner(color string) {
	m.spinner = spinner.New()
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	m.spinner.Spinner = spinner.Points
}

type CatFact struct {
	Fact   string
	Lenght int
}

type Response struct {
	StatusMsg  int
	rawBody    string
	objectBody CatFact
}
