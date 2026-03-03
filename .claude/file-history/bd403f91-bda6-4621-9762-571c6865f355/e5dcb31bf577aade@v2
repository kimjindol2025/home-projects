// 🔗 Test Mouse Audit Log - Hash Chain Verification
// Layer: Shared across all layers (Immutability guarantee)

import { HashChainLink, sha256 } from './types';

/**
 * HashChain: Test Mouse의 감사 로그
 * 모든 중요 작업을 SHA256 체인으로 기록하여 개조 불가능하도록 보증
 *
 * Unforgiving Rule: Chain Integrity = 100%
 *   - 모든 링크는 이전 해시에 의존 (부인 불가능)
 *   - 1비트 변조도 즉시 감지
 *   - 1000+ 링크 지원
 */
export class HashChain {
  private links: HashChainLink[] = [];
  private lastHash: string = '0'.repeat(64);

  /**
   * 새로운 링크를 체인에 추가
   * @param operationId 작업 ID (RDMA/SemanticSync 등)
   * @param operationType 작업 타입 (read/write/verify 등)
   * @param signature HMAC 서명
   * @returns 새로 추가된 링크
   */
  public addLink(
    operationId: bigint,
    operationType: string,
    signature: string
  ): HashChainLink {
    const content = JSON.stringify({
      type: operationType,
      operationId: operationId.toString(),
      signature,
    });

    const hashInput = this.lastHash + content;
    const newHash = sha256(hashInput);

    const link: HashChainLink = {
      index: this.links.length,
      timestamp: BigInt(Date.now()) * 1_000_000n,
      previousHash: this.lastHash,
      content: {
        type: operationType,
        operationId,
        signature,
      },
      hash: newHash,
    };

    this.links.push(link);
    this.lastHash = newHash;

    return link;
  }

  /**
   * 체인 무결성 검증
   * @returns { isValid: boolean, totalLinks: number, verifiedLinks: number, brokenAt: number }
   */
  public verify(): {
    isValid: boolean;
    totalLinks: number;
    verifiedLinks: number;
    brokenAt?: number;
  } {
    if (this.links.length === 0) {
      return { isValid: true, totalLinks: 0, verifiedLinks: 0 };
    }

    let verifiedLinks = 0;
    let previousHash = '0'.repeat(64);

    for (let i = 0; i < this.links.length; i++) {
      const link = this.links[i];

      // 이전 해시 검증
      if (link.previousHash !== previousHash) {
        return {
          isValid: false,
          totalLinks: this.links.length,
          verifiedLinks,
          brokenAt: i,
        };
      }

      // 현재 해시 재계산 및 검증
      const content = JSON.stringify(link.content, (key, value) =>
        typeof value === 'bigint' ? value.toString() : value
      );
      const expectedHash = sha256(previousHash + content);

      if (expectedHash !== link.hash) {
        return {
          isValid: false,
          totalLinks: this.links.length,
          verifiedLinks,
          brokenAt: i,
        };
      }

      verifiedLinks++;
      previousHash = link.hash;
    }

    return {
      isValid: true,
      totalLinks: this.links.length,
      verifiedLinks,
    };
  }

  /**
   * 특정 링크 찾기 (이진 검색)
   */
  public findLink(operationId: bigint): HashChainLink | undefined {
    return this.links.find((link) => link.content.operationId === operationId);
  }

  /**
   * 지정된 범위의 링크 조회
   */
  public getLinksInRange(startIndex: number, endIndex: number): HashChainLink[] {
    return this.links.slice(startIndex, Math.min(endIndex + 1, this.links.length));
  }

  /**
   * 체인 직렬화 (저장용)
   */
  public serialize(): string {
    return JSON.stringify(
      {
        links: this.links,
        lastHash: this.lastHash,
        timestamp: Date.now(),
      },
      (key, value) => {
        if (typeof value === 'bigint') {
          return value.toString();
        }
        return value;
      }
    );
  }

  /**
   * 체인 역직렬화 (복원용)
   */
  public deserialize(data: string): void {
    const parsed = JSON.parse(data);
    this.links = parsed.links.map((link: any) => ({
      ...link,
      timestamp: BigInt(link.timestamp),
      content: {
        ...link.content,
        operationId: BigInt(link.content.operationId),
      },
    }));
    this.lastHash = parsed.lastHash;
  }

  /**
   * 체인 통계
   */
  public getStatistics(): {
    totalLinks: number;
    chainDepth: number;
    oldestTimestamp: bigint;
    newestTimestamp: bigint;
    averageHashDistance: number;
  } {
    if (this.links.length === 0) {
      return {
        totalLinks: 0,
        chainDepth: 0,
        oldestTimestamp: 0n,
        newestTimestamp: 0n,
        averageHashDistance: 0,
      };
    }

    const oldestTimestamp = this.links[0].timestamp;
    const newestTimestamp = this.links[this.links.length - 1].timestamp;
    const chainDepth = this.links.length; // 각 링크는 이전 해시에 의존하므로 깊이 = 길이

    return {
      totalLinks: this.links.length,
      chainDepth,
      oldestTimestamp,
      newestTimestamp,
      averageHashDistance: chainDepth / 10, // 평균 10링크 = 1칸
    };
  }

  /**
   * 현재 체인 상태 (공개)
   */
  public getLastHash(): string {
    return this.lastHash;
  }

  public getLinks(): ReadonlyArray<HashChainLink> {
    return Object.freeze([...this.links]);
  }
}

/**
 * HashChainVerifier: 외부 검증용 유틸리티
 * 두 체인을 비교하여 동등성 검증
 */
export class HashChainVerifier {
  /**
   * 두 체인이 동일한지 검증
   */
  public static verifyEquivalence(chain1: HashChain, chain2: HashChain): boolean {
    if (chain1.getLastHash() !== chain2.getLastHash()) {
      return false;
    }

    const links1 = chain1.getLinks();
    const links2 = chain2.getLinks();

    if (links1.length !== links2.length) {
      return false;
    }

    for (let i = 0; i < links1.length; i++) {
      if (links1[i].hash !== links2[i].hash) {
        return false;
      }
    }

    return true;
  }

  /**
   * 체인에서 변조된 부분 찾기 (이진 검색)
   */
  public static findTampering(chain: HashChain): number | null {
    const verification = chain.verify();

    if (verification.isValid) {
      return null;
    }

    return verification.brokenAt ?? 0;
  }

  /**
   * 체인의 일부를 새로운 체인으로 추출
   */
  public static extractSubChain(chain: HashChain, startIndex: number, length: number): HashChain {
    const newChain = new HashChain();
    const links = chain.getLinks().slice(startIndex, startIndex + length);

    for (const link of links) {
      newChain.addLink(link.content.operationId, link.content.type, link.content.signature);
    }

    return newChain;
  }
}
