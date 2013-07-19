package util

import (
	"github.com/msbranco/goconfig"
)

type Config interface {
	GetString(section string, option string) (string, error)
	GetBool(section string, option string) (bool, error)
	GetFloat(section string, option string) (float64, error)
	GetInt64(section string, option string) (int64, error)
}

func ReadConfigFile(configFileName string) *goconfig.ConfigFile {
	conf, err := goconfig.ReadConfigFile(configFileName)
	if err != nil {
		panic("Error reading config file")
	}
	return conf
}
