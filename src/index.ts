import { readFileSync, existsSync } from "fs";
import { DnsEntry, DnsQueryResponse, ZoneSettings } from "./types.js";

const API_BASE = "https://api.cloudflare.com/client/v4";
const REPEAT_INTERVAL_MS = parseInt(process.env.CHECK_INTERVAL_SECONDS ?? "120") * 1000 ?? 2 * 60 * 1000;

function log(msg: string, params?: Record<string, unknown>) {
  console.log(`[${new Date().toISOString()}] ${msg}`, { ...params });
}

function readSettings(): ZoneSettings[] {
  const settingsFile = process.env.SETTINGS_FILE ?? "/settings.json";
  if (!existsSync(settingsFile)) {
    throw new Error(`Settings file ${settingsFile} does not exist!`);
  }
  return JSON.parse(readFileSync(settingsFile).toString()) as ZoneSettings[];
}

async function getDnsEntries(zoneId: string, token: string): Promise<DnsEntry[]> {
  const req = await fetch(`${API_BASE}/zones/${zoneId}/dns_records?type=A&per_page=5000`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  const data: DnsQueryResponse = (await req.json()) as DnsQueryResponse;
  return data.result;
}

async function updateDnsEntry(zoneId: string, token: string, entry: DnsEntry): Promise<void> {
  const payload = {
    name: entry.name,
    content: entry.content,
    ttl: entry.ttl,
    type: "A",
  };

  const headers = { Authorization: `Bearer ${token}`, "Content-type": "application/json" };

  if (entry.id == null) {
    log(`Creating DNS entry: ${entry.name} -> ${entry.content}`);
    await fetch(`${API_BASE}/zones/${zoneId}/dns_records`, { method: "POST", headers, body: JSON.stringify(payload) });
  } else {
    log(`Updating DNS entry: ${entry.name} -> ${entry.content}`);
    await fetch(`${API_BASE}/zones/${zoneId}/dns_records/${entry.id}`, {
      method: "PUT",
      headers,
      body: JSON.stringify(payload),
    });
  }
}

async function removeDnsEntry(zoneId: string, token: string, entry: DnsEntry): Promise<void> {
  const headers = { Authorization: `Bearer ${token}` };

  log(`Deleting DNS entry: ${entry.name}`);
  await fetch(`${API_BASE}/zones/${zoneId}/dns_records/${entry.id}`, { method: "DELETE", headers });
}

function domainDeletionAllowed(
  autoDelete: boolean,
  autoDeleteAllowList: string[],
  autoDeleteBlockList: string[],
  domain: string,
): boolean {
  // semantics:
  // - if auto delete is disabled, we obviously CANNOT delete
  // - if any item on the block list matches, we CANNOT delete
  // - if an allow list is set:
  //   - if any item on the list matches, we CAN delete
  //   - if no item on the list matches, we CANNOT delete
  // - if there is no allow list, we CAN delete

  if (!autoDelete) {
    log(`Not deleting ${domain} because auto-delete is disabled`);
    return false;
  }

  for (const regex of autoDeleteBlockList ?? []) {
    if (new RegExp(regex).test(domain)) {
      log(`Not deleting ${domain} because it matches an entry on the auto-delete block list`);
      return false;
    }
  }

  if (autoDeleteAllowList && autoDeleteAllowList.length > 0) {
    for (const regex of autoDeleteAllowList ?? []) {
      if (new RegExp(regex).test(domain)) {
        return true;
      }
    }
    log(`Not deleting ${domain} because it does not match any entry on the auto-delete allow list`);
    return false;
  } else {
    return true;
  }
}

async function getCurrentIp(): Promise<string> {
  let tries = 0;
  const maxTries = 3;
  while (tries < maxTries) {
    ++tries;
    try {
      const req = await fetch("https://ipecho.net/plain");
      return req.text();
    } catch (error) {
      log(`Failed to get current IP on attempt #${tries}, waiting 5s before trying again`, { error });
      await new Promise((resolve) => setTimeout(resolve, 5000));
    }
  }
  throw new Error(`Failed to get current IP after ${maxTries} attempts`);
}

async function updateDomains() {
  const zones = readSettings();
  const currentIp = await getCurrentIp();
  log(`Current IP is ${currentIp}`);
  for (const zone of zones) {
    const {
      zoneId,
      token,
      ttlSeconds,
      domains,
      autoWww,
      autoDelete: autoDeleteRaw,
      autoDeleteAllowList: autoDeleteAllowListRaw,
      autoDeleteBlockList: autoDeleteBlockListRaw,
    } = zone;
    const autoDelete = autoDeleteRaw ?? false;
    const autoDeleteAllowList = autoDeleteAllowListRaw ?? [];
    const autoDeleteBlockList = autoDeleteBlockListRaw ?? [];

    const dnsEntires = await getDnsEntries(zoneId, token);

    // create missing domains
    for (const domain of domains) {
      if (!dnsEntires.some((e) => e.name === domain)) {
        await updateDnsEntry(zoneId, token, { id: null, name: domain, content: currentIp, ttl: ttlSeconds ?? 120 });
      }
      if (autoWww && !dnsEntires.some((e) => e.name === "www." + domain)) {
        await updateDnsEntry(zoneId, token, {
          id: null,
          name: "www." + domain,
          content: currentIp,
          ttl: ttlSeconds ?? 120,
        });
      }
    }

    for (const entry of dnsEntires) {
      // remove undeclared domains
      if ((!autoWww || !entry.name.startsWith("www.")) && !domains.includes(entry.name)) {
        if (domainDeletionAllowed(autoDelete, autoDeleteAllowList, autoDeleteBlockList, entry.name)) {
          await removeDnsEntry(zoneId, token, entry);
        }
        continue;
      }
      if (autoWww && entry.name.startsWith("www.") && !domains.includes(entry.name.replace(/^www\./, ""))) {
        if (domainDeletionAllowed(autoDelete, autoDeleteAllowList, autoDeleteBlockList, entry.name)) {
          await removeDnsEntry(zoneId, token, entry);
        }
        continue;
      }

      // update out of date domains
      if (entry.content !== currentIp) {
        await updateDnsEntry(zoneId, token, { ...entry, content: currentIp, ttl: ttlSeconds ?? 120 });
      }
    }
  }
}

setInterval(() => {
  updateDomains().catch((error) => {
    log("Encountered an error during update - exiting");
    throw error;
  });
}, REPEAT_INTERVAL_MS);
