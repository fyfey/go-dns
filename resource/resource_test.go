package resource

import (
	"bytes"
	"encoding/hex"
	"testing"

	"fyfe.io/dns/query"
)

func TestNewResource(t *testing.T) {
	rr := "db4281800001000100000000037777770c6e6f7274686561737465726e036564750000010001c00c000100010000025800049b211144"

	b, _ := hex.DecodeString(rr)
	resource := NewResource(b, 38)

	if resource.Name != "www.northeastern.edu" {
		t.Errorf("Expected www.northeastern.edu; got %s\n", resource.Name)
	}
	if resource.Type != int(query.A) {
		t.Errorf("Expected type A (0x01); got %d", resource.Type)
	}
	if resource.Class != int(query.IN) {
		t.Errorf("Expected type IN (0x01); got %d", resource.Class)
	}
	if resource.TTL != 600 {
		t.Errorf("Expected TTL 600; got %d", resource.TTL)
	}
	if resource.RDataLength != 4 {
		t.Errorf("Expected RDataLength 4; got %d", resource.RDataLength)
	}
	if resource.RData != "155.33.17.68" {
		t.Errorf("Expected RData \"155.33.17.68\"; got %s", resource.RData)
	}
}

func TestMarshal(t *testing.T) {
	resource := Resource{
		Name:        "fyfe.io",
		Type:        int(query.A),
		Class:       query.IN,
		TTL:         300,
		RDataLength: 4,
		RData:       "178.62.55.224",
	}
	expected := []byte{0x04, 0x66, 0x79, 0x66, 0x65, 0x02, 0x69, 0x6f, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x01, 0x2c, 0x04, 0xb2, 0x3e, 0x37, 0xe0}

	b := resource.Marshal()

	if !bytes.Equal(b, expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s\n", hex.Dump(expected), hex.Dump(b))
	}
}
