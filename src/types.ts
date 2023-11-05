type ZoneSettings = {
  readonly zoneId: string;
  readonly token: string;
  readonly ttlSeconds?: number;
  readonly autoWww?: boolean;
  readonly autoDelete?: boolean;
  readonly autoDeleteAllowList?: string[];
  readonly autoDeleteBlockList?: string[];
  readonly domains: string[];
};

type DnsEntry = {
  readonly id: string | null; // null = create new
  readonly name: string; // domain
  readonly content: string; // IP address
  readonly ttl: number; // seconds
};

type DnsQueryResponse = {
  readonly result: DnsEntry[];
  readonly success: boolean;
  readonly result_info: {
    readonly page: number;
    readonly total_pages: number;
  };
};

export { ZoneSettings, DnsEntry, DnsQueryResponse };
