package googleplaces

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var hexPattern = regexp.MustCompile(`(0x[0-9a-fA-F]+):(0x[0-9a-fA-F]+)`)

// HexToPlaceID converts Google's internal feature hex identifiers (0x...:0x...) into a Place ID (ChIJ...)
func HexToPlaceID(hex1, hex2 string) (string, error) {
	v1, err := strconv.ParseUint(strings.TrimPrefix(strings.ToLower(hex1), "0x"), 16, 64)
	if err != nil {
		return "", fmt.Errorf("parsing hex1: %w", err)
	}

	v2, err := strconv.ParseUint(strings.TrimPrefix(strings.ToLower(hex2), "0x"), 16, 64)
	if err != nil {
		return "", fmt.Errorf("parsing hex2: %w", err)
	}

	buf := make([]byte, 20)
	buf[0] = 0x0a
	buf[1] = 0x12
	buf[2] = 0x09
	binary.LittleEndian.PutUint64(buf[3:11], v1)
	buf[11] = 0x11
	binary.LittleEndian.PutUint64(buf[12:20], v2)

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// ToDirectReviewURL converts a Google Maps URL, short link, or place URL into a direct 5-star review modal URL.
// If the URL is already a direct review link or cannot be resolved, it returns the original URL.
func ToDirectReviewURL(ctx context.Context, rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}

	// 1. If already a direct review link, return as is
	if strings.Contains(rawURL, "writereview?placeid=") || strings.HasSuffix(rawURL, "/review") {
		return rawURL
	}

	// 2. Check if already has placeid query parameter
	if parsed, err := url.Parse(rawURL); err == nil {
		if pid := parsed.Query().Get("placeid"); pid != "" {
			return fmt.Sprintf("https://search.google.com/local/writereview?placeid=%s", pid)
		}
	}

	// 3. Check if rawURL already contains hex pattern
	if m := hexPattern.FindStringSubmatch(rawURL); len(m) == 3 {
		if pid, err := HexToPlaceID(m[1], m[2]); err == nil && pid != "" {
			return fmt.Sprintf("https://search.google.com/local/writereview?placeid=%s", pid)
		}
	}

	// 4. Resolve redirects (for maps.app.goo.gl, goo.gl/maps, etc.)
	targetURL := rawURL
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err == nil {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		resp, err := client.Do(req)
		if err == nil {
			targetURL = resp.Request.URL.String()
			resp.Body.Close()
		}
	}

	// 5. Look for hex pattern in the resolved URL
	if m := hexPattern.FindStringSubmatch(targetURL); len(m) == 3 {
		if pid, err := HexToPlaceID(m[1], m[2]); err == nil && pid != "" {
			return fmt.Sprintf("https://search.google.com/local/writereview?placeid=%s", pid)
		}
	}

	// Fallback to original URL
	return rawURL
}
