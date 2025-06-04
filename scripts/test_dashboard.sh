#!/bin/bash

# Quick test script for DNS proxy dashboard
# Tests a mix of blocked and allowed domains to populate stats

DNS_SERVER="localhost"
DNS_PORT="5354"
DASHBOARD_URL="http://localhost:8080"

echo "🚀 DNS Proxy Dashboard Test Script"
echo "=================================="
echo "DNS Server: ${DNS_SERVER}:${DNS_PORT}"
echo "Dashboard: ${DASHBOARD_URL}"
echo ""

# Function to test a domain
test_domain() {
    local domain=$1
    local expected=$2
    
    result=$(dig @$DNS_SERVER -p $DNS_PORT +short +time=2 +tries=1 $domain 2>/dev/null)
    
    if [ -z "$result" ]; then
        status="BLOCKED"
        icon="🚫"
    else
        status="ALLOWED"
        icon="✅"
    fi
    
    echo "$icon $domain -> $status"
}

# Show initial stats
echo "📊 Initial Stats:"
curl -s $DASHBOARD_URL/api/stats | jq '{totalQueries, blockedQueries, allowedQueries, uniqueBlockedDomains: (.blockedDomains | length)}'
echo ""

# Test blocked domains (from your blocklist)
echo "🔴 Testing Blocked Domains:"
blocked_domains="facebook.com instagram.com doubleclick.net googleadservices.com googlesyndication.com amazon-adsystem.com google-analytics.com googletagmanager.com 180hits.de ads.yahoo.com advertising.com"

for domain in $blocked_domains; do
    test_domain $domain "blocked"
    sleep 0.1
done

echo ""

# Test allowed domains
echo "🟢 Testing Allowed Domains:"
allowed_domains="google.com github.com stackoverflow.com wikipedia.org cloudflare.com mozilla.org kernel.org gnu.org debian.org ubuntu.com"

for domain in $allowed_domains; do
    test_domain $domain "allowed"
    sleep 0.1
done

echo ""
echo "📊 Final Stats:"
curl -s $DASHBOARD_URL/api/stats | jq '{totalQueries, blockedQueries, allowedQueries, uniqueBlockedDomains: (.blockedDomains | length)}'

echo ""
echo "🏆 Top 5 Blocked Domains:"
curl -s "$DASHBOARD_URL/api/top-blocked?limit=5" | jq -r '.[] | "• \(.domain): \(.count) times (last: \(.lastSeen | strptime("%Y-%m-%dT%H:%M:%S") | strftime("%H:%M:%S")))"'

echo ""
echo "🎯 Dashboard Summary:"
stats=$(curl -s $DASHBOARD_URL/api/stats)
total=$(echo $stats | jq '.totalQueries')
blocked=$(echo $stats | jq '.blockedQueries')
allowed=$(echo $stats | jq '.allowedQueries')

if [ "$total" -gt 0 ]; then
    block_rate=$(echo "scale=1; $blocked * 100 / $total" | bc -l 2>/dev/null || echo "0")
    echo "   Total Queries: $total"
    echo "   Blocked: $blocked ($block_rate%)"
    echo "   Allowed: $allowed"
    echo ""
    echo "🌐 View live dashboard: $DASHBOARD_URL"
else
    echo "   No queries processed yet"
fi