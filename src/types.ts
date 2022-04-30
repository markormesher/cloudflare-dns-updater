interface IZoneSettings {
  readonly zoneId: string;
  readonly token: string;
  readonly ttlSeconds?: number;
  readonly autoWww?: boolean;
  readonly autoDelete?: boolean;
  readonly domains: string[];
}

interface IDnsEntry {
  readonly id: string;
  readonly name: string; // domain
  readonly content: string; // IP address
  readonly ttl: number; // seconds
}

interface IDnsQueryResponse {
  readonly result: IDnsEntry[];
  readonly success: boolean;
  readonly result_info: {
    readonly page: number;
    readonly total_pages: number;
  };
}

export { IZoneSettings, IDnsEntry, IDnsQueryResponse };
