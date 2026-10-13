CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL DEFAULT '',
    channel TEXT NOT NULL DEFAULT 'linkedin',
    status TEXT NOT NULL DEFAULT 'idea',
    sent_at TIMESTAMPTZ,
    next_follow_up TIMESTAMPTZ,
    contact TEXT NOT NULL DEFAULT '',
    location TEXT NOT NULL DEFAULT '',
    remote TEXT NOT NULL DEFAULT '',
    stack TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes TEXT NOT NULL DEFAULT ''
);