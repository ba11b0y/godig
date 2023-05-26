package main

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strings"
)

// DNSHeader ...
type DNSHeader struct {
	ID             uint16
	Flags          uint16
	NumQuestion    uint16
	NumAnswers     uint16
	NumAuthorities uint16
	NumAdditionals uint16
}

// ToBytes converts all field values of a DNS header to big endian two byte integers and concatenates each
// field's two byte integer representation.
func (header DNSHeader) ToBytes() []byte {
	var byteData bytes.Buffer

	v := reflect.ValueOf(header)

	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i).Interface().(uint16)
		// converting an integer(base10) to a 2-byte integer
		// for example: 23 is converted to 17
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, val)
		byteData.Write(b)
	}

	return byteData.Bytes()
}

type DNSQuestion struct {
	Name  []byte
	Type  uint16
	Class uint16
}

// ToBytes converts all field values of a DNS question to big endian two byte integers and concatenates each
// field's two byte integer representation.
func (question DNSQuestion) ToBytes() []byte {
	var byteData bytes.Buffer

	v := reflect.ValueOf(question)

	byteData.Write(question.Name)

	for i := 1; i < v.NumField(); i++ {
		val := v.Field(i).Interface().(uint16)
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, val)
		byteData.Write(b)
	}

	return byteData.Bytes()
}

// DomainNameEncoder encodes domain name to bytes
// input:  www.example.com
// output: '[3 119 119 119 7 101 120 97 109 112 108 101 3 99 111 109 0]'
func DomainNameEncoder(domainName string) []byte {
	var encodedDomainName bytes.Buffer
	parts := strings.Split(domainName, ".")
	for _, part := range parts {
		encodedDomainName.WriteByte(byte(len(part)))
		encodedDomainName.Write([]byte(part))
	}
	emptyByte := make([]byte, 1)
	encodedDomainName.Write(emptyByte)

	return encodedDomainName.Bytes()
}
