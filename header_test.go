package main

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestNewHeader(t *testing.T) {
	raw := []byte{0x73, 0x48, 0x01, 0x20, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

	header := NewHeader(raw)

	if header.QR != 0 {
		t.Errorf("Incorrect QR. Expected 0; got %d", header.QR)
	}
	if header.Opcode != 0 {
		t.Errorf("Incorrect Opcode. Expected 0; got %d", header.Opcode)
	}
	if header.AuthorativeAnswer != false {
		t.Errorf("Incorrect AuthoriativeAnswer. Expected false; got true")
	}
	if header.Truncated != false {
		t.Errorf("Incorrect Truncated. Expected false; got true")
	}
	if header.RecursionDesired != true {
		t.Errorf("Incorrect RecursionDesired. Expected true; got false")
	}
	if header.RecursionAvailable != false {
		t.Errorf("Incorrect RecursionAvailable. Expected false; got true")
	}
	if header.Z != false {
		t.Errorf("Incorrect Z. Expected false; got true")
	}
	if header.AuthenticatedData != false {
		t.Errorf("Incorrect AuthenticatedData. Expected false; got true")
	}
	if header.CheckingDisabled != false {
		t.Errorf("Incorrect CheckingDisabled. Expected false; got true")
	}
	if header.ReturnCode != 0 {
		t.Errorf("Incorrect ReturnCode. Expected 0; got %d", header.ReturnCode)
	}
	if header.TotalQuestions != 1 {
		t.Errorf("Incorrect TotalQuestion. Expected 1; got %d", header.TotalQuestions)
	}
	if header.TotalAnswerRR != 0 {
		t.Errorf("Incorrect TotalAnswerRR. Expected 0; got %d", header.TotalAnswerRR)
	}
	if header.TotalAuthorityRR != 0 {
		t.Errorf("Incorrect TotalAuthorityRR. Expected 0; got %d", header.TotalAnswerRR)
	}
	if header.TotalAdditionalRR != 1 {
		t.Errorf("Incorrect TotalAdditionalRR. Expected 1; got %d", header.TotalQuestions)
	}
}

func TestToBytes(t *testing.T) {
	header := &Header{
		ID:                 0x7348,
		QR:                 1,
		Opcode:             0,
		AuthorativeAnswer:  false,
		Truncated:          false,
		RecursionDesired:   true,
		RecursionAvailable: true,
		Z:                  false,
		AuthenticatedData:  false,
		CheckingDisabled:   false,
		ReturnCode:         0,
		TotalQuestions:     1,
		TotalAnswerRR:      0,
		TotalAuthorityRR:   0,
		TotalAdditionalRR:  1,
	}
	expected := []byte{0x73, 0x48, 0x81, 0x80, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

	b := MarshalHeader(header)

	if !bytes.Equal(b, expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", hex.Dump(expected), hex.Dump(b))
	}
}
