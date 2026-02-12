-- BookCommunity PostgreSQL 初始化脚本
-- 此脚本会在 PostgreSQL 容器首次启动时自动执行

-- 创建扩展（可选，用于增强功能）
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";      -- UUID 生成
CREATE EXTENSION IF NOT EXISTS "pg_trgm";        -- 三元组相似度搜索
CREATE EXTENSION IF NOT EXISTS "unaccent";       -- 去除重音符号

-- 设置默认搜索路径
SET search_path TO public;

-- 创建自定义类型（示例）
-- DO $$ BEGIN
--     CREATE TYPE user_role AS ENUM ('user', 'admin', 'moderator');
-- EXCEPTION
--     WHEN duplicate_object THEN null;
-- END $$;

-- 优化配置（可选）
-- ALTER SYSTEM SET shared_buffers = '256MB';
-- ALTER SYSTEM SET effective_cache_size = '1GB';
-- ALTER SYSTEM SET maintenance_work_mem = '64MB';
-- ALTER SYSTEM SET checkpoint_completion_target = 0.9;
-- ALTER SYSTEM SET wal_buffers = '16MB';
-- ALTER SYSTEM SET default_statistics_target = 100;
-- ALTER SYSTEM SET random_page_cost = 1.1;
-- ALTER SYSTEM SET effective_io_concurrency = 200;

-- 日志记录初始化完成
DO $$
BEGIN
    RAISE NOTICE 'BookCommunity database initialized successfully!';
END $$;
