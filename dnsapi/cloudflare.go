package dnsapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseUrl = "https://api.cloudflare.com/client/v4/zones/%s/dns_records/"
)

type Cloudflare struct {
}

type DnsRecord struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IP        string `json:"content"`
	Proxiable bool   `json:"proxiable"`
	Proxied   bool   `json:"proxied"`
	TTL       int    `json:"ttl"`
	Settings  struct {
	} `json:"settings"`
	Meta struct {
	} `json:"meta"`
	Comment    any       `json:"comment"`
	Tags       []any     `json:"tags"`
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
	Priority   int       `json:"priority,omitempty"`
}

type ZoneInfo struct {
	Result     []DnsRecord `json:"result"`
	Success    bool        `json:"success"`
	Errors     []any       `json:"errors"`
	Messages   []any       `json:"messages"`
	ResultInfo struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		Count      int `json:"count"`
		TotalCount int `json:"total_count"`
		TotalPages int `json:"total_pages"`
	} `json:"result_info"`
}

type DomainIds map[string]DnsRecord

func (c *Cloudflare) GetName() string {
	return "cloudflare"
}

func (c *Cloudflare) setHeader(req *http.Request, zone ConfigZone) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+zone.Secret)
}

func (c *Cloudflare) callApi(method string, url string, body []byte, zone ConfigZone) (response string, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	c.setHeader(req, zone)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	rbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(rbody), nil
}

func (c *Cloudflare) getRecords(zone ConfigZone) (DomainIds, error) {
	domainIds := make(map[string]DnsRecord, 0)
	url := fmt.Sprintf(baseUrl, zone.Zone)
	res, err := c.callApi("GET", url, nil, zone)
	if err != nil {
		return nil, err
	}
	var result ZoneInfo
	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		return nil, err
	}
	for _, result := range result.Result {
		if result.Type == "A" {
			domainIds[result.Name] = result
		}
	}

	return domainIds, nil
}

func (c *Cloudflare) Update(zone ConfigZone, ip string) (code UpdateResult, err error) {

	records, err := c.getRecords(zone)
	if err != nil {
		return ResultError, fmt.Errorf("error retrieving records for zone: %v", err)
	}
	if len(records) == 0 {
		return ResultError, errors.New("zone contains zero A records")
	}

	for _, record := range zone.Domains {
		_, ok := records[record]
		if !ok {
			return ResultSkip, fmt.Errorf("zone not updated, type A record for %s not found", record)
		}
	}

	for _, record := range zone.Domains {
		if records[record].IP == ip {
			continue
		}
		id := records[record].ID
		url := fmt.Sprintf(baseUrl, zone.Zone) + id
		payload := map[string]interface{}{
			"type":    "A", // TODO: configurable
			"name":    record,
			"content": ip,
		}
		body, err := json.Marshal(payload)
		if err != nil {
			return ResultError, fmt.Errorf("json marshal error: %v", err)
		}
		res, err := c.callApi("PUT", url, body, zone)
		if err != nil {
			return ResultError, fmt.Errorf("api call error: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(res), &result); err != nil {
			return ResultError, fmt.Errorf("json unmarshal error: %v", err)
		}

		if success, ok := result["success"].(bool); !ok || !success {
			return ResultError, fmt.Errorf("cloudflare api error: %v", result["errors"])
		}

	}

	return ResultOk, nil
}
