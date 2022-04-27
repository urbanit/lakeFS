BEGIN;

CREATE TABLE IF NOT EXISTS expired_tokens
(
    token_id         text UNIQUE NOT NULL,
    token_expires_at timestamptz NOT NULL
);

COMMIT;
