ALTER TABLE posts
ADD COLUMN featured TINYINT(1) NOT NULL DEFAULT 0
AFTER status;

CREATE INDEX idx_posts_featured
ON posts(featured);