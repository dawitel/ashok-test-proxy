package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/dawitel/Ashok-reverse-proxy-test/internal/config"
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

	// Attach cookies to the request
	for _, cookie := range cookies {
		r.AddCookie(cookie)
	}

	// Modify the request to look like it's coming from Semrush
	r.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.URL.Host = targetURL.Host

	r.Header.Set("User-Agent", userAgent)
	r.Header.Set("Host", targetURL.Host)
	r.Header.Set("Origin", "https://www.semrush.com")
	r.Header.Set("Referer", "https://www.semrush.com")

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "https://www.semrush.com")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Custom-Header")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Proxy response modification
	proxy.ModifyResponse = func(resp *http.Response) error {
		// Allow cross-origin requests in the response
		resp.Header.Set("Access-Control-Allow-Origin", "https://www.semrush.com")
		resp.Header.Set("Access-Control-Allow-Credentials", "true")

		// Remove security headers that could block the request
		resp.Header.Del("Strict-Transport-Security")
		resp.Header.Del("X-Frame-Options")

		// Optionally modify the body or headers as necessary
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "application/javascript") {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			modifiedBody := strings.ReplaceAll(string(bodyBytes), "mydomain.com", "semrush.com")

			resp.Body = io.NopCloser(strings.NewReader(modifiedBody))
			resp.ContentLength = int64(len(modifiedBody))
			resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(modifiedBody)))
		}

		return nil
	}

	// Serve the proxied request
	proxy.ServeHTTP(w, r)
}
