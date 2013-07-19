package handlers

import (
	"adsrv/msg"
	"adsrv/util"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type OpenSessionMsg struct {
	ver, sku           string
	sessionId, gamerId uint32
	etoken             []byte
	unkByte            byte
	uuid               string
}

func (m *OpenSessionMsg) Encode(w *msg.MessageWriter) error {
	if err := w.WriteString(m.ver); err != nil {
		return err
	}
	if err := w.WriteString(m.sku); err != nil {
		return err
	}
	if err := w.WriteInt(m.sessionId); err != nil {
		return err
	}
	if err := w.WriteInt(m.gamerId); err != nil {
		return err
	}
	if err := w.WriteBytes(m.etoken); err != nil {
		return err
	}
	if err := w.WriteByte(m.unkByte); err != nil {
		return err
	}
	if err := w.WriteString(m.uuid); err != nil {
		return err
	}
	return nil
}

func (m *OpenSessionMsg) Decode(r *msg.MessageReader) error {
	var err error
	if m.ver, err = r.ReadString(); err != nil {
		return err
	}
	if m.sku, err = r.ReadString(); err != nil {
		return err
	}
	if m.sessionId, err = r.ReadInt(); err != nil {
		return err
	}
	if m.gamerId, err = r.ReadInt(); err != nil {
		return err
	}
	if m.etoken, err = r.ReadBytes(); err != nil {
		return err
	}
	if m.unkByte, err = r.ReadByte(); err != nil {
		return err
	}
	if m.uuid, err = r.ReadString(); err != nil {
		return err
	}
	return nil
}

type SessionOpenMsg struct {
	sessionId, gamerId uint32
	baseTimeMs         uint64
}

func (m *SessionOpenMsg) Encode(w *msg.MessageWriter) error {
	if err := w.WriteInt(m.sessionId); err != nil {
		return err
	}
	if err := w.WriteInt(m.gamerId); err != nil {
		return err
	}
	if err := w.WriteLongLong(m.baseTimeMs); err != nil {
		return err
	}
	return nil
}

func (m *SessionOpenMsg) Decode(r *msg.MessageReader) error {
	var err error
	if m.sessionId, err = r.ReadInt(); err != nil {
		return err
	}
	if m.gamerId, err = r.ReadInt(); err != nil {
		return err
	}
	if m.baseTimeMs, err = r.ReadLongLong(); err != nil {
		return err
	}
	return nil
}

var rsaPrivKey *rsa.PrivateKey = LoadOrGenerateRsaKey("private.key", "public.key")

func OpenSessionHandler(conf util.Config, db *sql.DB, wr *msg.MessageWriter, rdr *msg.MessageReader, remoteAddr string) (statusCode int, tokenPtr *string, err error) {
	statusCode = http.StatusBadRequest

	reqMsg := OpenSessionMsg{}
	err = reqMsg.Decode(rdr)
	if err != nil {
		return
	}

	var realTokenBytes []byte
	realTokenBytes, err = rsa.DecryptOAEP(sha1.New(), rand.Reader, rsaPrivKey, reqMsg.etoken, nil)
	if err != nil {
		return
	}
	realToken := string(realTokenBytes)

	fmt.Printf("%s: Serving /adsrv/openSession with ver:%s sku:%s sessionId:%d gamerId:%d token:%s unkByte:%d uuid:%s\n", remoteAddr,
		reqMsg.ver, reqMsg.sku, reqMsg.sessionId, reqMsg.gamerId, realToken, uint32(reqMsg.unkByte), reqMsg.uuid)

	err = rdr.VerifyHash(&realToken)
	if err != nil {
		return
	}

	statusCode = http.StatusInternalServerError

	var sess AdSession
	var exists bool
	if reqMsg.sessionId != 0 && reqMsg.gamerId != 0 {
		sess, exists, err = GetSessionByIds(db, reqMsg.sessionId, reqMsg.gamerId)
	} else {
		reqMsg.gamerId, err = GetGamerIdForUuid(db, reqMsg.uuid)
		if err == nil {
			sess, exists, err = GetSessionByStrings(db, realToken, reqMsg.uuid)
		}
	}
	if err != nil {
		return
	}

	if exists == true {
		if realToken != sess.token || reqMsg.uuid != sess.uuid {
			statusCode = http.StatusBadRequest
			err = errors.New(fmt.Sprintf("Token/UUID mismatch for existing session %d", reqMsg.sessionId))
			fmt.Println("%s: %s", remoteAddr, err.Error())
			return
		}
	} else {
		var remoteIP net.IP
		{
			var remoteHost string
			remoteHost, _, err = net.SplitHostPort(remoteAddr)
			if err != nil {
				return
			}
			remoteIP = net.ParseIP(remoteHost)
			if remoteIP == nil {
				err = errors.New("Error parsing IP address: " + remoteHost)
				return
			}
		}

		//set the other members
		sess.gamerId = reqMsg.gamerId
		sess.token = realToken
		sess.uuid = reqMsg.uuid

		if err = InsertNewSession(db, &sess, remoteIP.String()); err != nil {
			return
		}
	}

	respMsg := SessionOpenMsg{sessionId: sess.sessionId, gamerId: sess.gamerId}
	err = respMsg.Encode(wr)
	if err != nil {
		return
	}

	fmt.Printf("%s: Sent response to /adsrv/openSession with sessionId:%d gamerId:%d\n", remoteAddr, respMsg.sessionId, respMsg.gamerId)
	statusCode, tokenPtr, err = http.StatusOK, &realToken, nil
	return
}
