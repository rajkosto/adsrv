package handlers

import (
	"database/sql"
)

type AdSession struct {
	sessionId, gamerId uint32
	token, uuid        string
}

func GetGamerIdForUuid(uuid string) (gamerId uint32, err error) {
	var rows *sql.Rows
	var exists bool

	rows, err = mainDatabase.Query("SELECT gamerId FROM gamers WHERE uuid = ?", uuid)
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&gamerId)
		if err != nil {
			return
		}
		exists = true
	}
	err = rows.Err()
	if err != nil {
		return
	}
	rows.Close()

	if exists == false {
		var res sql.Result
		res, err = mainDatabase.Exec("INSERT INTO gamers (uuid) VALUES (?)", uuid)
		if err != nil {
			return
		}

		var bigId int64
		bigId, err = res.LastInsertId()
		if err != nil {
			return
		}

		gamerId = uint32(bigId)
	}

	return
}

func touchGameSession(sessionId, gamerId uint32) error {
	_, err := mainDatabase.Exec("UPDATE sessions SET latest = CURRENT_TIMESTAMP WHERE sessionId = ? AND gamerId = ?", sessionId, gamerId)
	return err
}

func GetSessionByIds(sessionId, gamerId uint32) (sess AdSession, exists bool, err error) {
	var rows *sql.Rows
	rows, err = mainDatabase.Query("SELECT token, uuid FROM sessions JOIN gamers ON sessions.gamerId = gamers.gamerId WHERE sessions.sessionId = ? AND sessions.gamerId = ?", sessionId, gamerId)
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&sess.token, &sess.uuid)
		if err != nil {
			return
		}
		exists = true
	}
	err = rows.Err()
	if err != nil {
		return
	}
	rows.Close()

	if exists == true {
		if err = touchGameSession(sessionId, gamerId); err != nil {
			return
		}

		sess.sessionId = sessionId
		sess.gamerId = gamerId
	}

	return
}

func GetSessionByStrings(token, uuid string) (sess AdSession, exists bool, err error) {
	var rows *sql.Rows
	rows, err = mainDatabase.Query("SELECT sessions.sessionId, sessions.gamerId FROM sessions JOIN gamers ON sessions.gamerId = gamers.gamerId WHERE token = ? AND uuid = ?", token, uuid)
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&sess.sessionId, &sess.gamerId)
		if err != nil {
			return
		}
		exists = true
	}
	err = rows.Err()
	if err != nil {
		return
	}
	rows.Close()

	if exists == true {
		if err = touchGameSession(sess.sessionId, sess.gamerId); err != nil {
			return
		}

		sess.token = token
		sess.uuid = uuid
	}

	return
}

func InsertNewSession(sess *AdSession, ipAddr string) error {
	res, err := mainDatabase.Exec("INSERT INTO sessions (gamerId, token, ip) VALUES (?, ?, ?)", sess.gamerId, sess.token, ipAddr)
	if err != nil {
		return err
	}

	var bigId int64
	bigId, err = res.LastInsertId()
	if err != nil {
		return err
	}

	sess.sessionId = uint32(bigId)
	return nil
}
