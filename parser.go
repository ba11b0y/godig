package main

import (
	"bytes"
	"encoding/binary"
)

type DNSRecord struct {
	// Name is the domain name
	Name []byte
	// Type is the record type,ex: A, AAAA, MX
	Type  uint16
	Class uint16
	TTL   uint16
	Data  []byte
}

func parseHeader(reader *bytes.Reader) DNSHeader {
	var header DNSHeader

	err := binary.Read(reader, binary.BigEndian, &header)
	if err != nil {
		panic(err)
	}

	return header
}
