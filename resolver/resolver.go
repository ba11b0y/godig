package resolver

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"

	"godig/parser"
)

// root server by Verisign, Inc. https://www.iana.org/domains/root/servers
const rootNameServer = "198.41.0.4"

// BuildQuery builds a DNS query
func BuildQuery(domainName string, recordType uint16) []byte {
	encodedDomainName := parser.DomainNameEncoder(domainName)
	id := uint16(rand.Intn(65535))
	//recursionDesired := uint16(1 << 8)
	header := parser.DNSHeader{
		ID:          id,
		Flags:       0, // Set Flags to 0, since we don't need recursion.
		NumQuestion: 1,
	}
	question := parser.DNSQuestion{
		Name:  encodedDomainName,
		Type:  recordType,
		Class: parser.ClassIn,
	}

	var query bytes.Buffer
	query.Write(header.ToBytes())
	query.Write(question.ToBytes())

	return query.Bytes()
}

func SendQuery(ip, domainName string, recordType uint16) parser.DNSPacket {
	query := BuildQuery(domainName, recordType)

	// create a UDP socket
	conn, err := net.Dial("udp", fmt.Sprintf("%s:53", ip))
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
	packet := parser.ParseDNSPacket(responseReader)

	return packet
}

func Resolve(domainName string, recordType uint16) []byte {
	var nameServer string
	nameServer = rootNameServer
	for {
		fmt.Printf("Querying %s for %s\n", nameServer, domainName)
		responsePacket := SendQuery(nameServer, domainName, recordType)
		if ip := responsePacket.GetAnswer(); ip != nil {
			return ip
		}

		if nsIP := responsePacket.GetNameServerIP(); nsIP != nil {
			nameServer = parser.ParseIP(nsIP)
			continue
		}

		if nsDomain := responsePacket.GetNameServer(); nsDomain != "" {
			nameServer = parser.ParseIP(Resolve(nsDomain, parser.TypeA))
		}
	}
}
