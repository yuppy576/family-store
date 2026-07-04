-- 供应商表
CREATE TABLE IF NOT EXISTS "suppliers" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(200) NOT NULL,
    "contact_person" varchar(100),
    "phone" varchar(20),
    "address" text,
    "memo" text,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_suppliers_name ON "suppliers" ("name");

-- 进货入库记录
CREATE TYPE purchase_status_enum AS ENUM (
    'PENDING',     -- 待入库
    'COMPLETED',   -- 已入库
    'CANCELLED'    -- 已取消
);

CREATE TABLE IF NOT EXISTS "purchases" (
    "id" bigserial PRIMARY KEY,
    "supplier_id" bigint NOT NULL REFERENCES "suppliers"("id"),
    "operator" varchar(50) NOT NULL,
    "total_amount" decimal(18,2) NOT NULL DEFAULT 0,
    "status" purchase_status_enum NOT NULL DEFAULT 'PENDING',
    "remark" text,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_purchases_supplier ON "purchases" ("supplier_id");
CREATE INDEX idx_purchases_date ON "purchases" ("created_at");

-- 进货明细（入库商品）
CREATE TABLE IF NOT EXISTS "purchase_items" (
    "id" bigserial PRIMARY KEY,
    "purchase_id" bigint NOT NULL REFERENCES "purchases"("id"),
    "product_id" bigint NOT NULL REFERENCES "products"("id"),
    "quantity" bigint NOT NULL,
    "unit_price" decimal(18,2) NOT NULL,
    "total_price" decimal(18,2) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_purchase_items_purchase ON "purchase_items" ("purchase_id");
