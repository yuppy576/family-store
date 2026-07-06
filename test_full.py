import requests
import uuid
import time

BASE_URL = "https://store.yuppy576.top/v1"

print("=== 家族门店管理系统 - 完整功能测试 ===\n")

# 1. 登录测试
print("1. 登录测试")
resp = requests.post(f"{BASE_URL}/users/login", json={"email": "admin@test.com", "password": "test123456"})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    token = data["data"]["token"]
    headers = {"Authorization": f"Bearer {token}"}
    print("   ✅ 登录成功")
else:
    print(f"   ❌ 登录失败: {resp.text}")
    exit(1)

# 2. 用户信息测试
print("\n2. 用户信息测试")
resp = requests.get(f"{BASE_URL}/users/me", headers=headers)
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    print(f"   用户: {data['data']['email']}")
    print("   ✅ 用户信息正常")
else:
    print(f"   ❌ 用户信息获取失败")

# 3. 商品列表测试
print("\n3. 商品列表测试")
resp = requests.get(f"{BASE_URL}/products/", headers=headers, params={"limit": 5})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    products = data.get("data", {}).get("products", []) or []
    print(f"   商品数量: {len(products)}")
    print("   ✅ 商品列表接口正常")
else:
    print(f"   ❌ 商品列表接口异常")

# 4. 分类列表测试
print("\n4. 分类列表测试")
resp = requests.get(f"{BASE_URL}/categories/", headers=headers, params={"limit": 5})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    categories = data.get("data", {}).get("categories", []) or []
    print(f"   分类数量: {len(categories)}")
    print("   ✅ 分类列表接口正常")
else:
    print(f"   ❌ 分类列表接口异常")

# 5. 寄卖品列表测试
print("\n5. 寄卖品列表测试")
resp = requests.get(f"{BASE_URL}/consignment/items", headers=headers, params={"limit": 5})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    items = data.get("data", {}).get("items", []) or []
    print(f"   寄卖品数量: {len(items)}")
    print("   ✅ 寄卖品列表接口正常")
else:
    print(f"   ❌ 寄卖品列表接口异常")

# 6. 供应商列表测试
print("\n6. 供应商列表测试")
resp = requests.get(f"{BASE_URL}/suppliers", headers=headers, params={"limit": 5})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    suppliers = data.get("data", {}).get("suppliers", []) or []
    print(f"   供应商数量: {len(suppliers)}")
    print("   ✅ 供应商列表接口正常")
else:
    print(f"   ❌ 供应商列表接口异常")

# 7. 寄卖人列表测试
print("\n7. 寄卖人列表测试")
resp = requests.get(f"{BASE_URL}/consignment/consignors", headers=headers, params={"limit": 5})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    consignors = data.get("data", {}).get("consignors", []) or []
    print(f"   寄卖人数量: {len(consignors)}")
    print("   ✅ 寄卖人列表接口正常")
else:
    print(f"   ❌ 寄卖人列表接口异常")

# 8. 支付方式列表测试
print("\n8. 支付方式列表测试")
resp = requests.get(f"{BASE_URL}/payments?skip=0&limit=5", headers=headers)
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    payments = data.get("data", {}).get("payments", []) or []
    print(f"   支付方式数量: {len(payments)}")
    print("   ✅ 支付方式列表接口正常")
else:
    print(f"   ❌ 支付方式列表接口异常")

# 9. 订单列表测试
print("\n9. 订单列表测试")
resp = requests.get(f"{BASE_URL}/orders?skip=0&limit=5", headers=headers)
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    orders = data.get("data", {}).get("orders", []) or []
    print(f"   订单数量: {len(orders)}")
    print("   ✅ 订单列表接口正常")
else:
    print(f"   ❌ 订单列表接口异常")

# 10. 创建分类测试
print("\n10. 创建分类测试")
cat_name = f"测试分类_{uuid.uuid4().hex[:8]}"
resp = requests.post(f"{BASE_URL}/categories", headers=headers, json={"name": cat_name})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    print("   ✅ 创建分类正常")
else:
    print(f"   ❌ 创建分类失败: {resp.text}")

# 11. 审计日志测试
print("\n11. 审计日志测试")
time.sleep(2)
resp = requests.get(f"{BASE_URL}/audit-logs", headers=headers, params={"limit": 5})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    logs = data.get("data", {}).get("audit_logs", []) or []
    print(f"   日志数量: {len(logs)}")
    if logs:
        latest = logs[0]
        print(f"   最新日志: {latest.get('action')} - {latest.get('resource_type')}")
        print("   ✅ 审计日志记录正常")
    else:
        print("   ⚠️ 审计日志为空")
else:
    print(f"   ❌ 审计日志接口异常")

# 12. 订阅信息测试
print("\n12. 订阅信息测试")
resp = requests.get(f"{BASE_URL}/subscription", headers=headers)
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    sub = data["data"]
    print(f"   计划: {sub['plan']}")
    print(f"   状态: {sub['status']}")
    print(f"   到期日: {sub['end_date']}")
    print("   ✅ 订阅信息接口正常")
else:
    print(f"   ❌ 订阅信息接口异常: {resp.text}")

# 13. 注册API测试（公开接口）
print("\n13. 注册API测试")
test_email = f"test_register_{uuid.uuid4().hex[:6]}@example.com"
test_name = f"测试租户_{uuid.uuid4().hex[:4]}"
resp = requests.post(f"{BASE_URL}/stores/register", json={"name": test_name, "email": test_email, "password": "12345678"})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    print(f"   域名: {data.get('data', {}).get('domain', 'N/A')}")
    print("   ✅ 注册成功")
else:
    print(f"   ⚠️ 注册返回: {resp.text}")

# 14. 进货列表测试
print("\n14. 进货列表测试")
resp = requests.get(f"{BASE_URL}/purchases", headers=headers, params={"limit": 5})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    purchases = data.get("data", {}).get("purchases", []) or []
    print(f"   进货记录数量: {len(purchases)}")
    print("   ✅ 进货列表接口正常")
else:
    print(f"   ❌ 进货列表接口异常")

# 15. POS创建订单测试
print("\n15. POS创建订单测试")
resp = requests.post(f"{BASE_URL}/orders", headers=headers, json={
    "customer_name": "测试客户",
    "payment_id": 1,
    "total_paid": 100,
    "products": []
})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    print(f"   订单ID: {data['data'].get('id', 'N/A')}")
    print("   ✅ POS订单创建正常")
else:
    print(f"   ⚠️ POS订单创建: {resp.text}")

print("\n=== 所有功能测试完成 ===")