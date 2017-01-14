package godaddy

func PutDomainsRecords(domain, record, newIp, key, secret string) error {
	return putDomainsRecords(domain, record, newIp, key, secret)
}