# SSH 단축키 설정 가이드

**Status**: ✅ **이미 설정됨**

---

## 🚀 빠른 사용법

### 253 서버 접속
```bash
# 단축키로 간단하게 접속
ssh 253

# 또는 별칭 사용
ssh linux-dev
```

### Phase B 테스트 실행
```bash
# 1. SSH 접속
ssh 253

# 2. 프로젝트 디렉토리로 이동
cd ~/freelang-vm

# 3. 테스트 실행
cargo test
```

### 원라이너 (한 줄로 실행)
```bash
# 전체 테스트 실행
ssh 253 "cd ~/freelang-vm && cargo test"

# 릴리스 빌드로 테스트
ssh 253 "cd ~/freelang-vm && cargo test --release"

# 벤치마크 실행
ssh 253 "cd ~/freelang-vm && cargo bench"

# 결과를 파일로 저장
ssh 253 "cd ~/freelang-vm && cargo test 2>&1" > test_results.log
```

---

## 🔧 SSH 설정 내용

**파일 위치**: `~/.ssh/config`

```
Host 253
    HostName 192.168.45.253
    User kimjin
    IdentityFile ~/.ssh/id_ed25519
    Port 22
    ServerAliveInterval 60
    ServerAliveCountMax 3
    StrictHostKeyChecking accept-new
```

**설정 항목**:
- **Host**: 단축키 이름 (`ssh 253`으로 사용)
- **HostName**: 실제 IP 주소 (192.168.45.253)
- **User**: 로그인 사용자 (kimjin)
- **IdentityFile**: SSH 개인키 위치 (~/.ssh/id_ed25519)
- **Port**: SSH 포트 (기본값 22)
- **ServerAliveInterval**: 연결 유지 주기 (60초)
- **ServerAliveCountMax**: 재연결 시도 횟수

---

## ✅ 검증 체크리스트

- [x] SSH 설정 파일 존재 ✅
- [x] Host 253 설정됨 ✅
- [x] SSH 개인키 존재 (~/.ssh/id_ed25519) ✅
- [x] 키 권한 올바름 (600) ✅
- [x] SSH config 문법 정상 ✅
- [x] 별칭 설정 (linux-dev) ✅
- [ ] 네트워크 연결 확인 ⏳ (현재 타임아웃)

---

## 🔐 SSH 키 정보

```
Key Type:     Ed25519 (Modern, secure)
Key Location: ~/.ssh/id_ed25519
Permissions:  rw------- (600)
Key Size:     411 bytes
Created:      Feb 26 19:02
Status:       ✅ 설정됨
```

---

## 📡 네트워크 연결 방법

### 방법 1: Direct SSH (권장) 🌟
```bash
# 직접 접속 (단축키 사용)
ssh 253

# 또는 명시적 지정
ssh -i ~/.ssh/id_ed25519 -p 22 kimjin@192.168.45.253
```

### 방법 2: Console Access
```bash
# 하이퍼바이저 콘솔로 직접 접속
# 웹 브라우저에서 https://[hypervisor]:8443
```

### 방법 3: VPN 연결
```bash
# VPN 클라이언트로 연결
vpn connect 192.168.45.253
ssh 253
```

### 방법 4: Port Forwarding
```bash
# 로컬 포트로 포워딩
ssh -L 2222:192.168.45.253:22 bastion
ssh -p 2222 kimjin@localhost
```

### 방법 5: SSH Jump Host
```bash
# 중간 서버를 통한 접속
ssh -J jumphost 253
```

---

## 🛠️ 로컬 테스트 대안 (SSH 불가 시)

### Option A: Rust 설치 (로컬 Termux)
```bash
# 한 번만 실행
rustup default stable

# 테스트 실행
cd ~/freelang-vm
cargo test

# 결과 저장
cargo test 2>&1 | tee local_test_results.log
```

### Option B: Docker 사용
```bash
# Docker 이미지로 테스트
docker run --rm -v ~/freelang-vm:/work rust:latest \
  bash -c "cd /work && cargo test"
```

### Option C: GitHub Actions (CI/CD)
```bash
# 자동 테스트 (푸시 시 자동 실행)
# .github/workflows/test.yml 추가
```

---

## 🐛 문제 해결

### SSH 타임아웃 발생 시
```bash
# 1. 네트워크 확인
ping 192.168.45.253

# 2. 포트 확인
telnet 192.168.45.253 22

# 3. SSH 디버그 모드
ssh -vvv 253

# 4. 키 권한 확인
ls -la ~/.ssh/
chmod 600 ~/.ssh/id_ed25519
chmod 700 ~/.ssh

# 5. SSH 에이전트 확인
ssh-add -l
ssh-add ~/.ssh/id_ed25519
```

### 권한 오류 발생 시
```bash
# SSH 키 권한 수정
chmod 600 ~/.ssh/id_ed25519
chmod 644 ~/.ssh/id_ed25519.pub
chmod 700 ~/.ssh
```

### 호스트 키 검증 실패 시
```bash
# known_hosts에 추가
ssh-keyscan -H 192.168.45.253 >> ~/.ssh/known_hosts

# 또는 설정에서 자동 허용
# StrictHostKeyChecking accept-new
```

---

## 📊 로그 수집

### SSH 연결 로그
```bash
# 상세 디버그 로그 저장
ssh -vvv 253 > ssh_debug.log 2>&1

# 또는 SSH 로그
ssh 253 "tail -100 ~/.ssh/ssh_debug.log"
```

### 테스트 결과 저장
```bash
# 254 서버에서 테스트 실행 후 결과 다운로드
ssh 253 "cd ~/freelang-vm && cargo test 2>&1" > phase_b_test_results.log

# 또는 SCP로 파일 전송
scp 253:~/freelang-vm/test_results.log ./test_results.log
```

---

## 🎯 권장 작업 흐름

```
1. SSH 단축키 확인
   └─ ssh 253

2. 프로젝트 디렉토리 이동
   └─ cd ~/freelang-vm

3. Phase B 테스트 실행
   └─ cargo test

4. 결과 확인
   └─ 99/99 tests PASS

5. 벤치마크 실행 (선택)
   └─ cargo bench

6. 결과 저장
   └─ 로그 파일 다운로드
```

---

## 🚀 원라이너 명령어

```bash
# 전체 테스트 + 벤치마크 + 저장
ssh 253 "cd ~/freelang-vm && \
  echo '=== TEST ===' && cargo test && \
  echo '=== BENCHMARK ===' && cargo bench && \
  echo '=== DONE ===' " | tee full_results.log

# 또는 더 간단하게
ssh 253 'cd ~/freelang-vm && cargo test --release' > results.log 2>&1

# 스트레스 테스트
ssh 253 'cd ~/freelang-vm && cargo test stress'

# 특정 테스트만
ssh 253 'cd ~/freelang-vm && cargo test test_compact_'
```

---

## 📝 SSH 설정 추가 방법

### 새 호스트 추가 (필요시)
```bash
# SSH 설정 파일 편집
nano ~/.ssh/config

# 내용 추가 (예: prod 서버)
Host prod
    HostName 192.168.45.254
    User deploy
    IdentityFile ~/.ssh/id_ed25519
    Port 22

# 저장 후 사용
ssh prod
```

---

## 🎓 SSH 단축키 고급 팁

### 키 체인 등록
```bash
# 키를 SSH 에이전트에 등록 (매번 비번 안 물음)
ssh-add ~/.ssh/id_ed25519

# 등록된 키 확인
ssh-add -l
```

### 설정 테스트
```bash
# SSH 설정 문법 검증
ssh -G 253  # 실제 적용되는 설정 보기
```

### SSH Config 예제 모음
```bash
# 호스트별 다른 포트
Host github
    HostName github.com
    Port 22
    User git

# VPN 연결 필요
Host internal
    HostName 192.168.1.10
    ProxyCommand ssh -W %h:%p bastion

# 자동 로그인 설정
Host gitlab
    HostName gitlab.example.com
    User deploy
    IdentityFile ~/.ssh/deploy_key
    IdentitiesOnly yes
```

---

## ✅ 최종 체크리스트

- [x] SSH 단축키 설정 완료 (`ssh 253`)
- [x] SSH 키 설정 확인 (Ed25519)
- [x] SSH config 문법 정상
- [x] 별칭 설정 (linux-dev)
- [x] 테스트 실행 명령어 준비
- [ ] 네트워크 연결 (현재 미연결)
- [ ] Phase B 테스트 실행 (대기)

---

## 📞 다음 단계

1. **네트워크 연결 확인**
   ```bash
   ping 192.168.45.253
   ```

2. **SSH 접속 시도**
   ```bash
   ssh 253
   ```

3. **테스트 실행**
   ```bash
   cd ~/freelang-vm && cargo test
   ```

4. **결과 저장**
   ```bash
   cargo test 2>&1 | tee test_results.log
   ```

---

**Status**: ✅ SSH 단축키 설정 완료
**다음**: 네트워크 연결 후 `ssh 253` 접속

