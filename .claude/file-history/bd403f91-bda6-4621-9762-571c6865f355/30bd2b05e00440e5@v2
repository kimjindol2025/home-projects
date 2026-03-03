/**
 * Basic transpiler tests
 */

import ZLangCodeGen from "../src/transpiler";

describe("ZLangCodeGen", () => {
  const codegen = new ZLangCodeGen();

  test("empty program", () => {
    const ast = { stmts: [] };
    const result = codegen.generate(ast);
    expect(result).toBe("");
  });

  test("simple variable declaration", () => {
    const ast = {
      stmts: [
        {
          kind: "var_decl",
          name: "x",
          mutable: true,
          type: { kind: "i32" },
          init: { kind: "int_lit", value: 42 },
        },
      ],
    };
    const result = codegen.generate(ast);
    expect(result).toContain("let x: i32 = 42;");
  });

  test("integer literal", () => {
    const ast = {
      stmts: [
        {
          kind: "expr_stmt",
          expr: { kind: "int_lit", value: 123 },
        },
      ],
    };
    const result = codegen.generate(ast);
    expect(result).toContain("123");
  });

  test("string literal with escaping", () => {
    const ast = {
      stmts: [
        {
          kind: "expr_stmt",
          expr: { kind: "string_lit", value: 'Hello "World"' },
        },
      ],
    };
    const result = codegen.generate(ast);
    expect(result).toContain('Hello \\"World\\"');
  });

  test("binary operation", () => {
    const ast = {
      stmts: [
        {
          kind: "expr_stmt",
          expr: {
            kind: "binary",
            op: "+",
            left: { kind: "int_lit", value: 10 },
            right: { kind: "int_lit", value: 20 },
          },
        },
      ],
    };
    const result = codegen.generate(ast);
    expect(result).toContain("10 + 20");
  });

  test("return statement", () => {
    const ast = {
      stmts: [
        {
          kind: "return_stmt",
          value: { kind: "int_lit", value: 42 },
        },
      ],
    };
    const result = codegen.generate(ast);
    expect(result).toContain("return 42;");
  });

  test("function declaration", () => {
    const ast = {
      stmts: [
        {
          kind: "fn_decl",
          name: "add",
          params: [
            { name: "x", type: { kind: "i32" } },
            { name: "y", type: { kind: "i32" } },
          ],
          returnType: { kind: "i32" },
          body: [
            {
              kind: "return_stmt",
              value: {
                kind: "binary",
                op: "+",
                left: { kind: "ident", name: "x" },
                right: { kind: "ident", name: "y" },
              },
            },
          ],
        },
      ],
    };
    const result = codegen.generate(ast);
    expect(result).toContain("fn add(x: i32, y: i32) -> i32");
    expect(result).toContain("x + y");
  });

  test("if statement", () => {
    const ast = {
      stmts: [
        {
          kind: "if_stmt",
          condition: { kind: "bool_lit", value: true },
          then: [
            {
              kind: "expr_stmt",
              expr: { kind: "int_lit", value: 1 },
            },
          ],
          else_: [
            {
              kind: "expr_stmt",
              expr: { kind: "int_lit", value: 0 },
            },
          ],
        },
      ],
    };
    const result = codegen.generate(ast);
    expect(result).toContain("if true");
    expect(result).toContain("else");
  });

  test("for-in with range", () => {
    const ast = {
      stmts: [
        {
          kind: "for_stmt",
          variable: "i",
          iterable: {
            kind: "call",
            callee: { kind: "ident", name: "range" },
            args: [
              { kind: "int_lit", value: 1 },
              { kind: "int_lit", value: 5 },
            ],
          },
          body: [
            {
              kind: "expr_stmt",
              expr: { kind: "ident", name: "i" },
            },
          ],
        },
      ],
    };
    const result = codegen.generate(ast);
    expect(result).toContain("let i: i64 = 1");
    expect(result).toContain("while i <= 5");
    expect(result).toContain("i = i + 1");
  });

  test("type conversion", () => {
    const ast = {
      stmts: [
        {
          kind: "var_decl",
          name: "a",
          mutable: true,
          type: { kind: "i64" },
          init: { kind: "int_lit", value: 10 },
        },
        {
          kind: "var_decl",
          name: "b",
          mutable: true,
          type: { kind: "f64" },
          init: { kind: "float_lit", value: 3.14 },
        },
        {
          kind: "var_decl",
          name: "c",
          mutable: true,
          type: { kind: "bool" },
          init: { kind: "bool_lit", value: true },
        },
      ],
    };
    const result = codegen.generate(ast);
    expect(result).toContain("let a: i64 = 10");
    expect(result).toContain("let b: f64 = 3.14");
    expect(result).toContain("let c: bool = true");
  });
});
