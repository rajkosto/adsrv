package ads

import (
	"adsrv/msg"
)

//type 4
type ZoneElem struct {
	Id      uint32
	Name    string
	Targets []uint32
}

func (m *ZoneElem) Encode(w *msg.MessageWriter) error {
	if err := w.WriteInt(m.Id); err != nil {
		return err
	}
	if err := w.WriteString(m.Name); err != nil {
		return err
	}
	if err := w.WriteShort(uint16(len(m.Targets))); err != nil {
		return err
	}
	for _, targetId := range m.Targets {
		if err := w.WriteInt(targetId); err != nil {
			return err
		}
	}

	return nil
}

func (m *ZoneElem) Decode(r *msg.MessageReader) error {
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
	m.Targets = make([]uint32, numElems)
	for i := range m.Targets {
		if m.Targets[i], err = r.ReadInt(); err != nil {
			return err
		}
	}
	return nil
}
