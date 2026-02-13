-- =====================================================
-- BookCommunity 数据库迁移脚本
-- 从视频社交平台迁移到图书社交平台
-- =====================================================

-- 注意：此脚本假设数据库为 PostgreSQL
-- 执行前请备份数据库！

BEGIN;

-- =====================================================
-- Step 1: 重命名主表
-- =====================================================

-- 将 videos_models 重命名为 book_reviews
ALTER TABLE IF EXISTS videos_models RENAME TO book_reviews;

COMMENT ON TABLE book_reviews IS '图书书评表：用户发布的书评（视频/文字/混合）';


-- =====================================================
-- Step 2: 为 book_reviews 添加新字段
-- =====================================================

-- 添加图书相关字段
ALTER TABLE book_reviews ADD COLUMN IF NOT EXISTS book_isbn VARCHAR(20);
ALTER TABLE book_reviews ADD COLUMN IF NOT EXISTS book_title VARCHAR(200);

-- 添加评分字段
ALTER TABLE book_reviews ADD COLUMN IF NOT EXISTS rating DECIMAL(3,1) DEFAULT 0.0;

-- 添加统计字段
ALTER TABLE book_reviews ADD COLUMN IF NOT EXISTS view_count INTEGER DEFAULT 0;
ALTER TABLE book_reviews ADD COLUMN IF NOT EXISTS collect_count INTEGER DEFAULT 0;

-- 添加文字内容字段（必填）
ALTER TABLE book_reviews ADD COLUMN IF NOT EXISTS content TEXT NOT NULL DEFAULT '';

-- 添加图片字段（JSON 数组，最多9张）
ALTER TABLE book_reviews ADD COLUMN IF NOT EXISTS images TEXT;

-- 添加标签字段
ALTER TABLE book_reviews ADD COLUMN IF NOT EXISTS tags VARCHAR(500);

-- 删除视频相关字段（如果存在）
ALTER TABLE book_reviews DROP COLUMN IF EXISTS url;
ALTER TABLE book_reviews DROP COLUMN IF EXISTS video_url;
ALTER TABLE book_reviews DROP COLUMN IF EXISTS review_type;

COMMENT ON COLUMN book_reviews.book_isbn IS '关联的图书 ISBN';
COMMENT ON COLUMN book_reviews.book_title IS '图书标题（冗余存储）';
COMMENT ON COLUMN book_reviews.rating IS '用户评分 (0.0-10.0)';
COMMENT ON COLUMN book_reviews.view_count IS '浏览次数';
COMMENT ON COLUMN book_reviews.collect_count IS '收藏次数';
COMMENT ON COLUMN book_reviews.content IS '书评文字内容（必填）';
COMMENT ON COLUMN book_reviews.images IS '图片URL列表（JSON数组，最多9张）';
COMMENT ON COLUMN book_reviews.tags IS '标签列表（JSON数组）';


-- =====================================================
-- Step 3: 更新评论表
-- =====================================================

-- 重命名字段：video_id → review_id
ALTER TABLE comment_models RENAME COLUMN video_id TO review_id;

COMMENT ON COLUMN comment_models.review_id IS '关联的书评ID（原 video_id）';


-- =====================================================
-- Step 4: 更新点赞表（多对多关系）
-- =====================================================

-- 重命名字段：video_id → review_id
ALTER TABLE user_like RENAME COLUMN video_id TO review_id;

COMMENT ON TABLE user_like IS '用户点赞书评关系表';
COMMENT ON COLUMN user_like.review_id IS '书评ID（原 video_id）';


-- =====================================================
-- Step 5: 更新收藏表（多对多关系）
-- =====================================================

-- 重命名字段：video_id → review_id
ALTER TABLE user_collection RENAME COLUMN video_id TO review_id;

COMMENT ON TABLE user_collection IS '用户收藏书评关系表';
COMMENT ON COLUMN user_collection.review_id IS '书评ID（原 video_id）';


-- =====================================================
-- Step 6: 添加索引以提升查询性能
-- =====================================================

-- book_reviews 表索引
CREATE INDEX IF NOT EXISTS idx_book_reviews_isbn ON book_reviews(book_isbn);
CREATE INDEX IF NOT EXISTS idx_book_reviews_author ON book_reviews(author_id);
CREATE INDEX IF NOT EXISTS idx_book_reviews_created ON book_reviews(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_book_reviews_rating ON book_reviews(rating DESC);
CREATE INDEX IF NOT EXISTS idx_book_reviews_view_count ON book_reviews(view_count DESC);
CREATE INDEX IF NOT EXISTS idx_book_reviews_like_count ON book_reviews(like_count DESC);

-- comment_models 表索引
CREATE INDEX IF NOT EXISTS idx_comments_review_id ON comment_models(review_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comment_models(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_created ON comment_models(created_at DESC);

-- user_like 表索引
CREATE INDEX IF NOT EXISTS idx_user_like_review ON user_like(review_id);
CREATE INDEX IF NOT EXISTS idx_user_like_user ON user_like(user_id);

-- user_collection 表索引
CREATE INDEX IF NOT EXISTS idx_user_collection_review ON user_collection(review_id);
CREATE INDEX IF NOT EXISTS idx_user_collection_user ON user_collection(user_id);


-- =====================================================
-- Step 7: 添加约束
-- =====================================================

-- book_reviews 表约束
ALTER TABLE book_reviews
  ADD CONSTRAINT chk_rating_range
  CHECK (rating >= 0.0 AND rating <= 10.0);

ALTER TABLE book_reviews
  ADD CONSTRAINT chk_content_not_empty
  CHECK (LENGTH(TRIM(content)) > 0);


-- =====================================================
-- Step 8: 数据清理和迁移（可选）
-- =====================================================

-- 设置默认值：没有 view_count 的记录设为 0
UPDATE book_reviews SET view_count = 0 WHERE view_count IS NULL;

-- 设置默认值：没有 collect_count 的记录设为 0
UPDATE book_reviews SET collect_count = 0 WHERE collect_count IS NULL;

-- 设置默认值：没有 rating 的记录设为 0.0
UPDATE book_reviews SET rating = 0.0 WHERE rating IS NULL;

-- 设置默认值：没有 content 的记录设为空字符串（必填字段）
UPDATE book_reviews SET content = '' WHERE content IS NULL;


-- =====================================================
-- Step 9: 验证数据迁移
-- =====================================================

-- 检查表是否存在
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'book_reviews') THEN
        RAISE EXCEPTION '表 book_reviews 不存在，迁移失败！';
    END IF;

    RAISE NOTICE '✓ 表 book_reviews 已创建';
END $$;

-- 检查字段是否存在
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'book_reviews' AND column_name = 'book_isbn'
    ) THEN
        RAISE EXCEPTION '字段 book_isbn 不存在，迁移失败！';
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'comment_models' AND column_name = 'review_id'
    ) THEN
        RAISE EXCEPTION '字段 review_id 不存在，迁移失败！';
    END IF;

    RAISE NOTICE '✓ 所有新字段已添加';
END $$;

-- 显示统计信息
SELECT
    'book_reviews' AS table_name,
    COUNT(*) AS total_records,
    COUNT(DISTINCT author_id) AS total_authors,
    SUM(like_count) AS total_likes,
    SUM(comment_count) AS total_comments
FROM book_reviews;

COMMIT;

-- =====================================================
-- 迁移完成提示
-- =====================================================

SELECT '
╔════════════════════════════════════════════════════╗
║  ✓ 数据库迁移完成！                                ║
║                                                    ║
║  已完成以下操作：                                  ║
║  1. videos_models → book_reviews                   ║
║  2. 添加图书相关字段（ISBN, Title, Rating）        ║
║  3. video_id → review_id（所有关联表）             ║
║  4. 添加索引和约束                                 ║
║                                                    ║
║  下一步：                                          ║
║  - 更新 Go 代码中的 Handler 和 Service            ║
║  - 更新 API 路由                                   ║
║  - 重新生成 Swagger 文档                           ║
╚════════════════════════════════════════════════════╝
' AS migration_status;
