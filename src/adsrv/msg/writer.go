package msg

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"io"
	"strings"
)

type MessageWriter struct {
	io.Writer
	buf bytes.Buffer
}

func NewWriter() *MessageWriter {
	out := new(MessageWriter)
	out.buf = bytes.Buffer{}
	out.Writer = io.Writer(&out.buf)

	return out
}

func (w *MessageWriter) WriteByte(d uint8) error {
	err := binary.Write(w, binary.BigEndian, d)
	return err
}

func (w *MessageWriter) WriteShort(d uint16) error {
	err := binary.Write(w, binary.BigEndian, d)
	return err
}

func (w *MessageWriter) WriteInt(d uint32) error {
	err := binary.Write(w, binary.BigEndian, d)
	return err
}

func (w *MessageWriter) WriteLongLong(d uint64) error {
	err := binary.Write(w, binary.BigEndian, d)
	return err
}

func (w *MessageWriter) WriteBytes(b []byte) error {
	err := w.WriteInt(uint32(len(b)))
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func (w *MessageWriter) WriteString(s string) error {
	err := w.WriteInt(uint32(len(s)))
	if err != nil {
		return err
	}
	_, err = strings.NewReader(s).WriteTo(w)
	return err
}

func (w *MessageWriter) Finalize(sessionToken *string) (out []byte, err error) {
	if sessionToken != nil {
		hash := md5.New()
		_, err = hash.Write(w.buf.Bytes())
		if err != nil {
			return
		}
		_, err = strings.NewReader(*sessionToken).WriteTo(hash)
		if err != nil {
			return
		}
		err = w.WriteString(hex.EncodeToString(hash.Sum(nil)))
		if err != nil {
			return
		}
	}

	out = w.buf.Bytes()
	return
}
