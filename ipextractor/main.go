package ipextractor

import (
	"fmt"
	"io"
	"net/http"
)

type IpExtractor interface {
	GetIp() (ip string, err error)
}

func GetIpExtractor(name string) (IpExtractor, error) {
	switch name {
	case "aws":
		return &Aws{}, nil
	case "cloudflare":
		return &Cloudflare{}, nil
	case "parkyourdomain":
		return &Parkyourdomain{}, nil
	default:
		return nil, fmt.Errorf("unknown ip extractor type %s", name)
	}
}

func simpleGetRequest(url string) (body string, err error) {
	var res *http.Response
	var data []byte

	res, err = http.Get(url)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("response error: %v", err)
		return
	}

	data, err = io.ReadAll(res.Body)
	if err == nil {
		body = string(data)
	}
	return
}
