package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"
	"strings"
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

	name := decodeName(reader)

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

func decodeName(reader *bytes.Reader) []byte {
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

		if length&0b1100_0000 != 0 {
			name = append(decodeCompressedName(length, reader))
			break
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

func decodeCompressedName(length byte, reader *bytes.Reader) []byte {
	b, _ := reader.ReadByte()
	pointerBytes := []byte{length & 0b0011_1111, b}
	pointer := binary.BigEndian.Uint16(pointerBytes)
	currentPos, _ := reader.Seek(0, io.SeekCurrent)
	reader.Seek(int64(pointer), io.SeekStart)
	result := decodeName(reader)
	reader.Seek(currentPos, io.SeekStart)
	return result
}

func parseRecord(reader *bytes.Reader) DNSRecord {

	name := decodeName(reader)

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

func parseIP(data []byte) string {
	var ip strings.Builder

	for _, b := range data[0:3] {
		// each byte is an IP segment, convert it to string and write to ip
		ip.WriteString(strconv.Itoa(int(b)))
		ip.WriteString(".")
	}

	ip.WriteString(strconv.Itoa(int(data[3])))
	return ip.String()
}
