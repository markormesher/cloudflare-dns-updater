import { readFileSync, existsSync } from "fs";
import axios from "axios";
import { IDnsEntry, IDnsQueryResponse, IZoneSettings } from "./types";

const API_BASE = "https://api.cloudflare.com/client/v4";
const REPEAT_INTERVAL_MS = parseInt(process.env.CHECK_INTERVAL_SECONDS) * 1000 || 2 * 60 * 1000;

function log(msg: string) {
  console.log(`[${new Date().toISOString()}] ${msg}`);
}

function readSettings(): IZoneSettings[] {
  const settingsFile = process.env.SETTINGS_FILE || "/settings.json";
  if (!existsSync(settingsFile)) {
    throw new Error(`Settings file ${settingsFile} does not exist!`);
  }
  return JSON.parse(readFileSync(settingsFile).toString()) as IZoneSettings[];
}

async function getDnsEntries(zoneId: string, token: string): Promise<IDnsEntry[]> {
  const req = await axios.get(`${API_BASE}/zones/${zoneId}/dns_records?type=A&per_page=5000`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  const data: IDnsQueryResponse = req.data as IDnsQueryResponse;
  return data.result;
}

async function updateDnsEntry(zoneId: string, token: string, entry: IDnsEntry): Promise<void> {
  const payload = {
    name: entry.name,
    content: entry.content,
    ttl: entry.ttl,
    type: "A",
  };

  const headers = { Authorization: `Bearer ${token}` };

  if (entry.id == null) {
    log(`Creating DNS entry: ${entry.name} -> ${entry.content}`);
    await axios.post(`${API_BASE}/zones/${zoneId}/dns_records`, payload, { headers });
  } else {
    log(`Updating DNS entry: ${entry.name} -> ${entry.content}`);
    await axios.put(`${API_BASE}/zones/${zoneId}/dns_records/${entry.id}`, payload, { headers });
  }
}

async function removeDnsEntry(zoneId: string, token: string, entry: IDnsEntry): Promise<void> {
  const headers = { Authorization: `Bearer ${token}` };

  log(`Deleting DNS entry: ${entry.name}`);
  await axios.delete(`${API_BASE}/zones/${zoneId}/dns_records/${entry.id}`, { headers });
}

async function updateDomains() {
  const zones = readSettings();
  const currentIp = (await axios.get("https://ipecho.net/plain")).data;
  log(`Current IP is ${currentIp}`);
  for (const zone of zones) {
    const { zoneId, token, ttlSeconds, domains, autoWww } = zone;
    const dnsEntires = await getDnsEntries(zoneId, token);

    // create missing domains
    for (const domain of domains) {
      if (!dnsEntires.some((e) => e.name === domain)) {
        await updateDnsEntry(zoneId, token, { id: null, name: domain, content: currentIp, ttl: ttlSeconds });
      }
      if (autoWww && !dnsEntires.some((e) => e.name === "www." + domain)) {
        await updateDnsEntry(zoneId, token, { id: null, name: "www." + domain, content: currentIp, ttl: ttlSeconds });
      }
    }

    for (const entry of dnsEntires) {
      // remove undeclared domains
      if ((!autoWww || !entry.name.startsWith("www.")) && !domains.includes(entry.name)) {
        await removeDnsEntry(zoneId, token, entry);
        continue;
      }
      if (autoWww && entry.name.startsWith("www.") && !domains.includes(entry.name.replace(/^www\./, ""))) {
        await removeDnsEntry(zoneId, token, entry);
        continue;
      }

      // update out of date domains
      if (entry.content !== currentIp) {
        await updateDnsEntry(zoneId, token, { ...entry, content: currentIp, ttl: ttlSeconds });
      }
    }
  }

  setTimeout(updateDomains, REPEAT_INTERVAL_MS);
}

updateDomains();
