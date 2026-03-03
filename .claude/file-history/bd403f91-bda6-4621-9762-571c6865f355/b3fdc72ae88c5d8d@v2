#!/usr/bin/env ts-node

/**
 * FreeLang → Z-Lang Transpiler CLI
 *
 * 사용법:
 *   npx ts-node src/index.ts <input.fl> [-o <output.z>]
 */

import * as fs from "fs";
import * as path from "path";
import ZLangCodeGen from "./transpiler";

interface CliOptions {
  inputFile: string;
  outputFile: string;
  verbose: boolean;
}

function parseArgs(): CliOptions {
  const args = process.argv.slice(2);

  if (args.length === 0) {
    printUsage();
    process.exit(1);
  }

  const options: CliOptions = {
    inputFile: "",
    outputFile: "",
    verbose: false,
  };

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];

    if (arg === "-o" || arg === "--output") {
      options.outputFile = args[++i];
    } else if (arg === "-v" || arg === "--verbose") {
      options.verbose = true;
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
    // Generate output filename from input
    options.outputFile = options.inputFile.replace(/\.fl$/, ".z");
  }

  return options;
}

function printUsage(): void {
  console.log(`
FreeLang → Z-Lang Transpiler v1.0

Usage:
  npx ts-node src/index.ts <input.fl> [-o <output.z>] [-v]

Options:
  -o, --output FILE    Output Z-Lang file (default: input.z)
  -v, --verbose        Verbose output
  -h, --help          Show this help

Examples:
  npx ts-node src/index.ts hello.fl
  npx ts-node src/index.ts factorial.fl -o /tmp/factorial.z
`);
}

async function transpile(options: CliOptions): Promise<void> {
  try {
    // 1. 입력 파일 읽기
    if (!fs.existsSync(options.inputFile)) {
      throw new Error(`Input file not found: ${options.inputFile}`);
    }

    const sourceCode = fs.readFileSync(options.inputFile, "utf-8");
    if (options.verbose) {
      console.log(`[*] Read ${options.inputFile} (${sourceCode.length} bytes)`);
    }

    // 2. FreeLang 파서로 AST 생성
    // NOTE: 실제로는 freelang-v4의 파서를 임포트해야 하지만,
    // 이 예제에서는 간단한 파서를 사용하거나 더미 구현을 사용
    const ast = parseFreeLang(sourceCode);
    if (options.verbose) {
      console.log(`[*] Parsed ${ast.stmts.length} statements`);
    }

    // 3. Z-Lang 코드 생성
    const codegen = new ZLangCodeGen();
    const zlangCode = codegen.generate(ast);
    if (options.verbose) {
      console.log(`[*] Generated ${zlangCode.length} bytes of Z-Lang code`);
    }

    // 4. 출력 파일 저장
    fs.writeFileSync(options.outputFile, zlangCode, "utf-8");
    console.log(`✓ Transpiled to ${options.outputFile}`);

  } catch (error) {
    console.error(`Error: ${(error as Error).message}`);
    process.exit(1);
  }
}

/**
 * 간단한 FreeLang 파서 (매우 기본적인 구현)
 * 실제로는 freelang-v4의 완전한 파서를 사용해야 함
 */
function parseFreeLang(source: string): any {
  // 이것은 플레이스홀더입니다
  // 실제 구현에서는 freelang-v4의 Lexer + Parser를 사용합니다

  // 예제: 간단한 토큰화
  const lines = source.split("\n");
  const stmts: any[] = [];

  for (const line of lines) {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith("//")) {
      continue;
    }

    // 매우 기본적인 패턴 매칭 (실제로는 proper parser 필요)
    if (trimmed.startsWith("fn ")) {
      // 함수 선언
      stmts.push(parseFunctionDecl(trimmed));
    } else if (trimmed.startsWith("var ") || trimmed.startsWith("let ") || trimmed.startsWith("const ")) {
      // 변수 선언
      stmts.push(parseVarDecl(trimmed));
    } else if (trimmed.startsWith("return ")) {
      // Return 문
      stmts.push(parseReturn(trimmed));
    }
  }

  return { stmts };
}

function parseFunctionDecl(line: string): any {
  // 매우 간단한 구현
  return {
    kind: "fn_decl",
    name: "main",
    params: [],
    returnType: { kind: "i64" },
    body: [
      {
        kind: "return_stmt",
        value: { kind: "int_lit", value: 0 },
      },
    ],
  };
}

function parseVarDecl(line: string): any {
  return {
    kind: "var_decl",
    name: "x",
    mutable: true,
    type: { kind: "i64" },
    init: { kind: "int_lit", value: 0 },
  };
}

function parseReturn(line: string): any {
  return {
    kind: "return_stmt",
    value: { kind: "int_lit", value: 42 },
  };
}

// Main
const options = parseArgs();
transpile(options).catch((error) => {
  console.error("Transpilation failed:", error);
  process.exit(1);
});
