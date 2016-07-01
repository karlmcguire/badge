package badge

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	type user struct {
		username []byte
		id       uint32
		err      error
	}

	newTests := []user{
		user{
			[]byte("joe"),
			uint32(0),
			nil,
		},
		user{
			[]byte("....................................................................................................................................................................................................................................................................................................."),
			uint32(0),
			ErrInvalidUsername,
		},
		user{
			[]byte(""),
			uint32(0),
			ErrInvalidUsername,
		},
		user{
			nil,
			uint32(0),
			ErrInvalidUsername,
		},
	}

	for _, v := range newTests {
		_, err := New(v.username, v.id, []byte("key"))
		if err != v.err {
			t.Fatal("unexpected error")
		}
	}
}

func TestGet(t *testing.T) {
	type user struct {
		username []byte
		id       uint32
		key      []byte
		auth     bool
	}

	goodKey := []byte("goodKey")
	badKey := []byte("badKey")

	getTests := []user{
		user{
			username: []byte("joe"),
			id:       uint32(0),
			key:      goodKey,
			auth:     true,
		},
		user{
			username: []byte("joe"),
			id:       uint32(0),
			key:      badKey,
			auth:     false,
		},
	}

	scratchBadge := make([]byte, 0)

	for _, v := range getTests {
		b, err := New(v.username, v.id, v.key)
		if err != nil {
			t.Fatal(err)
		}

		username, id, auth := Get(b, goodKey)
		if v.auth != auth {
			t.Fatal("unexpected auth")
		}

		if auth && !bytes.Equal(username, v.username) {
			t.Fatal("different username when decoded")
		}

		if auth && id != v.id {
			t.Fatal("different id when decoded")
		}

		if auth {
			scratchBadge = b
		}
	}

	_, _, auth := Get(scratchBadge[:20], goodKey)
	if auth {
		t.Fatal("badge too short, shouldn't auth")
	}

	scratchBadge = append(scratchBadge, byte(0x00))

	_, _, auth = Get(scratchBadge, goodKey)
	if auth {
		t.Fatal("badge too long, shouldn't auth")
	}
}
