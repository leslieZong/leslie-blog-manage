CREATE TABLE projects (
    id CHAR(26) NOT NULL,

    name VARCHAR(200) NOT NULL,

    slug VARCHAR(200) NOT NULL,

    description VARCHAR(1000) DEFAULT NULL,

    cover VARCHAR(500) DEFAULT NULL,

    github_url VARCHAR(500) DEFAULT NULL,

    demo_url VARCHAR(500) DEFAULT NULL,

    featured TINYINT(1) NOT NULL DEFAULT 0,

    status TINYINT NOT NULL DEFAULT 1,

    sort INT NOT NULL DEFAULT 0,

    created_at DATETIME NOT NULL
        DEFAULT CURRENT_TIMESTAMP,

    updated_at DATETIME NOT NULL
        DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,

    deleted_at DATETIME DEFAULT NULL,

    PRIMARY KEY (id),

    UNIQUE KEY uk_projects_name (name),

    UNIQUE KEY uk_projects_slug (slug),

    KEY idx_projects_status (status),

    KEY idx_projects_featured (featured),

    KEY idx_projects_sort (sort),

    KEY idx_projects_deleted_at (deleted_at)

) ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci
COMMENT='项目表';
