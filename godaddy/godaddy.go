package godaddy

import (
	"net/http"
	"crypto/tls"
	"fmt"
	"encoding/json"
	"bytes"

	"github.com/7joe7/godaddyactualizer/resources"
)

func putDomainsRecords(domain, record, newIp, key, secret string) error {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	c := http.Client{Transport: tr}

	reqBody := resources.PutDomainRecordRequestBody{Data: newIp, Name: record, Type:"A"}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/%s", domain, record), bytes.NewReader(reqBodyBytes))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("sso-key %s:%s", key, secret))

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Something went wrong. %s", resp.Status)
	}
	return nil
}
