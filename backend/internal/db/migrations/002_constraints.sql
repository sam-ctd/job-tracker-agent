ALTER TABLE applications
    ADD CONSTRAINT applications_status_check
    CHECK (status IN ('idea', 'enviada', 'relance', 'entretien', 'offre', 'refus'));

ALTER TABLE applications
    ADD CONSTRAINT applications_channel_check
    CHECK (channel IN ('linkedin', 'wttj', 'indeed', 'hellowork', 'apec', 'jobteaser', 'lesjeudis', 'chooseyourboss', 'alumni_estia'));

CREATE INDEX IF NOT EXISTS idx_applications_status ON applications (status);
CREATE INDEX IF NOT EXISTS idx_applications_next_follow_up ON applications (next_follow_up);
CREATE INDEX IF NOT EXISTS idx_applications_company ON applications (company);