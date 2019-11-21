package main

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestReadBytes(t *testing.T) {
	//             4     f     y     f     e     2     i     o     0     1           1
	b := []byte{0x04, 0x66, 0x79, 0x66, 0x65, 0x02, 0x69, 0x6f, 0x00, 0x00, 0x01, 0x00, 0x01}

	_, query := NewQuestion(b)

	if query.Name != "fyfe.io" {
		t.Errorf("Expected fyfe.io; got %s\n", query.Name)
	}
	if query.Type != A {
		t.Errorf("Expected A; got %d\n", query.Type)
	}
	if query.Class != IN {
		t.Errorf("Expected IN; got %d\n", query.Class)
	}
}

func TestMarshalQuestion(t *testing.T) {
	q := &Question{
		Name:  "fyfe.io",
		Type:  A,
		Class: IN,
	}
	e := []byte{0x04, 0x66, 0x79, 0x66, 0x65, 0x02, 0x69, 0x6f, 0x00, 0x00, 0x01, 0x00, 0x01}

	b := q.Marshal()

	if !bytes.Equal(b, e) {
		t.Errorf("Expected\n%s\nGot\n%s", hex.Dump(e), hex.Dump(b))
	}
}
