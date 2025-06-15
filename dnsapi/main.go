package dnsapi

import "errors"

type ConfigZone struct {
	DnsApi  string   `toml:"dns_api"`
	Zone    string   `toml:"zone_id"`
	Secret  string   `toml:"secret"`
	Domains []string `toml:"domains"`
}

type UpdateResult int

const (
	ResultOk UpdateResult = iota
	ResultError
	ResultSkip
)

type DnsApi interface {
	GetName() string
	Update(zone ConfigZone, ip string) (UpdateResult, error)
}

func GetDnsApi(name string) (DnsApi, error) {
	switch name {
	case "cloudflare":
		return &Cloudflare{}, nil
	default:
		return nil, errors.New("unknown dns api type")
	}
}
