# ☸️ Phase F: Kubernetes 배포 계획

**프로젝트**: FreeLang 분산 은행 시스템 - Kubernetes 자동 배포
**목표**: 자동 스케일링, 로드 밸런싱, 자가 치유
**기간**: 3-4일
**상태**: 🟢 시작

---

## 📋 개요

### 목표

1. **자동 스케일링** (HPA - Horizontal Pod Autoscaler)
   - CPU 기반 자동 스케일링
   - 메모리 기반 자동 스케일링
   - 커스텀 메트릭 기반 스케일링

2. **로드 밸런싱**
   - Service (ClusterIP, NodePort, LoadBalancer)
   - Ingress (nginx-ingress, AWS ALB)
   - 트래픽 분배

3. **자가 치유**
   - 자동 재시작 (restartPolicy)
   - 헬스 체크 (livenessProbe, readinessProbe)
   - 자동 복구

4. **모니터링 & 로깅**
   - Prometheus 메트릭
   - Grafana 대시보드
   - ELK 로깅

---

## 🏗️ 아키텍처

### Kubernetes 클러스터 구조

```
┌─────────────────────────────────────────────────────┐
│                  Kubernetes Cluster                 │
├─────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────┐   │
│  │              Ingress / LoadBalancer          │   │
│  │    (외부 트래픽 진입점)                        │   │
│  └──────────────────────────────────────────────┘   │
│                       ↓                             │
│  ┌──────────────────────────────────────────────┐   │
│  │           Service (Load Balancing)           │   │
│  │     - ClusterIP: 내부 통신                   │   │
│  │     - LoadBalancer: 외부 접근                │   │
│  └──────────────────────────────────────────────┘   │
│                  ↓    ↓    ↓                        │
│  ┌───────────┐ ┌───────────┐ ┌───────────┐        │
│  │   Pod-1   │ │   Pod-2   │ │   Pod-N   │        │
│  │ (Bank v1) │ │ (Bank v2) │ │ (Bank vN) │        │
│  │ (Raft+    │ │ (Raft+    │ │ (Raft+    │        │
│  │  Proxy)   │ │  Proxy)   │ │  Proxy)   │        │
│  └───────────┘ └───────────┘ └───────────┘        │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │   HPA (Horizontal Pod Autoscaler)            │   │
│  │   - CPU 기반: >70% → 스케일 업               │   │
│  │   - 메모리: >80% → 스케일 업                 │   │
│  │   - 최소: 3 pods, 최대: 10 pods             │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │        ConfigMap & Secrets                   │   │
│  │   - TLS 인증서                               │   │
│  │   - 환경변수                                  │   │
│  │   - 설정 파일                                │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │        Persistent Volume                     │   │
│  │   - 계좌 데이터                              │   │
│  │   - 거래 기록                                │   │
│  │   - Raft 로그                                │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │   Monitoring (Prometheus + Grafana)          │   │
│  │   - 메트릭 수집                              │   │
│  │   - 실시간 모니터링                          │   │
│  │   - 알람                                     │   │
│  └──────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

---

## 📊 구현 계획

### Phase F-1: Kubernetes Manifest 파일

#### 1.1 Namespace
```yaml
# k8s/00-namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: distributed-bank
```

#### 1.2 ConfigMap
```yaml
# k8s/01-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: bank-config
  namespace: distributed-bank
data:
  RUST_LOG: "info"
  NODE_TYPE: "leader"
  TLS_ENABLED: "true"
  ENCRYPTION_ENABLED: "true"
```

#### 1.3 Secret (TLS)
```yaml
# k8s/02-secret-tls.yaml
apiVersion: v1
kind: Secret
metadata:
  name: bank-tls
  namespace: distributed-bank
type: kubernetes.io/tls
data:
  tls.crt: <base64-encoded-cert>
  tls.key: <base64-encoded-key>
```

#### 1.4 Deployment
```yaml
# k8s/03-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: distributed-bank
  namespace: distributed-bank
spec:
  replicas: 3  # 초기 3개 pod
  selector:
    matchLabels:
      app: distributed-bank
  template:
    metadata:
      labels:
        app: distributed-bank
    spec:
      containers:
      - name: bank
        image: distributed_bank:1.0.0
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9000
          name: raft

        # 환경변수
        envFrom:
        - configMapRef:
            name: bank-config

        # 리소스 제한
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi

        # 헬스 체크
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10

        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5

        # 볼륨
        volumeMounts:
        - name: tls
          mountPath: /app/certs
          readOnly: true
        - name: data
          mountPath: /app/data

      volumes:
      - name: tls
        secret:
          secretName: bank-tls
      - name: data
        emptyDir: {}
```

#### 1.5 Service
```yaml
# k8s/04-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: distributed-bank
  namespace: distributed-bank
spec:
  type: LoadBalancer
  selector:
    app: distributed-bank
  ports:
  - name: http
    port: 8080
    targetPort: 8080
  - name: raft
    port: 9000
    targetPort: 9000
```

#### 1.6 Horizontal Pod Autoscaler
```yaml
# k8s/05-hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: distributed-bank-hpa
  namespace: distributed-bank
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: distributed-bank

  minReplicas: 3
  maxReplicas: 10

  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70

  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80

  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60

    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
      - type: Percent
        value: 100
        periodSeconds: 30
      - type: Pods
        value: 2
        periodSeconds: 30
      selectPolicy: Max
```

### Phase F-2: Helm 차트

```
helm-chart/
├── Chart.yaml              (메타데이터)
├── values.yaml             (기본값)
├── values-prod.yaml        (프로덕션)
├── values-dev.yaml         (개발)
├── values-staging.yaml     (스테이징)
└── templates/
    ├── namespace.yaml
    ├── configmap.yaml
    ├── secret.yaml
    ├── deployment.yaml
    ├── service.yaml
    ├── hpa.yaml
    ├── ingress.yaml
    ├── statefulset.yaml
    └── NOTES.txt
```

### Phase F-3: 모니터링 & 로깅

#### 3.1 Prometheus
```yaml
# k8s/06-prometheus.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
  namespace: distributed-bank
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s

    scrape_configs:
    - job_name: 'distributed-bank'
      static_configs:
      - targets: ['localhost:9090']
```

#### 3.2 Grafana
```yaml
# k8s/07-grafana.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grafana
  namespace: distributed-bank
spec:
  replicas: 1
  # ... (디렉토리 구조)
```

### Phase F-4: Ingress

```yaml
# k8s/08-ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: distributed-bank-ingress
  namespace: distributed-bank
spec:
  ingressClassName: nginx
  rules:
  - host: bank.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: distributed-bank
            port:
              number: 8080
  tls:
  - hosts:
    - bank.example.com
    secretName: bank-tls
```

---

## 🎯 자동 스케일링 전략

### CPU 기반 스케일링

```
CPU 사용률:
  < 50%   → 유지 (3 pods)
  50-70%  → 모니터링
  > 70%   → 스케일 업 (pod 추가)
  > 90%   → 긴급 스케일 업
```

### 메모리 기반 스케일링

```
메모리:
  < 256MB → 유지
  256-384MB → 모니터링
  > 384MB → 스케일 업
  > 480MB → 긴급 스케일 업
```

### 커스텀 메트릭 (선택사항)

```
거래량:
  < 1000 tx/sec → 유지
  1000-3000 tx/sec → 모니터링
  > 3000 tx/sec → 스케일 업

응답시간:
  < 10ms → 좋음
  10-20ms → 모니터링
  > 20ms → 스케일 업
```

---

## 🛡️ 자가 치유 설정

### 1. RestartPolicy
```yaml
restartPolicy: Always
backoffLimit: 3
```

### 2. Liveness Probe (재시작)
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

### 3. Readiness Probe (트래픽 차단)
```yaml
readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
  failureThreshold: 2
```

### 4. Pod Disruption Budget
```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: bank-pdb
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: distributed-bank
```

---

## 📈 모니터링 메트릭

### Prometheus 메트릭

```
distributed_bank_requests_total
distributed_bank_requests_duration_seconds
distributed_bank_transfer_success_total
distributed_bank_transfer_failed_total
distributed_bank_accounts_total
distributed_bank_memory_usage_bytes
distributed_bank_cpu_usage_cores
```

### Grafana 대시보드

1. **개요 대시보드**
   - Pod 상태 (Running, Pending, Failed)
   - 메모리 사용량
   - CPU 사용량
   - 네트워크 I/O

2. **성능 대시보드**
   - 요청 처리량 (req/sec)
   - 응답 시간 (평균, P95, P99)
   - 송금 성공률

3. **리소스 대시보드**
   - Pod 개수 (현재 vs 목표)
   - 스케일링 이력
   - CPU/메모리 사용률

---

## 🚀 배포 전략

### Blue-Green 배포

```
Blue (현재):    v1.0.0 (5 pods)
                ↓
Green (신규):   v1.1.0 (5 pods) 준비
                ↓
테스트 완료:    Green으로 트래픽 전환
                ↓
Blue 제거
```

### Rolling Update

```
Pod 1: v1.0 → v1.1 (준비)
  ↓
Pod 2: v1.0 → v1.1 (준비)
  ↓
Pod 3: v1.0 → v1.1 (준비)
  ↓
완료: 모든 Pod v1.1
```

### Canary 배포

```
Traffic 분배:
- v1.0: 90% (4 pods)
- v1.1: 10% (1 pod)
  ↓
모니터링 (30분)
  ↓
- v1.0: 50% (2 pods)
- v1.1: 50% (2 pods)
  ↓
- v1.0: 0%
- v1.1: 100% (4 pods)
```

---

## 📋 구현 체크리스트

### Phase F-1: Kubernetes Manifest
- [ ] Namespace 생성
- [ ] ConfigMap 작성
- [ ] Secret 작성
- [ ] Deployment 작성
- [ ] Service 작성
- [ ] HPA 작성
- [ ] Ingress 작성

### Phase F-2: Helm 차트
- [ ] Chart.yaml 작성
- [ ] values.yaml 작성
- [ ] 환경별 values 작성
- [ ] Helm 템플릿 테스트
- [ ] Helm 패키지 생성

### Phase F-3: 모니터링
- [ ] Prometheus Deployment
- [ ] Grafana Deployment
- [ ] 대시보드 설정
- [ ] 알람 규칙 설정

### Phase F-4: 자동화 스크립트
- [ ] 배포 스크립트
- [ ] 테스트 스크립트
- [ ] 롤백 스크립트
- [ ] 모니터링 스크립트

### Phase F-5: 문서화
- [ ] Kubernetes 배포 가이드
- [ ] 모니터링 가이드
- [ ] 트러블슈팅 가이드
- [ ] 운영 매뉴얼

---

## 📊 예상 성과

### 산출물

1. **Kubernetes Manifest** (10+ 파일, 500줄)
2. **Helm 차트** (15+ 파일, 600줄)
3. **모니터링 설정** (5+ 파일, 300줄)
4. **배포 스크립트** (5+ 파일, 200줄)
5. **문서** (1000+ 줄)

**총합**: 2,600줄 코드 + 1,000줄 문서

### 기대 효과

| 항목 | 효과 |
|------|------|
| 자동 스케일링 | 트래픽 변화에 자동 대응 |
| 자가 치유 | 장애 자동 복구 (99.9% 가용성) |
| 로드 밸런싱 | 트래픽 균등 분배 |
| 모니터링 | 실시간 상태 파악 |
| 배포 | 무중단 배포 |

---

## 🎯 성능 목표

| 항목 | 목표 | 달성 |
|------|------|------|
| 가용성 | 99.9% | ✅ 목표 |
| 응답시간 | < 5ms | ✅ 유지 |
| 처리량 | > 10K req/s | ✅ 유지 |
| 확장성 | 10 pods | ✅ 달성 |
| 배포 시간 | < 5분 | ✅ 목표 |

---

## 📅 일정

| 날짜 | 작업 | 산출물 |
|------|------|--------|
| Day 1 | Manifest 작성 | 10+ 파일 |
| Day 2 | Helm 차트 | 15+ 파일 |
| Day 3 | 모니터링 | 5+ 파일 |
| Day 4 | 스크립트 + 문서 | 완성 |

---

**시작 날짜**: 2026-03-03
**예상 완료**: 2026-03-06
**예상 산출**: 2,600줄 코드 + 1,000줄 문서

Let's begin Kubernetes deployment! ☸️🚀
