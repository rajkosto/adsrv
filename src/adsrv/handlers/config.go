package handlers

import (
	"github.com/msbranco/goconfig"
)

func ReadConfigFile() *goconfig.ConfigFile {
	conf, err := goconfig.ReadConfigFile("adsrv.conf")
	if err != nil {
		panic("Error reading config file")
	}
	return conf
}

var configFile *goconfig.ConfigFile = ReadConfigFile()
