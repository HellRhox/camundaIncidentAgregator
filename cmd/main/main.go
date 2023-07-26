package main

import (
	configuration "camundaIncidentAggregator/pkg/config"
	"camundaIncidentAggregator/pkg/tui"
	"camundaIncidentAggregator/pkg/utils/constants"
	"github.com/charmbracelet/log"
	"os"
	"time"
)

func main() {
	var configurationError error
	config, configurationError := configuration.LoadConfig()
	constants.Config = &config
	if configurationError != nil {
		log.With(configurationError).Fatal("Error loading configuration file")
	}
	dirError := os.MkdirAll(config.LogPath, 0777)
	if dirError != nil {
		log.With(dirError).Error("Error creating logging directory")
	}
	logFile, fileError := os.OpenFile(config.LogPath+constants.APPNAME+(time.Now()).Format("02-01-2006")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if fileError != nil {
		log.With(fileError).Error("Error writing to log file. (program will work but not as intended)")
	}
	log.SetOutput(logFile)
	log.SetLevel(log.ParseLevel(config.LogLevel))
	log.Debug(config)
	tui.StartTea()
}
