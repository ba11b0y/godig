package main

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

func parseDNSPacket(reader *bytes.Reader) DNSPacket {
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
