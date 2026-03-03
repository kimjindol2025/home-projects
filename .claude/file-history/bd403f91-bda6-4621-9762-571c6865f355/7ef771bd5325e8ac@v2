#!/usr/bin/env node

/**
 * Benchmark Comparison Tool
 * Compares current vs previous measurements
 * Detects performance regressions
 */

const fs = require('fs');
const path = require('path');

const HISTORY_DIR = path.join(__dirname, '..', '.benchmark-history');
const METRICS_FILE = path.join(HISTORY_DIR, 'metrics.json');

if (!fs.existsSync(METRICS_FILE)) {
  console.log('No benchmark history found. Run: npm run benchmark');
  process.exit(0);
}

const allMetrics = JSON.parse(fs.readFileSync(METRICS_FILE, 'utf8'));

if (allMetrics.length < 2) {
  console.log('Need at least 2 measurements for comparison');
  process.exit(0);
}

const latest = allMetrics[allMetrics.length - 1];
const previous = allMetrics[allMetrics.length - 2];

console.log('\n📊 Benchmark Comparison\n');
console.log('═'.repeat(60));

// Test results comparison
console.log('\n✅ Test Results\n');
console.log(`Previous: ${previous.tests.passed}/${previous.tests.total} (${previous.tests.passRate}%)`);
console.log(`Current:  ${latest.tests.passed}/${latest.tests.total} (${latest.tests.passRate}%)`);

if (latest.tests.passed < previous.tests.passed) {
  console.log('❌ REGRESSION: Tests failing');
}

// Performance comparison
console.log('\n⏱️  Execution Time\n');
const prevTime = previous.performance.executionTimeMs;
const currTime = latest.performance.executionTimeMs;
const timeDiff = currTime - prevTime;
const timePercent = ((timeDiff / prevTime) * 100).toFixed(2);

console.log(`Previous: ${prevTime}ms`);
console.log(`Current:  ${currTime}ms`);
console.log(`Change:   ${timeDiff > 0 ? '+' : ''}${timeDiff}ms (${timeDiff > 0 ? '+' : ''}${timePercent}%)`);

if (timeDiff > prevTime * 0.1) { // More than 10% slower
  console.log('⚠️  WARNING: Performance degradation >10%');
}

// Trend analysis
console.log('\n📈 Trend (Last 5 measurements)\n');
const recent = allMetrics.slice(-5);
recent.forEach((m, i) => {
  const ts = new Date(m.timestamp).toLocaleString();
  console.log(`${i + 1}. [${ts}] ${m.performance.executionTimeMs}ms`);
});

// Statistical analysis
console.log('\n📋 Statistics\n');
const times = allMetrics.map(m => m.performance.executionTimeMs);
const avgTime = times.reduce((a, b) => a + b, 0) / times.length;
const minTime = Math.min(...times);
const maxTime = Math.max(...times);

console.log(`Average:  ${avgTime.toFixed(1)}ms`);
console.log(`Min:      ${minTime}ms`);
console.log(`Max:      ${maxTime}ms`);
console.log(`Variance: ${(maxTime - minTime)}ms`);

console.log('\n═'.repeat(60));

// Exit code based on regression
if (latest.tests.passRate < previous.tests.passRate || timeDiff > prevTime * 0.1) {
  console.log('\n⚠️  Performance regression detected');
  process.exit(1);
}

console.log('\n✅ All checks passed');
process.exit(0);
