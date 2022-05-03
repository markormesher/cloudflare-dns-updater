![CircleCI](https://img.shields.io/circleci/build/github/markormesher/cloudflare-dns-updater)

# Cloudflare DNS Updater

A utility for updating A records on [Cloudflare DNS](https://www.cloudflare.com/en-gb/dns/) to point at your current IP, which is ideal if you self-host Internet-facing services but do not have a static IP address. Once configured with a list of domains, this tool will:

- Check your IP address every 2 minutes (configurable).
- Add A records for any domains on your list that aren't configured on Cloudflare.
  - `www.` prefixes are handled automatically if you want them to be.
- Update existing A records when your IP address changes.
- Optionally remove any A records for domains you haven't configured.

:rocket: Jump to [quick-start example](#quick-start-docker-compose-example).

:whale: See releases on [ghcr.io](https://ghcr.io/markormesher/cloudflare-dns-updater).

_This work is heavily inspired by [timothymiller/cloudflare-ddns](https://github.com/timothymiller/cloudflare-ddns) and for the most part is just a slightly simplified and more opinionated version._

## Configuration

Configuration is split into two parts: a settings file with your list of domains (required) and environment variables for extra options (optional).

### Settings File (Domain List)

The main configuration - the list of domains - is stored as a JSON file containing a list of zones to manage. Each zone can have the following properties:

- `zoneId` - your Cloudflare zone ID
- `token` - a service token with permission to update DNS records
- `ttlSeconds` - the TTL value to set when creating new domains (optional, default 120)
- `autoWww` - whether or not to automatically create and update `www`-prefixed versions of the domains listed (optional, default false)
- `autoDelete` - whether or not to automatically remove DNS records for domains you have not specified (optional, default false)
- `autoDeleteAllowList` - list of regexes for domains that eligible for auto-removal (optional, default empty = all domains are eligible)
- `autoDeleteBlockList` - list of regexes for domains that ineligible for auto-removal (optional, default empty = all domains are eligible)
- `domains` - the domains you wish to manage; each should include the TLD and have no leading `.`

For example:

```json
[
  {
    "zoneId": "xzy-zone1",
    "token": "abracadabra",
    "ttlSeconds": 120,
    "autoWww": true,
    "autoDelete": true,
    "autoDeleteBlockList": ["legacy.example.com"],
    "domains": [
      "example.com",
      "sub1.example.com",
      "sub2.example.com",
      "sub3.example.com"
    ]
  },
  {
    "zoneId": "xzy-zone2",
    "token": "abracadabra",
    "ttlSeconds": 120,
    "autoWww": true,
    "autoDelete": true,
    "autoDeleteAllowList": ["sub.*.example.net"],
    "domains": [
      "example.net",
      "sub1.example.net",
      "sub2.example.net",
      "sub3.example.net"
    ]
  }
]
```

See [types.ts](./src/types.ts) for the full details.

### Environment Variables

| Variable                 | Required? | Description                                               | Default          |
| ------------------------ | --------- | --------------------------------------------------------- | ---------------- |
| `SETTINGS_FILE`          | no        | Path to where your JSON settings file is.                 | `/settings.json` |
| `CHECK_INTERVAL_SECONDS` | no        | How often to re-check your IP address and update records. | 120              |

## Quick-Start Docker-Compose Example

```yaml
version: "3.8"

services:
  cloudflare-dns-updater:
    image: ghcr.io/markormesher/cloudflare-dns-updater:VERSION
    restart: unless-stopped
    volumes:
      - ./settings.json:/settings.json:ro
```
