DROP INDEX idx_posts_featured
ON posts;

ALTER TABLE posts
DROP COLUMN featured;