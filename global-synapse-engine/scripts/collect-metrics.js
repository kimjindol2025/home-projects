#!/usr/bin/env node

/**
 * Benchmark Metrics Collector
 * Runs tests and extracts performance metrics
 * Stores results in .benchmark-history/
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

const HISTORY_DIR = path.join(__dirname, '..', '.benchmark-history');
const METRICS_FILE = path.join(HISTORY_DIR, 'metrics.json');

// Ensure history directory exists
if (!fs.existsSync(HISTORY_DIR)) {
  fs.mkdirSync(HISTORY_DIR, { recursive: true });
}

console.log('🏃 Running benchmarks...\n');

try {
  // Run tests and capture output
  const output = execSync('npm test 2>&1', { encoding: 'utf8', stdio: 'pipe' });

  // Parse metrics from output
  const metrics = parseMetrics(output);

  // Add timestamp and commit info
  const result = {
    timestamp: new Date().toISOString(),
    commit: process.env.GITHUB_SHA || 'local',
    branch: process.env.GITHUB_REF || 'unknown',
    ...metrics
  };

  // Save individual result
  const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, -5);
  const resultFile = path.join(HISTORY_DIR, `result_${timestamp}_${result.commit.slice(0, 7)}.json`);
  fs.writeFileSync(resultFile, JSON.stringify(result, null, 2));
  console.log(`✅ Saved: ${resultFile}`);

  // Update metrics.json with latest
  const allMetrics = loadAllMetrics();
  allMetrics.push(result);

  // Keep only last 100 results
  if (allMetrics.length > 100) {
    allMetrics.shift();
  }

  fs.writeFileSync(METRICS_FILE, JSON.stringify(allMetrics, null, 2));
  console.log(`✅ Updated metrics: ${METRICS_FILE}`);

  // Generate summary
  generateSummary(allMetrics);

  process.exit(0);
} catch (error) {
  console.error('❌ Benchmark failed:', error.message);
  process.exit(1);
}

function parseMetrics(output) {
  const lines = output.split('\n');

  // Extract test results
  const passedTests = (output.match(/✅/g) || []).length;
  const failedTests = (output.match(/❌/g) || []).length;

  // Extract timing
  const timeMatch = output.match(/⏱️\s+Total Time:\s+(\d+)ms/);
  const executionTimeMs = timeMatch ? parseInt(timeMatch[1]) : null;

  // Extract pass rate
  const passRateMatch = output.match(/Pass Rate:\s+([\d.]+)%/);
  const passRate = passRateMatch ? parseFloat(passRateMatch[1]) : null;

  return {
    tests: {
      passed: passedTests,
      failed: failedTests,
      total: passedTests + failedTests,
      passRate: passRate
    },
    performance: {
      executionTimeMs: executionTimeMs
    }
  };
}

function loadAllMetrics() {
  if (fs.existsSync(METRICS_FILE)) {
    return JSON.parse(fs.readFileSync(METRICS_FILE, 'utf8'));
  }
  return [];
}

function generateSummary(allMetrics) {
  if (allMetrics.length === 0) return;

  const latest = allMetrics[allMetrics.length - 1];
  const previous = allMetrics.length > 1 ? allMetrics[allMetrics.length - 2] : null;

  console.log('\n📊 Benchmark Summary\n');
  console.log(`Current: ${latest.timestamp}`);
  console.log(`Tests: ${latest.tests.passed}/${latest.tests.total} (${latest.tests.passRate}%)`);
  console.log(`Execution: ${latest.performance.executionTimeMs}ms`);

  if (previous) {
    const timeDiff = latest.performance.executionTimeMs - previous.performance.executionTimeMs;
    const timePercentage = ((timeDiff / previous.performance.executionTimeMs) * 100).toFixed(1);
    const trend = timeDiff > 0 ? '⬆️ slower' : '⬇️ faster';
    console.log(`Trend: ${trend} (${Math.abs(timePercentage)}%)`);
  }

  console.log(`\nHistory: ${allMetrics.length} measurements`);
}
