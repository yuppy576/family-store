import requests

BASE_URL = "https://store.yuppy576.top/v1"

# 登录
print("1. 登录测试...")
resp = requests.post(f"{BASE_URL}/users/login", json={"email": "admin@test.com", "password": "test123456"})
data = resp.json()
token = data["data"]["token"]
headers = {"Authorization": f"Bearer {token}"}
print(f"   状态码: {resp.status_code}")
print("   ✅ 登录成功")

# 获取订阅信息
print("\n2. 获取订阅信息...")
resp = requests.get(f"{BASE_URL}/subscription", headers=headers)
print(f"   状态码: {resp.status_code}")
print(f"   响应: {resp.text}")

# 续费订阅
print("\n3. 续费订阅（个人版，12个月）...")
resp = requests.post(f"{BASE_URL}/subscription/renew", headers=headers, json={"plan": "PERSONAL", "months": 12})
print(f"   状态码: {resp.status_code}")
print(f"   响应: {resp.text}")

# 再次获取订阅信息
print("\n4. 再次获取订阅信息...")
resp = requests.get(f"{BASE_URL}/subscription", headers=headers)
print(f"   状态码: {resp.status_code}")
print(f"   响应: {resp.text}")

print("\n=== 订阅管理功能测试完成 ===")