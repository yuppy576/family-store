-- 清理旧数据
DELETE FROM consignment_settlements;
DELETE FROM consignment_transfer_progress;
DELETE FROM consignment_vehicles WHERE consignment_id NOT IN (2);
DELETE FROM consignments WHERE id > 2;

-- 重置寄卖人（保留前4个，删掉旧的测试数据）
DELETE FROM consignors WHERE id > 4;

-- 本地化寄卖人
INSERT INTO consignors (id,name,phone,id_card,address) VALUES
(5,'陈建国','15274658888','432924198505061234','永州市冷水滩区凤凰路'),
(6,'刘翠英','15173659999','432924197803152345','永州市零陵区中山路'),
(7,'周大军','18874651234','432924199012013456','永州市东安县白牙市镇'),
(8,'唐秀英','15874668888','432924196509102345','永州市祁阳县浯溪镇'),
(9,'李满华','13774667890','432924198812053456','永州市双牌县紫金路'),
(10,'王春生','15574651234','432924199503201234','永州市道县濂溪路');

-- 真实寄卖品（永州县城风格）
INSERT INTO consignments (id,consignor_id,name,category,expected_price,commission_rate,status,contract_end,is_vehicle,memo,created_at) VALUES
(3,5,'九成新雅迪电动车','电动车',1800,10,'SOLD','2026-05-20',false,'买菜接娃神器，已售','2026-04-08'),
(4,5,'华为Mate 60 256G','数码',3200,10,'SOLD','2026-05-15',false,'国行蓝色，已售','2026-04-12'),
(5,6,'格力1.5匹变频空调','家电',2200,10,'SOLD','2026-05-25',false,'用了两年，制冷好，已售','2026-04-18'),
(6,7,'五羊本田摩托车','车辆',4500,8,'SOLD','2026-06-01',true,'125cc，手续齐全，已售','2026-04-22'),
(7,8,'老凤祥黄金项链','珠宝',3600,5,'SOLD','2026-05-30',false,'24K约8克，带发票，已售','2026-04-25'),
(8,9,'海尔双门冰箱','家电',1200,10,'SOLD','2026-06-10',false,'用了三年，制冷正常，已售','2026-05-02'),
(9,5,'iPhone 14 128G','数码',2800,10,'SOLD','2026-06-20',false,'国行紫色，已售','2026-05-08'),
(10,6,'联想笔记本电脑','数码',1500,10,'EXPIRED','2026-06-15',false,'i5-8代，到期未售已取回','2026-05-15'),
(11,10,'立马电动车','电动车',1600,10,'ON_SALE','2026-08-10',false,'60V20A，九成新','2026-05-20'),
(12,7,'三轮摩托车','车辆',3200,8,'ON_SALE','2026-08-20',true,'载重1吨，柴油机','2026-05-25'),
(13,9,'红米K70 12+256G','数码',1200,10,'ON_SALE','2026-08-25',false,'用了半年，性价比高','2026-06-01'),
(14,5,'全友布艺沙发','家具',800,10,'ON_SALE','2026-09-01',false,'七成新，自提优先','2026-06-05'),
(15,6,'美的大1.5匹空调','家电',1800,10,'ON_SALE','2026-09-05',false,'2024年买，搬家出售','2026-06-10'),
(16,10,'宗申三轮摩托车','车辆',5800,5,'ON_SALE','2026-09-15',true,'带雨棚，适合拉货','2026-06-15'),
(17,7,'索尼75寸电视','家电',2800,8,'ON_SALE','2026-09-20',false,'2023年款，画质好','2026-06-20'),
(18,8,'本田190cc摩托车','车辆',6800,5,'ON_SALE','2026-10-01',true,'2022年上牌，1.2万公里','2026-06-25'),
(19,9,'vivo X100 Pro','数码',2600,10,'ON_SALE','2026-10-05',false,'12+256G，在保','2026-07-01'),
(20,10,'台铃电动车','电动车',1400,10,'ON_SALE','2026-10-10',false,'48V20A，代步利器','2026-07-03');

-- 车辆信息
INSERT INTO consignment_vehicles (id,consignment_id,vin,plate_number,brand,model,year,mileage,color) VALUES
(4,6,'LBBPEK5C9AB123456','湘M4X567','五羊本田','WH125-7',2020,15000,'黑色'),
(5,12,'XYZPBK5C9AB789012','湘M8Y890','宗申','ZS200ZH',2021,8000,'红色'),
(6,16,'JLEPBK5C9AB345678','湘M6W234','本田','CBF190X',2022,12000,'红色'),
(7,18,'JLEPBK5C9AB901234','湘M7X345','本田','CBF190X',2020,12000,'黑色');

-- 过户进度
INSERT INTO consignment_transfer_progress (vehicle_id,status,remark,operator,created_at) VALUES
(4,'PENDING_INSPECTION','五羊本田到店，等待验车','管理员','2026-04-23'),
(4,'INSPECTED','验车通过，手续齐全','管理员','2026-04-25'),
(4,'TRANSFERRING','买家已付款，办理过户','管理员','2026-05-28'),
(4,'TRANSFERRED','过户完成，已交车','管理员','2026-06-01');

-- 已售寄卖品的结算
INSERT INTO consignment_settlements (consignment_id,type,sale_price,commission_amount,settlement_amount,remark,created_at) VALUES
(3,'SOLD_SETTLEMENT',1800,180,1620,'已结算','2026-05-10'),
(4,'SOLD_SETTLEMENT',3200,320,2880,'已结算','2026-05-12'),
(5,'SOLD_SETTLEMENT',2200,220,1980,'已结算','2026-05-20'),
(6,'SOLD_SETTLEMENT',4500,360,4140,'已结算','2026-05-28'),
(7,'SOLD_SETTLEMENT',3600,180,3420,'已结算','2026-05-28'),
(8,'SOLD_SETTLEMENT',1200,120,1080,'已结算','2026-06-08'),
(9,'SOLD_SETTLEMENT',2800,280,2520,'已结算','2026-06-15');

-- 更新序列
SELECT setval('consignors_id_seq', (SELECT max(id) FROM consignors));
SELECT setval('consignments_id_seq', (SELECT max(id) FROM consignments));
SELECT setval('consignment_vehicles_id_seq', (SELECT max(id) FROM consignment_vehicles));
