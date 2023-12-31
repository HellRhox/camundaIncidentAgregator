package configuration

import (
	camunda "camundaIncidentAggregator/pkg/utils"
	"fmt"
	"github.com/spf13/viper"
)

const (
	MainConfigPath      = "./resources/config/"
	SecondaryConfigPath = "../resources/config/"
	TertiaryConfigPath  = "."
)

type Config struct {
	Camundas   []camunda.Camunda `mapstructue:"CAMUNDAS"`
	LogLevel   string            `mapstructure:"LOG_LEVEL"`
	LogPath    string            `mapstructure:"LOG_PATH"`
	ExportPath string            `mapstructure:"EXPORT_PATH"`
}

func (config Config) String() string {
	output := "\nConfig:\n" + "\tCamundas:\n"
	for index, element := range config.Camundas {
		output += "\t\t" + fmt.Sprintf("%d", index) + fmt.Sprintf("%+v", element) + "\n"
	}
	output += "\tLogLevel: " + config.LogLevel
	output += "\tLogPath: " + config.LogPath
	output += "\tExportPath: " + config.ExportPath
	return output
}

func LoadConfig(customDir string) (config Config, err error) {
	if customDir != "" {
		viper.AddConfigPath(customDir)
	} else {
		viper.AddConfigPath(MainConfigPath)
		viper.AddConfigPath(SecondaryConfigPath)
		viper.AddConfigPath(TertiaryConfigPath)
	}
	viper.SetConfigName("environment")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
