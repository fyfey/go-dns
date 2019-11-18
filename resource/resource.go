package resource

import (
	bio "bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	"fyfe.io/dns/bufio"
)

type Resource struct {
	Name        string `json:"name"`
	Type        int    `json:"type"`
	Class       int    `json:"class"`
	TTL         int    `json:"ttl"`
	RDataLength int    `json:"rdataLength"`
	RData       string `json:"rdata"`
}

func (r *Resource) ReadBytes(reader *bio.Reader) *Resource {
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

	resource := &Resource{Name: string(name)}

	t, err := bufio.ReadNBytes(reader, 2)
	if err != nil {
		panic(err)
	}
	resource.Type = int(binary.BigEndian.Uint16(t))

	c, err := bufio.ReadNBytes(reader, 2)
	if err != nil {
		panic(err)
	}
	resource.Class = int(binary.BigEndian.Uint16(c))

	ttl, err := bufio.ReadNBytes(reader, 4)
	if err != nil {
		panic(err)
	}
	resource.TTL = int(binary.BigEndian.Uint16(ttl))

	return resource
}

func packDomainName(b []byte, name string) {
	parts := strings.Split(name, ".")

	for i := 0; i < len(parts); i++ {
		b = append(b, byte(len(parts[i])))
		b = append(b, []byte(parts[i])...)
	}
}

func unpackDomainName(b []byte) string {
	offset := 0
	parts := [][]byte{}
	for {
		len := int(b[offset])
		fmt.Printf("Len :%d\n", len)
		if len == 0 {
			break
		}
		parts = append(parts, b[offset+1:offset+len+1])
		offset += int(len) + 1
	}

	return string(bytes.Join(parts, []byte{0x2e}))
}
