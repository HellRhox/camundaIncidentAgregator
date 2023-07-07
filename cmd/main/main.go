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
	url     = "https://catfact.ninja/fact"
)

var (
	config configuration.Config
)

type model struct {
	status  int
	text    string
	err     error
	spinner spinner.Model
}

type statusMsg Response

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
	log.SetLevel(log.ParseLevel(config.LogLevel))
	log.Debug(config)
	m.resetSpinner("69")
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(checkServer, m.spinner.Tick)
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
		m.status = msg.StatusMsg
		m.text = msg.objectBody.Fact
		m.resetSpinner("70")
		return m, tea.Quit

	case errMsg:
		m.err = msg
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
	s := fmt.Sprintf("Checking %s %s%s", url, m.spinner.View(), " ")
	if m.err != nil {
		s += fmt.Sprintf("something went wrong: %s", m.err)
	} else if m.status != 0 {
		s += fmt.Sprintf("%d %s \n%s", m.status, http.StatusText(m.status), m.text)
	}
	return s + "\n"
}

func checkServer() tea.Msg {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	log.Debug("doing api call")
	res, err := c.Get(config.Camundas[0].URL)
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
