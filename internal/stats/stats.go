package stats

import (
	"sync"
	"time"
)

type BlockedDomain struct {
	Domain    string    `json:"domain"`
	Count     int       `json:"count"`
	LastSeen  time.Time `json:"lastSeen"`
}

type RecentQuery struct {
	Domain    string    `json:"domain"`
	Blocked   bool      `json:"blocked"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "BLOCKED" or "ALLOWED"
}

type Stats struct {
	mu              sync.RWMutex
	TotalQueries    int64                    `json:"totalQueries"`
	BlockedQueries  int64                    `json:"blockedQueries"`
	AllowedQueries  int64                    `json:"allowedQueries"`
	BlockedDomains  map[string]*BlockedDomain `json:"blockedDomains"`
	RecentQueries   []*RecentQuery           `json:"recentQueries"`
	StartTime       time.Time                `json:"startTime"`
}

func NewStats() *Stats {
	return &Stats{
		BlockedDomains: make(map[string]*BlockedDomain),
		RecentQueries:  make([]*RecentQuery, 0),
		StartTime:      time.Now(),
	}
}

func (s *Stats) RecordQuery(domain string, blocked bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.TotalQueries++
	
	// Add to recent queries (keep last 50)
	queryType := "ALLOWED"
	if blocked {
		queryType = "BLOCKED"
	}
	
	recentQuery := &RecentQuery{
		Domain:    domain,
		Blocked:   blocked,
		Timestamp: time.Now(),
		Type:      queryType,
	}
	
	s.RecentQueries = append([]*RecentQuery{recentQuery}, s.RecentQueries...)
	if len(s.RecentQueries) > 50 {
		s.RecentQueries = s.RecentQueries[:50]
	}
	
	if blocked {
		s.BlockedQueries++
		
		if existing, ok := s.BlockedDomains[domain]; ok {
			existing.Count++
			existing.LastSeen = time.Now()
		} else {
			s.BlockedDomains[domain] = &BlockedDomain{
				Domain:   domain,
				Count:    1,
				LastSeen: time.Now(),
			}
		}
	} else {
		s.AllowedQueries++
	}
}

func (s *Stats) GetStats() *Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Create a copy for safe concurrent access
	statsCopy := &Stats{
		TotalQueries:   s.TotalQueries,
		BlockedQueries: s.BlockedQueries,
		AllowedQueries: s.AllowedQueries,
		BlockedDomains: make(map[string]*BlockedDomain),
		RecentQueries:  make([]*RecentQuery, len(s.RecentQueries)),
		StartTime:      s.StartTime,
	}
	
	for k, v := range s.BlockedDomains {
		statsCopy.BlockedDomains[k] = &BlockedDomain{
			Domain:   v.Domain,
			Count:    v.Count,
			LastSeen: v.LastSeen,
		}
	}
	
	copy(statsCopy.RecentQueries, s.RecentQueries)
	
	return statsCopy
}

func (s *Stats) GetTopBlockedDomains(limit int) []*BlockedDomain {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	domains := make([]*BlockedDomain, 0, len(s.BlockedDomains))
	for _, domain := range s.BlockedDomains {
		domains = append(domains, domain)
	}
	
	// Sort by count (simple bubble sort for small datasets)
	for i := 0; i < len(domains)-1; i++ {
		for j := 0; j < len(domains)-i-1; j++ {
			if domains[j].Count < domains[j+1].Count {
				domains[j], domains[j+1] = domains[j+1], domains[j]
			}
		}
	}
	
	if limit > 0 && limit < len(domains) {
		domains = domains[:limit]
	}
	
	return domains
}