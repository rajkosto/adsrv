package handlers

import (
	"database/sql"
	"errors"
	"fmt"
)

type AdSession struct {
	sessionId, gamerId uint32
	token, uuid        string
	closed             bool
}

func GetGamerIdForUuid(db *sql.DB, uuid string) (gamerId uint32, err error) {
	var rows *sql.Rows
	var exists bool

	rows, err = db.Query("SELECT gamerId FROM gamers WHERE uuid = ?", uuid)
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
		res, err = db.Exec("INSERT INTO gamers (uuid) VALUES (?)", uuid)
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

func touchGameSession(db *sql.DB, sessionId, gamerId uint32) error {
	_, err := db.Exec("UPDATE sessions SET latest = CURRENT_TIMESTAMP WHERE sessionId = ? AND gamerId = ?", sessionId, gamerId)
	return err
}

func GetSessionByIds(db *sql.DB, sessionId, gamerId uint32) (sess AdSession, exists bool, err error) {
	var rows *sql.Rows
	rows, err = db.Query("SELECT token, uuid, (end IS NOT NULL) AS closed FROM sessions JOIN gamers ON sessions.gamerId = gamers.gamerId WHERE sessionId = ? AND sessions.gamerId = ?", sessionId, gamerId)
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&sess.token, &sess.uuid, &sess.closed)
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
		sess.sessionId = sessionId
		sess.gamerId = gamerId

		if sess.closed == false {
			if err = touchGameSession(db, sessionId, gamerId); err != nil {
				return
			}
		}
	}
	return
}

func GetSessionForConsumption(db *sql.DB, sessionId, gamerId uint32) (sess AdSession, err error) {
	var exists bool
	sess, exists, err = GetSessionByIds(db, sessionId, gamerId)
	if err != nil {
		return
	}
	if exists == false {
		err = errors.New(fmt.Sprintf("Session %d doesn't exist", sess.sessionId))
		return
	}
	if sess.closed == true {
		err = errors.New(fmt.Sprintf("Session %d is closed", sess.sessionId))
		return
	}
	err = nil
	return
}

func GetSessionByStrings(db *sql.DB, token, uuid string) (sess AdSession, exists bool, err error) {
	var rows *sql.Rows
	rows, err = db.Query("SELECT sessionId, sessions.gamerId, (end IS NOT NULL) as closed FROM sessions JOIN gamers ON sessions.gamerId = gamers.gamerId WHERE token = ? AND uuid = ?", token, uuid)
	defer rows.Close()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&sess.sessionId, &sess.gamerId, &sess.closed)
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
		sess.token = token
		sess.uuid = uuid

		if sess.closed == false {
			if err = touchGameSession(db, sess.sessionId, sess.gamerId); err != nil {
				return
			}
		}
	}
	return
}

func InsertNewSession(db *sql.DB, sess *AdSession, ipAddr string) error {
	res, err := db.Exec("INSERT INTO sessions (gamerId, token, ip) VALUES (?, ?, ?)", sess.gamerId, sess.token, ipAddr)
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

func CloseSession(db *sql.DB, sess *AdSession, timestamp uint64) error {
	if sess.closed == true {
		return nil
	}

	_, err := db.Exec("UPDATE sessions SET end = CURRENT_TIMESTAMP, durationMs = ? WHERE sessionId = ? AND gamerId = ?", timestamp, sess.sessionId, sess.gamerId)
	return err
}
