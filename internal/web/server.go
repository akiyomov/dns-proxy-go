package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/akiyomov/dns-proxy-go/internal/stats"
)

type WebServer struct {
	stats *stats.Stats
	port  string
}

func NewWebServer(statsManager *stats.Stats, port string) *WebServer {
	return &WebServer{
		stats: statsManager,
		port:  port,
	}
}

func (ws *WebServer) Start() {
	// Serve static files (dashboard)
	http.Handle("/", http.FileServer(http.Dir("./dashboard/")))
	
	// API endpoints
	http.HandleFunc("/api/stats", ws.handleStats)
	http.HandleFunc("/api/top-blocked", ws.handleTopBlocked)
	
	address := fmt.Sprintf(":%s", ws.port)
	log.Printf("Starting web dashboard on http://localhost%s", address)
	
	go func() {
		if err := http.ListenAndServe(address, nil); err != nil {
			log.Fatalf("Failed to start web server: %v", err)
		}
	}()
}

func (ws *WebServer) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	statsData := ws.stats.GetStats()
	
	if err := json.NewEncoder(w).Encode(statsData); err != nil {
		log.Printf("Error encoding stats: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (ws *WebServer) handleTopBlocked(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	// Get limit from query parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	topDomains := ws.stats.GetTopBlockedDomains(limit)
	
	if err := json.NewEncoder(w).Encode(topDomains); err != nil {
		log.Printf("Error encoding top blocked domains: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}