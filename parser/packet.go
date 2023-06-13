package parser

import (
	"bytes"
)

type DNSPacket struct {
	Header      DNSHeader
	Questions   []DNSQuestion
	Answers     []DNSRecord
	Authorities []DNSRecord
	Additionals []DNSRecord
}

// GetAnswer returns the first A record in the Answer section
func (packet DNSPacket) GetAnswer() []byte {
	for _, a := range packet.Answers {
		if a.Type == TypeA {
			return a.Data
		}
	}

	return nil
}

// GetNameServerIP returns the first A record in the Additionals section
func (packet DNSPacket) GetNameServerIP() []byte {
	for _, a := range packet.Additionals {
		if a.Type == TypeA {
			return a.Data
		}
	}

	return nil
}

// GetNameServer returns the first NS record in the Authority section
func (packet DNSPacket) GetNameServer() string {
	for _, a := range packet.Authorities {
		if a.Type == TypeNS {
			nameServerDomain := string(a.Data)
			return nameServerDomain
		}
	}

	return ""
}

func ParseDNSPacket(reader *bytes.Reader) DNSPacket {
	var (
		header                            DNSHeader
		questions                         []DNSQuestion
		answers, authorities, additionals []DNSRecord
	)
	header = parseHeader(reader)
	for i := 0; i < int(header.NumQuestion); i++ {
		questions = append(questions, parseQuestion(reader))
	}

	for i := 0; i < int(header.NumAnswers); i++ {
		answers = append(answers, parseRecord(reader))
	}

	for i := 0; i < int(header.NumAuthorities); i++ {
		authorities = append(authorities, parseRecord(reader))
	}

	for i := 0; i < int(header.NumAdditionals); i++ {
		additionals = append(additionals, parseRecord(reader))
	}

	return DNSPacket{
		Header:      header,
		Questions:   questions,
		Answers:     answers,
		Authorities: authorities,
		Additionals: additionals,
	}
}
