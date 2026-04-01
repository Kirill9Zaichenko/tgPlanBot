CREATE TABLE IF NOT EXISTS organization_memberships (
                                                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                                                        organization_id INTEGER NOT NULL,
                                                        user_id INTEGER NOT NULL,
                                                        role TEXT NOT NULL DEFAULT 'member',
                                                        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

                                                        FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

CREATE UNIQUE INDEX IF NOT EXISTS idx_org_memberships_org_user
    ON organization_memberships(organization_id, user_id);

CREATE INDEX IF NOT EXISTS idx_org_memberships_user_id
    ON organization_memberships(user_id);

CREATE INDEX IF NOT EXISTS idx_org_memberships_org_id
    ON organization_memberships(organization_id);