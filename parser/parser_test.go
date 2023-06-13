package parser_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	main2 "godig/parser"
)

func TestDomainNameEncoder(t *testing.T) {
	t.Run("test domain name encoder", func(t *testing.T) {
		expectedHexString := "06676f6f676c6503636f6d00"
		out := main2.DomainNameEncoder("google.com")
		fmt.Println(out)
		hexString := hex.EncodeToString(out)
		if hexString != expectedHexString {
			t.Fatalf("domain name encoding failed\nexpected: %s\nactual  : %s\n", hexString, expectedHexString)
		}
	})
}
