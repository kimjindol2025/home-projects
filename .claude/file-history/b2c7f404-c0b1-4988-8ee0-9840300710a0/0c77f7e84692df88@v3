# v11: 블록체인 & DPoS 합의 엔진 구현 계획

## Context
v10 Federated Learning 완료 후, 분산 합의 알고리즘 심화 연구.
v9.3 Raft(CFT, 로그 복제)와 달리 DPoS는 경제적 인센티브 기반 BFT 합의.
hashlib(표준 라이브러리)으로 SHA-256 해시 체인 + Merkle Tree + 스택 VM 구현.

## 파일 위치
- **메인**: `university_v11_BLOCKCHAIN_DPOS.py` (~1,350줄)
- **테스트**: `test_v11_blockchain_dpos.py` (~420줄)

## 참조 파일
- `university_v10_DISTRIBUTED_ML.py` — PART/SECTION 구조, dataclass 패턴
- `test_v10_distributed_ml.py` — Group A~E, test_NN_ 번호 체계
- `university_v9_3_PLANETARY_CONSENSUS.py` — 합의 알고리즘 비교 기준

## Import (표준 라이브러리만)
```python
import hashlib
import hmac
import json
import math
import random
import time
from dataclasses import dataclass, field
from enum import Enum
from typing import Dict, List, Optional, Tuple, Any
from collections import defaultdict
```

---

## PART 0: Enum 4개 + Dataclass 5개

### Enum (소문자 문자열 값)
| 이름 | 값 |
|------|---|
| `TransactionType` | TRANSFER / STAKE / DELEGATE / CONTRACT_CALL |
| `BlockStatus` | PENDING / CONFIRMED / FINALIZED |
| `NodeRole` | VALIDATOR / DELEGATOR / OBSERVER |
| `ConsensusPhase` | ELECTION / PRODUCING / VERIFYING / FINALIZING |

### Dataclass
| 이름 | 핵심 필드 |
|------|---------|
| `Transaction` | `tx_hash`, `sender`, `receiver`, `amount`, `tx_type`, `timestamp`, `data: Dict`, `fee: float = 0.001` |
| `BlockHeader` | `block_id`, `prev_hash`, `merkle_root`, `timestamp`, `validator`, `slot_number` |
| `Block` | `header: BlockHeader`, `transactions: List[Transaction]`, `block_hash`, `size_bytes`, `status: BlockStatus` |
| `ValidatorInfo` | `address`, `total_stake`, `delegator_count`, `blocks_produced`, `is_active`, `reward_balance` |
| `ConsensusResult` | `slot`, `producer`, `block_hash`, `confirmed`, `vote_count`, `round_time_ms` |

---

## PART 1~8: 클래스 목록

### PART 1: CryptoUtils
역할: SHA-256 해시, Merkle Tree, 주소 생성

| 메서드 | 설명 |
|--------|------|
| `sha256(data)` | `hashlib.sha256(data.encode()).hexdigest()` → 64자 hex |
| `double_sha256(data)` | sha256(sha256(data)) |
| `compute_tx_hash(tx)` | 트랜잭션 필드 JSON → sha256 |
| `compute_block_hash(header)` | 헤더 필드 JSON → sha256 |
| `build_merkle_tree(tx_hashes)` | 이진 Merkle Tree, 홀수 시 마지막 복사 |
| `generate_address(seed)` | sha256(seed) 앞 40자 + "0x" prefix |
| `compute_hmac(key, msg)` | hmac.new(sha256) → hex digest |

핵심 알고리즘:
```
Merkle Root:
  while len(leaves) > 1:
      if len(leaves) % 2 == 1: leaves.append(leaves[-1])
      leaves = [sha256(a+b) for a,b in zip(leaves[::2], leaves[1::2])]
  return leaves[0]
```

---

### PART 2: TransactionPool
역할: Mempool 관리, 수수료 우선순위 정렬

| 메서드 | 설명 |
|--------|------|
| `add_transaction(tx)` | 중복 거부, 음수 amount 거부, max_size 초과 시 낮은 fee 제거 |
| `select_for_block(max_tx=50)` | fee 내림차순 상위 반환 |
| `get_pending_count()` | 미확인 트랜잭션 수 |
| `get_pool_stats()` | 전체 수, 평균/최고/최저 fee, 타입별 분포 |
| `flush_old_transactions(max_age_sec=300)` | 오래된 트랜잭션 제거 |

---

### PART 3: BlockFactory
역할: 블록 생성, 구조 검증

| 메서드 | 설명 |
|--------|------|
| `create_genesis_block()` | block_id=0, prev_hash="0"*64, 빈 tx |
| `create_block(txs, prev_block, validator, slot)` | Merkle Root → Header → Hash → Block |
| `validate_block_structure(block, prev)` | prev_hash, block_id, merkle_root, timestamp 검증 |
| `serialize_block(block)` | JSON 직렬화 (P2P 전파용) |

---

### PART 4: Blockchain
역할: 체인 상태 관리, 잔액 추적, 포크 해결

| 메서드 | 설명 |
|--------|------|
| `initialize()` | 제네시스 블록 생성, 초기 잔액 배분 |
| `append_block(block)` | 검증 → 잔액 업데이트 → CONFIRMED 설정 |
| `get_balance(address)` | 주소 잔액 (없으면 0.0) |
| `validate_chain()` | 전체 prev_hash + merkle_root 검증 → (bool, msg) |
| `resolve_fork(fork_chain)` | 최장 체인 규칙, _rebuild_balances() |
| `get_chain_stats()` | height, total_transactions, average_block_size |

포크 해결:
```
if len(fork) <= len(main): reject
if genesis 불일치: reject
fork 전체 구조 검증 → 교체 → _rebuild_balances()
```

---

### PART 5: DPoSEngine ← 핵심
역할: 스테이킹, 위임 투표, 검증자 선출, 보상, 슬래싱

| 메서드 | 설명 |
|--------|------|
| `register_validator(address, stake)` | ValidatorInfo 등록 |
| `stake(address, amount)` | total_stake += amount |
| `delegate(delegator, validator, amount)` | 위임 지분 누적 |
| `run_election()` | total_stake 내림차순 → 상위 21명 → 슬롯 스케줄 |
| `_generate_slot_schedule()` | Fisher-Yates 셔플로 슬롯 순서 결정 |
| `get_slot_producer(slot)` | slot % len(schedule) 위치 검증자 |
| `produce_block(slot, blockchain, tx_pool)` | 슬롯 생산자가 블록 생성 |
| `collect_votes(block, validators)` | 85% 참여율 시뮬레이션, quorum = ceil(n×2/3)+1 |
| `compute_reward(validator, total_stake)` | BASE_REWARD × (v_stake / total) |
| `distribute_rewards(result)` | 생산자 reward_balance 증가 |
| `slash_validator(address, reason)` | stake × (1 - SLASH_RATE), is_active=False |
| `advance_epoch()` | epoch++, 재선출 |

상수: `NUM_VALIDATORS=21`, `SLASH_RATE=0.05`, `BASE_REWARD=100.0`

핵심 수학:
```
선출: sorted(validators, key=total_stake, reverse=True)[:21]
보상: reward = 100.0 × (validator_stake / total_network_stake)
슬래싱: new_stake = stake × 0.95
합의: quorum = math.ceil(21 × 2/3) + 1 = 15표
```

---

### PART 6: SmartContract (StackVM 포함)
역할: 스택 기반 VM + 컨트랙트 배포/실행

StackVM 옵코드(14개):
```
PUSH(3gas), POP(3gas), ADD/SUB(5gas), MUL/DIV(10gas),
GT/LT/EQ(5gas), STORE/LOAD(20gas), EMIT(15gas), HALT(0gas)
```

| 메서드 | 설명 |
|--------|------|
| `StackVM.execute(bytecode)` | 스택 기반 실행 → {result, gas_used, logs, storage} |
| `SmartContract.deploy_contract(owner, bytecode)` | address = generate_address(owner+time) |
| `SmartContract.call_contract(addr, caller, args)` | StackVM 실행 → state 업데이트 |
| `SmartContract.get_contract_state(addr)` | 저장소 상태 반환 |

GAS 소진: `gas_used >= gas_limit` 시 `{"success": False, "error": "gas_exhausted"}`

---

### PART 7: P2PNetwork
역할: 피어 관리, Gossip 전파, 체인 동기화

| 메서드 | 설명 |
|--------|------|
| `join_network(peer_id, role, blockchain)` | 중복 거부, peers 등록 |
| `broadcast_block(block, sender_id)` | 95% 성공률 Gossip 시뮬레이션 |
| `broadcast_transaction(tx)` | 모든 피어 tx_pool에 추가 |
| `sync_chain(requester_id)` | 최장 체인 피어 찾기 → resolve_fork |
| `simulate_partition(isolated, duration)` | 분리 → 재연결 → 통계 반환 |

---

### PART 8: BlockExplorer
역할: 읽기 전용 쿼리 인터페이스

| 메서드 | 설명 |
|--------|------|
| `get_transaction(tx_hash)` | 모든 블록 순회 검색 |
| `get_address_history(address)` | 송수신 전체 트랜잭션 |
| `get_validator_ranking()` | total_stake 내림차순 |
| `get_rich_list(top_n=10)` | 잔액 상위 주소 |
| `get_chain_analytics()` | 평균 tx/블록, 상위 검증자, 총 수수료 |
| `search(query)` | 숫자→블록, 64자→tx, 42자→주소 자동 판별 |

---

## PART 9: main() SECTION 1~7

| SECTION | 내용 |
|---------|------|
| 1 | 제네시스 블록 + 30명 검증자 등록 → DPoS 선출 (21명) |
| 2 | 트랜잭션 생성 (TRANSFER 100개, STAKE 10개) + Mempool |
| 3 | DPoS 블록 생성 10슬롯 (합의 성공률, 평균 tx/블록) |
| 4 | 슬래싱 3명 + 보상 분배 상위 5명 출력 |
| 5 | 스마트 컨트랙트 배포 & VM 실행 2회 (gas 출력) |
| 6 | P2P 5노드 네트워크 + 블록 전파 + 파티션 복구 |
| 7 | 블록 탐색기 통계 + validate_chain() + 최종 박스 |

---

## 테스트 계획 (총 20개)

| Group | 클래스 | 테스트 | 수 |
|-------|--------|--------|---|
| A | `TestCryptoAndTransaction` | test_01~04 | 4개 |
| B | `TestBlock` | test_05~08 | 4개 |
| C | `TestDPoS` | test_09~12 | 4개 |
| D | `TestSmartContract` | test_13~16 | 4개 |
| E | `TestNetwork` | test_17~20 | 4개 |

| 번호 | 검증 내용 |
|------|---------|
| 01 | SHA-256 결정론적 출력, double_sha256 ≠ sha256 |
| 02 | Merkle Root 결정론, 빈 목록 처리 |
| 03 | 트랜잭션 해시 고유성, 음수 amount 거부 |
| 04 | 주소 형식 (42자, 0x prefix), seed별 고유 주소 |
| 05 | 제네시스: block_id=0, prev_hash="0"×64 |
| 06 | 블록 해시 체인 (prev_hash 일치) |
| 07 | Merkle Root 재계산 일치, 변조 시 불일치 |
| 08 | 잔액 추적 정확성, validate_chain() → True |
| 09 | 30명 등록 → run_election → 21명 선발 |
| 10 | 위임 후 total_stake 누적, delegator_count 갱신 |
| 11 | reward = 100.0 × (10000/100000) = 10.0 수치 검증 |
| 12 | slash: new_stake = old × 0.95, amount = old × 0.05 |
| 13 | StackVM ADD: [PUSH 5, PUSH 3, ADD, HALT] → 8 |
| 14 | StackVM STORE/LOAD: storage["balance"] = 100 |
| 15 | GAS 소진: gas_used ≤ gas_limit, success=False |
| 16 | 컨트랙트 배포 → 주소 42자, call → success=True |
| 17 | 피어 가입/탈퇴, 중복 join → False |
| 18 | 블록 전파: propagated ≥ 2 (5노드 중) |
| 19 | 블록 탐색기: get_transaction, get_validator_ranking 정렬 |
| 20 | 통합: 제네시스 → 21명 선출 → 5슬롯 → validate_chain=True |

---

## 검증 방법
```bash
python university_v11_BLOCKCHAIN_DPOS.py
python test_v11_blockchain_dpos.py
# 기대: Ran 20 tests in X.XXXs — OK
git add university_v11_BLOCKCHAIN_DPOS.py test_v11_blockchain_dpos.py
git commit -m "feat: v11 블록체인 & DPoS 합의 엔진"
git push origin master
```

## v9.3 Raft vs v11 DPoS 핵심 차이
| 항목 | v9.3 Raft | v11 DPoS |
|------|-----------|---------|
| 합의 방식 | 로그 복제 + 과반수 | 위임 지분 투표 |
| 리더 선출 | 타임아웃 + term | 지분 순위 상위 21명 |
| 경제 인센티브 | 없음 | 보상 + 슬래싱 |
| 포크 해결 | term + 인덱스 | 최장 체인 규칙 |
| 스마트 컨트랙트 | 없음 | 스택 VM 14 옵코드 |
