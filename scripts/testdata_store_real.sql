-- 更新商品为永州县城风格（保留ID 1-4的烟酒基础，替换其他）
UPDATE products SET name='白沙(精品)',price=100,stock=50,unit='条',base_unit='包',conversion_rate=10 WHERE id=1;
UPDATE products SET name='芙蓉王(硬)',price=250,stock=40,unit='条',base_unit='包',conversion_rate=10 WHERE id=2;
UPDATE products SET name='椰岛椰汁1L',price=60,stock=80,unit='箱',base_unit='瓶',conversion_rate=12 WHERE id=3;
UPDATE products SET name='旺旺大礼包',price=300,stock=30,unit='箱',base_unit='包',conversion_rate=20 WHERE id=4;

-- 清理旧的额外商品和供应商
DELETE FROM order_products;
DELETE FROM orders;
DELETE FROM purchase_items;
DELETE FROM purchases;
DELETE FROM products WHERE id > 4;
DELETE FROM suppliers WHERE id > 2;

-- 永州县城风格的商品
INSERT INTO products (id,category_id,name,price,stock,unit,base_unit,conversion_rate) VALUES
(5,1,'双喜(软经典)',80,60,'条','包',10),
(6,1,'黄果树(长征)',50,100,'条','包',10),
(7,2,'永州凉茶500ml',30,120,'箱','瓶',15),
(8,2,'农夫矿泉水550ml',24,200,'箱','瓶',24),
(9,3,'辣条大礼包',60,100,'箱','包',30),
(10,3,'康师傅方便面',36,150,'箱','桶',12),
(11,3,'洽洽瓜子',80,60,'箱','包',20),
(12,4,'雕牌洗洁精',50,40,'箱','瓶',6),
(13,4,'立白洗衣粉',90,35,'箱','袋',10);

-- 重置商品序列
SELECT setval('products_id_seq', (SELECT max(id) FROM products));

-- 永州本地供应商
INSERT INTO suppliers (id,name,contact_person,phone,address) VALUES
(3,'永州副食批发部','唐老板','0746-8412345','永州市冷水滩区零陵路'),
(4,'零陵烟酒商行','张经理','0746-6223456','永州市零陵区潇水路'),
(5,'东安百货批发','刘总','0746-4213456','永州市东安县建设中路');
SELECT setval('suppliers_id_seq', (SELECT max(id) FROM suppliers));

-- 4月-7月销售订单（真实风格）
INSERT INTO orders (user_id,payment_id,customer_name,total_price,total_paid,total_return,created_at,updated_at) VALUES
(1,1,'王老板',1000,1000,0,'2026-04-05 09:30:00+08','2026-04-05 09:30:00+08'),
(1,2,'张翠花',250,250,0,'2026-04-10 15:20:00+08','2026-04-10 15:20:00+08'),
(1,3,'李建国',140,140,0,'2026-04-18 10:15:00+08','2026-04-18 10:15:00+08'),
(1,1,'陈师傅',200,200,0,'2026-04-25 08:00:00+08','2026-04-25 08:00:00+08'),
(1,2,'周老板',500,500,0,'2026-05-08 11:30:00+08','2026-05-08 11:30:00+08'),
(1,3,'刘奶奶',36,36,0,'2026-05-15 16:45:00+08','2026-05-15 16:45:00+08'),
(1,1,'王老师',250,250,0,'2026-05-22 10:00:00+08','2026-05-22 10:00:00+08'),
(1,2,'赵师傅',300,300,0,'2026-06-02 09:00:00+08','2026-06-02 09:00:00+08'),
(1,4,'唐老板',1000,1000,0,'2026-06-10 14:00:00+08','2026-06-10 14:00:00+08'),
(1,1,'散客',36,36,0,'2026-06-18 11:30:00+08','2026-06-18 11:30:00+08'),
(1,2,'杨老板',250,250,0,'2026-06-25 10:00:00+08','2026-06-25 10:00:00+08'),
(1,3,'肖师傅',80,80,0,'2026-07-02 16:00:00+08','2026-07-02 16:00:00+08');

INSERT INTO order_products (order_id,product_id,quantity,total_price,created_at,updated_at) VALUES
(2,2,1,250,'2026-04-10 15:20:00+08','2026-04-10 15:20:00+08'),
(3,1,1,100,'2026-04-18 10:15:00+08','2026-04-18 10:15:00+08'),
(3,5,1,40,'2026-04-18 10:15:00+08','2026-04-18 10:15:00+08'),
(4,2,1,200,'2026-04-25 08:00:00+08','2026-04-25 08:00:00+08'),
(5,2,2,500,'2026-05-08 11:30:00+08','2026-05-08 11:30:00+08'),
(6,8,24,36,'2026-05-15 16:45:00+08','2026-05-15 16:45:00+08'),
(7,2,1,250,'2026-05-22 10:00:00+08','2026-05-22 10:00:00+08'),
(8,6,6,300,'2026-06-02 09:00:00+08','2026-06-02 09:00:00+08'),
(9,2,4,1000,'2026-06-10 14:00:00+08','2026-06-10 14:00:00+08'),
(10,8,24,36,'2026-06-18 11:30:00+08','2026-06-18 11:30:00+08'),
(11,2,1,250,'2026-06-25 10:00:00+08','2026-06-25 10:00:00+08'),
(12,5,2,80,'2026-07-02 16:00:00+08','2026-07-02 16:00:00+08');

-- 注意：orders表使用序列，需要知道实际ID
-- 先插入orders，再查询实际ID，再插入order_products
-- 但上面的INSERT已经用固定ID插入了（因为我们重建了数据）

-- 库存调整
UPDATE products SET stock=stock-1 WHERE id=2;
UPDATE products SET stock=stock-1 WHERE id=1;
UPDATE products SET stock=stock-4 WHERE id=5;
UPDATE products SET stock=stock-1 WHERE id=2;
UPDATE products SET stock=stock-2 WHERE id=2;
UPDATE products SET stock=stock-24 WHERE id=8;
UPDATE products SET stock=stock-1 WHERE id=2;
UPDATE products SET stock=stock-6 WHERE id=6;
UPDATE products SET stock=stock-4 WHERE id=2;
UPDATE products SET stock=stock-24 WHERE id=8;
UPDATE products SET stock=stock-2 WHERE id=5;
