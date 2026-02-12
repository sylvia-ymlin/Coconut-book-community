#!/bin/bash

# BookCommunity Kubernetes 部署脚本

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}==>${NC} $1"
}

# 检查必要的工具
check_prerequisites() {
    print_step "检查必要工具..."
    
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl 未安装，请先安装 kubectl"
        exit 1
    fi
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    print_info "✅ 所有必要工具已就绪"
}

# 构建 Docker 镜像
build_image() {
    print_step "构建 Docker 镜像..."
    
    DOCKER_USERNAME=${DOCKER_USERNAME:-yourusername}
    IMAGE_TAG=${IMAGE_TAG:-latest}
    
    docker build -t ${DOCKER_USERNAME}/bookcommunity:${IMAGE_TAG} .
    
    print_info "✅ Docker 镜像构建完成: ${DOCKER_USERNAME}/bookcommunity:${IMAGE_TAG}"
    
    # 询问是否推送到 Docker Hub
    read -p "是否推送镜像到 Docker Hub? (y/n): " push_image
    if [ "$push_image" = "y" ]; then
        docker push ${DOCKER_USERNAME}/bookcommunity:${IMAGE_TAG}
        print_info "✅ 镜像已推送到 Docker Hub"
    fi
}

# 部署到 Kubernetes
deploy_to_k8s() {
    print_step "部署到 Kubernetes..."
    
    # 创建 namespace
    print_info "创建 namespace..."
    kubectl apply -f k8s/base/namespace.yaml
    
    # 部署 Secret
    print_info "部署 Secret..."
    kubectl apply -f k8s/base/secret.yaml
    
    # 部署 ConfigMap
    print_info "部署 ConfigMap..."
    kubectl apply -f k8s/base/configmap.yaml
    
    # 部署 PostgreSQL
    print_info "部署 PostgreSQL..."
    kubectl apply -f k8s/base/postgres.yaml
    
    # 部署 Redis
    print_info "部署 Redis..."
    kubectl apply -f k8s/base/redis.yaml
    
    # 部署 RabbitMQ
    print_info "部署 RabbitMQ..."
    kubectl apply -f k8s/base/rabbitmq.yaml
    
    # 等待数据库就绪
    print_info "等待 PostgreSQL 就绪..."
    kubectl wait --for=condition=ready pod -l app=postgres -n bookcommunity --timeout=300s
    
    print_info "等待 Redis 就绪..."
    kubectl wait --for=condition=ready pod -l app=redis -n bookcommunity --timeout=300s
    
    print_info "等待 RabbitMQ 就绪..."
    kubectl wait --for=condition=ready pod -l app=rabbitmq -n bookcommunity --timeout=300s
    
    # 部署应用
    print_info "部署 BookCommunity 应用..."
    kubectl apply -f k8s/base/deployment.yaml
    kubectl apply -f k8s/base/service.yaml
    kubectl apply -f k8s/base/hpa.yaml
    
    # 等待应用就绪
    print_info "等待应用就绪..."
    kubectl wait --for=condition=ready pod -l app=bookcommunity -n bookcommunity --timeout=300s
    
    print_info "✅ 部署完成！"
}

# 查看部署状态
check_status() {
    print_step "查看部署状态..."
    
    echo ""
    print_info "Pods 状态:"
    kubectl get pods -n bookcommunity
    
    echo ""
    print_info "Services 状态:"
    kubectl get svc -n bookcommunity
    
    echo ""
    print_info "HPA 状态:"
    kubectl get hpa -n bookcommunity
    
    echo ""
    print_info "Ingress 状态:"
    kubectl get ingress -n bookcommunity
}

# 获取访问地址
get_access_url() {
    print_step "获取访问地址..."
    
    # 检查是否使用 Minikube
    if kubectl config current-context | grep -q "minikube"; then
        print_info "检测到 Minikube 环境"
        echo ""
        echo "访问应用:"
        echo "  minikube service bookcommunity-service -n bookcommunity"
        echo ""
        echo "或者使用端口转发:"
        echo "  kubectl port-forward svc/bookcommunity-service 8080:80 -n bookcommunity"
        echo "  然后访问: http://localhost:8080"
    else
        SERVICE_IP=$(kubectl get svc bookcommunity-service -n bookcommunity -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
        if [ -z "$SERVICE_IP" ]; then
            print_warn "LoadBalancer IP 尚未分配，使用端口转发:"
            echo "  kubectl port-forward svc/bookcommunity-service 8080:80 -n bookcommunity"
        else
            print_info "访问地址: http://${SERVICE_IP}"
        fi
    fi
}

# 清理部署
cleanup() {
    print_step "清理部署..."
    
    read -p "确认删除所有资源? (yes/no): " confirm
    if [ "$confirm" != "yes" ]; then
        print_info "取消清理"
        exit 0
    fi
    
    kubectl delete namespace bookcommunity
    
    print_info "✅ 清理完成"
}

# 主菜单
show_menu() {
    echo ""
    echo "====================================="
    echo "  BookCommunity K8s 部署工具"
    echo "====================================="
    echo "1. 构建 Docker 镜像"
    echo "2. 部署到 Kubernetes"
    echo "3. 查看部署状态"
    echo "4. 获取访问地址"
    echo "5. 完整部署 (构建 + 部署)"
    echo "6. 清理部署"
    echo "7. 退出"
    echo "====================================="
    read -p "请选择操作 [1-7]: " choice
    
    case $choice in
        1)
            build_image
            ;;
        2)
            deploy_to_k8s
            check_status
            get_access_url
            ;;
        3)
            check_status
            ;;
        4)
            get_access_url
            ;;
        5)
            build_image
            deploy_to_k8s
            check_status
            get_access_url
            ;;
        6)
            cleanup
            ;;
        7)
            print_info "退出"
            exit 0
            ;;
        *)
            print_error "无效选择"
            show_menu
            ;;
    esac
}

# 主逻辑
main() {
    check_prerequisites
    
    if [ $# -eq 0 ]; then
        show_menu
    else
        case "$1" in
            build)
                build_image
                ;;
            deploy)
                deploy_to_k8s
                check_status
                get_access_url
                ;;
            status)
                check_status
                ;;
            url)
                get_access_url
                ;;
            clean)
                cleanup
                ;;
            full)
                build_image
                deploy_to_k8s
                check_status
                get_access_url
                ;;
            *)
                echo "用法: $0 {build|deploy|status|url|clean|full}"
                echo ""
                echo "命令:"
                echo "  build   - 构建 Docker 镜像"
                echo "  deploy  - 部署到 Kubernetes"
                echo "  status  - 查看部署状态"
                echo "  url     - 获取访问地址"
                echo "  clean   - 清理部署"
                echo "  full    - 完整部署"
                echo ""
                echo "或直接运行 $0 使用交互式菜单"
                exit 1
                ;;
        esac
    fi
}

main "$@"
