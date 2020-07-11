package handlers

import (
	"adsrv/ads"
	"adsrv/msg"
	"adsrv/util"
	"database/sql"
	"fmt"
	"net/http"
)

type EnterZoneMsg struct {
	sessionId, gamerId uint32
	zoneName           string
	timestamp          uint64
}

func (m *EnterZoneMsg) Encode(w *msg.MessageWriter) error {
	if err := w.WriteInt(m.sessionId); err != nil {
		return err
	}
	if err := w.WriteInt(m.gamerId); err != nil {
		return err
	}
	if err := w.WriteString(m.zoneName); err != nil {
		return err
	}
	if err := w.WriteLongLong(m.timestamp); err != nil {
		return err
	}
	return nil
}

func (m *EnterZoneMsg) Decode(r *msg.MessageReader) error {
	var err error
	if m.sessionId, err = r.ReadInt(); err != nil {
		return err
	}
	if m.gamerId, err = r.ReadInt(); err != nil {
		return err
	}
	if m.zoneName, err = r.ReadString(); err != nil {
		return err
	}
	if m.timestamp, err = r.ReadLongLong(); err != nil {
		return err
	}
	return nil
}

func EnterZoneHandler(conf util.Config, db *sql.DB, wr *msg.MessageWriter, rdr *msg.MessageReader, remoteAddr string) (statusCode int, tokenPtr *string, err error) {
	statusCode = http.StatusBadRequest

	reqMsg := EnterZoneMsg{}
	if err = reqMsg.Decode(rdr); err != nil {
		return
	}

	fmt.Printf("%s: Serving /adsrv/enterZone sessionId:%d gamerId:%d zoneName:%s timestamp:%d\n", remoteAddr, reqMsg.sessionId, reqMsg.gamerId, reqMsg.zoneName, reqMsg.timestamp)
	var sess AdSession
	if sess, err = GetSessionForConsumption(db, reqMsg.sessionId, reqMsg.gamerId); err != nil {
		return
	}

	if err = rdr.VerifyHash(&sess.token); err != nil {
		return
	}

	statusCode = http.StatusInternalServerError

	respMsg := ads.ZonePackage{}
	if respMsg.Zone, err = ads.GetZoneByName(db, reqMsg.zoneName); err != nil {
		return
	}
	if err = ads.AddZoneVisit(db, sess.sessionId, &respMsg.Zone, reqMsg.timestamp); err != nil {
		return
	}
	if respMsg.Targets, err = ads.GetTargetsByZoneId(db, respMsg.Zone.Id); err != nil {
		return
	}
	var mediaIds map[uint32]bool = make(map[uint32]bool)
	for _, target := range respMsg.Targets {
		for _, mediaId := range target.Medias {
			mediaIds[mediaId] = true
		}
	}
	respMsg.Medias = make([]ads.MediaElem, 0, len(mediaIds))
	for k, v := range mediaIds {
		if v != true {
			continue
		}

		var media ads.MediaElem
		if media, err = ads.GetMediaById(db, k); err != nil	{
			return
		}
		respMsg.Medias = append(respMsg.Medias, media)
	}

	if err = respMsg.Encode(wr); err != nil {
		return
	}

	fmt.Printf("%s: Sent response to /adsrv/enterZone id:%d name:%s\n", remoteAddr, respMsg.Zone.Id, respMsg.Zone.Name)
	statusCode, tokenPtr, err = http.StatusOK, &sess.token, nil
	return
}
