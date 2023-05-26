package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

type DNSRecord struct {
	// Name is the domain name
	Name []byte
	// Type is the record type,ex: A, AAAA, MX
	Type  uint16
	Class uint16
	TTL   uint32
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

func parseQuestion(reader *bytes.Reader) DNSQuestion {

	// parse question type and class
	var typeAndClass struct {
		Type  uint16
		Class uint16
	}

	name := nameParser(reader)

	err := binary.Read(reader, binary.BigEndian, &typeAndClass)
	if err != nil {
		panic(err)
	}

	return DNSQuestion{
		Name:  name,
		Type:  typeAndClass.Type,
		Class: typeAndClass.Class,
	}
}

func nameParser(reader *bytes.Reader) []byte {
	var (
		length byte
		name   []byte
	)

	// parse domain name.
	for {
		length, _ = reader.ReadByte()
		if length == 0 {
			break
		}

		if length&128 == 1 && length&64 == 1 {
			panic("implement compression decoder")
		} else {
			// find a way to read multiple bytes at once, possible reader.ReadAt()
			for i := 0; i < int(length); i++ {
				b, _ := reader.ReadByte()
				name = append(name, b)
			}
		}
	}

	return name
}

func parseRecord(reader *bytes.Reader) DNSRecord {

	name := nameParser(reader)

	var recordData struct {
		Type    uint16
		Class   uint16
		TTL     uint32
		DataLen uint16
	}

	err := binary.Read(reader, binary.BigEndian, &recordData)
	if err != nil {
		panic(err)
	}

	var data = make([]byte, 10)

	_, err = io.ReadFull(reader, data)
	if err != nil {
		panic(err)
	}

	return DNSRecord{
		Name:  name,
		Type:  recordData.Type,
		Class: recordData.Class,
		TTL:   recordData.TTL,
		Data:  data,
	}
}
