/**
 * Example RPC Client for Explorer
 * Copy this to your Explorer repo: lib/rpcClient.ts
 */

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

