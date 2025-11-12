/**
 * Example SDK RPC Client with Multi-RPC Failover
 * Copy this to your SDK repo: src/rpc/client.ts
 */

type RpcInit = {
  baseUrl?: string; // Deprecated: use baseUrls instead
  baseUrls?: string[];
  timeoutMs?: number;
};

export class Rpc {
  private hosts: string[];
  private timeoutMs: number;

  constructor(opts: RpcInit) {
    const list =
      opts.baseUrls?.filter(Boolean) ?? (opts.baseUrl ? [opts.baseUrl] : []);

    if (list.length === 0) {
      throw new Error('Rpc requires baseUrl or baseUrls');
    }

    this.hosts = list;
    this.timeoutMs = opts.timeoutMs ?? 3000;
  }

  private async fetchAny<T = any>(
    path: string,
    init?: RequestInit
  ): Promise<T> {
    let lastErr: unknown;

    for (let i = 0; i < this.hosts.length; i++) {
      const host = this.hosts[i];
      const ac = new AbortController();
      const timeout = setTimeout(() => ac.abort(), this.timeoutMs);

      try {
        const res = await fetch(`${host}${path}`, {
          ...init,
          signal: ac.signal,
          headers: {
            'Content-Type': 'application/json',
            ...init?.headers,
          },
        });

        clearTimeout(timeout);

        if (!res.ok) {
          throw new Error(`HTTP ${res.status}: ${res.statusText}`);
        }

        const data = await res.json();

        // Simple host rotation on success (move successful host to front)
        if (i !== 0) {
          this.hosts.unshift(...this.hosts.splice(i, 1));
        }

        return data;
      } catch (e) {
        clearTimeout(timeout);
        lastErr = e;
        continue;
      }
    }

    throw lastErr ?? new Error('All RPCs failed');
  }

  // Chain endpoints
  async getChainTip() {
    return this.fetchAny<{ height: string; hash: string; difficulty: string }>(
      '/chainTip'
    );
  }

  async getRecentBlocks(n = 20) {
    return this.fetchAny<{ blocks: any[]; count: number }>(
      `/recentBlocks?count=${n}`
    );
  }

  async getBlockByHeight(height: number | string) {
    return this.fetchAny(`/block/${height}`);
  }

  // Challenge endpoints
  async getChallenge() {
    return this.fetchAny<{
      challenge: number[];
      difficulty: number;
      height: number;
    }>('/challenge');
  }

  // Health endpoints
  async getHealth() {
    return this.fetchAny<{ ok: boolean; height: number; peers: number }>(
      '/healthz'
    );
  }

  // Account endpoints
  async getBalance(address: string) {
    return this.fetchAny<{ address: string; balance: number; nonce: number }>(
      `/balance/${address}`
    );
  }

  async getAccounts() {
    return this.fetchAny<{ count: number; accounts: any[] }>('/accounts');
  }

  // Transaction endpoints
  async submitTx(tx: any) {
    return this.fetchAny('/submitTx', {
      method: 'POST',
      body: JSON.stringify(tx),
    });
  }
}

// Usage example:
/*
import { Rpc } from '@archivas/sdk';

const rpc = new Rpc({
  baseUrls: [
    'https://seed.archivas.ai',
    'https://seed2.archivas.ai',
    'https://seed3.archivas.ai',
  ],
  timeoutMs: 2500,
});

const tip = await rpc.getChainTip();
console.log('Current height:', tip.height);
*/

