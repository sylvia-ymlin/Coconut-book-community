# Kubernetes 部署指南

## 📋 概述

本文档介绍如何将 BookCommunity 部署到 Kubernetes 集群。

### 架构图

```
┌─────────────────────────────────────────────────────────┐
│                    Ingress Controller                    │
│                  (nginx-ingress)                         │
└────────────────────────┬────────────────────────────────┘
                         ↓
┌────────────────────────────────────────────────────────┐
│              BookCommunity Service (ClusterIP)          │
└────────────────────────┬───────────────────────────────┘
                         ↓
        ┌────────────────┼────────────────┐
        ↓                ↓                ↓
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│BookCommunity │  │BookCommunity │  │BookCommunity │
│   Pod 1      │  │   Pod 2      │  │   Pod 3      │
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘
       │                 │                 │
       └─────────────────┼─────────────────┘
                         ↓
        ┌────────────────┼────────────────┐
        ↓                ↓                ↓
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ PostgreSQL   │  │    Redis     │  │  RabbitMQ    │
│  Service     │  │   Service    │  │   Service    │
└──────────────┘  └──────────────┘  └──────────────┘
```

---

## 🚀 快速开始

### 前置要求

1. **Kubernetes 集群**
   - Minikube（本地开发）
   - GKE / EKS / AKS（生产环境）
   - K3s / Kind（轻量级）

2. **工具**
   ```bash
   # 必须
   - kubectl
   - docker
   
   # 可选
   - helm
   - k9s (K8s 可视化工具)
   ```

3. **资源要求**
   - 最小：2 CPU, 4GB RAM
   - 推荐：4 CPU, 8GB RAM

---

## 📦 方式一：使用部署脚本（推荐）

### 1. 使用交互式菜单

```bash
cd /Users/ymlin/Downloads/003-Study/137-Projects/BookCommunity

./scripts/k8s-deploy.sh
```

**菜单选项：**
```
1. 构建 Docker 镜像
2. 部署到 Kubernetes
3. 查看部署状态
4. 获取访问地址
5. 完整部署 (构建 + 部署)  ← 推荐首次使用
6. 清理部署
7. 退出
```

### 2. 使用命令行参数

```bash
# 完整部署
./scripts/k8s-deploy.sh full

# 仅部署
./scripts/k8s-deploy.sh deploy

# 查看状态
./scripts/k8s-deploy.sh status

# 清理
./scripts/k8s-deploy.sh clean
```

---

## 🔧 方式二：手动部署

### 1. 构建 Docker 镜像

```bash
# 构建镜像
docker build -t sylvia-ymlin/bookcommunity:latest .

# 推送到 Docker Hub
docker push sylvia-ymlin/bookcommunity:latest

# 或使用 Minikube 的 Docker 环境
eval $(minikube docker-env)
docker build -t bookcommunity:latest .
```

### 2. 部署到 Kubernetes

```bash
# 创建 namespace
kubectl apply -f k8s/base/namespace.yaml

# 部署配置和密钥
kubectl apply -f k8s/base/secret.yaml
kubectl apply -f k8s/base/configmap.yaml

# 部署数据库
kubectl apply -f k8s/base/postgres.yaml
kubectl apply -f k8s/base/redis.yaml
kubectl apply -f k8s/base/rabbitmq.yaml

# 等待数据库就绪
kubectl wait --for=condition=ready pod -l app=postgres -n bookcommunity --timeout=300s
kubectl wait --for=condition=ready pod -l app=redis -n bookcommunity --timeout=300s
kubectl wait --for=condition=ready pod -l app=rabbitmq -n bookcommunity --timeout=300s

# 部署应用
kubectl apply -f k8s/base/deployment.yaml
kubectl apply -f k8s/base/service.yaml
kubectl apply -f k8s/base/hpa.yaml

# 可选：部署 Ingress
kubectl apply -f k8s/base/ingress.yaml
```

### 3. 查看部署状态

```bash
# 查看所有资源
kubectl get all -n bookcommunity

# 查看 Pods
kubectl get pods -n bookcommunity

# 查看日志
kubectl logs -f deployment/bookcommunity -n bookcommunity

# 查看 HPA 状态
kubectl get hpa -n bookcommunity
```

---

## 🎯 方式三：使用 Helm（推荐生产环境）

### 1. 安装 Helm Chart

```bash
# 安装
helm install bookcommunity ./helm/bookcommunity -n bookcommunity --create-namespace

# 查看状态
helm status bookcommunity -n bookcommunity

# 查看值
helm get values bookcommunity -n bookcommunity
```

### 2. 自定义配置

创建 `values-custom.yaml`：

```yaml
image:
  repository: sylvia-ymlin/bookcommunity
  tag: v1.0.0

replicaCount: 5

resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "1Gi"
    cpu: "1000m"

autoscaling:
  minReplicas: 3
  maxReplicas: 20

ingress:
  enabled: true
  hosts:
    - host: bookcommunity.example.com
      paths:
        - path: /
          pathType: Prefix
```

```bash
# 使用自定义配置安装
helm install bookcommunity ./helm/bookcommunity \
  -n bookcommunity --create-namespace \
  -f values-custom.yaml
```

### 3. 升级和回滚

```bash
# 升级
helm upgrade bookcommunity ./helm/bookcommunity -n bookcommunity

# 回滚
helm rollback bookcommunity -n bookcommunity

# 卸载
helm uninstall bookcommunity -n bookcommunity
```

---

## 🌐 访问应用

### Minikube 环境

```bash
# 方式1：使用 Minikube service
minikube service bookcommunity-service -n bookcommunity

# 方式2：端口转发
kubectl port-forward svc/bookcommunity-service 8080:80 -n bookcommunity

# 访问
curl http://localhost:8080/health
```

### 云环境（GKE/EKS/AKS）

```bash
# 获取 LoadBalancer IP
kubectl get svc bookcommunity-service -n bookcommunity

# 或使用 Ingress
kubectl get ingress -n bookcommunity
```

---

## 📊 监控和调试

### 查看 Pods 状态

```bash
# 列出所有 Pods
kubectl get pods -n bookcommunity

# 查看详细信息
kubectl describe pod <pod-name> -n bookcommunity

# 查看日志
kubectl logs -f <pod-name> -n bookcommunity

# 进入 Pod
kubectl exec -it <pod-name> -n bookcommunity -- /bin/sh
```

### 查看 HPA 自动扩展

```bash
# 查看 HPA 状态
kubectl get hpa -n bookcommunity

# 详细信息
kubectl describe hpa bookcommunity-hpa -n bookcommunity

# 模拟负载测试
kubectl run -it --rm load-generator --image=busybox /bin/sh
# 在容器内执行
while true; do wget -q -O- http://bookcommunity-service.bookcommunity/health; done
```

### 查看资源使用

```bash
# 查看节点资源
kubectl top nodes

# 查看 Pod 资源
kubectl top pods -n bookcommunity

# 查看事件
kubectl get events -n bookcommunity --sort-by='.lastTimestamp'
```

---

## 🔒 生产环境配置

### 1. 更新 Secret

**⚠️ 重要：不要在生产环境使用默认密码！**

```bash
# 生成安全密码
DB_PASSWORD=$(openssl rand -base64 32)
REDIS_PASSWORD=$(openssl rand -base64 32)
JWT_SIGN_KEY=$(openssl rand -hex 16)
JWT_SECRET=$(openssl rand -hex 16)

# 创建 Secret
kubectl create secret generic bookcommunity-secret \
  --from-literal=DB_PASSWORD=${DB_PASSWORD} \
  --from-literal=REDIS_PASSWORD=${REDIS_PASSWORD} \
  --from-literal=JWT_SIGN_KEY_HEX=${JWT_SIGN_KEY} \
  --from-literal=JWT_SECRET_HEX=${JWT_SECRET} \
  -n bookcommunity
```

### 2. 配置持久化存储

使用云提供商的存储类：

```yaml
# storage-class.yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/gce-pd  # GKE
# provisioner: kubernetes.io/aws-ebs  # EKS
# provisioner: kubernetes.io/azure-disk  # AKS
parameters:
  type: pd-ssd  # SSD 类型
  replication-type: regional-pd  # 区域复制
```

### 3. 配置资源限制

```yaml
# resource-quota.yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: bookcommunity-quota
  namespace: bookcommunity
spec:
  hard:
    requests.cpu: "10"
    requests.memory: 20Gi
    limits.cpu: "20"
    limits.memory: 40Gi
    persistentvolumeclaims: "10"
```

### 4. 配置网络策略

```yaml
# network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: bookcommunity-netpol
  namespace: bookcommunity
spec:
  podSelector:
    matchLabels:
      app: bookcommunity
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - protocol: TCP
      port: 5432
  - to:
    - podSelector:
        matchLabels:
          app: redis
    ports:
    - protocol: TCP
      port: 6379
```

---

## 🐛 常见问题

### Q1: Pod 一直处于 Pending 状态

**原因：** 资源不足

**解决方案：**
```bash
# 查看节点资源
kubectl top nodes

# 查看 Pod 事件
kubectl describe pod <pod-name> -n bookcommunity

# 减少资源请求或增加节点
```

### Q2: 镜像拉取失败

**原因：** ImagePullBackOff

**解决方案：**
```bash
# 检查镜像是否存在
docker pull sylvia-ymlin/bookcommunity:latest

# 或使用 Minikube 本地镜像
eval $(minikube docker-env)
docker build -t bookcommunity:latest .

# 更新 deployment 使用本地镜像
kubectl set image deployment/bookcommunity bookcommunity=bookcommunity:latest -n bookcommunity
```

### Q3: 数据库连接失败

**原因：** PostgreSQL 未就绪

**解决方案：**
```bash
# 检查 PostgreSQL Pod
kubectl get pods -n bookcommunity -l app=postgres

# 查看日志
kubectl logs -l app=postgres -n bookcommunity

# 检查 Service
kubectl get svc postgres-service -n bookcommunity
```

### Q4: HPA 不工作

**原因：** Metrics Server 未安装

**解决方案：**
```bash
# Minikube 启用 metrics-server
minikube addons enable metrics-server

# 其他集群安装 Metrics Server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# 验证
kubectl get apiservice v1beta1.metrics.k8s.io
```

---

## 📚 参考资源

- [Kubernetes 官方文档](https://kubernetes.io/docs/)
- [Helm 文档](https://helm.sh/docs/)
- [Minikube 文档](https://minikube.sigs.k8s.io/docs/)
- [K8s 最佳实践](https://kubernetes.io/docs/concepts/configuration/overview/)

---

## 🆘 获取帮助

```bash
# 查看所有资源
kubectl get all -n bookcommunity

# 查看事件
kubectl get events -n bookcommunity --sort-by='.lastTimestamp'

# 查看日志
kubectl logs -f deployment/bookcommunity -n bookcommunity

# 使用 k9s 可视化管理
k9s -n bookcommunity
```

---

**部署成功后，您可以访问：**
- 应用 API：http://<service-ip>:80/health
- Prometheus 监控：http://<service-ip>:80/metrics
- Swagger 文档：http://<service-ip>:80/swagger/index.html（如已配置）

**祝部署顺利！** 🚀
