package badge

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

var (
	ErrInvalidUsername = errors.New("usernames must be <= 12 bytes")
	ErrInvalidBadge    = errors.New("invalid badge")
)

func New(username []byte, id uint32, key []byte) ([]byte, error) {
	if len(username) > 12 {
		return nil, ErrInvalidUsername
	}

	raw_id := convert(id)

	username_padded := make([]byte, 12)

	for i := 0; i < 12; i++ {
		if i >= len(username) {
			username_padded[i] = 0x3d
			continue
		}
		username_padded[i] = username[i]
	}

	token := make([]byte, 72)
	token[17] = 0x2e
	token[27] = 0x2e

	base64.URLEncoding.Encode(token[0:16], username_padded)
	base64.URLEncoding.Encode(token[18:26], raw_id)

	h := hmac.New(sha256.New, key)
	h.Write(token[0:27])

	base64.URLEncoding.Encode(token[28:], h.Sum(nil))

	return token, nil
}

func Get(badge, key []byte) ([]byte, uint32, bool) {
	if !Check(badge, key) {
		return nil, 0, false
	}

	ue := make([]byte, 16)
	ie := make([]byte, 8)

	for i := 0; i < 16; i++ {
		ue[i] = badge[i]
	}

	for i := 18; i < 26; i++ {
		ie[i-18] = badge[i]
	}

	up := make([]byte, 12)
	id := make([]byte, 8)

	base64.URLEncoding.Decode(up, ue)
	base64.URLEncoding.Decode(id, ie)

	ul := 0

	for i := 0; i < 12; i++ {
		if up[i] == 0x3d {
			break
		}
		ul++
	}

	return up[:ul], deconvert(id), true
}

func Check(badge, key []byte) bool {
	if len(badge) != 72 {
		return false
	}

	v := make([]byte, 27)
	s := make([]byte, 46)
	vs := make([]byte, 46)

	for i := 0; i < 26; i++ {
		v[i] = badge[i]
	}

	for i := 28; i < 72; i++ {
		s[i-28] = badge[i]
	}

	h := hmac.New(sha256.New, key)
	h.Write(v)
	base64.URLEncoding.Encode(vs, h.Sum(nil))

	return bytes.Equal(s, vs)
}

func deconvert(id []byte) uint32 {
	var i uint32

	i |= uint32(id[3]) << 24
	i |= uint32(id[2]) << 16
	i |= uint32(id[1]) << 8
	i |= uint32(id[0])

	return i
}

func convert(id uint32) []byte {
	i := make([]byte, 4)

	i[3] = byte((id & 0xff000000) >> 24)
	i[2] = byte((id & 0x00ff0000) >> 16)
	i[1] = byte((id & 0x0000ff00) >> 8)
	i[0] = byte((id & 0x000000ff))

	return i
}
