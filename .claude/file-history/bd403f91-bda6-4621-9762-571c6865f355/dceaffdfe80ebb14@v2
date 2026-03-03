/**
 * FreeLang → Z-Lang Transpiler
 *
 * AST를 받아서 Z-Lang 소스 코드 문자열로 변환
 */

// FreeLang AST 타입 임포트 (심볼릭 링크로 연결)
type TypeAnnotation = any;
type Expr = any;
type Stmt = any;
type Program = any;
type Param = any;

export class ZLangCodeGen {
  private indentLevel: number = 0;
  private indent(): string {
    return "  ".repeat(this.indentLevel);
  }

  /**
   * 프로그램 전체 변환
   */
  generate(ast: Program): string {
    if (!ast || !ast.stmts) {
      return "";
    }

    const statements = ast.stmts
      .map((stmt: Stmt) => this.genStmt(stmt))
      .filter((s: string) => s.length > 0)
      .join("\n");

    return statements;
  }

  /**
   * 문(Statement) 변환
   */
  private genStmt(stmt: Stmt): string {
    if (!stmt) return "";

    switch (stmt.kind) {
      case "var_decl":
        return this.genVarDecl(stmt);
      case "fn_decl":
        return this.genFunction(stmt);
      case "if_stmt":
        return this.genIf(stmt);
      case "for_stmt":
        return this.genForIn(stmt);
      case "return_stmt":
        return this.genReturn(stmt);
      case "expr_stmt":
        return this.genExprStmt(stmt);
      case "spawn_stmt":
        // Phase 2: spawn 미지원, 주석 처리
        return `// TODO: spawn not yet supported`;
      case "match_stmt":
        // Phase 2: match 미지원
        return `// TODO: match not yet supported`;
      default:
        return `// TODO: unknown statement kind: ${(stmt as any).kind}`;
    }
  }

  /**
   * 변수 선언: var/let/const → let + semicolon
   */
  private genVarDecl(stmt: Stmt): string {
    const name = stmt.name;
    const type = this.convertType(stmt.type);
    const init = this.genExpr(stmt.init);
    return `${this.indent()}let ${name}: ${type} = ${init};`;
  }

  /**
   * 함수 정의
   */
  private genFunction(stmt: Stmt): string {
    const name = stmt.name;
    const params = stmt.params
      .map((p: Param) => `${p.name}: ${this.convertType(p.type)}`)
      .join(", ");

    const returnType = stmt.returnType
      ? ` -> ${this.convertType(stmt.returnType)}`
      : "";

    const body = this.genBlock(stmt.body, stmt.returnType);
    return `${this.indent()}fn ${name}(${params})${returnType} {\n${body}\n${this.indent()}}`;
  }

  /**
   * If 문
   */
  private genIf(stmt: Stmt): string {
    const condition = this.genExpr(stmt.condition);
    this.indentLevel++;
    const thenBody = stmt.then
      .map((s: Stmt) => this.genStmt(s))
      .filter((s: string) => s.length > 0)
      .join("\n");
    this.indentLevel--;

    let result = `${this.indent()}if ${condition} {\n${thenBody}\n${this.indent()}}`;

    if (stmt.else_ && stmt.else_.length > 0) {
      this.indentLevel++;
      const elseBody = stmt.else_
        .map((s: Stmt) => this.genStmt(s))
        .filter((s: string) => s.length > 0)
        .join("\n");
      this.indentLevel--;
      result += ` else {\n${elseBody}\n${this.indent()}}`;
    }

    return result;
  }

  /**
   * For-in 루프 → While 변환
   */
  private genForIn(stmt: Stmt): string {
    const variable = stmt.variable;
    const iterable = stmt.iterable;

    // range(start, end) 패턴 인식
    if (iterable.kind === "call" &&
        iterable.callee &&
        iterable.callee.name === "range" &&
        iterable.args &&
        iterable.args.length >= 2) {

      const start = this.genExpr(iterable.args[0]);
      const end = this.genExpr(iterable.args[1]);

      this.indentLevel++;
      const body = stmt.body
        .map((s: Stmt) => this.genStmt(s))
        .filter((s: string) => s.length > 0)
        .join("\n");
      this.indentLevel--;

      return `${this.indent()}let ${variable}: i64 = ${start};
${this.indent()}while ${variable} <= ${end} {
${body}
${this.indent()}  ${variable} = ${variable} + 1;
${this.indent()}}`;
    }

    // 배열 루프: Phase 2에서 지원
    return `${this.indent()}// TODO: for-in array not yet supported`;
  }

  /**
   * Return 문
   */
  private genReturn(stmt: Stmt): string {
    if (stmt.value) {
      const value = this.genExpr(stmt.value);
      return `${this.indent()}return ${value};`;
    } else {
      return `${this.indent()}return;`;
    }
  }

  /**
   * 표현식 문
   */
  private genExprStmt(stmt: Stmt): string {
    const expr = this.genExpr(stmt.expr);
    return `${this.indent()}${expr};`;
  }

  /**
   * 블록(여러 문)
   */
  private genBlock(stmts: Stmt[], returnType?: TypeAnnotation): string {
    if (!stmts || stmts.length === 0) {
      return "";
    }

    this.indentLevel++;
    const lines: string[] = [];

    // 마지막 문 이전까지 처리
    for (let i = 0; i < stmts.length - 1; i++) {
      const line = this.genStmt(stmts[i]);
      if (line.length > 0) {
        lines.push(line);
      }
    }

    // 마지막 문 처리: 암묵적 반환
    const lastStmt = stmts[stmts.length - 1];
    if (lastStmt) {
      if (returnType && returnType.kind !== "void" && lastStmt.kind === "expr_stmt") {
        // 마지막이 표현식 문이고 반환 타입이 있으면 return으로 변환
        const expr = this.genExpr(lastStmt.expr);
        lines.push(`${this.indent()}return ${expr};`);
      } else {
        const line = this.genStmt(lastStmt);
        if (line.length > 0) {
          lines.push(line);
        }
      }
    }

    this.indentLevel--;
    return lines.join("\n");
  }

  /**
   * 표현식(Expression) 변환
   */
  private genExpr(expr: Expr): string {
    if (!expr) return "null";

    switch (expr.kind) {
      case "int_lit":
        return expr.value.toString();
      case "float_lit":
        return expr.value.toString();
      case "string_lit":
        return `"${this.escapeString(expr.value)}"`;
      case "bool_lit":
        return expr.value ? "true" : "false";
      case "ident":
        return expr.name;
      case "binary":
        return this.genBinaryOp(expr);
      case "unary":
        return this.genUnaryOp(expr);
      case "call":
        return this.genCall(expr);
      case "assign":
        return this.genAssignment(expr);
      case "if_expr":
        return this.genIfExpr(expr);
      case "array_lit":
        return this.genArrayLit(expr);
      default:
        return `/* TODO: ${expr.kind} */`;
    }
  }

  /**
   * 이항 연산자
   */
  private genBinaryOp(expr: Expr): string {
    const left = this.genExpr(expr.left);
    const right = this.genExpr(expr.right);
    const op = expr.op;
    return `${left} ${op} ${right}`;
  }

  /**
   * 단항 연산자
   */
  private genUnaryOp(expr: Expr): string {
    const operand = this.genExpr(expr.operand);
    const op = expr.op;
    return `${op}${operand}`;
  }

  /**
   * 함수 호출
   */
  private genCall(expr: Expr): string {
    const callee = this.genExpr(expr.callee);
    const args = expr.args
      .map((arg: Expr) => this.genExpr(arg))
      .join(", ");
    return `${callee}(${args})`;
  }

  /**
   * 할당
   */
  private genAssignment(expr: Expr): string {
    const target = this.genExpr(expr.target);
    const value = this.genExpr(expr.value);
    return `${target} = ${value}`;
  }

  /**
   * If 표현식
   */
  private genIfExpr(expr: Expr): string {
    const condition = this.genExpr(expr.condition);
    const thenExpr = expr.then.length > 0
      ? this.genExpr(expr.then[expr.then.length - 1])
      : "null";
    const elseExpr = expr.else_ && expr.else_.length > 0
      ? this.genExpr(expr.else_[expr.else_.length - 1])
      : "null";

    return `if ${condition} { ${thenExpr} } else { ${elseExpr} }`;
  }

  /**
   * 배열 리터럴
   */
  private genArrayLit(expr: Expr): string {
    const elements = expr.elements
      .map((e: Expr) => this.genExpr(e))
      .join(", ");
    return `[${elements}]`;
  }

  /**
   * 타입 변환: FreeLang → Z-Lang
   */
  private convertType(type: TypeAnnotation | null): string {
    if (!type) return "i64";

    switch (type.kind) {
      case "i32":
      case "i64":
      case "f64":
      case "bool":
      case "string":
        return type.kind;
      case "void":
        return "void";
      case "array":
        const elementType = this.convertType(type.element);
        return `[${elementType}]`;
      case "channel":
        const channelType = this.convertType(type.element);
        return `channel<${channelType}>`;
      case "option":
        const optionType = this.convertType(type.element);
        return `option<${optionType}>`;
      case "result":
        const okType = this.convertType(type.ok);
        const errType = this.convertType(type.err);
        return `result<${okType}, ${errType}>`;
      default:
        return "i64";
    }
  }

  /**
   * 문자열 이스케이핑
   */
  private escapeString(str: string): string {
    return str
      .replace(/\\/g, "\\\\")
      .replace(/"/g, '\\"')
      .replace(/\n/g, "\\n")
      .replace(/\t/g, "\\t")
      .replace(/\r/g, "\\r");
  }
}

export default ZLangCodeGen;
