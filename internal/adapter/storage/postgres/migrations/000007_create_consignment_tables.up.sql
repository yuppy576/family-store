-- 寄卖人表
CREATE TABLE IF NOT EXISTS "consignors" (
    "id" bigserial PRIMARY KEY,
    "name" varchar(100) NOT NULL,
    "phone" varchar(20) NOT NULL,
    "id_card" varchar(18),
    "address" text,
    "memo" text,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_consignors_name ON "consignors" ("name");
CREATE INDEX idx_consignors_phone ON "consignors" ("phone");

-- 寄卖品主表
CREATE TYPE consignment_status_enum AS ENUM (
    'ON_SALE',       -- 在售
    'SOLD',          -- 已售出
    'EXPIRED',       -- 到期未售
    'RETURNED',      -- 已取回
    'CANCELLED'      -- 已取消
);

CREATE TABLE IF NOT EXISTS "consignments" (
    "id" bigserial PRIMARY KEY,
    "consignor_id" bigint NOT NULL REFERENCES "consignors"("id"),
    "name" varchar(200) NOT NULL,
    "description" text,
    "images" text[],                     -- 图片URL数组
    "category" varchar(50),             -- 分类：electronics/vehicle/furniture/etc
    "expected_price" decimal(18,2),     -- 寄卖人期望价
    "recommended_price" decimal(18,2),  -- 我方建议价
    "final_price" decimal(18,2),        -- 实际售价
    "commission_rate" decimal(5,2),     -- 佣金比例(%)，每件灵活设
    "commission_amount" decimal(18,2),  -- 佣金金额（自动算）
    "status" consignment_status_enum NOT NULL DEFAULT 'ON_SALE',
    "contract_end" date,                -- 寄卖到期日
    "is_vehicle" boolean NOT NULL DEFAULT false,
    "memo" text,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_consignments_consignor ON "consignments" ("consignor_id");
CREATE INDEX idx_consignments_status ON "consignments" ("status");
CREATE INDEX idx_consignments_enddate ON "consignments" ("contract_end");

-- 车辆寄卖专属信息
CREATE TABLE IF NOT EXISTS "consignment_vehicles" (
    "id" bigserial PRIMARY KEY,
    "consignment_id" bigint NOT NULL UNIQUE REFERENCES "consignments"("id"),
    "vin" varchar(17),                  -- 车架号
    "plate_number" varchar(10),         -- 车牌号
    "brand" varchar(50),               -- 品牌
    "model" varchar(100),              -- 型号
    "year" integer,                     -- 初次上牌年份
    "mileage" integer,                  -- 里程数(km)
    "displacement" varchar(20),        -- 排量
    "color" varchar(20),
    "inspection_expire" date,          -- 年检到期
    "insurance_expire" date,           -- 保险到期
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now()
);

-- 车辆过户进度
CREATE TYPE transfer_status_enum AS ENUM (
    'PENDING_INSPECTION',   -- 待验车
    'INSPECTED',            -- 验车完成
    'TRANSFERRING',         -- 过户办理中
    'TRANSFERRED',          -- 过户完成
    'SETTLED'               -- 已结算
);

CREATE TABLE IF NOT EXISTS "consignment_transfer_progress" (
    "id" bigserial PRIMARY KEY,
    "vehicle_id" bigint NOT NULL REFERENCES "consignment_vehicles"("id"),
    "status" transfer_status_enum NOT NULL DEFAULT 'PENDING_INSPECTION',
    "remark" text,                      -- 备注
    "attachment" text,                  -- 凭证图片
    "operator" varchar(50),            -- 经办人
    "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_transfer_vehicle ON "consignment_transfer_progress" ("vehicle_id");

-- 结算记录
CREATE TYPE settlement_type_enum AS ENUM (
    'SOLD_SETTLEMENT',      -- 卖出结算
    'RETURN_SETTLEMENT',    -- 到期取回结算
    'RENEWAL'               -- 续费
);

CREATE TABLE IF NOT EXISTS "consignment_settlements" (
    "id" bigserial PRIMARY KEY,
    "consignment_id" bigint NOT NULL REFERENCES "consignments"("id"),
    "type" settlement_type_enum NOT NULL,
    "sale_price" decimal(18,2),                              -- 成交价
    "commission_amount" decimal(18,2),                       -- 佣金
    "settlement_amount" decimal(18,2),                       -- 实际结算给寄卖人
    "renewal_fee" decimal(18,2),                             -- 续费金额
    "renewal_months" integer,                                -- 续费月数
    "new_end_date" date,                                     -- 续费后新到期日
    "remark" text,
    "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_settlements_consignment ON "consignment_settlements" ("consignment_id");
