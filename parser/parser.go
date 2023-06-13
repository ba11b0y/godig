package parser

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

const (
	// TypeA is the infamous A record type
	TypeA   uint16 = 1
	TypeNS  uint16 = 2
	ClassIn        = 1
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
			// add period('.') to the domain name after adding this domain name part
			name = append(name, 46)
			break
		} else {
			// find a way to read multiple bytes at once, possible reader.ReadAt()
			part := make([]byte, length)
			io.ReadFull(reader, part)
			name = append(name, part...)
			// add period('.') to the domain name after adding this domain name part
			name = append(name, 46)
		}
	}

	// strip the extra period('.') from the end
	name = name[0 : len(name)-1]

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

	var record = DNSRecord{
		Name:  name,
		Type:  recordData.Type,
		Class: recordData.Class,
		TTL:   recordData.TTL,
	}

	switch recordData.Type {
	//case TypeA:
	case TypeNS:
		data := decodeName(reader)
		record.Data = data
	default:
		var data = make([]byte, recordData.DataLen)

		_, err = io.ReadFull(reader, data)
		if err != nil {
			panic(err)
		}

		record.Data = data
	}

	return record
}

func ParseIP(data []byte) string {
	var ip strings.Builder

	for _, b := range data[0:3] {
		// each byte is an IP segment, convert it to string and write to ip
		ip.WriteString(strconv.Itoa(int(b)))
		ip.WriteString(".")
	}

	ip.WriteString(strconv.Itoa(int(data[3])))
	return ip.String()
}
