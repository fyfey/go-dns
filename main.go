package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"fyfe.io/dns/query"
)

func main() {
	// ip := &net.UDPAddr{Port: 1234}

	conn, err := net.ListenPacket("udp", ":1234")
	// conn, err := net.ListenUDP("udp4", ip)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Listening on :1234\n")
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		l, addr, err := conn.ReadFrom(buf)
		if err != nil {
			break
		}
		fmt.Printf("packet-received: bytes=%d from=%s\n",
			l, addr.String())

		header := NewHeader(buf[0:12])

		headerJson, _ := json.MarshalIndent(header, "", " ")
		fmt.Printf("Header: %s", headerJson)
		if header.QR == 1 {
			fmt.Errorf("Received response?!\n")
			continue
		}
		questions := []*query.Query{}
		offset := 12
		for i := 0; i < header.TotalQuestions; i++ {
			o, q := query.ReadBytes(buf[offset:])
			fmt.Printf("New offset: %d", o)
			offset += o
			fmt.Printf("Q: %v\n", q)
			questions = append(questions, q)
		}

		fmt.Printf("New offset: %d", offset)

		additionalRRs := []*query.Query{}
		for i := 0; i < header.TotalAdditionalRR; i++ {
			o, q := query.ReadBytes(buf[offset:])
			offset += o
			fmt.Printf("Q: %v\n", q)
			additionalRRs = append(additionalRRs, q)
		}
		fmt.Printf("Questions: %v", questions)
		fmt.Printf("AdditionalRRs: %v", additionalRRs)
	}

}
