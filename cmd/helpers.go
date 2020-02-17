package cmd

import (
	"fmt"
	"net/url"
	"strings"
)

// ParseURL parses a url and splits the result in scheme, host and port
func parseURL(URL string) (scheme, host, port string, err error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", "", "", fmt.Errorf("url requires scheme://host:port formatting. Received: %s", u)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", "", "", fmt.Errorf("url scheme needs to be http or https. Received: %s", u.Scheme)
	}
	scheme = u.Scheme

	parts := strings.Split(u.Host, ":")
	host = parts[0]
	if len(parts) > 1 {
		port = parts[1]
	} else {
		// set default scheme if port not set explicitly
		if scheme == "http" {
			port = "80"
		} else {
			port = "443"
		}
	}
	return scheme, host, port, nil
}

// ensureGatewayClient returns a client implementing gateway.Handler
// This exists because of the way cobra inititalizes configuration only on execution and
// not during registration of flags and commands
func ensureClient(c client, api_url string) client {
	if c != nil {
		return c
	}

	scheme, host, port, _ := parseURL(api_url)
	return NewClient(scheme, host, port)
}
