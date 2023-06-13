package main

import (
	"fmt"

	"godig/parser"
	"godig/resolver"
)

func main() {
	domainName := "twitter.com"
	recordType := parser.TypeA
	ipData := resolver.Resolve(domainName, recordType)
	fmt.Printf("Resolved IP for %s is %s\n", domainName, parser.ParseIP(ipData))
}
