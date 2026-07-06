CREATE TABLE IF NOT EXISTS "subscriptions" (
    "id" bigserial PRIMARY KEY,
    "store_id" bigint NOT NULL REFERENCES "stores"("id") ON DELETE CASCADE,
    "plan" varchar(20) NOT NULL DEFAULT 'TRIAL',
    "status" varchar(20) NOT NULL DEFAULT 'TRIAL',
    "start_date" timestamptz NOT NULL DEFAULT now(),
    "end_date" timestamptz,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_subscriptions_store ON "subscriptions" ("store_id");
CREATE INDEX idx_subscriptions_status ON "subscriptions" ("status");
CREATE INDEX idx_subscriptions_end_date ON "subscriptions" ("end_date");

INSERT INTO "subscriptions" ("store_id", "plan", "status", "start_date") VALUES (1, 'TRIAL', 'ACTIVE', now());