package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

var client = resty.New()

const API_BASE = "https://api.cloudflare.com/client/v4"

func getIpAddress() (string, error) {
	req := client.NewRequest()
	res, err := req.Get("https://ipecho.net/plain")
	if err != nil {
		return "", err
	}

	return string(res.Body()), nil
}

func getDnsEntries(zoneId string, token string) ([]DnsEntry, error) {
	req := client.NewRequest()
	req.Header.Add("authorization", "Bearer "+token)
	res, err := req.Get(fmt.Sprintf("%s/zones/%s/dns_records?type=A&per_page=5000", API_BASE, zoneId))
	if err != nil {
		return nil, err
	}

	var output DnsQueryResponse
	err = json.Unmarshal(res.Body(), &output)
	if err != nil {
		return nil, err
	}

	if !output.Success {
		return nil, fmt.Errorf("API returned non-success response")
	}

	return output.Result, nil
}

func updateDnsEntry(zoneId string, token string, entry DnsEntry) error {
	req := client.NewRequest()
	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("content-type", "application/json")

	body, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("error preparing request body: %w", err)
	}
	req.SetBody(body)

	if entry.ID == nil {
		l.Info("creating DNS entry", "domain", entry.Name)
		_, err := req.Post(fmt.Sprintf(`%s/zones/%s/dns_records`, API_BASE, zoneId))
		if err != nil {
			return fmt.Errorf("error creating DNS entry: %w", err)
		}
	} else {
		l.Info("updating DNS entry", "domain", entry.Name)
		_, err := req.Put(fmt.Sprintf(`%s/zones/%s/dns_records/%s`, API_BASE, zoneId, *entry.ID))
		if err != nil {
			return fmt.Errorf("error updating DNS entry: %w", err)
		}
	}

	return nil
}

func deleteDnsEntry(zoneId string, token string, entry DnsEntry) error {
	req := client.NewRequest()
	req.Header.Add("authorization", "Bearer "+token)

	l.Info("deleting DNS entry", "domain", entry.Name)
	_, err := req.Delete(fmt.Sprintf(`%s/zones/%s/dns_records/%s`, API_BASE, zoneId, *entry.ID))
	if err != nil {
		return fmt.Errorf("error creating DNS entry: %w", err)
	}

	return nil
}
