package utils

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func ReadConfig(paramsMap map[string]string) (map[string]interface{}, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	viper.AddConfigPath(wd)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		GenerateConfig(paramsMap)
		return nil, err
	}

	configMap := make(map[string]interface{})

	for key, value := range paramsMap {
		switch value {
		case "string":
			configMap[key] = viper.GetString(key)
		case "int":
			configMap[key] = viper.GetInt(key)
		case "bool":
			configMap[key] = viper.GetBool(key)
		case "uint64":
			configMap[key] = viper.GetUint64(key)
		case "int64":
			configMap[key] = viper.GetInt64(key)
		}
	}

	return configMap, nil
}

func GenerateConfig(paramsMap map[string]string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	for k, v := range paramsMap {
		viper.Set(k, v)
	}

	path := filepath.Join(wd, "config.yaml")
	return viper.WriteConfigAs(path)
}
