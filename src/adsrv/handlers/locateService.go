package handlers

import (
	"fmt"
	"adsrv/msg"
	"net/http"
)

type LocateServiceMsg struct {
	ver, sku string
}

func (m *LocateServiceMsg) Encode(w *msg.MessageWriter) error {
	err := w.WriteString(m.ver)
	if err != nil {
		return err
	}
	err = w.WriteString(m.sku)
	return err
}

func (m *LocateServiceMsg) Decode(r *msg.MessageReader) error {
	var err error
	m.ver, err = r.ReadString()
	if err != nil {
		return err
	}
	m.sku, err = r.ReadString()
	return err
}

type ServiceLocationMsg struct {
	zoneSrv, impSrv, mediaSrv string
}

func (m *ServiceLocationMsg) Encode(w *msg.MessageWriter) error {
	err := w.WriteString(m.zoneSrv)
	if err != nil {
		return err
	}
	err = w.WriteString(m.impSrv)
	if err != nil {
		return err
	}
	err = w.WriteString(m.mediaSrv)
	return err
}

func (m *ServiceLocationMsg) Decode(r *msg.MessageReader) error {
	var err error
	m.zoneSrv, err = r.ReadString()
	if err != nil {
		return err
	}
	m.impSrv, err = r.ReadString()
	if err != nil {
		return err
	}
	m.mediaSrv, err = r.ReadString()
	return err
}

func LocateServiceHandler(wr *msg.MessageWriter, rdr *msg.MessageReader, remoteAddr string) (statusCode int, tokenPtr *string, err error) {
	reqMsg := LocateServiceMsg{}
	err = reqMsg.Decode(rdr)
	if err != nil {
		statusCode = http.StatusBadRequest
		return
	}

	fmt.Printf("%s: Serving /adsrv/locateService with ver:%s sku:%s\n", remoteAddr, reqMsg.ver, reqMsg.sku)

	respMsg := ServiceLocationMsg{}
	respMsg.zoneSrv, err = configFile.GetString("servers", "zone")
	if err != nil {
		respMsg.zoneSrv = "DISABLED"
	}
	respMsg.impSrv, err = configFile.GetString("servers", "impression")
	if err != nil {
		respMsg.impSrv = "DISABLED"
	}
	respMsg.mediaSrv, err = configFile.GetString("servers", "media")
	if err != nil {
		respMsg.mediaSrv = "DISABLED"
	}

	err = respMsg.Encode(wr)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	fmt.Printf("%s: Sent response to /adsrv/locateService with zone:%s imp:%s media:%s\n", remoteAddr, respMsg.zoneSrv, respMsg.impSrv, respMsg.mediaSrv)
	statusCode, err = http.StatusOK, nil
	return
}
