package configuration

import (
	camunda "camundaIncidentAggregator/pkg/utils"
	"camundaIncidentAggregator/pkg/utils/constants"
	"fmt"
	"github.com/spf13/viper"
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
	viper.AddConfigPath(constants.CONFIG_PATH)
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
