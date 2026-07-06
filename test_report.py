import requests

BASE_URL = "https://store.yuppy576.top/v1"

resp = requests.post(f"{BASE_URL}/users/login", json={"email": "admin@test.com", "password": "test123456"})
token = resp.json()["data"]["token"]
headers = {"Authorization": f"Bearer {token}"}

print("16. 销售统计测试")
resp = requests.get(f"{BASE_URL}/reports/sales/stats", headers=headers, params={"start_date": "2026-01-01", "end_date": "2026-12-31"})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    stats = data.get("data", {})
    print(f"   订单数: {stats.get('total_orders', 0)}")
    print(f"   销售额: {stats.get('total_revenue', 0)}")
    print("   ✅ 销售统计接口正常")
else:
    print(f"   ❌ 销售统计接口异常: {resp.text}")

print("\n17. 每日销售测试")
resp = requests.get(f"{BASE_URL}/reports/sales/daily", headers=headers, params={"start_date": "2026-01-01", "end_date": "2026-12-31"})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    sales = data.get("data", [])
    print(f"   数据点数: {len(sales)}")
    print("   ✅ 每日销售接口正常")
else:
    print(f"   ❌ 每日销售接口异常: {resp.text}")

print("\n18. 库存预警测试")
resp = requests.get(f"{BASE_URL}/products/low-stock", headers=headers, params={"threshold": "10"})
print(f"   状态码: {resp.status_code}")
if resp.status_code == 200:
    data = resp.json()
    products = data.get("data", {}).get("products", [])
    print(f"   低库存商品数量: {len(products)}")
    print("   ✅ 库存预警接口正常")
else:
    print(f"   ❌ 库存预警接口异常: {resp.text}")

print("\n=== 报表和库存预警功能测试完成 ===")
