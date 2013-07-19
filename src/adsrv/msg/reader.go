package msg

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"strings"
)

type MessageReader struct {
	io.Reader
	hs hash.Hash
}

func NewReader(r io.Reader) *MessageReader {
	out := new(MessageReader)
	out.hs = md5.New()
	out.Reader = io.TeeReader(r, out.hs)
	return out
}

func (r *MessageReader) ReadByte() (out uint8, err error) {
	err = binary.Read(r, binary.BigEndian, &out)
	return
}

func (r *MessageReader) ReadShort() (out uint16, err error) {
	err = binary.Read(r, binary.BigEndian, &out)
	return
}

func (r *MessageReader) ReadInt() (out uint32, err error) {
	err = binary.Read(r, binary.BigEndian, &out)
	return
}

func (r *MessageReader) ReadLongLong() (out uint64, err error) {
	err = binary.Read(r, binary.BigEndian, &out)
	return
}

func (r *MessageReader) ReadBytes() (out []byte, err error) {
	byteLen, err := r.ReadInt()
	if err != nil {
		return
	}

	out = make([]byte, byteLen)
	var bytesRead int
	bytesRead, err = r.Read(out)
	if bytesRead != int(byteLen) && err != nil {
		out = nil
		return
	}

	err = nil
	return
}

func (r *MessageReader) ReadString() (out string, err error) {
	slice, err := r.ReadBytes()
	if err != nil {
		return
	}

	out = string(slice)
	return
}

func (r *MessageReader) VerifyHash(sessionToken *string) error {
	if sessionToken != nil {
		_, err := strings.NewReader(*sessionToken).WriteTo(r.hs)
		if err != nil {
			return err
		}
	}
	ourHash := hex.EncodeToString(r.hs.Sum(nil))

	theirHash, err := r.ReadString()
	if err != nil {
		return err
	}

	hashLen := hex.EncodedLen(md5.Size)
	if len(theirHash) != hashLen {
		return errors.New("Invalid hash length")
	}

	if ourHash != theirHash {
		return errors.New("Message hash verification mismatch")
	}

	return nil
}
