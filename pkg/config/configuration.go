package configuration

import (
	camunda "camundaIncidentAggregator/pkg/utils"
	"fmt"
	"github.com/spf13/viper"
)

const (
	CONFIG_PATH = "./resources/config/"
)

type Config struct {
	Camundas []camunda.Camunda `mapstructue:"CAMUNDAS"`
	LogLevel string            `mapstructure:"LOG_LEVEL"`
	LogPath  string            `mapstructure:"LOG_PATH"`
}

func (config Config) String() string {
	output := "\nConfig:\n" + "\tCamundas:\n"
	for index, element := range config.Camundas {
		output += "\t\t" + fmt.Sprintf("%d", index) + fmt.Sprintf("%+v", element) + "\n"
	}
	output += "\tLogLevel: " + config.LogLevel
	return output
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(CONFIG_PATH)
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
