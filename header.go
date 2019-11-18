package main

import (
	"encoding/binary"
	"fmt"
)

type Header struct {
	ID                 int  `json:"id"`
	QR                 int  `json:"qr"`
	Opcode             int  `json:"opcode"`
	AuthorativeAnswer  bool `json:"authorativeAnswer"`
	Truncated          bool `json:"truncated"`
	RecursionDesired   bool `json:"recursionDesired"`
	RecursionAvailable bool `json:"recursionAvailable"`
	Z                  bool `json:"z"`
	AuthenticatedData  bool `json:"authenticatedData"`
	CheckingDisabled   bool `json:"checkingDisabled"`
	ReturnCode         int  `json:"returnCode"`
	TotalQuestions     int  `json:"totalQuestions"`
	TotalAnswerRR      int  `json:"totalAnswerRR"`
	TotalAuthorityRR   int  `json:"totalAuthorityRR"`
	TotalAdditionalRR  int  `json:"totalAdditionalRR"`
}

func NewHeader(b []byte) *Header {
	qr := b[2] >> 7
	opCode := b[2] & (0xf << 4)
	aa := b[2] & 0x01 << 2
	tc := b[2] & 0x01 << 1
	rd := b[2] & 0x01
	ra := b[3] & 0x01 << 7
	z := b[3] & 0x01 << 6
	ad := b[3] & 0x01 << 5
	cd := b[3] & 0x01 << 4
	rcode := b[3] & 0x0f
	totalQuestions := binary.BigEndian.Uint16(b[4:6])
	totalAnswerRR := binary.BigEndian.Uint16(b[6:8])
	totalAuthorityRR := binary.BigEndian.Uint16(b[8:10])
	totalAdditionalRR := binary.BigEndian.Uint16(b[10:12])

	header := &Header{
		int(binary.BigEndian.Uint16(b[0:2])),
		int(qr),
		int(binary.BigEndian.Uint16([]byte{0, 0, 0, opCode})),
		int(binary.BigEndian.Uint16([]byte{0, 0, 0, 0, aa})) == 1,
		int(tc) == 1,
		int(rd) == 1,
		int(ra) == 1,
		int(z) == 1,
		int(ad) == 1,
		int(cd) == 1,
		int(rcode),
		int(totalQuestions),
		int(totalAnswerRR),
		int(totalAuthorityRR),
		int(totalAdditionalRR),
	}

	return header
}

// Sets the bit at pos in the integer n.
func setBit(n int, pos uint) int {
	n |= (1 << pos)
	return n
}

// Clears the bit at pos in n.
func clearBit(n int, pos uint) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}

func MarshalHeader(h *Header) []byte {
	b := make([]byte, 12)

	fmt.Printf("ID: %x", uint16(h.ID))
	binary.BigEndian.PutUint16(b[0:4], uint16(h.ID))

	header0 := 0
	binary.BigEndian.PutUint16(b[2:], uint16(h.Opcode<<3))
	if h.QR == 1 {
		header0 = setBit(header0, 7)
	}
	if h.AuthorativeAnswer {
		header0 = setBit(header0, 2)
	}
	if h.Truncated {
		header0 = setBit(header0, 1)
	}
	if h.RecursionDesired {
		header0 = setBit(header0, 0)
	}

	b[2] = byte(header0)

	header3 := 0
	header3 &= h.ReturnCode
	if h.RecursionAvailable {
		header3 = setBit(header3, 7)
	}
	if h.Z {
		header3 = setBit(header3, 6)
	}
	if h.AuthenticatedData {
		header3 = setBit(header3, 5)
	}
	if h.CheckingDisabled {
		header3 = setBit(header3, 4)
	}

	b[3] = byte(header3)

	binary.BigEndian.PutUint16(b[4:6], uint16(h.TotalQuestions))
	binary.BigEndian.PutUint16(b[6:8], uint16(h.TotalAnswerRR))
	binary.BigEndian.PutUint16(b[8:10], uint16(h.TotalAuthorityRR))
	binary.BigEndian.PutUint16(b[10:12], uint16(h.TotalAdditionalRR))

	return b
}
