#!/bin/bash
# ============================================
# Kubernetes 배포 스크립트
# ============================================
# FreeLang 분산 은행 시스템을 Kubernetes에 배포

set -e

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 설정
NAMESPACE="distributed-bank"
DEPLOYMENT_NAME="distributed-bank"
IMAGE="${1:-distributed_bank:1.0.0}"
KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}"

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Kubernetes 배포 스크립트              ║${NC}"
echo -e "${BLUE}║  FreeLang 분산 은행 시스템             ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo

# 1. 클러스터 확인
echo -e "${YELLOW}1️⃣  Kubernetes 클러스터 확인...${NC}"
if kubectl cluster-info &> /dev/null; then
  echo -e "${GREEN}✅ 클러스터 연결 성공${NC}"
  kubectl cluster-info
else
  echo -e "${RED}❌ 클러스터에 연결할 수 없습니다${NC}"
  exit 1
fi

echo

# 2. 네임스페이스 생성
echo -e "${YELLOW}2️⃣  네임스페이스 생성...${NC}"
if kubectl get namespace "$NAMESPACE" &> /dev/null; then
  echo -e "${GREEN}✅ 네임스페이스 이미 존재${NC}"
else
  kubectl create namespace "$NAMESPACE"
  echo -e "${GREEN}✅ 네임스페이스 생성 완료${NC}"
fi

echo

# 3. Secret 생성 (TLS 인증서)
echo -e "${YELLOW}3️⃣  TLS Secret 생성...${NC}"
if [ -f "certs/server.crt" ] && [ -f "certs/server.key" ]; then
  kubectl create secret tls distributed-bank-tls \
    --cert=certs/server.crt \
    --key=certs/server.key \
    --namespace="$NAMESPACE" \
    --dry-run=client -o yaml | kubectl apply -f -
  echo -e "${GREEN}✅ TLS Secret 생성 완료${NC}"
else
  echo -e "${YELLOW}⚠️  인증서 파일이 없습니다 (선택사항)${NC}"
fi

echo

# 4. Manifest 파일 적용
echo -e "${YELLOW}4️⃣  Kubernetes Manifest 적용...${NC}"

echo "  - Namespace..."
kubectl apply -f k8s/00-namespace.yaml

echo "  - ConfigMap..."
kubectl apply -f k8s/01-configmap.yaml

echo "  - Deployment..."
sed "s|distributed_bank:1.0.0|$IMAGE|g" k8s/02-deployment.yaml | kubectl apply -f -

echo "  - Service..."
kubectl apply -f k8s/03-service.yaml

echo "  - Autoscaling..."
kubectl apply -f k8s/04-autoscaling.yaml

echo "  - Ingress..."
kubectl apply -f k8s/05-ingress.yaml

echo "  - Monitoring..."
kubectl apply -f k8s/06-monitoring.yaml

echo -e "${GREEN}✅ Manifest 적용 완료${NC}"

echo

# 5. 배포 상태 대기
echo -e "${YELLOW}5️⃣  배포 상태 대기 중...${NC}"
kubectl rollout status deployment/"$DEPLOYMENT_NAME" \
  --namespace="$NAMESPACE" \
  --timeout=5m

echo -e "${GREEN}✅ 배포 완료${NC}"

echo

# 6. 배포 정보 조회
echo -e "${YELLOW}6️⃣  배포 정보 조회...${NC}"
echo
echo -e "${BLUE}📊 Pod 상태:${NC}"
kubectl get pods --namespace="$NAMESPACE" -o wide

echo
echo -e "${BLUE}📊 Service:${NC}"
kubectl get svc --namespace="$NAMESPACE" -o wide

echo
echo -e "${BLUE}📊 HPA 상태:${NC}"
kubectl get hpa --namespace="$NAMESPACE"

echo

# 7. 접근 정보
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  배포 완료!                           ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo

echo -e "${GREEN}✅ 배포 성공!${NC}"
echo
echo "📋 접근 정보:"
echo "  - 네임스페이스: $NAMESPACE"
echo "  - Deployment: $DEPLOYMENT_NAME"
echo "  - Service ClusterIP: $(kubectl get svc distributed-bank -n $NAMESPACE -o jsonpath='{.spec.clusterIP}')"

EXTERNAL_IP=$(kubectl get svc distributed-bank-lb -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "Pending")
echo "  - Service LoadBalancer IP: $EXTERNAL_IP"

echo
echo "🔗 다음 단계:"
echo "  1. Pod 로그 확인:"
echo "     kubectl logs -f deployment/$DEPLOYMENT_NAME -n $NAMESPACE"
echo
echo "  2. Pod 진입:"
echo "     kubectl exec -it <pod-name> -n $NAMESPACE -- /bin/bash"
echo
echo "  3. 포워딩:"
echo "     kubectl port-forward svc/distributed-bank 8080:8080 -n $NAMESPACE"
echo
echo "  4. 모니터링:"
echo "     kubectl get pods -n $NAMESPACE -w"
echo

echo "✅ 배포 완료!"
