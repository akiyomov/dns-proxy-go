# DNS Proxy Go

A high-performance DNS proxy server with ad blocking capabilities and a real-time web dashboard for monitoring blocked domain statistics.

## Features

- 🚫 DNS request proxying with ad/tracker blocking (3000+ domains)
- 📊 Real-time statistics tracking and monitoring
- 🌐 Beautiful web dashboard with live updates
- 📋 Configurable blocklist with comment support
- 🏆 Top blocked domains ranking with counts
- ⏰ Recent queries tracking (last 50 queries)
- ⏱️ Server uptime monitoring
- 🔄 Auto-refresh dashboard (every 2 seconds)
- 🛡️ High-performance concurrent DNS resolution

## Quick Start

### Prerequisites

- Go 1.21+ installed
- Linux/macOS/Windows
- Port 53 access (for production) or use alternative ports

### Installation & Setup

1. **Clone and enter the project directory:**
   ```bash
   git clone <repository-url>
   cd dns-proxy-go
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Build the project:**
   ```bash
   mkdir -p bin
   go build -o bin/dns-proxy ./cmd/dns-proxy
   ```

4. **Run the DNS proxy:**
   
   **Option A: Development mode (non-privileged ports):**
   ```bash
   ./bin/dns-proxy -port 5354 -web-port 8080 -blocklist ./config/blocklist.txt -upstream 8.8.8.8:53
   ```
   
   **Option B: Production mode (requires sudo for port 53):**
   ```bash
   sudo ./bin/dns-proxy -port 53 -web-port 8080 -blocklist ./config/blocklist.txt -upstream 8.8.8.8:53
   ```

5. **Access the dashboard:**
   Open your browser and go to: `http://localhost:8080`

## Project Structure

```
dns-proxy-go/
├── cmd/dns-proxy/          # Main application entry point
│   └── main.go            # Application startup and configuration
├── internal/              # Private application packages
│   ├── stats/             # Statistics tracking and management
│   │   └── stats.go       # Stats collection, storage, and retrieval
│   └── web/               # Web server and API endpoints
│       └── server.go      # HTTP server, dashboard, and REST API
├── resolver/              # DNS resolution logic
│   └── dns_proxy.go       # DNS request handling and blocking logic
├── dashboard/             # Static web dashboard files
│   └── index.html         # Interactive dashboard with real-time updates
├── config/                # Configuration files
│   └── blocklist.txt      # Domain blocklist (3000+ ad/tracker domains)
├── scripts/               # Utility and test scripts
│   └── test_dashboard.sh  # Comprehensive testing script
├── bin/                   # Compiled binaries (created after build)
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
└── README.md              # This file
```

## Configuration Options

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `53` | DNS server port (use 5354+ for non-root) |
| `-web-port` | `8080` | Web dashboard port |
| `-blocklist` | `./config/blocklist.txt` | Path to domain blocklist file |
| `-upstream` | `1.1.1.1:53` | Upstream DNS server address |

### Example Configurations

**Development Setup:**
```bash
./bin/dns-proxy -port 5354 -web-port 8080 -upstream 8.8.8.8:53
```

**Production Setup:**
```bash
sudo ./bin/dns-proxy -port 53 -web-port 80 -upstream 1.1.1.1:53
```

**Custom Blocklist:**
```bash
./bin/dns-proxy -port 5354 -blocklist /path/to/custom/blocklist.txt
```

## Dashboard Features

Visit `http://localhost:8080` after starting the server to access:

### 📊 Real-time Statistics
- **Total DNS Queries**: Complete count of all processed requests
- **Blocked Queries**: Number and percentage of blocked requests
- **Allowed Queries**: Number of legitimate requests forwarded
- **Server Uptime**: How long the server has been running

### 🏆 Top Blocked Domains
- Most frequently blocked domains
- Request counts per domain
- Last seen timestamps
- Auto-updating list

### ⏰ Recent DNS Queries
- Live feed of the last 50 DNS queries
- Color-coded status (🚫 Blocked / ✅ Allowed)
- Timestamps for each query
- Real-time updates

### 🔄 Auto-refresh Controls
- Toggle auto-refresh on/off
- Manual refresh button
- 2-second update interval
- Connection status indicator

## Testing Your Setup

### 1. Test DNS Resolution

**Test a blocked domain (should return NXDOMAIN):**
```bash
dig @localhost -p 5354 facebook.com
# Expected: NXDOMAIN response
```

**Test an allowed domain (should return IP):**
```bash
dig @localhost -p 5354 google.com  
# Expected: Valid IP address
```

### 2. Test Dashboard API

**Get current statistics:**
```bash
curl -s http://localhost:8080/api/stats | jq
```

**Get top blocked domains:**
```bash
curl -s "http://localhost:8080/api/top-blocked?limit=5" | jq
```

### 3. Run Comprehensive Tests

**Execute the test script:**
```bash
chmod +x scripts/test_dashboard.sh
./scripts/test_dashboard.sh
```

This script will:
- Test multiple blocked domains
- Test multiple allowed domains  
- Show real-time statistics updates
- Display dashboard summary

## Development

### Building from Source

```bash
# Install dependencies
go mod tidy

# Build for current platform
go build -o bin/dns-proxy ./cmd/dns-proxy

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o bin/dns-proxy-linux ./cmd/dns-proxy
GOOS=windows GOARCH=amd64 go build -o bin/dns-proxy.exe ./cmd/dns-proxy
GOOS=darwin GOARCH=amd64 go build -o bin/dns-proxy-mac ./cmd/dns-proxy
```

### Code Structure

- **`cmd/dns-proxy/main.go`**: Application entry point and configuration
- **`internal/stats/`**: Thread-safe statistics collection and management
- **`internal/web/`**: HTTP server with REST API endpoints
- **`resolver/`**: DNS request processing and domain blocking logic
- **`dashboard/`**: Static web assets for the monitoring dashboard

### Adding Custom Domains to Blocklist

Edit `config/blocklist.txt`:
```bash
# Comments start with # and are ignored
facebook.com
instagram.com
doubleclick.net
your-custom-domain.com
```

## API Reference

### GET /api/stats
Returns complete DNS proxy statistics including recent queries.

**Response:**
```json
{
  "totalQueries": 1523,
  "blockedQueries": 892,
  "allowedQueries": 631,
  "blockedDomains": {
    "facebook.com": {
      "domain": "facebook.com",
      "count": 45,
      "lastSeen": "2025-06-04T14:30:15Z"
    }
  },
  "recentQueries": [
    {
      "domain": "google.com",
      "blocked": false,
      "timestamp": "2025-06-04T14:30:20Z",
      "type": "ALLOWED"
    }
  ],
  "startTime": "2025-06-04T10:15:30Z"
}
```

### GET /api/top-blocked?limit=N
Returns the top N most blocked domains (default: 10).

**Response:**
```json
[
  {
    "domain": "facebook.com",
    "count": 45,
    "lastSeen": "2025-06-04T14:30:15Z"
  }
]
```

## Troubleshooting

### Common Issues

**1. "Permission denied" on port 53:**
```bash
# Use sudo for privileged ports or use alternative ports
sudo ./bin/dns-proxy -port 53
# OR
./bin/dns-proxy -port 5354
```

**2. "Address already in use":**
```bash
# Check what's using the port
sudo lsof -i :53
# Kill the process or use a different port
./bin/dns-proxy -port 5354
```

**3. Dashboard not loading:**
- Ensure the `dashboard/` directory exists in the current working directory
- Check web server logs for errors
- Verify port 8080 is not blocked by firewall

**4. DNS queries timing out:**
- Check upstream DNS server connectivity
- Verify firewall allows UDP traffic on the DNS port
- Test with: `dig @8.8.8.8 google.com`

### Logging

The application logs DNS queries in real-time:
```
2025/06/04 14:30:15 [BLOCKED] facebook.com
2025/06/04 14:30:16 [ALLOWED] google.com
2025/06/04 14:30:17 [BLOCKED] doubleclick.net
```

## Performance

- **Concurrent Request Handling**: Uses goroutines for high-performance DNS resolution
- **Thread-safe Statistics**: Lock-based concurrent access to statistics
- **Memory Efficient**: Recent queries limited to 50 entries
- **Fast Domain Lookup**: Hash map-based blocklist for O(1) domain checking

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:
1. Check the troubleshooting section above
2. Review the logs for error messages
3. Test with the provided test script
4. Open an issue with detailed information about your setup

---

**Happy DNS blocking! 🛡️**
