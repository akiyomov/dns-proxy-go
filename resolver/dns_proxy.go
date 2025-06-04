package resolver

import (
	"log"
	"os"
	"strings"

	"github.com/miekg/dns"
)

type StatsRecorder interface {
	RecordQuery(domain string, blocked bool)
}

type DNSProxy struct {
	blockedDomains map[string]struct{}
	upstreamServer string
	stats          StatsRecorder
}

func NewDNSProxy(upstreamServer string, blocklistPath string, stats StatsRecorder) *DNSProxy {
	// Load blocked domains
	blockedDomains := make(map[string]struct{})
	if blocklistPath != "" {
		data, err := os.ReadFile(blocklistPath)
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				domain := strings.TrimSpace(line)
				if domain != "" && !strings.HasPrefix(domain, "//") && !strings.HasPrefix(domain, "#") {
					blockedDomains[domain] = struct{}{}
				}
			}
		}
	}

	log.Printf("Loaded %d blocked domains", len(blockedDomains))

	return &DNSProxy{
		blockedDomains: blockedDomains,
		upstreamServer: upstreamServer,
		stats:          stats,
	}
}

func (p *DNSProxy) HandleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Compress = false

	for _, q := range r.Question {
		domain := strings.TrimSuffix(q.Name, ".")
		
		// Check if domain is blocked
		if _, blocked := p.blockedDomains[domain]; blocked {
			// Record blocked query
			if p.stats != nil {
				p.stats.RecordQuery(domain, true)
			}
			
			// Blocked response - return NXDOMAIN
			msg.SetRcode(r, dns.RcodeNameError)
			log.Printf("[BLOCKED] %s", domain)
		} else {
			// Record allowed query
			if p.stats != nil {
				p.stats.RecordQuery(domain, false)
			}
			
			// Forward to upstream
			upstreamResp, err := dns.Exchange(r, p.upstreamServer)
			if err != nil {
				log.Printf("[ERROR] Failed to forward request for %s: %v", domain, err)
				msg.SetRcode(r, dns.RcodeServerFailure)
			} else {
				msg.Answer = upstreamResp.Answer
				msg.Ns = upstreamResp.Ns
				msg.Extra = upstreamResp.Extra
				msg.SetRcode(r, upstreamResp.Rcode)
				log.Printf("[ALLOWED] %s", domain)
			}
		}
	}

	w.WriteMsg(msg)
}
