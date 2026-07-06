ALTER TABLE "users" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "payments" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "orders" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "order_products" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "categories" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "products" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "consignors" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "consignments" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "vehicles" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "settlements" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "transfer_progress" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "suppliers" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "purchases" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "purchase_items" DROP COLUMN IF EXISTS "store_id";
ALTER TABLE "audit_logs" DROP COLUMN IF EXISTS "store_id";

DROP TABLE IF EXISTS "stores";