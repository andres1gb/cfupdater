package main

import (
	"cfupdater/dnsapi"
	"cfupdater/ipextractor"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

func main() {
	var ip string

	config, err := readConfig()
	if err != nil {
		panic(err.Error())
	}

	extractor, err := ipextractor.GetIpExtractor(config.Main.IpExtractor)
	if err != nil {
		panic(err.Error())
	}

	ip, err = extractor.GetIp()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("current ip is", ip)

	for _, zone := range config.Zones {
		api, err := dnsapi.GetDnsApi(zone.DnsApi)
		if err != nil {
			panic(err.Error())
		}
		_, err = api.Update(zone, ip)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}

func readConfig() (*Config, error) {
	var config Config

	file, err := os.Open("config.toml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = toml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
