package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

const (
	Query     = 0
	Response  = 1
	NoError   = 0
	FormError = 1
	ServFail  = 2
	NxDomain  = 3
)

func main() {
	// ip := &net.UDPAddr{Port: 1234}

	records := map[string]*Resource{}
	records["fyfe.io"] = &Resource{
		Name:        "fyfe.io",
		Class:       IN,
		Type:        A,
		TTL:         300,
		RData:       "178.62.55.224",
		RDataLength: 4,
	}

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

		queryHeader := NewHeader(buf[0:12])

		headerJson, _ := json.MarshalIndent(queryHeader, "", " ")
		fmt.Printf("Header: %s", headerJson)
		if queryHeader.QR == Response {
			fmt.Errorf("Received response?!\n")
			continue
		}
		questions := []*Question{}
		offset := 12
		for i := 0; i < queryHeader.TotalQuestions; i++ {
			o, q := NewQuestion(buf[offset:])
			fmt.Printf("New offset: %d", o)
			offset += o
			fmt.Printf("Q: %v\n", q)
			questions = append(questions, q)
		}

		answers := []*Resource{}
		for i := 0; i < queryHeader.TotalAnswerRR; i++ {
			o, resource := NewResource(buf, offset)
			offset += o
			fmt.Printf("Q: %v\n", resource)
			answers = append(answers, resource)
		}

		authorityRRs := []*Resource{}
		for i := 0; i < queryHeader.TotalAnswerRR; i++ {
			o, resource := NewResource(buf, offset)
			offset += o
			fmt.Printf("Q: %v\n", resource)
			authorityRRs = append(authorityRRs, resource)
		}

		additionalRRs := []*Resource{}
		for i := 0; i < queryHeader.TotalAdditionalRR; i++ {
			o, resource := NewResource(buf, offset)
			offset += o
			fmt.Printf("Q: %v\n", resource)
			additionalRRs = append(additionalRRs, resource)
		}
		fmt.Printf("Questions: %v", questions)
		fmt.Printf("Answers: %v", answers)
		fmt.Printf("AuthorityRRs: %v", authorityRRs)
		fmt.Printf("AdditionalRRs: %v", additionalRRs)

		foundAnswers := 1
		responseCode := NoError
		record, ok := records[questions[0].Name]
		if !ok {
			foundAnswers = 0
			responseCode = NxDomain
		}
		response := make([]byte, 0)

		header := Header{
			ID:                 queryHeader.ID,
			RecursionDesired:   false,
			RecursionAvailable: false,
			QR:                 Response,
			ReturnCode:         responseCode,
			TotalQuestions:     1,
			TotalAnswerRR:      foundAnswers,
			AuthorativeAnswer:  true,
			AuthenticatedData:  true,
			Z:                  false,
			Truncated:          false,
			CheckingDisabled:   false,
		}
		response = append(response, header.Marshal()...)
		response = append(response, questions[0].Marshal()...)
		if foundAnswers > 0 {
			response = append(response, record.Marshal()...)
		}

		fmt.Println(hex.Dump(response))

		conn.WriteTo(response, addr)
	}

}
