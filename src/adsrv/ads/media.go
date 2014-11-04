package ads

import (
	"adsrv/msg"
)

type Crex struct {
	id          uint32
	minSize     uint16
	rotDuration uint16
	minAngleDeg uint8
}

func (m *Crex) Encode(w *msg.MessageWriter) error {
	if err := w.WriteInt(m.id); err != nil {
		return err
	}
	if err := w.WriteShort(m.minSize); err != nil {
		return err
	}
	if err := w.WriteShort(m.rotDuration); err != nil {
		return err
	}
	if err := w.WriteByte(m.minAngleDeg); err != nil {
		return err
	}
	return nil
}

func (m *Crex) Decode(r *msg.MessageReader) error {
	var err error
	if m.id, err = r.ReadInt(); err != nil {
		return err
	}
	if m.minSize, err = r.ReadShort(); err != nil {
		return err
	}
	if m.rotDuration, err = r.ReadShort(); err != nil {
		return err
	}
	if m.minAngleDeg, err = r.ReadByte(); err != nil {
		return err
	}
	return nil
}

//type 1
type MediaElem struct {
	id       uint32
	name     string
	unkShort uint16 //probably mimeType
	fileSize uint32
	fileMd5  string
	filePath string
	crex     Crex
}

func (m *MediaElem) Encode(w *msg.MessageWriter) error {
	if err := w.WriteInt(m.id); err != nil {
		return err
	}
	if err := w.WriteString(m.name); err != nil {
		return err
	}
	if err := w.WriteShort(m.unkShort); err != nil {
		return err
	}
	if err := w.WriteInt(m.fileSize); err != nil {
		return err
	}
	if err := w.WriteString(m.fileMd5); err != nil {
		return err
	}
	if err := w.WriteString(m.filePath); err != nil {
		return err
	}
	if err := m.crex.Encode(w); err != nil {
		return err
	}
	return nil
}

func (m *MediaElem) Decode(r *msg.MessageReader) error {
	var err error
	if m.id, err = r.ReadInt(); err != nil {
		return err
	}
	if m.name, err = r.ReadString(); err != nil {
		return err
	}
	if m.unkShort, err = r.ReadShort(); err != nil {
		return err
	}
	if m.fileSize, err = r.ReadInt(); err != nil {
		return err
	}
	if m.fileMd5, err = r.ReadString(); err != nil {
		return err
	}
	if m.filePath, err = r.ReadString(); err != nil {
		return err
	}
	if err = m.crex.Decode(r); err != nil {
		return err
	}
	return nil
}
