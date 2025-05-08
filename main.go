// main.go
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/akiyomov/dns-proxy-go/resolver"
	"github.com/miekg/dns"
)

func main() {
	port := flag.String("port", "53", "Port to run the DNS server on")
	blocklistPath := flag.String("blocklist", "./config/blocklist.txt", "Path to the blocklist file")
	upstream := flag.String("upstream", "1.1.1.1:53", "Upstream DNS server address")
	flag.Parse()

	proxy := resolver.NewDNSProxy(*upstream, *blocklistPath)

	address := fmt.Sprintf(":%s", *port)
	server := &dns.Server{Addr: address, Net: "udp"}
	server.Handler = dns.HandlerFunc(proxy.HandleDNSRequest)

	log.Printf("Starting DNS proxy on %s...", address)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start DNS server: %v", err)
	}
}
