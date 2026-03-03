/**
 * FreeLang Ownership Checker
 * 소유권 규칙 검증 엔진
 *
 * 역할:
 * 1. 변수 소유권 추적
 * 2. 이동(move) 감지
 * 3. 스코프 기반 해제 확인
 * 4. 함수 소유권 이전 검증
 */

import { ASTNode } from './types';

// ============================================================
// 데이터 구조
// ============================================================

/**
 * 변수의 소유권 상태
 */
interface OwnershipState {
  variable: string;
  definedAt: number;              // 정의 라인
  ownedBy: string;                // 현재 소유자
  moved: boolean;                 // 소유권 이동됨?
  movedAt?: number;               // 이동 라인
  movedTo?: string;               // 이동 대상
  copies: Array<{
    line: number;
    target: string;
  }>;
  scopes: Array<{
    start: number;
    end: number;
    valid: boolean;
  }>;
}

/**
 * 소유권 에러
 */
interface OwnershipError {
  type: 'use_after_move' | 'use_after_free' | 'invalid_scope';
  variable: string;
  line: number;
  message: string;
  suggestion?: string;
}

/**
 * 소유권 맵 (변수 → 상태)
 */
type OwnershipMap = Map<string, OwnershipState>;

// ============================================================
// Ownership Checker
// ============================================================

export class OwnershipChecker {
  private ownershipMap: OwnershipMap = new Map();
  private errors: OwnershipError[] = [];
  private scopeStack: Array<{ start: number; variables: Set<string> }> = [];
  private currentScope: number = 0;

  /**
   * 메인 검증 메서드
   */
  check(ast: ASTNode): OwnershipError[] {
    this.ownershipMap.clear();
    this.errors = [];
    this.scopeStack = [];
    this.currentScope = 0;

    // AST 분석
    const statements: ASTNode[] = ast.type === 'block' ? (ast as any).statements : [ast];
    this.analyzeStatements(statements);

    return this.errors;
  }

  /**
   * Statements 분석
   */
  private analyzeStatements(statements: ASTNode[]): void {
    for (const stmt of statements) {
      this.analyzeStatement(stmt);
    }
  }

  /**
   * 단일 Statement 분석
   */
  private analyzeStatement(node: ASTNode): void {
    switch (node.type) {
      case 'assignment': {
        const n = node as any;
        // 변수 정의
        if (!this.ownershipMap.has(n.variable)) {
          this.ownershipMap.set(n.variable, {
            variable: n.variable,
            definedAt: 0,
            ownedBy: n.variable,
            moved: false,
            copies: [],
            scopes: [{
              start: this.currentScope,
              end: this.currentScope + 1000, // 임시값
              valid: true
            }]
          });
        }

        // 우측값 분석
        this.analyzeExpression(n.value);
        break;
      }

      case 'block': {
        const n = node as any;
        this.currentScope++;
        this.analyzeStatements(n.statements);
        this.currentScope--;
        break;
      }

      case 'if': {
        const n = node as any;
        this.analyzeExpression(n.condition);
        this.analyzeStatement(n.thenBranch);
        if (n.elseBranch) {
          this.analyzeStatement(n.elseBranch);
        }
        break;
      }

      case 'while': {
        const n = node as any;
        this.analyzeExpression(n.condition);
        this.analyzeStatement(n.body);
        break;
      }

      case 'functionDef': {
        const n = node as any;
        this.currentScope++;
        this.analyzeStatement(n.body);
        this.currentScope--;
        break;
      }

      default:
        // 표현식
        this.analyzeExpression(node);
        break;
    }
  }

  /**
   * Expression 분석
   */
  private analyzeExpression(node: ASTNode): void {
    switch (node.type) {
      case 'identifier': {
        const n = node as any;
        const state = this.ownershipMap.get(n.name);
        if (state && state.moved) {
          this.errors.push({
            type: 'use_after_move',
            variable: n.name,
            line: 0,
            message: `Variable '${n.name}' was moved and cannot be used`,
            suggestion: `Use the new owner instead`
          });
        }
        break;
      }

      case 'binaryOp': {
        const n = node as any;
        this.analyzeExpression(n.left);
        this.analyzeExpression(n.right);
        break;
      }

      case 'unaryOp': {
        const n = node as any;
        this.analyzeExpression(n.operand);
        break;
      }

      case 'functionCall': {
        const n = node as any;
        for (const arg of n.args) {
          this.analyzeExpression(arg);
        }
        break;
      }

      case 'arrayLiteral': {
        const n = node as any;
        for (const elem of n.elements) {
          this.analyzeExpression(elem);
        }
        break;
      }

      case 'arrayAccess': {
        const n = node as any;
        this.analyzeExpression(n.array);
        this.analyzeExpression(n.index);
        break;
      }

      default:
        break;
    }
  }

  /**
   * Move 감지
   */
  private detectMove(variable: string, target: string): void {
    const state = this.ownershipMap.get(variable);
    if (state) {
      state.moved = true;
      state.movedTo = target;
    }
  }

  /**
   * Copy 감지
   */
  private detectCopy(variable: string, target: string, line: number): void {
    const state = this.ownershipMap.get(variable);
    if (state) {
      state.copies.push({ line, target });
    }
  }

  /**
   * 소유권 맵 반환
   */
  getOwnershipMap(): OwnershipMap {
    return this.ownershipMap;
  }

  /**
   * 에러 반환
   */
  getErrors(): OwnershipError[] {
    return this.errors;
  }
}

// ============================================================
// Type Definitions
// ============================================================

export type { OwnershipState, OwnershipError, OwnershipMap };
