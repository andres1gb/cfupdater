package ipextractor

const awsUrl = "https://checkip.amazonaws.com"

type Aws struct{}

func (a *Aws) GetIp() (ip string, err error) {
	return simpleGetRequest(awsUrl)
}
