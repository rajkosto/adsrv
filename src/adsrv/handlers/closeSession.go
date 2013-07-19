package handlers

import (
	"adsrv/msg"
	"adsrv/util"
	"database/sql"
	"fmt"
	"net/http"
)

type CloseSessionMsg struct {
	sessionId, gamerId uint32
	timestamp          uint64
}

func (m *CloseSessionMsg) Encode(w *msg.MessageWriter) error {
	if err := w.WriteInt(m.sessionId); err != nil {
		return err
	}
	if err := w.WriteInt(m.gamerId); err != nil {
		return err
	}
	if err := w.WriteLongLong(m.timestamp); err != nil {
		return err
	}
	return nil
}

func (m *CloseSessionMsg) Decode(r *msg.MessageReader) error {
	var err error
	if m.sessionId, err = r.ReadInt(); err != nil {
		return err
	}
	if m.gamerId, err = r.ReadInt(); err != nil {
		return err
	}
	if m.timestamp, err = r.ReadLongLong(); err != nil {
		return err
	}
	return nil
}

type SessionClosedMsg struct {
	sessionId, gamerId uint32
	message            string
}

func (m *SessionClosedMsg) Encode(w *msg.MessageWriter) error {
	if err := w.WriteInt(m.sessionId); err != nil {
		return err
	}
	if err := w.WriteInt(m.gamerId); err != nil {
		return err
	}
	if err := w.WriteString(m.message); err != nil {
		return err
	}
	return nil
}

func (m *SessionClosedMsg) Decode(r *msg.MessageReader) error {
	var err error
	if m.sessionId, err = r.ReadInt(); err != nil {
		return err
	}
	if m.gamerId, err = r.ReadInt(); err != nil {
		return err
	}
	if m.message, err = r.ReadString(); err != nil {
		return err
	}
	return nil
}

func CloseSessionHandler(conf util.Config, db *sql.DB, wr *msg.MessageWriter, rdr *msg.MessageReader, remoteAddr string) (statusCode int, tokenPtr *string, err error) {
	statusCode = http.StatusBadRequest

	reqMsg := CloseSessionMsg{}
	if err = reqMsg.Decode(rdr); err != nil {
		return
	}

	fmt.Printf("%s: Serving /adsrv/closeSession sessionId:%d gamerId:%d timestamp:%d\n", remoteAddr, reqMsg.sessionId, reqMsg.gamerId, reqMsg.timestamp)
	var sess AdSession
	if sess, err = GetSessionForConsumption(db, reqMsg.sessionId, reqMsg.gamerId); err != nil {
		return
	}

	if err = rdr.VerifyHash(&sess.token); err != nil {
		return
	}

	statusCode = http.StatusInternalServerError
	if err = CloseSession(db, &sess, reqMsg.timestamp); err != nil {
		return
	}

	const MASSIVE_SIGNOFF_MSG = "Thank you for playing Star Wars Galaxies!"
	respMsg := SessionClosedMsg{sessionId: sess.sessionId, gamerId: sess.gamerId, message: MASSIVE_SIGNOFF_MSG}
	if err = respMsg.Encode(wr); err != nil {
		return
	}

	fmt.Printf("%s: Sent response to /adsrv/closeSession sessionId:%d gamerId:%d\n", remoteAddr, respMsg.sessionId, respMsg.gamerId)
	statusCode, tokenPtr, err = http.StatusOK, &sess.token, nil
	return
}
