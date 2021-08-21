package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func GenerateConfigReader(defaults map[string]interface{}, configName string, configPaths []string) *viper.Viper {
	wiper := viper.New()

	wiper.SetConfigType("yaml")

	for k, v := range defaults {
		wiper.SetDefault(k, v)
	}

	var customConfigFilePath string
	pflag.StringVar(&customConfigFilePath, "config", "", "Custom config file path")
	pflag.Parse()

	if len(customConfigFilePath) > 0 {
		wiper.SetConfigFile(customConfigFilePath)
	} else {
		wiper.SetConfigName(configName)

		for _, path := range configPaths {
			wiper.AddConfigPath(path)
		}
	}

	return wiper
}
