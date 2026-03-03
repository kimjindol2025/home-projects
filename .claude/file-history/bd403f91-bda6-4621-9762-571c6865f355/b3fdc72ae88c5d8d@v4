#!/usr/bin/env ts-node

/**
 * FreeLang → Z-Lang Transpiler CLI v2.0
 *
 * 사용법:
 *   fl2z <input.fl> [-o <output.z>] [-v] [--batch] [--report]
 *   fl2z --help
 *   fl2z --version
 */

import * as fs from "fs";
import * as path from "path";
import { Lexer } from "./lexer";
import { Parser } from "./parser";
import ZLangCodeGen from "./transpiler";

interface CliOptions {
  inputFile: string;
  outputFile: string;
  verbose: boolean;
  batch: boolean;
  report: boolean;
  reportFile: string;
}

interface TranspilationResult {
  success: boolean;
  zlangCode: string;
  errors: string[];
  warnings: string[];
  stats: {
    inputLines: number;
    outputLines: number;
    statementsCount: number;
    lexTokens: number;
    parseTime: number;
    genTime: number;
  };
}

interface BatchReport {
  timestamp: string;
  totalFiles: number;
  successCount: number;
  failureCount: number;
  files: Array<{
    filename: string;
    success: boolean;
    inputLines: number;
    outputLines: number;
    error?: string;
  }>;
}

const VERSION = "2.0.0";

function parseArgs(): CliOptions {
  const args = process.argv.slice(2);

  if (args.length === 0 || args[0] === "-h" || args[0] === "--help") {
    printUsage();
    process.exit(args.length === 0 ? 1 : 0);
  }

  if (args[0] === "--version" || args[0] === "-v") {
    console.log(`FreeLang → Z-Lang Transpiler v${VERSION}`);
    process.exit(0);
  }

  const options: CliOptions = {
    inputFile: "",
    outputFile: "",
    verbose: false,
    batch: false,
    report: false,
    reportFile: "transpilation_report.json",
  };

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];

    if (arg === "-o" || arg === "--output") {
      options.outputFile = args[++i];
    } else if (arg === "-v" || arg === "--verbose") {
      options.verbose = true;
    } else if (arg === "--batch") {
      options.batch = true;
    } else if (arg === "--report") {
      options.report = true;
    } else if (arg === "--report-file") {
      options.reportFile = args[++i];
    } else if (!options.inputFile) {
      options.inputFile = arg;
    }
  }

  if (!options.inputFile) {
    console.error("Error: input file required");
    printUsage();
    process.exit(1);
  }

  if (!options.outputFile) {
    options.outputFile = options.inputFile.replace(/\.fl$/, ".z");
  }

  return options;
}

function printUsage(): void {
  console.log(`
FreeLang → Z-Lang Transpiler v${VERSION}

USAGE:
  fl2z <input.fl> [OPTIONS]

OPTIONS:
  -o, --output FILE       Output Z-Lang file (default: input.z)
  -v, --verbose           Verbose output
  --batch                 Batch mode (process directory)
  --report                Generate transpilation report
  --report-file FILE      Report output file (default: transpilation_report.json)
  --help                  Show this help
  --version               Show version

EXAMPLES:
  fl2z hello.fl
  fl2z factorial.fl -o /tmp/factorial.z
  fl2z examples/ --batch --report
  fl2z fizzbuzz.fl -v

FEATURES:
  ✓ FreeLang v4 → Z-Lang transpilation
  ✓ Automatic for-in → while conversion
  ✓ Match, Spawn, Try expression support
  ✓ Complete type system mapping
  ✓ Batch processing for directories
  ✓ JSON report generation
`);
}

function transpile(
  sourceCode: string,
  filename: string,
  options: CliOptions
): TranspilationResult {
  const result: TranspilationResult = {
    success: false,
    zlangCode: "",
    errors: [],
    warnings: [],
    stats: {
      inputLines: sourceCode.split("\n").length,
      outputLines: 0,
      statementsCount: 0,
      lexTokens: 0,
      parseTime: 0,
      genTime: 0,
    },
  };

  try {
    // 1. Lexing
    const lexStartTime = Date.now();
    const lexer = new Lexer(sourceCode);
    const lexResult = lexer.tokenize();
    result.stats.lexTokens = lexResult.tokens.length;
    result.stats.parseTime = Date.now() - lexStartTime;

    if (options.verbose) {
      console.log(`    [Lexer] ${lexResult.tokens.length} tokens`);
    }

    // 2. Parsing
    const parseStartTime = Date.now();
    const parser = new Parser(lexResult.tokens);
    const parseResult = parser.parse();
    result.stats.statementsCount = parseResult.program.stmts.length;
    result.stats.parseTime = Date.now() - parseStartTime;

    if (parseResult.errors.length > 0) {
      result.warnings = parseResult.errors.map(
        (e: any) => `Parse warning at ${e.line}:${e.col}: ${e.message}`
      );
    }

    if (options.verbose) {
      console.log(`    [Parser] ${result.stats.statementsCount} statements`);
    }

    // 3. Code Generation
    const genStartTime = Date.now();
    const codegen = new ZLangCodeGen();
    const zlangCode = codegen.generate(parseResult.program);
    result.stats.genTime = Date.now() - genStartTime;

    result.stats.outputLines = zlangCode.split("\n").length;

    if (options.verbose) {
      console.log(`    [CodeGen] ${result.stats.outputLines} lines`);
      console.log(
        `    [Time] Lex: ${result.stats.parseTime}ms, Gen: ${result.stats.genTime}ms`
      );
    }

    result.zlangCode = zlangCode;
    result.success = true;

  } catch (error) {
    result.success = false;
    result.errors.push(`Transpilation failed: ${(error as Error).message}`);
  }

  return result;
}

function processFile(
  inputPath: string,
  outputPath: string,
  options: CliOptions
): TranspilationResult {
  if (options.verbose) {
    console.log(`[*] Processing: ${inputPath}`);
  }

  if (!fs.existsSync(inputPath)) {
    return {
      success: false,
      zlangCode: "",
      errors: [`File not found: ${inputPath}`],
      warnings: [],
      stats: {
        inputLines: 0,
        outputLines: 0,
        statementsCount: 0,
        lexTokens: 0,
        parseTime: 0,
        genTime: 0,
      },
    };
  }

  const sourceCode = fs.readFileSync(inputPath, "utf-8");
  const result = transpile(sourceCode, path.basename(inputPath), options);

  if (result.success) {
    fs.writeFileSync(outputPath, result.zlangCode, "utf-8");
    if (options.verbose) {
      console.log(`    ✅ Saved to ${outputPath}`);
    }
  }

  return result;
}

function processBatch(dirPath: string, options: CliOptions): BatchReport {
  const report: BatchReport = {
    timestamp: new Date().toISOString(),
    totalFiles: 0,
    successCount: 0,
    failureCount: 0,
    files: [],
  };

  const files = fs
    .readdirSync(dirPath)
    .filter((f) => f.endsWith(".fl"))
    .sort();

  if (files.length === 0) {
    console.warn(`No .fl files found in ${dirPath}`);
    return report;
  }

  console.log(`📦 Batch processing ${files.length} files...`);

  files.forEach((file) => {
    const inputPath = path.join(dirPath, file);
    const outputPath = path.join(
      dirPath,
      file.replace(/\.fl$/, ".z")
    );

    const result = processFile(inputPath, outputPath, {
      ...options,
      verbose: false,
    });

    report.totalFiles++;
    if (result.success) {
      report.successCount++;
      console.log(`  ✅ ${file}`);
    } else {
      report.failureCount++;
      console.log(`  ❌ ${file}: ${result.errors[0]}`);
    }

    report.files.push({
      filename: file,
      success: result.success,
      inputLines: result.stats.inputLines,
      outputLines: result.stats.outputLines,
      error: result.errors[0],
    });
  });

  // Save report if requested
  if (options.report) {
    fs.writeFileSync(options.reportFile, JSON.stringify(report, null, 2));
    console.log(
      `\n📊 Report saved to ${options.reportFile}`
    );
  }

  // Print summary
  console.log(`\n📈 Summary:`);
  console.log(`  Total: ${report.totalFiles}`);
  console.log(`  Success: ${report.successCount} ✅`);
  console.log(`  Failure: ${report.failureCount} ❌`);
  console.log(
    `  Success rate: ${((report.successCount / report.totalFiles) * 100).toFixed(1)}%`
  );

  return report;
}

async function main(): Promise<void> {
  try {
    const options = parseArgs();

    if (options.verbose) {
      console.log(
        `🔄 FreeLang → Z-Lang Transpiler v${VERSION}\n`
      );
    }

    // Check if input is directory (batch mode)
    if (options.batch || fs.statSync(options.inputFile).isDirectory()) {
      processBatch(options.inputFile, options);
    } else {
      // Single file mode
      const result = processFile(
        options.inputFile,
        options.outputFile,
        options
      );

      // Print result
      if (!result.success) {
        console.error("❌ Transpilation failed:");
        result.errors.forEach((err) => console.error(`  ${err}`));
        process.exit(1);
      }

      console.log("✅ Transpilation successful!");
      console.log("");
      console.log("📊 Statistics:");
      console.log(`  Input:      ${result.stats.inputLines} lines`);
      console.log(`  Statements: ${result.stats.statementsCount}`);
      console.log(`  Output:     ${result.stats.outputLines} lines`);
      console.log(`  Tokens:     ${result.stats.lexTokens}`);
      console.log(`  Time:       ${result.stats.parseTime + result.stats.genTime}ms`);
      console.log("");
      console.log(`📁 Output: ${options.outputFile}`);

      // Generate report if requested
      if (options.report) {
        const reportData = {
          timestamp: new Date().toISOString(),
          version: VERSION,
          inputFile: options.inputFile,
          outputFile: options.outputFile,
          success: result.success,
          stats: result.stats,
          warnings: result.warnings,
          errors: result.errors,
        };
        fs.writeFileSync(
          options.reportFile,
          JSON.stringify(reportData, null, 2)
        );
        console.log(`📊 Report: ${options.reportFile}`);
      }
    }
  } catch (error) {
    console.error(`Error: ${(error as Error).message}`);
    process.exit(1);
  }
}

// Main
main().catch((error) => {
  console.error("Fatal error:", error);
  process.exit(1);
});
