# Multi-RPC Failover Guide

This guide explains how to implement multi-RPC failover for the Archivas Explorer and SDK.

## Overview

Archivas now supports multiple seed nodes:
- **Primary**: `https://seed.archivas.ai` (Server A)
- **Secondary**: `https://seed2.archivas.ai` (Server D)
- **Tertiary**: `https://seed3.archivas.ai` (Server E)

Clients (Explorer, SDK, Farmers) should try primary → secondary → tertiary with automatic failover.

## Explorer Implementation

### Environment Variables

Add to `.env`:

```env
NEXT_PUBLIC_RPC_PRIMARY=https://seed.archivas.ai
NEXT_PUBLIC_RPC_SECONDARY=https://seed2.archivas.ai
NEXT_PUBLIC_RPC_TERTIARY=https://seed3.archivas.ai
```

### RPC Client (`lib/rpcClient.ts`)

```typescript
const RPCS = [
  process.env.NEXT_PUBLIC_RPC_PRIMARY!,
  process.env.NEXT_PUBLIC_RPC_SECONDARY!,
  process.env.NEXT_PUBLIC_RPC_TERTIARY!,
].filter(Boolean);

interface RpcResponse<T = any> {
  data: T;
  host: string; // Which RPC host succeeded
}

async function tryOnce(
  url: string,
  path: string,
  signal: AbortSignal,
  init?: RequestInit
): Promise<Response> {
  const res = await fetch(`${url}${path}`, {
    ...init,
    signal,
    headers: {
      'Content-Type': 'application/json',
      ...init?.headers,
    },
  });

  if (!res.ok) {
    throw new Error(`HTTP ${res.status}: ${res.statusText}`);
  }

  return res;
}

export async function rpcGet<T = any>(
  path: string,
  timeoutMs = 3000
): Promise<RpcResponse<T>> {
  const ac = new AbortController();
  const timeout = setTimeout(() => ac.abort(), timeoutMs);

  try {
    let lastErr: unknown;

    for (const base of RPCS) {
      try {
        const res = await tryOnce(base, path, ac.signal);
        const data = await res.json();
        clearTimeout(timeout);
        return { data, host: base };
      } catch (e) {
        lastErr = e;
        // Small jitter before next host (avoid thundering herd)
        await new Promise((r) => setTimeout(r, 150));
      }
    }

    throw lastErr ?? new Error('All RPCs failed');
  } finally {
    clearTimeout(timeout);
  }
}

export async function rpcPost<T = any>(
  path: string,
  body: any,
  timeoutMs = 5000
): Promise<RpcResponse<T>> {
  const ac = new AbortController();
  const timeout = setTimeout(() => ac.abort(), timeoutMs);

  try {
    let lastErr: unknown;

    for (const base of RPCS) {
      try {
        const res = await tryOnce(base, path, ac.signal, {
          method: 'POST',
          body: JSON.stringify(body),
        });
        const data = await res.json();
        clearTimeout(timeout);
        return { data, host: base };
      } catch (e) {
        lastErr = e;
        await new Promise((r) => setTimeout(r, 150));
      }
    }

    throw lastErr ?? new Error('All RPCs failed');
  } finally {
    clearTimeout(timeout);
  }
}
```

### Usage in API Routes

```typescript
// app/api/chainTip/route.ts
import { rpcGet } from '@/lib/rpcClient';

export async function GET() {
  try {
    const { data, host } = await rpcGet('/chainTip');
    return Response.json(data, {
      headers: {
        'X-RPC-Host': host, // Optional: show which RPC was used
      },
    });
  } catch (error) {
    return Response.json(
      { error: 'Failed to fetch chain tip' },
      { status: 503 }
    );
  }
}
```

### Usage in Components

```typescript
// app/components/ChainTip.tsx
'use client';

import { useEffect, useState } from 'react';
import { rpcGet } from '@/lib/rpcClient';

export function ChainTip() {
  const [tip, setTip] = useState<any>(null);
  const [host, setHost] = useState<string>('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchTip() {
      try {
        const { data, host: rpcHost } = await rpcGet('/chainTip');
        setTip(data);
        setHost(rpcHost);
      } catch (error) {
        console.error('Failed to fetch chain tip:', error);
      } finally {
        setLoading(false);
      }
    }

    fetchTip();
    const interval = setInterval(fetchTip, 5000); // Poll every 5s
    return () => clearInterval(interval);
  }, []);

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <p>Height: {tip?.height}</p>
      <p className="text-sm text-gray-500">RPC: {host}</p>
    </div>
  );
}
```

### UI Badge (Optional)

Show which RPC is being used:

```typescript
// app/components/RpcStatus.tsx
'use client';

import { useEffect, useState } from 'react';
import { rpcGet } from '@/lib/rpcClient';

export function RpcStatus() {
  const [host, setHost] = useState<string>('');

  useEffect(() => {
    async function checkRpc() {
      try {
        const { host: rpcHost } = await rpcGet('/healthz');
        setHost(rpcHost);
      } catch (error) {
        setHost('offline');
      }
    }

    checkRpc();
    const interval = setInterval(checkRpc, 30000); // Check every 30s
    return () => clearInterval(interval);
  }, []);

  if (!host) return null;

  const hostName = host.replace('https://', '').replace('http://', '');
  const isPrimary = hostName === 'seed.archivas.ai';

  return (
    <span
      className={`text-xs px-2 py-1 rounded ${
        isPrimary ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
      }`}
      title={`Using RPC: ${hostName}`}
    >
      {isPrimary ? '🟢' : '🟡'} {hostName}
    </span>
  );
}
```

## SDK Implementation

### RPC Client (`src/rpc/client.ts`)

```typescript
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
}
```

### Usage

```typescript
import { Rpc } from '@archivas/sdk';

// Multi-RPC with failover
const rpc = new Rpc({
  baseUrls: [
    'https://seed.archivas.ai',
    'https://seed2.archivas.ai',
    'https://seed3.archivas.ai',
  ],
  timeoutMs: 2500,
});

// Get chain tip
const tip = await rpc.getChainTip();
console.log('Current height:', tip.height);

// Get recent blocks
const blocks = await rpc.getRecentBlocks(10);
console.log('Recent blocks:', blocks.blocks);

// Get challenge
const challenge = await rpc.getChallenge();
console.log('Current challenge:', challenge.challenge);
```

### Backward Compatibility

The SDK still supports single `baseUrl`:

```typescript
// Single RPC (backward compatible)
const rpc = new Rpc({
  baseUrl: 'https://seed.archivas.ai',
});

// Or use baseUrls with single entry
const rpc = new Rpc({
  baseUrls: ['https://seed.archivas.ai'],
});
```

## Farmer Configuration

Farmers can connect to any seed node:

```bash
# Use primary
./archivas-farmer farm \
  --node https://seed.archivas.ai \
  --plots ./plots \
  --farmer-key YOUR_PRIVKEY

# Use secondary (if primary is slow)
./archivas-farmer farm \
  --node https://seed2.archivas.ai \
  --plots ./plots \
  --farmer-key YOUR_PRIVKEY

# Use tertiary
./archivas-farmer farm \
  --node https://seed3.archivas.ai \
  --plots ./plots \
  --farmer-key YOUR_PRIVKEY
```

### Recommended: Local Node

For zero-latency farming, run a local node:

```bash
# Start local node
./archivas-node --rpc 127.0.0.1:8080 --p2p :9090

# Connect farmer to local node
./archivas-farmer farm \
  --node http://127.0.0.1:8080 \
  --plots ./plots \
  --farmer-key YOUR_PRIVKEY
```

## Testing

### Test RPC Failover

```bash
# Test primary
curl https://seed.archivas.ai/healthz

# Test secondary
curl https://seed2.archivas.ai/healthz

# Test tertiary
curl https://seed3.archivas.ai/healthz

# Test all endpoints
for rpc in seed.archivas.ai seed2.archivas.ai seed3.archivas.ai; do
  echo "Testing $rpc:"
  curl -s "https://$rpc/chainTip" | jq '.height'
done
```

### Test SDK

```typescript
import { Rpc } from '@archivas/sdk';

const rpc = new Rpc({
  baseUrls: [
    'https://seed.archivas.ai',
    'https://seed2.archivas.ai',
    'https://seed3.archivas.ai',
  ],
});

// This should work even if one RPC is down
try {
  const tip = await rpc.getChainTip();
  console.log('Success:', tip);
} catch (error) {
  console.error('All RPCs failed:', error);
}
```

## Best Practices

1. **Timeout**: Use short timeouts (2-3 seconds) to quickly fail over
2. **Retry**: Implement exponential backoff between retries
3. **Rotation**: Rotate successful hosts to front (better performance)
4. **Monitoring**: Log which RPC is being used for debugging
5. **Fallback**: Always have at least 2 RPCs configured
6. **Local First**: Use local node for farmers when possible

## Troubleshooting

### All RPCs Failing

1. Check DNS: `dig seed.archivas.ai`
2. Check TLS: `curl -v https://seed.archivas.ai/healthz`
3. Check firewall: Ensure ports 443 are open
4. Check logs: `sudo journalctl -u archivas-node-seed2 -n 50`

### Slow Failover

1. Reduce timeout: `timeoutMs: 2000`
2. Check network latency: `ping seed.archivas.ai`
3. Use local node for farmers

### Rate Limiting

1. Check rate limits: `sudo tail -f /var/log/nginx/seed2.archivas.ai.error.log`
2. Adjust rate limits in Nginx config
3. Use multiple RPCs to distribute load

## Rollout Plan

1. **Deploy seed2**: Run `deploy/seed2/setup.sh` on Server D
2. **Deploy seed3**: Run `deploy/seed3/setup.sh` on Server E
3. **Update Explorer**: Add multi-RPC support, deploy
4. **Update SDK**: Add `baseUrls` support, publish npm package
5. **Update Docs**: Add multi-RPC failover guide
6. **Test**: Verify failover works in production

## Success Criteria

- ✅ Explorer works when one RPC is down
- ✅ SDK automatically fails over to secondary RPC
- ✅ Farmers can use any seed node
- ✅ No single point of failure
- ✅ Load is distributed across multiple RPCs

