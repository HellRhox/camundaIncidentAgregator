package configuration

import (
	"fmt"
	"github.com/spf13/viper"
)

const CONFIG_PATH = "./resources/config/"

type Config struct {
	Camundas []Camunda `mapstructue:"CAMUNDAS"`
	LogLevel string    `mapstructure:"LOG_LEVEL"`
	LogPath  string    `mapstructure:"LOG_PATH"`
}

type Camunda struct {
	URL  string
	User string
	pw   string
}

func (camunda Camunda) String() string {
	return "{ URL:" + camunda.URL + " User:" + camunda.User + " pw" + camunda.pw + " }"
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
