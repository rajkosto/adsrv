package util

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

func ConnectToDatabase(dbHost, dbUser, dbPass, dbBase string) *sql.DB {
	var con *sql.DB
	var err error
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
