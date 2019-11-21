package main

import (
	"encoding/binary"
)

type QueryType int
type QueryClass int

const (
	A        QueryType = 1
	NS                 = 2
	CNAME              = 5
	SOA                = 6
	PTR                = 12
	MX                 = 15
	TXT                = 16
	SPF                = 99
	DNSSECTA           = 32768
	DNSSECLV           = 32769
)
const (
	Reserved QueryClass = iota
	IN                  = 1
)

type Question struct {
	Name  string     `json:"name"`
	Type  QueryType  `json:"type"`
	Class QueryClass `json:"class"`
}

func NewQuestion(b []byte) (offset int, q *Question) {
	offset, name := UnpackDomainName(b)

	q = &Question{Name: name}

	t := b[offset : offset+2]
	q.Type = QueryType(binary.BigEndian.Uint16(t))
	offset += 2

	c := b[offset : offset+2]
	q.Class = QueryClass(binary.BigEndian.Uint16(c))

	return
}

func (q *Question) Marshal() []byte {
	b := make([]byte, 0)

	offset := PackDomainName(&b, q.Name)
	b = append(b, 0x00, 0x00, 0x00, 0x00)
	binary.BigEndian.PutUint16(b[offset:], uint16(q.Type))
	binary.BigEndian.PutUint16(b[offset+2:], uint16(q.Class))

	return b
}
