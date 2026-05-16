package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

const CloudflareAPIBase = "https://api.cloudflare.com/client/v4"

func getIPAddress() (string, error) {
	req, err := http.NewRequest("GET", "https://ipecho.net/plain", nil)
	if err != nil {
		return "", fmt.Errorf("error building request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return "", fmt.Errorf("non-OK response: %v", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %w", err)
	}

	return string(body), nil
}

func getDNSEntries(zoneID string, token string) ([]DNSEntry, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/zones/%s/dns_records?type=A&per_page=5000", CloudflareAPIBase, zoneID), nil)
	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	req.Header.Add("authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("non-OK response: %v", res.Status)
	}

	var output DNSQueryResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&output)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if !output.Success {
		return nil, fmt.Errorf("non-success response")
	}

	return output.Result, nil
}

func updateDNSEntry(zoneID string, token string, entry DNSEntry) error {
	body, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("error preparing request body: %w", err)
	}

	var req *http.Request
	if entry.ID == nil {
		slog.Info("creating DNS entry", "domain", entry.Name)
		req, err = http.NewRequest("POST", fmt.Sprintf(`%s/zones/%s/dns_records`, CloudflareAPIBase, zoneID), bytes.NewReader(body))
	} else {
		slog.Info("updating DNS entry", "domain", entry.Name)
		req, err = http.NewRequest("PUT", fmt.Sprintf(`%s/zones/%s/dns_records/%s`, CloudflareAPIBase, zoneID, *entry.ID), bytes.NewReader(body))
	}
	if err != nil {
		return fmt.Errorf("error building request: %w", err)
	}

	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("non-OK response: %v", res.Status)
	}

	return nil
}

func deleteDNSEntry(zoneID string, token string, entry DNSEntry) error {
	slog.Info("deleting DNS entry", "domain", entry.Name)

	req, err := http.NewRequest("DELETE", fmt.Sprintf(`%s/zones/%s/dns_records/%s`, CloudflareAPIBase, zoneID, *entry.ID), nil)
	if err != nil {
		return fmt.Errorf("error building request: %w", err)
	}

	req.Header.Add("authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("non-OK response: %v", res.Status)
	}

	return nil
}
