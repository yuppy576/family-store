CREATE TABLE IF NOT EXISTS "audit_logs" (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint REFERENCES "users"("id"),
    "user_name" varchar(100),
    "action" varchar(50) NOT NULL,
    "resource_type" varchar(100) NOT NULL,
    "resource_id" bigint,
    "old_data" jsonb,
    "new_data" jsonb,
    "ip_address" varchar(50),
    "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_user ON "audit_logs" ("user_id");
CREATE INDEX idx_audit_logs_resource ON "audit_logs" ("resource_type", "resource_id");
CREATE INDEX idx_audit_logs_date ON "audit_logs" ("created_at");
CREATE INDEX idx_audit_logs_action ON "audit_logs" ("action");
