-- 创建后台用户表
CREATE TABLE users (

    -- 用户唯一 ID
    --
    -- 使用 CHAR(26) 是因为后面我们计划使用 ULID。
    --
    -- ULID 类似：
    -- 01JXXXXXXXXXXXXXXX
    --
    -- 比 UUID 更适合我们的项目排序需求。
    id CHAR(26) NOT NULL,

    -- 登录用户名
    username VARCHAR(50) NOT NULL,

    -- 用户邮箱
    email VARCHAR(255) DEFAULT NULL,

    -- 密码哈希
    --
    -- 注意：
    -- 这里绝对不能保存明文密码。
    password_hash VARCHAR(255) NOT NULL,

    -- 用户显示名称
    display_name VARCHAR(100) NOT NULL DEFAULT '',

    -- 用户头像地址
    avatar_url VARCHAR(500) NOT NULL DEFAULT '',

    -- 用户状态
    --
    -- 1 = 正常
    -- 0 = 禁用
    status TINYINT NOT NULL DEFAULT 1,

    -- 最后一次登录时间
    last_login_at DATETIME(3) DEFAULT NULL,

    -- 创建时间
    created_at DATETIME(3)
        NOT NULL
        DEFAULT CURRENT_TIMESTAMP(3),

    -- 更新时间
    updated_at DATETIME(3)
        NOT NULL
        DEFAULT CURRENT_TIMESTAMP(3)
        ON UPDATE CURRENT_TIMESTAMP(3),

    -- 软删除时间
    --
    -- NULL = 没有删除
    -- 有值 = 已经删除
    deleted_at DATETIME(3) DEFAULT NULL,

    -- 主键
    PRIMARY KEY (id),

    -- 用户名必须唯一
    UNIQUE KEY uk_users_username (username),

    -- 邮箱必须唯一
    UNIQUE KEY uk_users_email (email),

    -- 用户状态索引
    KEY idx_users_status (status),

    -- 删除时间索引
    KEY idx_users_deleted_at (deleted_at)

)
ENGINE = InnoDB

DEFAULT CHARSET = utf8mb4

COLLATE = utf8mb4_unicode_ci

COMMENT = '后台用户表';