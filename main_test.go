package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestDomainNameEncoder(t *testing.T) {
	t.Run("test domain name encoder", func(t *testing.T) {
		expectedHexString := "06676f6f676c6503636f6d00"
		out := DomainNameEncoder("google.com")
		fmt.Println(out)
		hexString := hex.EncodeToString(out)
		if hexString != expectedHexString {
			t.Fatalf("domain name encoding failed\nexpected: %s\nactual  : %s\n", hexString, expectedHexString)
		}
	})
}

func TestHeaderToBytes(t *testing.T) {
	t.Run("test header to bytes conversion", func(t *testing.T) {
		header := DNSHeader{
			ID:             4884,
			Flags:          0,
			NumQuestion:    1,
			NumAnswers:     0,
			NumAuthorities: 0,
			NumAdditionals: 0,
		}

		fmt.Println(header.ToBytes())
	})
}

func TestBuildQuery(t *testing.T) {
	t.Run("test build query", func(t *testing.T) {
		out := BuildQuery("www.example.com", TypeA)
		fmt.Println(out)
	})
}
