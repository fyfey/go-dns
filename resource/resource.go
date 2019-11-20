package resource

import (
	"encoding/binary"
	"fmt"

	"fyfe.io/dns/dnsutils"
)

// Resource is a resource record
type Resource struct {
	Name        string `json:"name"`
	Type        int    `json:"type"`
	Class       int    `json:"class"`
	TTL         int    `json:"ttl"`
	RDataLength int    `json:"rdataLength"`
	RData       string `json:"rdata"`
}

// NewResource creates a new resource record from raw bytes
func NewResource(b []byte, offset int) *Resource {
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
			name = append(name, []byte(dnsutils.UnpackDomainName(b[pointer:]))...)
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

	resource.Type = int(binary.BigEndian.Uint16(b[offset : offset+2]))
	offset += 2

	resource.Class = int(binary.BigEndian.Uint16(b[offset : offset+2]))
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

	return resource
}

func (r *Resource) Marshal() []byte {
	b := make([]byte, 0)

	offset := dnsutils.PackDomainName(&b, r.Name)
	b = append(b, make([]byte, 8)...)
	binary.BigEndian.PutUint16(b[offset:], uint16(r.Type))
	binary.BigEndian.PutUint16(b[offset+2:], uint16(r.Class))
	binary.BigEndian.PutUint32(b[offset+4:], uint32(r.TTL))
	b = append(b, byte(r.RDataLength))
	dnsutils.PackIPAddress(&b, r.RData)

	return b
}
