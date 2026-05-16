package main

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"
)

func main() {
	settings, err := getSettings()
	if err != nil {
		slog.Error("error getting settings", "error", err)
		os.Exit(1)
	}

	interval, err := getCheckInterval()
	if err != nil {
		slog.Error("error getting check interval", "error", err)
		os.Exit(1)
	}

	wrappedRun := func() {
		err := runUpdate(settings)
		if err != nil {
			slog.Error("error running update", "error", err)
			os.Exit(1)
		}
	}

	if interval == 0 {
		slog.Info("no check interval set - running once then exiting")
		wrappedRun()
	} else {
		slog.Info("running repeatedly", "interval", interval)
		for ; true; <-time.Tick(time.Duration(interval) * time.Second) {
			wrappedRun()
		}
	}
}

func runUpdate(settings []ZoneSettings) error {
	currentIP, err := getIPAddress()
	if err != nil {
		return fmt.Errorf("error getting current IP address: %w", err)
	}
	slog.Info("got current IP", "ip", currentIP)

	for _, zone := range settings {
		slog.Info("checking zone", "zone", zone.ZoneID)

		autoDeleteAllowList := make([]*regexp.Regexp, 0, len(zone.AutoDeleteAllowList))
		for _, patternStr := range zone.AutoDeleteAllowList {
			pattern, err := regexp.Compile(patternStr)
			if err != nil {
				return fmt.Errorf("could not compile allow-list regex: %w", err)
			}
			autoDeleteAllowList = append(autoDeleteAllowList, pattern)
		}

		autoDeleteBlockList := make([]*regexp.Regexp, 0, len(zone.AutoDeleteBlockList))
		for _, patternStr := range zone.AutoDeleteBlockList {
			pattern, err := regexp.Compile(patternStr)
			if err != nil {
				return fmt.Errorf("could not compile block-list regex: %w", err)
			}
			autoDeleteBlockList = append(autoDeleteBlockList, pattern)
		}

		entries, err := getDNSEntries(zone.ZoneID, zone.Token)
		if err != nil {
			return fmt.Errorf("error getting current entries: %w", err)
		}

		knownDomains := make([]string, 0, len(entries))
		for _, e := range entries {
			knownDomains = append(knownDomains, e.Name)
		}

		zoneTTL := zone.TTLSeconds
		if zoneTTL == 0 {
			zoneTTL = 120
		}

		// creating missing entries
		for _, domain := range zone.Domains {
			if !slices.Contains(knownDomains, domain) {
				err := updateDNSEntry(zone.ZoneID, zone.Token, DNSEntry{
					Name:    domain,
					Content: currentIP,
					TTL:     zoneTTL,
					Type:    "A",
				})
				if err != nil {
					return fmt.Errorf("error creating entry: %w", err)
				}

				if zone.AutoWWW {
					err := updateDNSEntry(zone.ZoneID, zone.Token, DNSEntry{
						Name:    "www." + domain,
						Content: currentIP,
						TTL:     zoneTTL,
						Type:    "A",
					})
					if err != nil {
						return fmt.Errorf("error creating entry: %w", err)
					}
				}
			}
		}

		// update existing entries
		for _, entry := range entries {
			// remove undeclared domains
			delete := false

			if (!zone.AutoWWW || !strings.HasPrefix(entry.Name, "www.")) && !slices.Contains(zone.Domains, entry.Name) {
				delete = true
			}

			if (zone.AutoWWW && strings.HasPrefix(entry.Name, "www.")) && !slices.Contains(zone.Domains, strings.TrimPrefix(entry.Name, "www.")) {
				delete = true
			}

			if delete && domainDeletionAllowed(zone.AutoDelete, autoDeleteAllowList, autoDeleteBlockList, entry.Name) {
				err := deleteDNSEntry(zone.ZoneID, zone.Token, entry)
				if err != nil {
					return fmt.Errorf("error deleting entry: %w", err)
				}
				continue
			}

			// update out of sync entities
			if entry.Content != currentIP || entry.TTL != zoneTTL {
				entry.Content = currentIP
				entry.TTL = zoneTTL

				err := updateDNSEntry(zone.ZoneID, zone.Token, entry)
				if err != nil {
					return fmt.Errorf("error updating entry: %w", err)
				}
			}
		}
	}

	return nil
}

func domainDeletionAllowed(autoDelete bool, autoDeleteAllowList []*regexp.Regexp, autoDeleteBlockList []*regexp.Regexp, domain string) bool {
	// semantics:
	// - if auto delete is disabled, we obviously CANNOT delete
	// - if any item on the block list matches, we CANNOT delete
	// - if an allow list is set:
	//   - if any item on the list matches, we CAN delete
	//   - if no item on the list matches, we CANNOT delete
	// - if there is no allow list, we CAN delete

	if !autoDelete {
		slog.Info("not deleting domain because auto-delete is disabled", "domain", domain)
		return false
	}

	for _, pattern := range autoDeleteBlockList {
		if pattern.MatchString(domain) {
			slog.Info("not deleting domain because it matches an entry on the auto-delete block list", "domain", domain)
			return false
		}
	}

	if len(autoDeleteAllowList) > 0 {
		for _, pattern := range autoDeleteAllowList {
			if pattern.MatchString(domain) {
				return true
			}
		}
		slog.Info("not deleting domain because it does not match any entry on the auto-delete allow list", "domain", domain)
		return false
	} else {
		return true
	}
}
