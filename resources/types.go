package resources

type ActualIpResponse struct {
	Ip string
}

type Config struct {
	ActualIp          string
	GoDaddyApiKey     string
	GoDaddySecret     string
	Domains           map[string]Domain
	EmailSmtpWithPort string
	EmailAddress      string
	EmailAdminAddress string
	EmailSmtp         string
	EmailPassword     string
}

type Domain struct {
	RecordsToActualize []string
}

type PutDomainRecordRequestBody struct {
	Data string `json:"data"`
	Name string `json:"name"`
	Type string `json:"type"`
}
