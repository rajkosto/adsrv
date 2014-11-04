// adsrv project main.go
package main

import (
	"adsrv/handlers"
	"adsrv/util"
	"fmt"
	"net/http"
)

func main() {
	config := util.ReadConfigFile("adsrv.conf")
	var dbHost, dbUser, dbPass, dbBase string
	var err error
	dbHost, err = config.GetString("database", "host")
	if err != nil {
		dbHost = "tcp(localhost:3306)"
	}
	dbUser, err = config.GetString("database", "username")
	if err != nil {
		dbUser = "root"
	}
	dbPass, err = config.GetString("database", "password")
	if err != nil {
		dbPass = ""
	}
	dbBase, err = config.GetString("database", "database")
	if err != nil {
		dbBase = "adsrv"
	}

	mainDb := util.ConnectToDatabase(dbHost, dbUser, dbPass, dbBase)

	http.HandleFunc("/adsrv/locateService", handlers.Make(handlers.LocateServiceHandler, config, mainDb))
	http.HandleFunc("/adsrv/openSession", handlers.Make(handlers.OpenSessionHandler, config, mainDb))
	http.HandleFunc("/adsrv/closeSession", handlers.Make(handlers.CloseSessionHandler, config, mainDb))
	http.HandleFunc("/adsrv/enterZone", handlers.Make(handlers.EnterZoneHandler, config, mainDb))

	var listenHost string
	var listenPort int64
	listenHost, err = config.GetString("host", "address")
	if err != nil {
		listenHost = "localhost"
	}
	listenPort, err = config.GetInt64("host", "port")
	if err != nil {
		listenPort = 8123
	}

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", listenHost, listenPort), nil)
	if err != nil {
		panic("Error listening: " + err.Error())
	}
}
