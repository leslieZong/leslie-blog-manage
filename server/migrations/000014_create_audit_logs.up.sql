CREATE TABLE audit_logs (
    id CHAR(26) NOT NULL,

    user_id CHAR(26) DEFAULT NULL,

    action VARCHAR(100) NOT NULL,

    resource VARCHAR(100) NOT NULL,

    resource_id CHAR(26) DEFAULT NULL,

    request_id VARCHAR(100) DEFAULT NULL,

    ip VARCHAR(64) DEFAULT NULL,

    user_agent VARCHAR(500) DEFAULT NULL,

    result VARCHAR(20) NOT NULL DEFAULT 'success',

    error_message VARCHAR(1000) DEFAULT NULL,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),

    KEY idx_audit_logs_user_id (user_id),

    KEY idx_audit_logs_action (action),

    KEY idx_audit_logs_resource (resource),

    KEY idx_audit_logs_resource_id (resource_id),

    KEY idx_audit_logs_request_id (request_id),

    KEY idx_audit_logs_created_at (created_at),

    CONSTRAINT fk_audit_logs_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL

) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='审计日志表';