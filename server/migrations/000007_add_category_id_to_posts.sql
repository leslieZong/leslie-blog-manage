ALTER TABLE posts
ADD COLUMN category_id CHAR(26) NOT NULL
AFTER author_id;
CREATE INDEX idx_posts_category_id
ON posts(category_id);
ALTER TABLE posts
ADD CONSTRAINT fk_posts_category
FOREIGN KEY (category_id)
REFERENCES categories(id);