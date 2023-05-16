package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
)

const (
	// TypeA is the infamous A record type
	TypeA   uint16 = 1
	ClassIn        = 1
)

// BuildQuery builds a DNS query
func BuildQuery(domainName string, recordType uint16) []byte {
	encodedDomainName := DomainNameEncoder(domainName)
	id := uint16(rand.Intn(65535))
	recursionDesired := uint16(1 << 8)
	header := DNSHeader{
		ID:          id,
		Flags:       recursionDesired,
		NumQuestion: 1,
	}
	question := DNSQuestion{
		Name:  encodedDomainName,
		Type:  recordType,
		Class: ClassIn,
	}

	var query bytes.Buffer
	query.Write(header.ToBytes())
	query.Write(question.ToBytes())

	return query.Bytes()
}

func main() {
	query := BuildQuery("www.example.com", TypeA)

	// create a UDP socket
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// send our query to 8.8.8.8, port 53. Port 53 is the DNS port.
	_, err = conn.Write(query)
	if err != nil {
		panic(err)
	}

	// read the response. UDP DNS responses are usually less than 512 bytes
	// so reading 1024 bytes is enough
	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		panic(err)
	}

	responseReader := bytes.NewReader(response)

	// Process the response here
	fmt.Println(parseHeader(responseReader))

}
