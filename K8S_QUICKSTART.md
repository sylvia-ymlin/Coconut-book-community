# Kubernetes 快速开始指南

## ⚡ 5分钟部署到 Kubernetes

### 前置准备

```bash
# 1. 安装 Minikube (macOS)
brew install minikube

# 2. 启动 Minikube
minikube start --cpus=4 --memory=8192

# 3. 验证集群
kubectl cluster-info
```

---

## 🚀 方式一：一键部署（推荐）

```bash
cd /Users/ymlin/Downloads/003-Study/137-Projects/BookCommunity

# 运行部署脚本
./scripts/k8s-deploy.sh full

# 或使用交互式菜单
./scripts/k8s-deploy.sh
```

**脚本会自动：**
1. ✅ 构建 Docker 镜像
2. ✅ 部署 PostgreSQL
3. ✅ 部署 Redis
4. ✅ 部署 RabbitMQ
5. ✅ 部署 BookCommunity 应用
6. ✅ 配置自动扩展（HPA）
7. ✅ 显示访问地址

---

## 📦 方式二：手动部署

### 1. 构建镜像（使用 Minikube Docker）

```bash
# 使用 Minikube 的 Docker 环境
eval $(minikube docker-env)

# 构建镜像
docker build -t bookcommunity:latest .

# 验证镜像
docker images | grep bookcommunity
```

### 2. 部署到 K8s

```bash
# 一键部署所有资源
kubectl apply -f k8s/base/

# 或逐步部署
kubectl apply -f k8s/base/namespace.yaml
kubectl apply -f k8s/base/secret.yaml
kubectl apply -f k8s/base/configmap.yaml
kubectl apply -f k8s/base/postgres.yaml
kubectl apply -f k8s/base/redis.yaml
kubectl apply -f k8s/base/rabbitmq.yaml
kubectl apply -f k8s/base/deployment.yaml
kubectl apply -f k8s/base/service.yaml
kubectl apply -f k8s/base/hpa.yaml
```

### 3. 等待 Pods 就绪

```bash
# 查看所有 Pods
kubectl get pods -n bookcommunity -w

# 等待所有 Pods Running
kubectl wait --for=condition=ready pod --all -n bookcommunity --timeout=300s
```

### 4. 访问应用

```bash
# 方式1：使用 Minikube service
minikube service bookcommunity-service -n bookcommunity

# 方式2：端口转发
kubectl port-forward svc/bookcommunity-service 8080:80 -n bookcommunity
# 然后访问: http://localhost:8080/health

# 测试API
curl http://localhost:8080/health
```

---

## 📊 常用命令

### 查看资源状态

```bash
# 查看所有资源
kubectl get all -n bookcommunity

# 查看 Pods
kubectl get pods -n bookcommunity

# 查看 Services
kubectl get svc -n bookcommunity

# 查看 HPA（自动扩展）
kubectl get hpa -n bookcommunity
```

### 查看日志

```bash
# 查看应用日志
kubectl logs -f deployment/bookcommunity -n bookcommunity

# 查看 PostgreSQL 日志
kubectl logs -f deployment/postgres -n bookcommunity

# 查看所有容器日志
kubectl logs -f -l app=bookcommunity -n bookcommunity --all-containers
```

### 进入容器调试

```bash
# 进入应用容器
kubectl exec -it deployment/bookcommunity -n bookcommunity -- /bin/sh

# 进入 PostgreSQL
kubectl exec -it deployment/postgres -n bookcommunity -- psql -U bookcommunity

# 进入 Redis
kubectl exec -it deployment/redis -n bookcommunity -- redis-cli
```

### 更新部署

```bash
# 更新镜像
kubectl set image deployment/bookcommunity bookcommunity=bookcommunity:v2 -n bookcommunity

# 滚动重启
kubectl rollout restart deployment/bookcommunity -n bookcommunity

# 查看滚动更新状态
kubectl rollout status deployment/bookcommunity -n bookcommunity

# 回滚
kubectl rollout undo deployment/bookcommunity -n bookcommunity
```

### 扩缩容

```bash
# 手动扩容到 5 个副本
kubectl scale deployment/bookcommunity --replicas=5 -n bookcommunity

# 查看 HPA 自动扩展
kubectl get hpa -n bookcommunity -w
```

---

## 🧪 测试部署

### 1. 健康检查

```bash
# 转发端口
kubectl port-forward svc/bookcommunity-service 8080:80 -n bookcommunity

# 测试健康检查
curl http://localhost:8080/health

# 预期输出:
# {"status":"healthy","service":"BookCommunity API"}
```

### 2. 注册用户

```bash
curl -X POST http://localhost:8080/douyin/user/register/ \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 3. 获取图书推荐

```bash
# 先登录获取 token
TOKEN=$(curl -X POST http://localhost:8080/douyin/user/login/ \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' | jq -r '.token')

# 获取推荐
curl -X GET "http://localhost:8080/douyin/recommend?top_k=5" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🔧 故障排查

### Pod 一直 Pending

```bash
# 查看 Pod 详情
kubectl describe pod <pod-name> -n bookcommunity

# 查看节点资源
kubectl top nodes

# 解决方案：增加 Minikube 资源
minikube stop
minikube delete
minikube start --cpus=4 --memory=8192
```

### ImagePullBackOff

```bash
# 检查镜像
docker images | grep bookcommunity

# 确保使用 Minikube Docker
eval $(minikube docker-env)

# 重新构建
docker build -t bookcommunity:latest .

# 更新 deployment
kubectl set image deployment/bookcommunity bookcommunity=bookcommunity:latest -n bookcommunity
```

### CrashLoopBackOff

```bash
# 查看日志
kubectl logs <pod-name> -n bookcommunity --previous

# 查看事件
kubectl get events -n bookcommunity --sort-by='.lastTimestamp'

# 检查配置
kubectl get configmap bookcommunity-config -n bookcommunity -o yaml
kubectl get secret bookcommunity-secret -n bookcommunity -o yaml
```

---

## 🗑️ 清理资源

```bash
# 删除所有资源
kubectl delete namespace bookcommunity

# 或使用脚本
./scripts/k8s-deploy.sh clean

# 停止 Minikube
minikube stop

# 删除 Minikube 集群
minikube delete
```

---

## 📚 下一步

1. **配置 Ingress**
   - 安装 NGINX Ingress Controller
   - 配置域名访问

2. **启用 HTTPS**
   - 安装 cert-manager
   - 配置 Let's Encrypt

3. **配置监控**
   - 部署 Prometheus
   - 配置 Grafana 仪表盘

4. **持久化存储**
   - 配置 StorageClass
   - 使用云提供商存储

5. **CI/CD 集成**
   - GitHub Actions 自动部署
   - ArgoCD GitOps

---

## 🎯 生产环境检查清单

- [ ] 更新所有 Secret 密码
- [ ] 配置资源限制和请求
- [ ] 启用持久化存储
- [ ] 配置备份策略
- [ ] 启用监控和告警
- [ ] 配置日志聚合
- [ ] 设置网络策略
- [ ] 配置 HTTPS/TLS
- [ ] 实施滚动更新策略
- [ ] 配置健康检查和就绪探针

---

**完整文档：** `docs/KUBERNETES_DEPLOYMENT.md`

**祝部署顺利！** 🚀
