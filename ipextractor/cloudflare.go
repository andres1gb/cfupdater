package ipextractor

import (
	"errors"
	"net"
	"regexp"
	"time"
)

const (
	cloudflareUrl  = "https://cloudflare.com/cdn-cgi/trace"
	requestTimeout = 5 * time.Second
)

var (
	ipRegex = regexp.MustCompile(`ip=([0-9\.]+)\n`)
)

type Cloudflare struct{}

// GetIp retrieves the public IP address from Cloudflare's trace endpoint
// Returns the IP address or an error if the request fails or IP cannot be extracted
func (c *Cloudflare) GetIp() (string, error) {
	body, err := simpleGetRequest(cloudflareUrl)
	if err != nil {
		return "", errors.New("failed to request Cloudflare trace: " + err.Error())
	}

	if body == "" {
		return "", errors.New("empty response body from Cloudflare trace")
	}

	ipFound := ipRegex.FindStringSubmatch(body)
	if len(ipFound) < 2 {
		return "", errors.New("IP address pattern not found in response")
	}

	ip := ipFound[1]
	if net.ParseIP(ip) == nil {
		return "", errors.New("invalid IP address format received")
	}

	return ip, nil
}
