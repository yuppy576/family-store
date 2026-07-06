import requests

BASE_URL = "https://store.yuppy576.top/v1"

# 登录
print("登录...")
resp = requests.post(f"{BASE_URL}/users/login", json={"email": "admin@test.com", "password": "test123456"})
print(f"登录状态码: {resp.status_code}")
data = resp.json()
token = data["data"]["token"]
headers = {"Authorization": f"Bearer {token}"}

print("\n测试支付方式列表接口...")
resp = requests.get(f"{BASE_URL}/payments?skip=0&limit=5", headers=headers)
print(f"状态码: {resp.status_code}")
print(f"响应: {resp.text}")

print("\n测试订单列表接口...")
resp = requests.get(f"{BASE_URL}/orders?skip=0&limit=5", headers=headers)
print(f"状态码: {resp.status_code}")
print(f"响应: {resp.text}")

print("\n测试分类列表接口...")
resp = requests.get(f"{BASE_URL}/categories?skip=0&limit=5", headers=headers)
print(f"状态码: {resp.status_code}")
print(f"响应: {resp.text}")