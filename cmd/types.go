package main

type ZoneSettings struct {
	ZoneID              string   `json:"zoneId"`
	Token               string   `json:"token"`
	TTLSeconds          int      `json:"ttlSeconds"`
	AutoWWW             bool     `json:"autoWww"`
	AutoDelete          bool     `json:"autoDelete"`
	AutoDeleteAllowList []string `json:"autoDeleteAllowList"`
	AutoDeleteBlockList []string `json:"autoDeleteBlockList"`
	Domains             []string `json:"domains"`
}

type DnsEntry struct {
	ID      *string `json:"id"` // nil = create new
	Name    string  `json:"name"`
	Content string  `json:"content"`
	TTL     int     `json:"ttl"`
	Type    string  `json:"type"`
}

type DnsQueryResponse struct {
	Result     []DnsEntry `json:"result"`
	Success    bool       `json:"success"`
	ResultInfo struct {
		Page       int `json:"page"`
		TotalPages int `json:"total_pages"`
	}
}
