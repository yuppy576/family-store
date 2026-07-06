import requests

BASE_URL = "https://store.yuppy576.top/v1"

# 登录
resp = requests.post(f"{BASE_URL}/users/login", json={"email": "admin@test.com", "password": "test123456"})
token = resp.json()["data"]["token"]
headers = {"Authorization": f"Bearer {token}"}

# 检查审计日志原始响应
resp = requests.get(f"{BASE_URL}/audit-logs", headers=headers, params={"skip": 0, "limit": 10})
print("审计日志响应:", resp.text)

# 创建一个新分类
import uuid
cat_name = f"测试分类_{uuid.uuid4().hex[:8]}"
resp = requests.post(f"{BASE_URL}/categories", headers=headers, json={"name": cat_name})
print(f"创建分类: {resp.status_code} - {resp.text}")

# 再次查询审计日志
import time
time.sleep(1)  # 等待异步写入
resp = requests.get(f"{BASE_URL}/audit-logs", headers=headers, params={"skip": 0, "limit": 5})
print("创建后审计日志:", resp.text)
