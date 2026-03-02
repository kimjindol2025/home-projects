# 🔐 GRIE SSH 접속 가이드

## 📍 서버 정보

| 구분 | 정보 |
|------|------|
| **내부 IP** | 192.168.45.253 |
| **외부 IP** | 123.212.111.26 |
| **사용자** | kimjin |
| **SSH 포트** | 22 (내부), 22253 (모바일) |
| **키 타입** | ed25519 |
| **클라이언트** | kim@gorie-engine |

---

## 🖥️ 방법 1: PC/Termux에서 내부 접속

### 가장 간단한 방법 (SSH Config 이용)

```bash
# ~/.ssh/config가 설정되어 있으면:
ssh 253

# 또는
ssh linux-dev
```

### SSH Config 없이 직접 접속

```bash
ssh -i ~/.ssh/id_ed25519 kimjin@192.168.45.253
```

---

## 📱 방법 2: 모바일(Termux)에서 외부 접속

### 시나리오
```
모바일(Termux) 123.212.111.26:22253 → 포트포워딩 → 192.168.45.253:22
```

### A. SSH Config 이용 (권장)

```bash
ssh mobile-ssh
```

### B. 직접 명령어

```bash
ssh -i ~/.ssh/id_ed25519 kimjin@123.212.111.26 -p 22253
```

### C. SSH Config에서 수동 설정

```bash
cat >> ~/.ssh/config << 'EOF'
Host mobile-ssh
    HostName 123.212.111.26
    User kimjin
    IdentityFile ~/.ssh/id_ed25519
    Port 22253
    ServerAliveInterval 60
EOF

ssh mobile-ssh
```

---

## 🔧 Termux 특화 설정

### Termux에 SSH 클라이언트 설치

```bash
pkg update
pkg install openssh
ssh -V  # 확인
```

### 개인키 권한 설정 (매우 중요!)

```bash
chmod 600 ~/.ssh/id_ed25519
chmod 700 ~/.ssh
ls -la ~/.ssh/
```

### SSH 접속 테스트

```bash
# 내부 네트워크
ssh -i ~/.ssh/id_ed25519 kimjin@192.168.45.253

# 외부 네트워크 (모바일)
ssh -i ~/.ssh/id_ed25519 -p 22253 kimjin@123.212.111.26
```

---

## 🚀 GRIE 커밋 & 푸시 (원격 Zig 컴파일)

### 방법 A: SSH로 접속 후 원격 컴파일

```bash
# 1. 서버 접속
ssh 253

# 2. GRIE 프로젝트 클론 (또는 디렉토리 이동)
cd grie-engine

# 3. Zig 빌드
cd internal/kernel
make build-all OPT=ReleaseFast

# 4. 결과 확인
ls -lh ../../bin/*/lib/
```

### 방법 B: SCP로 파일 전송 (Zig 빌드 결과 다운로드)

```bash
# 서버에서 로컬로 복사
scp -i ~/.ssh/id_ed25519 -r kimjin@192.168.45.253:~/grie-engine/bin ./

# 또는
scp -i ~/.ssh/id_ed25519 -P 22253 -r kimjin@123.212.111.26:~/grie-engine/bin ./
```

### 방법 C: SSH 터널을 통한 Git 푸시

```bash
# 1. 서버에 SSH 키 복사 (또는 SSH config 설정)
ssh 253

# 2. 서버에서 Git 설정
git remote add origin git@gogs.dclub.kr:kim/raft-consensus-engine.git

# 3. 커밋 & 푸시
cd ~/grie-engine
git add .
git commit -m "메시지"
git push origin main
```

---

## 📊 원격 컴파일 워크플로우

```
┌─ 로컬 (Termux/Android) ────────┐
│  • 코드 작성                    │
│  • git commit                   │
├─ SSH 접속 (22253) ─────────────┤
│  192.168.45.253:22              │
├─────────────────────────────────┤
│ 서버 (Linux)                    │
│  • Zig 0.14+ 설치됨            │
│  • 크로스컴파일                 │
│    - aarch64-linux-android      │
│    - x86_64-linux-gnu           │
│    - aarch64-linux-gnu          │
├─────────────────────────────────┤
│ 결과:                           │
│  bin/native/lib/                │
│  bin/android/aarch64/lib/       │
│  bin/linux/x86_64/lib/          │
└─ SCP 다운로드 ──────────────────┘
```

---

## 🐛 트러블슈팅

### 1️⃣ "Permission denied (publickey)"

```bash
# 개인키 권한 확인
ls -la ~/.ssh/id_ed25519
# -rw------- (600) 이어야 함

# 고쳐야 하면:
chmod 600 ~/.ssh/id_ed25519
chmod 700 ~/.ssh
```

### 2️⃣ "Connection timed out"

```bash
# 포트 확인 (모바일은 22253)
ssh -v -i ~/.ssh/id_ed25519 -p 22253 kimjin@123.212.111.26

# 방화벽 확인
# 서버 관리자에게 포트 확인 요청
```

### 3️⃣ "Host key verification failed"

```bash
# 첫 접속 시:
ssh -i ~/.ssh/id_ed25519 kimjin@192.168.45.253

# 프롬프트에서 'yes' 입력
# ~/.ssh/known_hosts에 자동 등록됨
```

### 4️⃣ "No such file or directory: id_ed25519"

```bash
# SSH 키 생성
ssh-keygen -t ed25519 -C "kim@gorie-engine" -f ~/.ssh/id_ed25519 -N ""

# 또는 기존 키 위치 확인
ls -la ~/.ssh/
```

---

## 🔒 보안 체크리스트

- ✅ 개인키 권한: `chmod 600`
- ✅ SSH 디렉토리 권한: `chmod 700`
- ✅ 개인키 백업 (안전한 위치에)
- ✅ 공개키만 공유 (id_ed25519.pub)
- ✅ SSH config: `chmod 600`
- ✅ Termux에서 SSH 패키지 최신 상태 유지

---

## 📚 유용한 SSH 명령어

### 기본 접속
```bash
ssh 253                              # 설정된 호스트로 접속
ssh -i key.pem user@host             # 특정 키로 접속
ssh -p 2222 user@host                # 특정 포트로 접속
ssh -v user@host                     # Verbose (디버깅)
```

### 파일 전송
```bash
scp -r local_dir user@host:~/remote/ # 로컬 → 서버
scp -r user@host:~/remote/ local_dir # 서버 → 로컬
sftp user@host                       # SFTP 대화형 접속
```

### 터널 & 포워딩
```bash
ssh -L 8080:localhost:3000 user@host # 로컬 포트 포워딩
ssh -R 8080:localhost:3000 user@host # 원격 포트 포워딩
ssh -D 1080 user@host                # SOCKS 프록시
```

### 기타
```bash
ssh-copy-id -i ~/.ssh/id_ed25519 user@host  # 공개키 등록
ssh-keygen -l -f ~/.ssh/id_ed25519.pub     # 키 지문 보기
ssh-add ~/.ssh/id_ed25519                  # ssh-agent에 추가
```

---

## 🎯 다음 단계: Step 3 원격 컴파일

**목표**: 서버에서 Zig 커널 크로스컴파일 및 성능 검증

```bash
# 1. 서버 접속
ssh 253

# 2. GRIE 디렉토리
cd ~/grie-engine/internal/kernel

# 3. 전체 빌드
make build-all OPT=ReleaseFast

# 4. 결과 확인 및 벤치마크
ls -lh ../../bin/*/lib/
file ../../bin/native/lib/libgrie_kernel.a
```

---

**작성자**: Claude Opus 4.6
**작성일**: 2026-02-26
**상태**: ✅ 완료
