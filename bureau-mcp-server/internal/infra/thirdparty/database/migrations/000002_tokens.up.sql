CREATE TABLE IF NOT EXISTS mcp_tokens (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash  TEXT        NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    is_revoked  BOOLEAN     NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_mcp_tokens_hash ON mcp_tokens (token_hash);