package main

import (
	"time"
	"net/http"
	"os"
	"runtime/debug"
	"fmt"
	"log"
	"crypto/tls"
	"io/ioutil"
	"encoding/json"

	"github.com/7joe7/godaddyactualizer/resources"
	"github.com/7joe7/godaddyactualizer/godaddy"
	"net/smtp"
)

var (
	config *resources.Config
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panicked. %v %s\n", r, string(debug.Stack()))
			log.Printf("Panicked. %v %s", r, string(debug.Stack()))
			err := sendEmail("To: You\nSubject: godaddyactualizer failed\nHello,\n\nsomething failed while running godaddyactualizer.\nPlease check it out.\n\nHave a nice day.\n\nJOT")
			if err != nil {
				log.Printf("Email sending failed. %v", err)
			}
			os.Exit(3)
		}
	}()

	// in cycle verify my actual IP address

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	c := http.Client{Transport: tr}

	confContent, err := ioutil.ReadFile(resources.CONF_STORE_PATH)
	if err != nil {
		panic(err)
	}
	config = &resources.Config{Domains: map[string]resources.Domain{}}
	err = json.Unmarshal(confContent, config)
	if err != nil {
		panic(err)
	}

	for {
		log.Printf("Going to verify actual IP address.")
		req, err := http.NewRequest("GET", "https://api.ipify.org?format=json", nil)
		if err != nil {
			panic(err)
		}
		resp, err := c.Do(req)
		if err != nil {
			panic(err)
		}
		res, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		air := &resources.ActualIpResponse{}
		err = json.Unmarshal(res, air)
		if err != nil {
			panic(err)
		}
		log.Printf("Real IP address verified as %s.", air.Ip)
		if air.Ip != "" && air.Ip != config.ActualIp {
			log.Printf("Actual IP address and real IP address don't comform.")
			err = sendEmail(fmt.Sprintf("To: You\nSubject: godaddyactualizer is changing IP records\nHello,\n\nnew IP address is detected. godaddyactualizer is going to change IP address to %s.\n\nHave a nice day.\n\nJOT", air.Ip))
			if err != nil {
				panic(err)
			}
			for domainName, domain := range config.Domains {
				for i := 0; i < len(domain.RecordsToActualize); i++ {
					err = godaddy.PutDomainsRecords(domainName, domain.RecordsToActualize[i], air.Ip, config.GoDaddyApiKey, config.GoDaddySecret)
					if err != nil {
						panic(err)
					}
				}
			}
			config.ActualIp = air.Ip
			newConfigContent, err := json.Marshal(config)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(resources.CONF_STORE_PATH, newConfigContent, os.FileMode(0777))
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(time.Minute)
	}
}

func sendEmail(message string) error {
	return smtp.SendMail(config.EmailSmtpWithPort, smtp.PlainAuth("", config.EmailAdminAddress, config.EmailPassword, config.EmailSmtp), config.EmailAdminAddress, []string{config.EmailAddress}, []byte(message))
}
