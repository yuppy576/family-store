CREATE TABLE IF NOT EXISTS "stores" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(100) NOT NULL,
    "domain" varchar(255) UNIQUE,
    "status" varchar(20) NOT NULL DEFAULT 'ACTIVE',
    "trial_end" timestamptz,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_stores_domain ON "stores" ("domain");
CREATE INDEX IF NOT EXISTS idx_stores_status ON "stores" ("status");

INSERT INTO "stores" ("id", "name", "domain", "status") VALUES (1, '默认租户', 'store.yuppy576.top', 'ACTIVE') ON CONFLICT DO NOTHING;

ALTER TABLE IF EXISTS "users" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "payments" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "orders" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "order_products" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "categories" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "products" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "consignors" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "consignments" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "vehicles" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "settlements" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "transfer_progress" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "suppliers" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "purchases" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "purchase_items" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");
ALTER TABLE IF EXISTS "audit_logs" ADD COLUMN IF NOT EXISTS "store_id" bigint DEFAULT 1 REFERENCES "stores"("id");

CREATE INDEX IF NOT EXISTS idx_users_store ON "users" ("store_id");
CREATE INDEX IF NOT EXISTS idx_payments_store ON "payments" ("store_id");
CREATE INDEX IF NOT EXISTS idx_orders_store ON "orders" ("store_id");
CREATE INDEX IF NOT EXISTS idx_categories_store ON "categories" ("store_id");
CREATE INDEX IF NOT EXISTS idx_products_store ON "products" ("store_id");
CREATE INDEX IF NOT EXISTS idx_consignors_store ON "consignors" ("store_id");
CREATE INDEX IF NOT EXISTS idx_consignments_store ON "consignments" ("store_id");
CREATE INDEX IF NOT EXISTS idx_suppliers_store ON "suppliers" ("store_id");
CREATE INDEX IF NOT EXISTS idx_purchases_store ON "purchases" ("store_id");
CREATE INDEX IF NOT EXISTS idx_audit_logs_store ON "audit_logs" ("store_id");
