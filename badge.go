package badge

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

var (
	ErrInvalidUsername = errors.New("usernames must be <= 255 bytes")
	ErrInvalidBadge    = errors.New("invalid badge")
	ErrInvalidKey      = errors.New("invalid key")
)

func New(username []byte, id uint32, key []byte) ([]byte, error) {
	if len(username) > 255 {
		return nil, ErrInvalidUsername
	}

	badge := make([]byte, (len(username))+53)
	badge[0] = uint8(len(username))

	for i := 1; i <= len(username); i++ {
		badge[i] = username[i-1]
	}

	idb := make([]byte, 4)

	idb[3] = byte((id & 0xff000000) >> 24)
	idb[2] = byte((id & 0x00ff0000) >> 16)
	idb[1] = byte((id & 0x0000ff00) >> 8)
	idb[0] = byte((id & 0x000000ff))

	hex.Encode(badge[1+(len(username)):], idb)

	h := hmac.New(sha256.New, key)
	h.Write(badge[:1+len(username)+8])

	base64.URLEncoding.Encode(badge[1+len(username)+8:], h.Sum(nil))

	return badge, nil
}

func Get(badge []byte, key []byte) ([]byte, uint32, error) {
	l := uint8(badge[0])
	if l > 255 || len(badge) != (int(l)+53) {
		return nil, 0, ErrInvalidBadge
	}

	username := make([]byte, l)
	for a := 0; a < int(l); a++ {
		username[a] = badge[a+1]
	}

	idb := make([]byte, 4)
	hex.Decode(idb, badge[1+l:(1+l)+8])

	var i uint32

	i |= uint32(idb[3]) << 24
	i |= uint32(idb[2]) << 16
	i |= uint32(idb[1]) << 8
	i |= uint32(idb[0])

	h := hmac.New(sha256.New, key)
	h.Write(badge[:9+l])

	t := make([]byte, 44)

	base64.URLEncoding.Encode(t, h.Sum(nil))

	if !bytes.Equal(t, badge[9+l:]) {
		return nil, 0, ErrInvalidKey
	}

	return username, i, nil
}
