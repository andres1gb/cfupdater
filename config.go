package main

import (
	"cfupdater/dnsapi"
)

type Config struct {
	Main struct {
		IpExtractor string  `toml:"ip_extractor"`
		LogFile     *string `toml:"log_file"`
		LogLevel    int     `toml:"log_level"`
	} `toml:"main"`
	Zones []dnsapi.ConfigZone `toml:"zones"`
}
