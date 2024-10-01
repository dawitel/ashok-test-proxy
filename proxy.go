package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/dawitel/Ashok-reverse-proxy-test/internal/config"
	"github.com/dawitel/Ashok-reverse-proxy-test/internal/utils"
)

// loadUserAgent reads the user agent string from the user agent file.
func loadUserAgent() string {
	ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"
	return ua
}

// loadCookies reads cookies from the file specified in the configuration.
func loadCookies(cfg *config.Config) (string, error) {
	cookieData, err := os.ReadFile(cfg.Proxy.CookieFile)
	if err != nil {
		return "", fmt.Errorf("failed to read cookies: %v", err)
	}
	return strings.TrimSpace(string(cookieData)), nil
}

// parseCookies parses the cookie string into a slice of HTTP cookies.
func parseCookies(cookieString string) []*http.Cookie {
	var cookies []*http.Cookie
	for _, line := range strings.Split(cookieString, "\n") {
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 7 {
			fmt.Printf("Skipping invalid cookie line: %s\n", line)
			continue
		}
		name := strings.Trim(parts[5], `"`)
		value := strings.Trim(parts[6], `"`)
		cookie := &http.Cookie{
			Name:   name,
			Value:  value,
			Path:   parts[2],
			Domain: parts[0],
		}
		cookies = append(cookies, cookie)
	}
	return cookies
}

// ProxyHandler handles incoming requests and proxies them to the target URL.
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		http.Error(w, "Failed to load configuration", http.StatusInternalServerError)
		return
	}

	userAgent := loadUserAgent()

	cookieString, err := loadCookies(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookies := parseCookies(cookieString)

	targetURL, err := url.Parse(cfg.Proxy.TargetURL)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	for _, cookie := range cookies {
		r.AddCookie(cookie)
	}
	
	// Modify the headers for Intercom API if it's a ping request
	if isPingRequest(r) {
		handlePingRequestHeaders(r)
		// Add payload for the ping request
		payload := generatePingPayload(r)
		r.Body = io.NopCloser(bytes.NewBuffer(payload))
		r.ContentLength = int64(len(payload))
	}

	// Modify the request to point to the target URL
	r.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.URL.Host = targetURL.Host

	r.Header.Set("User-Agent", userAgent)
	r.Header.Set("Host", targetURL.Host)
	r.Header.Set("Origin", targetURL.String())
	r.Header.Set("Referer", r.Referer())

	// Ensure we bypass security headers that would restrict local development
	w.Header().Set("Access-Control-Allow-Origin", targetURL.String())
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Del("Strict-Transport-Security")
	w.Header().Del("X-Frame-Options")

	// Serve the proxied request
	proxy.ServeHTTP(w, r)
}
func isPingRequest(req *http.Request) bool {
	return strings.Contains(req.URL.Path, "/messenger/web/ping")
}

// handlePingRequestHeaders sets the required headers for the Intercom API ping request.
func handlePingRequestHeaders(req *http.Request) {
	req.Header.Set("access-control-allow-origin", "https://www.semrush.com")
	req.Header.Set("access-control-allow-credentials", "true")
	req.Header.Set("access-control-allow-headers", "Content-Type, Idempotency-Key, X-INTERCOM-APP, X-INTERCOM-PAGE-TITLE, X-INTERCOM-USER-DATA")
	req.Header.Set("access-control-allow-methods", "POST, GET, OPTIONS")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("X-INTERCOM-APP", "cs07vi2k")
	req.Header.Set("X-INTERCOM-PAGE-TITLE", req.URL.String())
	req.Header.Set("Idempotency-Key", utils.GenerateIdempotencyKey()) // Call to generate a new Idempotency Key
}

// generatePingPayload generates the specific payload required for the ping request.
func generatePingPayload(req *http.Request) []byte {
	payload := map[string]interface{}{
		"app_id":                 "cs07vi2k",
		"v":                      "3",
		"g":                      "ddf69cd514575586eed88a06b20f18a9e3a3eb07",
		"s":                      "45e565db-81f8-4788-9fae-70fe6db83fb2",
		"r":                      fmt.Sprintf("https://semrush.com%s", req.URL.Path),
		"platform":               "web",
		"installation_type":      "js-snippet",
		"Idempotency-Key":        utils.GenerateIdempotencyKey(),
		"is_intersection_booted": false,
		"user_active_company_id": "undefined",
		"user_data":              generateUserData(), // Dynamic user data from request
		"page_title":             req.URL.String(),
		"source":                 "apiBoot",
		"sampling":               false,
		"referer": fmt.Sprintf("https://semrush.com%s", req.URL.Path),
	}

	payloadStr := encodePayload(payload) // A utility function to encode payload as form-urlencoded
	return []byte(payloadStr)
}

// generateUserData generates dynamic user data for the payload.
func generateUserData() map[string]interface{} {
	return map[string]interface{}{
		"email":        "dianeburms1.6.1.990@gmail.com",
		"user_id":      "22774131",
		"user_hash":    "a72d3c46c056ce84efd558247257c833d6d516384bc9b59827535e0ed9691cb7",
		"GA Client ID": "1536862937.1727372350",
		"name":         " ",
		"phone":        nil,
		"created_at":   1727327044,
		"Paid":         false,
		"Product":      "guru",
		"Expire Date":  "2024-10-03 01:09:07",
	}
}

// encodePayload encodes the payload into a URL-encoded format for form submissions.
func encodePayload(payload map[string]interface{}) string {
	var encodedPayload []string
	for key, value := range payload {
		encodedPayload = append(encodedPayload, fmt.Sprintf("%s=%v", key, value))
	}
	return strings.Join(encodedPayload, "&")
}
