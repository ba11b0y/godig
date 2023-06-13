package main

import (
	"fmt"
)

// root server by Verisign, Inc. https://www.iana.org/domains/root/servers
const rootNameServer = "198.41.0.4"

func resolve(domainName string, recordType uint16) []byte {
	var nameServer string
	nameServer = rootNameServer
	for {
		fmt.Printf("Querying %s for %s\n", nameServer, domainName)
		responsePacket := SendQuery(nameServer, domainName, recordType)
		if ip := responsePacket.getAnswer(); ip != nil {
			return ip
		}

		if nsIP := responsePacket.getNameServerIP(); nsIP != nil {
			nameServer = parseIP(nsIP)
			continue
		}

		if nsDomain := responsePacket.getNameServer(); nsDomain != "" {
			nameServer = parseIP(resolve(nsDomain, TypeA))
		}
	}
}
