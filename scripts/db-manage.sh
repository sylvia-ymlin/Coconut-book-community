#!/bin/bash

# BookCommunity 数据库管理脚本
# 用于快速启动、停止、重置 PostgreSQL 开发环境

set -e

COMPOSE_FILE="docker-compose.dev.yaml"
DB_CONTAINER="bookcommunity-postgres"
DB_NAME="bookcommunity"
DB_USER="bookcommunity"
DB_PASSWORD="dev_password_2024"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

# 启动数据库
start_db() {
    print_info "启动 PostgreSQL 开发环境..."
    docker-compose -f "$COMPOSE_FILE" up -d postgres

    print_info "等待 PostgreSQL 就绪..."
    sleep 5

    # 检查健康状态
    if docker exec "$DB_CONTAINER" pg_isready -U "$DB_USER" > /dev/null 2>&1; then
        print_info "✅ PostgreSQL 已就绪！"
        print_info "连接信息:"
        echo "  Host: localhost"
        echo "  Port: 5432"
        echo "  Database: $DB_NAME"
        echo "  Username: $DB_USER"
        echo "  Password: $DB_PASSWORD"
    else
        print_error "PostgreSQL 启动失败"
        exit 1
    fi
}

# 停止数据库
stop_db() {
    print_info "停止 PostgreSQL..."
    docker-compose -f "$COMPOSE_FILE" stop postgres
    print_info "✅ PostgreSQL 已停止"
}

# 重启数据库
restart_db() {
    print_info "重启 PostgreSQL..."
    docker-compose -f "$COMPOSE_FILE" restart postgres
    sleep 3
    print_info "✅ PostgreSQL 已重启"
}

# 查看日志
logs_db() {
    print_info "查看 PostgreSQL 日志 (Ctrl+C 退出)..."
    docker-compose -f "$COMPOSE_FILE" logs -f postgres
}

# 进入数据库 shell
shell_db() {
    print_info "连接到 PostgreSQL (输入 \\q 退出)..."
    docker exec -it "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME"
}

# 重置数据库（删除所有数据）
reset_db() {
    print_warn "⚠️  警告：此操作将删除所有数据！"
    read -p "确认重置数据库? (yes/no): " confirm

    if [ "$confirm" != "yes" ]; then
        print_info "取消操作"
        exit 0
    fi

    print_info "停止并删除容器..."
    docker-compose -f "$COMPOSE_FILE" down -v

    print_info "重新启动数据库..."
    start_db

    print_info "✅ 数据库已重置"
}

# 备份数据库
backup_db() {
    BACKUP_DIR="./backups"
    mkdir -p "$BACKUP_DIR"

    BACKUP_FILE="$BACKUP_DIR/bookcommunity_$(date +%Y%m%d_%H%M%S).sql"

    print_info "备份数据库到: $BACKUP_FILE"
    docker exec "$DB_CONTAINER" pg_dump -U "$DB_USER" -d "$DB_NAME" > "$BACKUP_FILE"

    print_info "✅ 备份完成: $BACKUP_FILE"
}

# 恢复数据库
restore_db() {
    if [ -z "$1" ]; then
        print_error "请指定备份文件: $0 restore <backup_file>"
        exit 1
    fi

    BACKUP_FILE="$1"

    if [ ! -f "$BACKUP_FILE" ]; then
        print_error "备份文件不存在: $BACKUP_FILE"
        exit 1
    fi

    print_warn "⚠️  警告：此操作将覆盖当前数据！"
    read -p "确认恢复数据库? (yes/no): " confirm

    if [ "$confirm" != "yes" ]; then
        print_info "取消操作"
        exit 0
    fi

    print_info "恢复数据库从: $BACKUP_FILE"
    docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" < "$BACKUP_FILE"

    print_info "✅ 恢复完成"
}

# 启动 pgAdmin
start_pgadmin() {
    print_info "启动 pgAdmin..."
    docker-compose -f "$COMPOSE_FILE" up -d pgadmin

    print_info "✅ pgAdmin 已启动！"
    print_info "访问地址: http://localhost:5050"
    print_info "登录信息:"
    echo "  Email: admin@bookcommunity.local"
    echo "  Password: admin"
    echo ""
    print_info "添加服务器连接:"
    echo "  Host: postgres (容器内) 或 localhost (容器外)"
    echo "  Port: 5432"
    echo "  Database: $DB_NAME"
    echo "  Username: $DB_USER"
    echo "  Password: $DB_PASSWORD"
}

# 帮助信息
show_help() {
    cat << EOF
BookCommunity 数据库管理脚本

用法: $0 [command]

命令:
  start       启动 PostgreSQL 数据库
  stop        停止 PostgreSQL 数据库
  restart     重启 PostgreSQL 数据库
  logs        查看 PostgreSQL 日志
  shell       进入 PostgreSQL 命令行
  reset       重置数据库 (删除所有数据)
  backup      备份数据库
  restore     恢复数据库
  pgadmin     启动 pgAdmin Web 管理工具
  help        显示此帮助信息

示例:
  $0 start                    # 启动数据库
  $0 shell                    # 进入数据库 shell
  $0 backup                   # 备份数据库
  $0 restore backups/backup.sql  # 恢复数据库

EOF
}

# 主逻辑
case "${1:-}" in
    start)
        start_db
        ;;
    stop)
        stop_db
        ;;
    restart)
        restart_db
        ;;
    logs)
        logs_db
        ;;
    shell)
        shell_db
        ;;
    reset)
        reset_db
        ;;
    backup)
        backup_db
        ;;
    restore)
        restore_db "$2"
        ;;
    pgadmin)
        start_pgadmin
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "未知命令: ${1:-}"
        echo ""
        show_help
        exit 1
        ;;
esac
