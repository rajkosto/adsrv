package ads

import (
	"adsrv/msg"
)

//type 3
type TargetElem struct {
	Id     uint32
	Name   string
	Medias []uint32
}

func (m *TargetElem) Encode(w *msg.MessageWriter) error {
	if err := w.WriteInt(m.Id); err != nil {
		return err
	}
	if err := w.WriteString(m.Name); err != nil {
		return err
	}
	if err := w.WriteShort(uint16(len(m.Medias))); err != nil {
		return err
	}
	for _, refId := range m.Medias {
		if err := w.WriteInt(refId); err != nil {
			return err
		}
	}

	return nil
}

func (m *TargetElem) Decode(r *msg.MessageReader) error {
	var err error
	if m.Id, err = r.ReadInt(); err != nil {
		return err
	}
	if m.Name, err = r.ReadString(); err != nil {
		return err
	}
	var numElems uint16
	if numElems, err = r.ReadShort(); err != nil {
		return err
	}
	m.Medias = make([]uint32, numElems)
	for i := range m.Medias {
		if m.Medias[i], err = r.ReadInt(); err != nil {
			return err
		}
	}
	return nil
}
