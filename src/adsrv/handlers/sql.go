package handlers

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

func ConnectToDatabase() *sql.DB {
	var dbHost, dbUser, dbPass, dbBase string
	var err error
	dbHost, err = configFile.GetString("database", "host")
	if err != nil {
		dbHost = "tcp(localhost:3306)"
	}
	dbUser, err = configFile.GetString("database", "username")
	if err != nil {
		dbUser = "root"
	}
	dbPass, err = configFile.GetString("database", "password")
	if err != nil {
		dbPass = ""
	}
	dbBase, err = configFile.GetString("database", "database")
	if err != nil {
		dbBase = "adsrv"
	}

	var con *sql.DB
	if len(dbPass) > 0 {
		con, err = sql.Open("mysql", dbUser+":"+dbPass+"@"+dbHost+"/"+dbBase)
	} else {
		con, err = sql.Open("mysql", dbUser+"@"+dbHost+"/"+dbBase)
	}

	if err != nil {
		panic("Error creating database: " + err.Error())
	}

	err = con.Ping()
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}

	return con
}

var mainDatabase *sql.DB = ConnectToDatabase()
