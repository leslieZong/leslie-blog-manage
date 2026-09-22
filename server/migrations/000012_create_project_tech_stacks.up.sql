CREATE TABLE project_tech_stacks (
    project_id CHAR(26) NOT NULL,

    tech_stack_id CHAR(26) NOT NULL,

    sort INT NOT NULL DEFAULT 0,

    PRIMARY KEY (project_id, tech_stack_id),

    KEY idx_project_tech_stacks_project_id (project_id),

    KEY idx_project_tech_stacks_tech_stack_id (tech_stack_id),

    CONSTRAINT fk_project_tech_stacks_project
        FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_project_tech_stacks_tech_stack
        FOREIGN KEY (tech_stack_id)
        REFERENCES tech_stacks(id)
        ON DELETE CASCADE
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci
COMMENT='项目技术栈表';
