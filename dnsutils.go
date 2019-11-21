package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// PackDomainName packs name into b
func PackDomainName(b *[]byte, name string) int {
	parts := strings.Split(name, ".")

	n := 0
	for i := 0; i < len(parts); i++ {
		*b = append(*b, byte(len(parts[i])))
		*b = append(*b, []byte(parts[i])...)
		n += len(parts[i]) + 1
	}
	*b = append(*b, 0x00)
	n++

	return n
}

// UnpackDomainName unpacks domain name from b
func UnpackDomainName(b []byte) (int, string) {
	offset := 0
	parts := [][]byte{}
	for {
		len := int(b[offset])
		fmt.Printf("Len :%d\n", len)
		offset++
		if len == 0 {
			break
		}
		parts = append(parts, b[offset:offset+len])
		offset += int(len)
	}

	return offset, string(bytes.Join(parts, []byte{0x2e}))
}

// PackIPAddress writes address to b
func PackIPAddress(b *[]byte, address string) {
	parts := strings.Split(address, ".")
	for _, p := range parts {
		n, _ := strconv.Atoi(p)
		*b = append(*b, byte(n))
	}
}
