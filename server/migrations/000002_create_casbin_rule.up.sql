CREATE TABLE casbin_rule (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    ptype VARCHAR(100) NOT NULL DEFAULT '',
    v0 VARCHAR(100) NOT NULL DEFAULT '',
    v1 VARCHAR(100) NOT NULL DEFAULT '',
    v2 VARCHAR(100) NOT NULL DEFAULT '',
    v3 VARCHAR(100) NOT NULL DEFAULT '',
    v4 VARCHAR(100) NOT NULL DEFAULT '',
    v5 VARCHAR(100) NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    KEY idx_casbin_rule_ptype (ptype),
    KEY idx_casbin_rule_v0 (v0),
    KEY idx_casbin_rule_v1 (v1)
)
ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_unicode_ci
COMMENT = 'Casbin 权限策略表';