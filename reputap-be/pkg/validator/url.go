package validator

import (
	"net/url"
	"strings"
)

var allowedGoogleDomains = []string{
	"search.google.com",
	"www.google.com",
	"google.com",
	"maps.google.com",
	"www.google.co.id",
	"google.co.id",
	"maps.app.goo.gl",
	"goo.gl",
	"g.page",
	"g.co",
}

func IsValidGoogleReviewURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	if parsed.Scheme != "https" {
		return false
	}

	host := strings.ToLower(parsed.Host)
	for _, domain := range allowedGoogleDomains {
		if host == domain {
			return true
		}
	}

	return false
}
