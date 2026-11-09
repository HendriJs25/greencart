ALTER TABLE users
DROP CONSTRAINT IF EXISTS uq_users_email;

CREATE UNIQUE INDEX idx_users_email_active ON users (email) WHERE deleted_at IS NULL;