package query

import (
	bio "bufio"
	"encoding/binary"
	"log"

	"fyfe.io/dns/bufio"
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

func ReadBytes(reader *bio.Reader) *Query {
	name := []byte{}
	for {
		lenByte, err := reader.ReadByte()
		if err != nil {
			log.Panic(err)
		}
		l := int(lenByte)
		if l == 0 {
			break
		}
		if len(name) > 0 {
			name = append(name, 0x2e)
		}
		buf, err := bufio.ReadNBytes(reader, l)
		if err != nil {
			panic(err)
		}
		name = append(name, buf...)
	}

	query := &Query{Name: string(name)}

	t, err := bufio.ReadNBytes(reader, 2)
	if err != nil {
		panic(err)
	}
	query.Type = QueryType(binary.BigEndian.Uint16(t))

	c, err := bufio.ReadNBytes(reader, 2)
	if err != nil {
		panic(err)
	}
	query.Class = QueryClass(binary.BigEndian.Uint16(c))

	return query
}
