-- 回滚 users 表。
--
-- 当执行：
--
-- migrate ... down
--
-- 就会执行这里。
DROP TABLE IF EXISTS users;