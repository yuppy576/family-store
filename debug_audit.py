import requests

BASE_URL = "https://store.yuppy576.top/v1"

# 登录
print("登录...")
resp = requests.post(f"{BASE_URL}/users/login", json={"email": "admin@test.com", "password": "test123456"})
data = resp.json()
token = data["data"]["token"]
headers = {"Authorization": f"Bearer {token}"}

# 获取用户信息
print("\n获取用户信息...")
resp = requests.get(f"{BASE_URL}/users/me", headers=headers)
data = resp.json()
print(f"用户信息: {data}")

# 查询审计日志
print("\n查询审计日志...")
resp = requests.get(f"{BASE_URL}/audit-logs", headers=headers, params={"limit": 5})
print(f"状态码: {resp.status_code}")
print(f"响应: {resp.text}")

# 创建分类
print("\n创建分类...")
resp = requests.post(f"{BASE_URL}/categories", headers=headers, json={"name": "调试分类"})
print(f"状态码: {resp.status_code}")
print(f"响应: {resp.text}")

# 等待2秒后再次查询
import time
time.sleep(2)

print("\n再次查询审计日志...")
resp = requests.get(f"{BASE_URL}/audit-logs", headers=headers, params={"limit": 5})
print(f"状态码: {resp.status_code}")
print(f"响应: {resp.text}")