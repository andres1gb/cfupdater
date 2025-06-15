package ipextractor

const parkyourdomainUrl = "https://dynamicdns.park-your-domain.com/getip"

type Parkyourdomain struct{}

func (p *Parkyourdomain) GetIp() (ip string, err error) {
	return simpleGetRequest(parkyourdomainUrl)
}
