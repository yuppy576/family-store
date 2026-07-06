import requests

BASE_URL = "https://store.yuppy576.top/v1"

resp = requests.post(f"{BASE_URL}/users/login", json={"email": "admin@test.com", "password": "test123456"})
token = resp.json()["data"]["token"]
headers = {"Authorization": f"Bearer {token}"}

print("测试审计日志接口:")
resp = requests.get(f"{BASE_URL}/audit-logs", headers=headers, params={"page": "1", "page_size": "10"})
print(f"状态码: {resp.status_code}")
print(f"响应: {resp.text}")

print("\n创建分类触发审计日志:")
import uuid
resp = requests.post(f"{BASE_URL}/categories", headers=headers, json={"name": f"测试分类-{uuid.uuid4().hex[:8]}"})
print(f"状态码: {resp.status_code}")

print("\n再次查询审计日志:")
resp = requests.get(f"{BASE_URL}/audit-logs", headers=headers, params={"page": "1", "page_size": "10"})
print(f"状态码: {resp.status_code}")
print(f"响应: {resp.text}")

print("\n使用skip/limit参数查询:")
resp = requests.get(f"{BASE_URL}/audit-logs", headers=headers, params={"skip": "0", "limit": "10"})
print(f"状态码: {resp.status_code}")
print(f"响应: {resp.text}")
