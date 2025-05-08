package resolver

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/miekg/dns"
)

type DNSProxy struct {
	blockedDomains map[string]struct{}
	upstreamServer string
}

func NewDNSProxy(upstreamServer string, blocklistPath string) *DNSProxy {
	// Load blocked domains
	blockedDomains := make(map[string]struct{})
	if blocklistPath != "" {
		data, err := os.ReadFile(blocklistPath)
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				domain := strings.TrimSpace(line)
				if domain != "" {
					blockedDomains[domain] = struct{}{}
				}
			}
		}
	}

	return &DNSProxy{
		blockedDomains: blockedDomains,
		upstreamServer: upstreamServer,
	}
}

func (p *DNSProxy) HandleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Compress = false

	for _, q := range r.Question {
		domain := strings.TrimSuffix(q.Name, ".")
		if _, blocked := p.blockedDomains[domain]; blocked {
			// Blocked response
			msg.Answer = []dns.RR{
				&dns.A{
					Hdr: dns.RR_Header{
						Name:   q.Name,
						Rrtype: dns.TypeA,
						Class:  dns.ClassINET,
						Ttl:    0,
					},
					A: net.IPv4(0, 0, 0, 0),
				},
			}
			log.Printf("[BLOCKED] %s", domain)
		} else {
			// Forward to upstream
			upstreamResp, err := dns.Exchange(r, p.upstreamServer)
			if err != nil {
				log.Printf("[ERROR] Failed to forward request: %v", err)
				return
			}
			msg.Answer = upstreamResp.Answer
		}
	}

	w.WriteMsg(msg)
}
