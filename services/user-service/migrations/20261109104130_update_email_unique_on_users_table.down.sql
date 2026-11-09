DROP INDEX IF EXISTS idx_users_email_active;

ALTER TABLE users
ADD CONSTRAINT uq_users_email UNIQUE (email);