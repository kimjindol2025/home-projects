# Helm 배포 가이드

분산 은행 시스템을 Helm을 통해 Kubernetes에 배포하는 완전한 가이드입니다.

---

## 📋 목차

1. [Helm 설치](#helm-설치)
2. [차트 구조](#차트-구조)
3. [환경별 배포](#환경별-배포)
4. [배포 명령어](#배포-명령어)
5. [사용자 정의 설정](#사용자-정의-설정)
6. [배포 검증](#배포-검증)
7. [문제 해결](#문제-해결)
8. [차트 업그레이드](#차트-업그레이드)

---

## 🚀 Helm 설치

### 1. Helm 설치 확인

```bash
# Helm 버전 확인
helm version

# 버전이 3.x 이상인지 확인
helm version --short
```

### 2. 저장소 추가 (선택사항)

```bash
# Helm Hub 저장소 추가
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

---

## 📂 차트 구조

```
helm/
├── Chart.yaml              # 차트 메타데이터
├── values.yaml             # 기본 설정값
├── values-dev.yaml         # 개발 환경 설정
├── values-staging.yaml     # 스테이징 환경 설정
├── values-prod.yaml        # 프로덕션 환경 설정
└── templates/              # Kubernetes 템플릿
    ├── _helpers.tpl        # 헬퍼 함수
    ├── namespace.yaml      # 네임스페이스
    ├── configmap.yaml      # 설정 맵
    ├── deployment.yaml     # 배포
    ├── service.yaml        # 서비스
    ├── autoscaling.yaml    # HPA & PDB
    ├── ingress.yaml        # Ingress & NetworkPolicy
    ├── pvc.yaml            # PersistentVolumeClaim
    ├── monitoring.yaml     # ServiceMonitor & PrometheusRule
    ├── NOTES.txt           # 배포 후 안내
```

---

## 🌍 환경별 배포

### 개발 환경 (Development)

```bash
# 최소 리소스로 빠른 반복 개발
helm install distributed-bank ./helm \
  -f helm/values-dev.yaml \
  -n distributed-bank-dev \
  --create-namespace
```

**특징**:
- 1개 Pod (자동 스케일 안함)
- 50m CPU, 64Mi 메모리 요청
- NodePort로 접근
- 로그 레벨: DEBUG
- pprof 디버깅 활성화

### 스테이징 환경 (Staging)

```bash
# 프로덕션과 유사한 설정으로 테스트
helm install distributed-bank ./helm \
  -f helm/values-staging.yaml \
  -n distributed-bank-staging \
  --create-namespace
```

**특징**:
- 2-5개 Pod (자동 스케일 활성화)
- 100m CPU, 128Mi 메모리 요청
- LoadBalancer로 접근
- 기본 모니터링 활성화
- 알람 규칙 활성화

### 프로덕션 환경 (Production)

```bash
# 고가용성, 보안, 성능 최적화
helm install distributed-bank ./helm \
  -f helm/values-prod.yaml \
  -n distributed-bank \
  --create-namespace
```

**특징**:
- 5-20개 Pod (적극적 자동 스케일)
- 200m CPU, 256Mi 메모리 요청
- 고성능 SSD 스토리지
- Ingress + TLS 활성화
- 완전한 모니터링 & 알람

---

## 🎯 배포 명령어

### 1. 차트 검증

```bash
# 템플릿 검증
helm lint ./helm

# 렌더링된 manifest 확인
helm template distributed-bank ./helm -f helm/values-dev.yaml

# Dry-run (실제 배포하지 않음)
helm install distributed-bank ./helm \
  -f helm/values-dev.yaml \
  --dry-run=client \
  --debug
```

### 2. 배포

```bash
# 개발 환경 배포
helm install my-release ./helm \
  -f helm/values-dev.yaml \
  -n distributed-bank-dev \
  --create-namespace

# 스테이징 환경 배포
helm install my-release ./helm \
  -f helm/values-staging.yaml \
  -n distributed-bank-staging \
  --create-namespace

# 프로덕션 환경 배포
helm install my-release ./helm \
  -f helm/values-prod.yaml \
  -n distributed-bank \
  --create-namespace
```

### 3. 배포 상태 확인

```bash
# 릴리즈 목록
helm list -n distributed-bank

# 릴리즈 상태
helm status my-release -n distributed-bank

# 릴리즈 값 확인
helm get values my-release -n distributed-bank
```

### 4. 업그레이드

```bash
# 새 버전으로 업그레이드
helm upgrade my-release ./helm \
  -f helm/values-prod.yaml \
  -n distributed-bank

# 이전 버전으로 롤백
helm rollback my-release 1 -n distributed-bank
```

### 5. 삭제

```bash
# 릴리즈 제거
helm uninstall my-release -n distributed-bank

# 네임스페이스도 함께 삭제 (주의!)
kubectl delete namespace distributed-bank
```

---

## 🔧 사용자 정의 설정

### 1. 단일 값 오버라이드

```bash
# 레플리카 수 변경
helm install my-release ./helm \
  -f helm/values-prod.yaml \
  --set deployment.replicas.initial=10 \
  -n distributed-bank

# 이미지 태그 변경
helm install my-release ./helm \
  --set image.tag=2.0.0 \
  -n distributed-bank
```

### 2. 사용자 정의 values 파일 생성

`values-custom.yaml`:
```yaml
# 기본값을 오버라이드하는 커스텀 설정
deployment:
  replicas:
    initial: 3
    min: 3
    max: 15

image:
  tag: "1.5.0"

configMap:
  rustLog: info

persistence:
  data:
    size: 20Gi
```

배포:
```bash
helm install my-release ./helm \
  -f helm/values-prod.yaml \
  -f helm/values-custom.yaml \
  -n distributed-bank
```

### 3. JSON으로 여러 값 설정

```bash
helm install my-release ./helm \
  -f helm/values-prod.yaml \
  --set 'deployment.resources.requests.cpu=500m,deployment.resources.limits.memory=2Gi' \
  -n distributed-bank
```

---

## ✅ 배포 검증

### 1. Pod 상태 확인

```bash
# Pod 조회
kubectl get pods -n distributed-bank

# Pod 상세 정보
kubectl describe pod <pod-name> -n distributed-bank

# Pod 로그
kubectl logs <pod-name> -n distributed-bank
```

### 2. 서비스 확인

```bash
# 서비스 목록
kubectl get svc -n distributed-bank

# LoadBalancer IP 확인
kubectl get svc distributed-bank-lb -n distributed-bank -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```

### 3. 헬스체크

```bash
# 포트 포워딩
kubectl port-forward svc/distributed-bank 8080:8080 -n distributed-bank &

# 헬스체크 요청
curl http://localhost:8080/health

# 준비 상태 확인
curl http://localhost:8080/ready

# 메트릭 확인
curl http://localhost:8888/metrics
```

### 4. 이벤트 확인

```bash
# 최근 이벤트
kubectl get events -n distributed-bank --sort-by='.lastTimestamp'

# 특정 Pod의 이벤트
kubectl describe pod <pod-name> -n distributed-bank
```

---

## 🔍 문제 해결

### Pod가 Running 상태가 아님

```bash
# Pod 상태 확인
kubectl get pods -n distributed-bank

# 상세 정보 확인
kubectl describe pod <pod-name> -n distributed-bank

# 로그 확인
kubectl logs <pod-name> -n distributed-bank

# 이전 로그 확인 (Crash했을 경우)
kubectl logs <pod-name> -n distributed-bank --previous
```

### ImagePullBackOff 에러

```bash
# 이미지 확인
kubectl describe pod <pod-name> -n distributed-bank

# 이미지 레지스트리 접근 확인
docker login <registry>

# ImagePullSecret 추가
kubectl create secret docker-registry regcred \
  --docker-server=<registry> \
  --docker-username=<username> \
  --docker-password=<password> \
  -n distributed-bank
```

### 리소스 부족

```bash
# 노드 상태 확인
kubectl get nodes
kubectl describe node <node-name>

# 클러스터 리소스 확인
kubectl top nodes
kubectl top pods -n distributed-bank

# HPA 상태 확인
kubectl get hpa -n distributed-bank
kubectl describe hpa distributed-bank-hpa -n distributed-bank
```

### Ingress가 작동하지 않음

```bash
# Ingress 상태 확인
kubectl get ingress -n distributed-bank
kubectl describe ingress distributed-bank-ingress -n distributed-bank

# Ingress Controller 확인
kubectl get pods -n ingress-nginx

# 인증서 확인
kubectl get secret distributed-bank-tls -n distributed-bank -o yaml
```

---

## 📈 차트 업그레이드

### 버전 업그레이드

```bash
# 현재 릴리즈 확인
helm list -n distributed-bank

# 차트 업그레이드
helm upgrade my-release ./helm \
  -f helm/values-prod.yaml \
  -n distributed-bank

# 업그레이드 진행 상황 확인
kubectl rollout status deployment/distributed-bank -n distributed-bank
```

### 롤백

```bash
# 릴리즈 히스토리 확인
helm history my-release -n distributed-bank

# 이전 버전으로 롤백
helm rollback my-release 1 -n distributed-bank

# 특정 리비전으로 롤백
helm rollback my-release 3 -n distributed-bank
```

### 매니페스트 확인

```bash
# 현재 릴리즈의 manifest 확인
helm get manifest my-release -n distributed-bank

# 렌더링된 values 확인
helm get values my-release -n distributed-bank

# 전체 정보 확인
helm get all my-release -n distributed-bank
```

---

## 📊 모니터링 & 로깅

### Prometheus 메트릭 수집

```bash
# ServiceMonitor 확인
kubectl get servicemonitor -n distributed-bank

# Prometheus 접근
kubectl port-forward -n monitoring svc/prometheus 9090:9090
# http://localhost:9090에서 메트릭 확인
```

### 로그 수집

```bash
# 로그 조회
kubectl logs -f deployment/distributed-bank -n distributed-bank

# 특정 시간대 로그
kubectl logs deployment/distributed-bank \
  --since=1h \
  -n distributed-bank

# 모든 Pod의 로그
kubectl logs -f -l app=distributed-bank -n distributed-bank
```

### 대시보드 접근

```bash
# Grafana 포워딩
kubectl port-forward -n monitoring svc/grafana 3000:3000
# http://localhost:3000에서 대시보드 확인
```

---

## 🛡️ 보안 체크리스트

- [ ] TLS 인증서 설정됨
- [ ] NetworkPolicy 활성화됨
- [ ] Non-root 사용자로 실행됨
- [ ] 읽기 전용 파일시스템 설정됨
- [ ] 리소스 제한 설정됨
- [ ] 보안 컨텍스트 설정됨
- [ ] RBAC 정책 구성됨
- [ ] Pod Disruption Budget 설정됨
- [ ] 모니터링 & 알람 활성화됨
- [ ] 로깅 수집 설정됨

---

## 📚 추가 리소스

- [Helm 공식 문서](https://helm.sh/docs/)
- [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)
- [Prometheus Operator](https://prometheus-operator.dev/)
- [RBAC 설정](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)

---

**마지막 업데이트**: 2026-03-02
**차트 버전**: 1.0.0
**앱 버전**: 1.0.0
