#!/usr/bin/env node
/**
 * Archivas Seed2 Relay - Health & Metrics Service
 * Provides monitoring endpoints for the RPC relay
 */

const fastify = require('fastify')({ logger: true });
const { register, Counter, Histogram, Gauge, collectDefaultMetrics } = require('prom-client');
const axios = require('axios');
const { execSync } = require('child_process');
const fs = require('fs');

// Configuration
const PORT = process.env.PORT || 9090;
const UPSTREAM_URL = process.env.UPSTREAM_URL || 'https://seed.archivas.ai';
const CACHE_DIR = process.env.CACHE_DIR || '/var/cache/nginx/archivas_rpc';

// Enable default metrics (CPU, memory, etc.)
collectDefaultMetrics({ prefix: 'seed2_' });

// Custom metrics
const cacheHits = new Counter({
  name: 'seed2_cache_hits_total',
  help: 'Total number of cache hits',
  labelNames: ['endpoint']
});

const cacheMisses = new Counter({
  name: 'seed2_cache_miss_total',
  help: 'Total number of cache misses',
  labelNames: ['endpoint']
});

const upstreamLatency = new Histogram({
  name: 'seed2_upstream_latency_ms',
  help: 'Upstream request latency in milliseconds',
  buckets: [10, 50, 100, 500, 1000, 2000, 5000]
});

const rateLimited = new Counter({
  name: 'seed2_rate_limited_total',
  help: 'Total number of rate-limited requests'
});

const txSubmitSuccess = new Counter({
  name: 'seed2_tx_submit_success_total',
  help: 'Total number of successful transaction submissions'
});

const txSubmitFailed = new Counter({
  name: 'seed2_tx_submit_failed_total',
  help: 'Total number of failed transaction submissions'
});

const cacheSize = new Gauge({
  name: 'seed2_cache_size_bytes',
  help: 'Current cache size in bytes'
});

const cacheFreeSpace = new Gauge({
  name: 'seed2_cache_free_space_percent',
  help: 'Free space percentage on cache disk'
});

// Helper: Get cache statistics
function getCacheStats() {
  try {
    // Get cache size
    const sizeOutput = execSync(`du -sb ${CACHE_DIR} 2>/dev/null | cut -f1`).toString().trim();
    const size = parseInt(sizeOutput) || 0;
    cacheSize.set(size);
    
    // Get disk free space percentage
    const dfOutput = execSync(`df ${CACHE_DIR} | tail -1 | awk '{print $5}'`).toString().trim();
    const usedPercent = parseInt(dfOutput.replace('%', '')) || 0;
    const freePercent = 100 - usedPercent;
    cacheFreeSpace.set(freePercent);
    
    // Get cache file count
    const fileCount = execSync(`find ${CACHE_DIR} -type f 2>/dev/null | wc -l`).toString().trim();
    
    return {
      sizeBytes: size,
      sizeMB: (size / 1024 / 1024).toFixed(2),
      freeSpacePercent: freePercent,
      fileCount: parseInt(fileCount) || 0
    };
  } catch (error) {
    return {
      sizeBytes: 0,
      sizeMB: '0',
      freeSpacePercent: 0,
      fileCount: 0,
      error: error.message
    };
  }
}

// Helper: Get cache hit ratio from Nginx logs
function getCacheHitRatio() {
  try {
    const logFile = '/var/log/nginx/seed2-access.log';
    if (!fs.existsSync(logFile)) {
      return { hits: 0, misses: 0, expired: 0, ratio: 0 };
    }
    
    // Count last 1000 requests
    const logs = execSync(`tail -1000 ${logFile} | grep -oP 'X-Cache-Status: \\K\\w+' | sort | uniq -c`).toString();
    
    let hits = 0, misses = 0, expired = 0;
    logs.split('\n').forEach(line => {
      const match = line.trim().match(/(\d+)\s+(\w+)/);
      if (match) {
        const count = parseInt(match[1]);
        const status = match[2];
        if (status === 'HIT') hits = count;
        else if (status === 'MISS') misses = count;
        else if (status === 'EXPIRED') expired = count;
      }
    });
    
    const total = hits + misses + expired;
    const ratio = total > 0 ? ((hits + expired) / total * 100).toFixed(2) : 0;
    
    return { hits, misses, expired, total, ratio: parseFloat(ratio) };
  } catch (error) {
    return { hits: 0, misses: 0, expired: 0, total: 0, ratio: 0 };
  }
}

// Helper: Check upstream health
async function checkUpstream() {
  const start = Date.now();
  try {
    const response = await axios.get(`${UPSTREAM_URL}/chainTip`, { timeout: 5000 });
    const latency = Date.now() - start;
    upstreamLatency.observe(latency);
    
    return {
      healthy: response.status === 200,
      latency,
      height: response.data.height || null
    };
  } catch (error) {
    const latency = Date.now() - start;
    upstreamLatency.observe(latency);
    
    return {
      healthy: false,
      latency,
      error: error.message
    };
  }
}

// Routes

/**
 * GET /health - Simple health check
 * Returns 200 if service is running
 */
fastify.get('/health', async (request, reply) => {
  return { ok: true, timestamp: Date.now() };
});

/**
 * GET /ready - Readiness check
 * Verifies upstream reachability and cache disk space
 * Returns 200 if ready, 503 if not
 */
fastify.get('/ready', async (request, reply) => {
  const upstream = await checkUpstream();
  const cache = getCacheStats();
  
  const ready = upstream.healthy && cache.freeSpacePercent > 10;
  
  if (!ready) {
    reply.code(503);
  }
  
  return {
    ready,
    checks: {
      upstream: {
        healthy: upstream.healthy,
        latency: upstream.latency
      },
      cache: {
        freeSpace: cache.freeSpacePercent,
        sufficient: cache.freeSpacePercent > 10
      }
    },
    timestamp: Date.now()
  };
});

/**
 * GET /status - Detailed status information
 * Returns comprehensive relay status
 */
fastify.get('/status', async (request, reply) => {
  const upstream = await checkUpstream();
  const cache = getCacheStats();
  const hitRatio = getCacheHitRatio();
  
  return {
    relay: upstream.healthy ? 'healthy' : 'degraded',
    cache: 'enabled',
    backend: 'seed1',
    upstream: {
      url: UPSTREAM_URL,
      healthy: upstream.healthy,
      latency_ms: upstream.latency,
      height: upstream.height
    },
    cache_stats: {
      size_mb: parseFloat(cache.sizeMB),
      free_space_percent: cache.freeSpacePercent,
      file_count: cache.fileCount,
      hit_ratio_5m: hitRatio.ratio,
      hits: hitRatio.hits,
      misses: hitRatio.misses,
      expired: hitRatio.expired
    },
    timestamp: Date.now()
  };
});

/**
 * GET /metrics - Prometheus metrics
 * Returns metrics in Prometheus format
 */
fastify.get('/metrics', async (request, reply) => {
  // Update cache metrics before exporting
  getCacheStats();
  getCacheHitRatio();
  
  reply.type('text/plain');
  return register.metrics();
});

/**
 * POST /internal/record-cache-hit
 * Internal endpoint for Nginx to record cache metrics
 */
fastify.post('/internal/record-cache-hit', async (request, reply) => {
  const { endpoint, status } = request.body || {};
  
  if (status === 'HIT' || status === 'EXPIRED') {
    cacheHits.inc({ endpoint: endpoint || 'unknown' });
  } else if (status === 'MISS') {
    cacheMisses.inc({ endpoint: endpoint || 'unknown' });
  }
  
  return { ok: true };
});

/**
 * POST /internal/record-tx-submit
 * Internal endpoint to record TX submission results
 */
fastify.post('/internal/record-tx-submit', async (request, reply) => {
  const { success } = request.body || {};
  
  if (success) {
    txSubmitSuccess.inc();
  } else {
    txSubmitFailed.inc();
  }
  
  return { ok: true };
});

/**
 * POST /internal/record-rate-limit
 * Internal endpoint to record rate limit events
 */
fastify.post('/internal/record-rate-limit', async (request, reply) => {
  rateLimited.inc();
  return { ok: true };
});

// Error handler
fastify.setErrorHandler((error, request, reply) => {
  fastify.log.error(error);
  reply.status(500).send({ error: 'Internal Server Error', message: error.message });
});

// Start server
const start = async () => {
  try {
    await fastify.listen({ port: PORT, host: '127.0.0.1' });
    console.log(`✅ Archivas Relay Service listening on port ${PORT}`);
    console.log(`📊 Endpoints:`);
    console.log(`   - GET  /health   - Simple health check`);
    console.log(`   - GET  /ready    - Readiness probe`);
    console.log(`   - GET  /status   - Detailed status`);
    console.log(`   - GET  /metrics  - Prometheus metrics`);
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
};

start();

// Graceful shutdown
process.on('SIGTERM', async () => {
  console.log('SIGTERM received, shutting down gracefully...');
  await fastify.close();
  process.exit(0);
});

process.on('SIGINT', async () => {
  console.log('SIGINT received, shutting down gracefully...');
  await fastify.close();
  process.exit(0);
});

