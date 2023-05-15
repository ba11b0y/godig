package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"strings"
)

const (
	// TypeA is the infamous A record type
	TypeA   = 1
	ClassIn = 1
)

// DNSHeader ...
type DNSHeader struct {
	ID             int
	Flags          int
	NumQuestion    int
	NumAnswers     int
	NumAuthorities int
	NumAdditionals int
}

// ToBytes converts all field values of a DNS header to big endian two byte integers and concatenates each
// field's two byte integer representation.
func (header DNSHeader) ToBytes() []byte {
	var byteData bytes.Buffer

	v := reflect.ValueOf(header)

	for i := 0; i < v.NumField(); i++ {
		val := uint16(v.Field(i).Interface().(int))
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
	Type  int
	Class int
}

// ToBytes converts all field values of a DNS question to big endian two byte integers and concatenates each
// field's two byte integer representation.
func (question DNSQuestion) ToBytes() []byte {
	var byteData bytes.Buffer

	v := reflect.ValueOf(question)

	byteData.Write(question.Name)

	for i := 1; i < v.NumField(); i++ {
		val := uint16(v.Field(i).Interface().(int))
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, val)
		byteData.Write(b)
	}

	return byteData.Bytes()
}

// DomainNameEncoder encodes domain name to bytes
// input:  google.com
// output: '[6 103 111 111 103 108 101 3 99 111 109 0]'
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

// BuildQuery builds a DNS query
func BuildQuery(domainName string, recordType int) []byte {
	encodedDomainName := DomainNameEncoder(domainName)
	id := 33432
	recursionDesired := 1 << 8
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

	// Process the response here
	fmt.Println(response)
}
