package ads

import (
	"adsrv/msg"
	"errors"
	"fmt"
)

type ZonePackage struct {
	Medias  []MediaElem
	Targets []TargetElem
	Zone    ZoneElem
}

const MEDIA_TAG_BYTE = 1
const TARGET_TAG_BYTE = 3
const ZONE_TAG_BYTE = 4
const FINISH_TAG_BYTE = 5

func (m *ZonePackage) Encode(w *msg.MessageWriter) error {
	for i := range m.Medias {
		if err := w.WriteByte(MEDIA_TAG_BYTE); err != nil {
			return err
		}
		if err := m.Medias[i].Encode(w); err != nil {
			return err
		}
	}
	for i := range m.Targets {
		if err := w.WriteByte(TARGET_TAG_BYTE); err != nil {
			return err
		}
		if err := m.Targets[i].Encode(w); err != nil {
			return err
		}
	}

	if err := w.WriteByte(ZONE_TAG_BYTE); err != nil {
		return err
	}
	if err := m.Zone.Encode(w); err != nil {
		return err
	}

	if err := w.WriteByte(FINISH_TAG_BYTE); err != nil {
		return err
	}

	return nil
}

func (m *ZonePackage) Decode(r *msg.MessageReader) error {
	m.Medias = make([]MediaElem, 0, 10)
	m.Targets = make([]TargetElem, 0, 40)
	m.Zone = ZoneElem{}

	var err error
infini:
	for {
		var tagType uint8
		tagType, err = r.ReadByte()
		if err != nil { //no more bytes to read
			break
		}
		switch tagType {
		case MEDIA_TAG_BYTE:
			var tmp MediaElem
			err = tmp.Decode(r)
			if err == nil {
				m.Medias = append(m.Medias, tmp)
			} else {
				return err
			}
		case TARGET_TAG_BYTE:
			var tmp TargetElem
			err = tmp.Decode(r)
			if err == nil {
				m.Targets = append(m.Targets, tmp)
			} else {
				return err
			}
		case ZONE_TAG_BYTE:
			if m.Zone.Id != 0 {
				return errors.New("Only one zone element allowed in package")
			} else {
				err = m.Zone.Decode(r)
				if err != nil {
					return err
				}
			}
		case FINISH_TAG_BYTE:
			break infini
		default:
			return errors.New(fmt.Sprintf("Unrecognized tag type: %d", uint32(tagType)))
		}
	}

	if m.Zone.Id == 0 {
		return err
	}

	return nil
}
