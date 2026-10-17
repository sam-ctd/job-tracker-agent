ALTER TABLE applications
    ADD COLUMN IF NOT EXISTS priority TEXT NOT NULL DEFAULT 'media';

ALTER TABLE applications
    DROP CONSTRAINT IF EXISTS applications_priority_check;

ALTER TABLE applications
    ADD CONSTRAINT applications_priority_check
    CHECK (priority IN ('alta', 'media', 'baja'));
