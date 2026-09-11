CREATE TABLE posts (
    id CHAR(26) NOT NULL,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL,
    summary VARCHAR(500) DEFAULT NULL,
    content LONGTEXT NOT NULL,
    cover VARCHAR(500) DEFAULT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'draft',

    author_id CHAR(26) NOT NULL,

    published_at DATETIME DEFAULT NULL,

    view_count BIGINT UNSIGNED NOT NULL DEFAULT 0,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    deleted_at DATETIME DEFAULT NULL,

    PRIMARY KEY (id),

    UNIQUE KEY uk_posts_slug (slug),

    KEY idx_posts_author_id (author_id),

    KEY idx_posts_status (status),

    KEY idx_posts_published_at (published_at),

    KEY idx_posts_deleted_at (deleted_at),

    CONSTRAINT fk_posts_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
) 
ENGINE=InnoDB 
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci
COMMENT='文章表';