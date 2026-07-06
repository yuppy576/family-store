import requests

BASE_URL = "https://store.yuppy576.top/v1"

resp = requests.post(f"{BASE_URL}/users/login", json={"email": "admin@test.com", "password": "test123456"})
token = resp.json()["data"]["token"]
headers = {"Authorization": f"Bearer {token}"}

resp = requests.get(f"{BASE_URL}/products/low-stock", headers=headers, params={"threshold": "10"})
print("状态码:", resp.status_code)
print("完整响应:", resp.text)
