# distributed-bank Helm Chart

FreeLang 분산 은행 시스템을 Kubernetes에 배포하기 위한 Helm 차트입니다.

## 개요

이 Helm 차트는 다음을 제공합니다:

- **다중 환경 지원**: 개발(dev), 스테이징(staging), 프로덕션(prod)
- **자동 스케일링**: CPU/메모리 기반 Horizontal Pod Autoscaler
- **자가 치유**: liveness/readiness/startup probe를 통한 자동 복구
- **고가용성**: Pod anti-affinity, Pod Disruption Budget
- **보안**: TLS 1.3, NetworkPolicy, non-root 사용자
- **모니터링**: Prometheus 메트릭 수집, 자동 알람
- **로드 밸런싱**: 내부/외부 서비스, Ingress 지원

## 설치

### 사전 요구사항

- Kubernetes 1.20+
- Helm 3.0+
- 2GB+ 메모리, 1+ CPU
- `distributed-bank` 이미지가 사용 가능한 레지스트리

### 빠른 시작

```bash
# 개발 환경 배포
helm install distributed-bank ./helm \
  -f helm/values-dev.yaml \
  -n distributed-bank-dev \
  --create-namespace

# 프로덕션 환경 배포
helm install distributed-bank ./helm \
  -f helm/values-prod.yaml \
  -n distributed-bank \
  --create-namespace
```

## 환경별 설정

### 개발 (values-dev.yaml)

```yaml
- 1개 Pod (자동 스케일 비활성화)
- 최소 리소스 (50m CPU, 64Mi 메모리)
- NodePort 서비스
- DEBUG 로그 레벨
- Ingress 비활성화
```

**사용처**: 로컬 개발, 테스트

### 스테이징 (values-staging.yaml)

```yaml
- 2-5개 Pod
- 균형잡힌 리소스 (100m CPU, 128Mi 메모리)
- LoadBalancer 서비스
- 모니터링 & 알람 활성화
```

**사용처**: 프로덕션 배포 전 테스트

### 프로덕션 (values-prod.yaml)

```yaml
- 5-20개 Pod
- 넉넉한 리소스 (200m CPU, 256Mi 메모리)
- Ingress + TLS
- 완전한 모니터링 & 알람
- Pod anti-affinity 필수
- 고성능 스토리지
```

**사용처**: 실제 운영 환경

## 커스터마이징

### 단일 값 오버라이드

```bash
helm install distributed-bank ./helm \
  -f helm/values-prod.yaml \
  --set deployment.replicas.initial=10 \
  --set image.tag=2.0.0 \
  -n distributed-bank
```

### 사용자 정의 파일

`values-custom.yaml`:
```yaml
deployment:
  replicas:
    initial: 5
    min: 5
    max: 20

persistence:
  data:
    size: 50Gi
```

```bash
helm install distributed-bank ./helm \
  -f helm/values-prod.yaml \
  -f helm/values-custom.yaml \
  -n distributed-bank
```

## 배포

### 배포 확인

```bash
# 템플릿 검증
helm lint ./helm

# Dry-run
helm install distributed-bank ./helm \
  -f helm/values-dev.yaml \
  --dry-run=client \
  --debug

# 실제 배포
helm install distributed-bank ./helm \
  -f helm/values-dev.yaml \
  -n distributed-bank-dev \
  --create-namespace
```

### 상태 확인

```bash
# 릴리즈 상태
helm status distributed-bank -n distributed-bank-dev

# Pod 조회
kubectl get pods -n distributed-bank-dev

# 로그 확인
kubectl logs -f deployment/distributed-bank -n distributed-bank-dev

# 이벤트 확인
kubectl get events -n distributed-bank-dev
```

## 업그레이드 & 롤백

### 업그레이드

```bash
helm upgrade distributed-bank ./helm \
  -f helm/values-prod.yaml \
  -n distributed-bank
```

### 롤백

```bash
# 이전 버전으로
helm rollback distributed-bank -n distributed-bank

# 특정 리비전으로
helm rollback distributed-bank 3 -n distributed-bank
```

## 모니터링

### 메트릭 확인

```bash
kubectl port-forward svc/distributed-bank 8888:8888 -n distributed-bank
curl http://localhost:8888/metrics
```

### Prometheus

```bash
kubectl get servicemonitor -n distributed-bank
```

### 알람

프로메테우스에서 설정된 6개의 알람 규칙:
- HighCPUUsage (>75% for 5m)
- HighMemoryUsage (>80% for 5m)
- HighPodRestartRate (>0.05 per 15m)
- HighErrorRate (>1% for 5m)
- HighResponseTime (P99 >100ms for 5m)
- PodNotReady (>5m)

## 문제 해결

### Pod가 Running 상태가 아님

```bash
kubectl describe pod <pod-name> -n distributed-bank
kubectl logs <pod-name> -n distributed-bank
```

### ImagePullBackOff

```bash
kubectl describe pod <pod-name> -n distributed-bank
# 이미지 레지스트리 접근 확인
```

### CrashLoopBackOff

```bash
kubectl logs <pod-name> -n distributed-bank --previous
```

### 리소스 부족

```bash
kubectl top nodes
kubectl top pods -n distributed-bank
```

## 안내사항

자세한 배포 방법은 [HELM_DEPLOYMENT_GUIDE.md](../HELM_DEPLOYMENT_GUIDE.md)를 참고하세요.

## 파일 구조

```
helm/
├── Chart.yaml                # 차트 메타데이터
├── README.md                 # 차트 README
├── values.yaml              # 기본 설정값
├── values-dev.yaml          # 개발 환경 설정
├── values-staging.yaml      # 스테이징 환경 설정
├── values-prod.yaml         # 프로덕션 환경 설정
└── templates/
    ├── _helpers.tpl         # 헬퍼 함수
    ├── namespace.yaml       # 네임스페이스
    ├── configmap.yaml       # 설정 맵
    ├── deployment.yaml      # 배포
    ├── service.yaml         # 서비스
    ├── autoscaling.yaml     # HPA & PDB
    ├── ingress.yaml         # Ingress & NetworkPolicy
    ├── pvc.yaml             # 볼륨 청구
    ├── monitoring.yaml      # 모니터링
    ├── NOTES.txt            # 배포 후 안내
```

## 설정 옵션

주요 설정 옵션:

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `image.repository` | `distributed-bank` | 컨테이너 이미지 |
| `image.tag` | `1.0.0` | 이미지 태그 |
| `deployment.replicas.initial` | `3` | 초기 Pod 수 |
| `deployment.replicas.min` | `3` | 최소 Pod 수 |
| `deployment.replicas.max` | `10` | 최대 Pod 수 |
| `autoscaling.enabled` | `true` | 자동 스케일링 활성화 |
| `persistence.enabled` | `true` | 영구 저장소 활성화 |
| `monitoring.enabled` | `true` | 모니터링 활성화 |
| `global.environment` | `development` | 환경 이름 |

전체 옵션은 `values.yaml`을 참고하세요.

## 라이센스

MIT License

## 지원

문제가 있으면 [GitHub Issues](https://gogs.dclub.kr/kim/freelang-distributed-system/issues)를 통해 보고해주세요.
