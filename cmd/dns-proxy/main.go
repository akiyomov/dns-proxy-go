package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/akiyomov/dns-proxy-go/internal/stats"
	"github.com/akiyomov/dns-proxy-go/internal/web"
	"github.com/akiyomov/dns-proxy-go/resolver"
	"github.com/miekg/dns"
)

func main() {
	port := flag.String("port", "53", "Port to run the DNS server on")
	blocklistPath := flag.String("blocklist", "./config/blocklist.txt", "Path to the blocklist file")
	upstream := flag.String("upstream", "1.1.1.1:53", "Upstream DNS server address")
	webPort := flag.String("web-port", "8080", "Port to run the web dashboard on")
	flag.Parse()

	// Initialize stats tracking
	statsManager := stats.NewStats()

	// Initialize DNS proxy with stats
	proxy := resolver.NewDNSProxy(*upstream, *blocklistPath, statsManager)

	// Start web server for dashboard
	webServer := web.NewWebServer(statsManager, *webPort)
	webServer.Start()

	// Setup DNS server
	address := fmt.Sprintf(":%s", *port)
	server := &dns.Server{Addr: address, Net: "udp"}
	server.Handler = dns.HandlerFunc(proxy.HandleDNSRequest)

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting DNS proxy on %s...", address)
		log.Printf("Dashboard available at http://localhost:%s", *webPort)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start DNS server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-c
	log.Println("Shutting down DNS proxy...")
	server.Shutdown()
}