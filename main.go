package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"fyfe.io/dns/opcode"
	"fyfe.io/dns/query"
)

func main() {
	ip := &net.UDPAddr{Port: 1234}
	conn, err := net.ListenUDP("udp4", ip)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Listening on :%d\n", ip.Port)
	defer conn.Close()

	buf := make([]byte, 0)
	tmp := make([]byte, 256)
	for {
		l, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Panic(err)
			}
			break
		}
		buf = append(buf, tmp[:l]...)
		fmt.Printf("Read %d\n%s\n", l, hex.Dump(buf))

		fmt.Printf("%x", buf[0:2])
		id := binary.BigEndian.Uint16(buf[0:2])
		fmt.Printf("ID BE: %d\n", id)

		qr := buf[2] & 0x01 << 7
		opCode := buf[2] & 0x0f << 6
		aa := buf[2] & 0x01 << 2
		tc := buf[2] & 0x01 << 1
		rd := buf[2] & 0x01
		ra := buf[3] & 0x01 << 7
		z := buf[3] & 0x01 << 6
		ad := buf[3] & 0x01 << 5
		cd := buf[3] & 0x01 << 4
		rcode := buf[3] & 0x0f
		totalQuestions := binary.BigEndian.Uint16(buf[4:6])
		totalAnswerRR := binary.BigEndian.Uint16(buf[6:8])
		totalAuthorityRR := binary.BigEndian.Uint16(buf[8:10])
		totalAdditionalRR := binary.BigEndian.Uint16(buf[10:12])

		header := &Header{
			int(id),
			int(qr),
			int(opcode.Opcode(opCode)),
			int(aa) == 1,
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

		headerJson, _ := json.MarshalIndent(header, "", "  ")
		fmt.Printf("Header: %s\n", headerJson)

		reader := bytes.NewReader(buf[12:])
		bufReader := bufio.NewReader(reader)

		query := query.ReadBytes(bufReader)
		queryJson, _ := json.MarshalIndent(query, "", " ")
		fmt.Printf("Query: %s\n", queryJson)
	}

}
