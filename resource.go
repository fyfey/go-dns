package main

import (
	"encoding/binary"
	"fmt"
)

// Resource is a resource record
type Resource struct {
	Name        string     `json:"name"`
	Type        QueryType  `json:"type"`
	Class       QueryClass `json:"class"`
	TTL         int        `json:"ttl"`
	RDataLength int        `json:"rdataLength"`
	RData       string     `json:"rdata"`
}

// NewResource creates a new resource record from raw bytes
func NewResource(b []byte, offset int) (int, *Resource) {
	startOffset := offset
	name := []byte{}
	for {
		lenByte := b[offset]
		l := int(lenByte)
		offset++
		if l == 0 {
			// End of name
			break
		}
		if l>>4 == 0xc {
			// It's a pointer!
			pointer := binary.BigEndian.Uint16(b[offset-1:offset+1]) & 0xfff
			_, domainName := UnpackDomainName(b[pointer:])
			name = append(name, []byte(domainName)...)
			fmt.Printf("Pointer to %d\n%s\n", pointer, name)
			offset++
			break
		}
		if len(name) > 0 {
			name = append(name, 0x2e)
		}
		name = append(name, b[offset:offset+l]...)
		offset += l
	}

	resource := &Resource{Name: string(name)}

	resource.Type = QueryType(binary.BigEndian.Uint16(b[offset : offset+2]))
	offset += 2

	resource.Class = QueryClass(binary.BigEndian.Uint16(b[offset : offset+2]))
	offset += 2

	resource.TTL = int(binary.BigEndian.Uint32(b[offset : offset+4]))
	offset += 4

	resource.RDataLength = int(binary.BigEndian.Uint16(b[offset : offset+2]))
	offset += 2

	address := ""
	for i := offset; i < offset+resource.RDataLength; i++ {
		address = fmt.Sprintf("%s%d", address, b[i])
		if i < offset+resource.RDataLength-1 {
			address += "."
		}
	}
	resource.RData = string(address)
	offset += resource.RDataLength

	return offset - startOffset, resource
}

func (r *Resource) Marshal() []byte {
	b := make([]byte, 0)

	offset := PackDomainName(&b, r.Name)
	b = append(b, make([]byte, 10)...)
	binary.BigEndian.PutUint16(b[offset:], uint16(r.Type))
	binary.BigEndian.PutUint16(b[offset+2:], uint16(r.Class))
	binary.BigEndian.PutUint32(b[offset+4:], uint32(r.TTL))
	binary.BigEndian.PutUint16(b[offset+8:], uint16(r.RDataLength))
	PackIPAddress(&b, r.RData)

	return b
}
