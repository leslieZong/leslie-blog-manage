CREATE TABLE tech_stacks (
    id CHAR(26) NOT NULL,

    name VARCHAR(100) NOT NULL,

    slug VARCHAR(100) NOT NULL,

    icon VARCHAR(500) DEFAULT NULL,

    description VARCHAR(255) DEFAULT NULL,

    official_url VARCHAR(500) DEFAULT NULL,

    category VARCHAR(50) DEFAULT NULL,

    sort INT NOT NULL DEFAULT 0,

    status TINYINT NOT NULL DEFAULT 1,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    deleted_at DATETIME DEFAULT NULL,

    PRIMARY KEY (id),

    UNIQUE KEY uk_tech_stacks_name (name),

    UNIQUE KEY uk_tech_stacks_slug (slug),

    KEY idx_tech_stacks_category (category),

    KEY idx_tech_stacks_sort (sort),

    KEY idx_tech_stacks_status (status),

    KEY idx_tech_stacks_deleted_at (deleted_at)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci
COMMENT='技术栈表';