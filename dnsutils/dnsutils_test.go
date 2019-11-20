package dnsutils

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestPackDomainName(t *testing.T) {
	b := make([]byte, 0)
	PackDomainName(&b, "fyfe.io")

	expected := []byte{0x04, 0x66, 0x79, 0x66, 0x65, 0x02, 0x69, 0x6f}

	if !bytes.Equal(b, expected) {
		t.Errorf("Incorrect packed domain. Expected %s, got %s", hex.Dump(expected), hex.Dump(b))
	}
}

func TestUnpackDomainName(t *testing.T) {
	b := []byte{0x04, 0x66, 0x79, 0x66, 0x65, 0x02, 0x69, 0x6f, 0x00}

	name := UnpackDomainName(b)

	if name != "fyfe.io" {
		t.Errorf("Unpack domain failed. Expected fyfe.io; got %s", name)
	}
}
