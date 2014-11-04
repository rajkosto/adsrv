package ads

import (
	"adsrv/msg"
)

//type 3
type TargetElem struct {
	id     uint32
	name   string
	medias []uint32
}

func (m *TargetElem) Encode(w *msg.MessageWriter) error {
	if err := w.WriteInt(m.id); err != nil {
		return err
	}
	if err := w.WriteString(m.name); err != nil {
		return err
	}
	if err := w.WriteShort(uint16(len(m.medias))); err != nil {
		return err
	}
	for _, refId := range m.medias {
		if err := w.WriteInt(refId); err != nil {
			return err
		}
	}

	return nil
}

func (m *TargetElem) Decode(r *msg.MessageReader) error {
	var err error
	if m.id, err = r.ReadInt(); err != nil {
		return err
	}
	if m.name, err = r.ReadString(); err != nil {
		return err
	}
	var numElems uint16
	if numElems, err = r.ReadShort(); err != nil {
		return err
	}
	m.medias = make([]uint32, numElems)
	for i := range m.medias {
		if m.medias[i], err = r.ReadInt(); err != nil {
			return err
		}
	}
	return nil
}
