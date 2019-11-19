package query

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

type Query struct {
	Name  string     `json:"name"`
	Type  QueryType  `json:"type"`
	Class QueryClass `json:"class"`
}

func ReadBytes(b []byte) (offset int, q *Query) {
	name := []byte{}
	for {
		lenByte := b[offset]
		offset++
		l := int(lenByte)
		if l == 0 {
			break
		}
		if len(name) > 0 {
			name = append(name, 0x2e)
		}
		buf := b[offset : offset+l]
		name = append(name, buf...)
		offset += l
	}

	q = &Query{Name: string(name)}

	t := b[offset : offset+2]
	q.Type = QueryType(binary.BigEndian.Uint16(t))
	offset += 2

	c := b[offset : offset+2]
	q.Class = QueryClass(binary.BigEndian.Uint16(c))

	return
}
