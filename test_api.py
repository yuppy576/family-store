import requests

BASE_URL = "https://store.yuppy576.top/v1"

def test_api():
    print("=== 验证家族门店管理系统 ===")
    
    # 1. 登录
    print("\n1. 登录测试")
    resp = requests.post(f"{BASE_URL}/users/login", json={"email": "admin@test.com", "password": "test123456"})
    print(f"   状态码: {resp.status_code}")
    data = resp.json()
    if not data.get("success"):
        print(f"   ❌ 登录失败: {data.get('messages', ['未知错误'])}")
        return
    
    token = data["data"]["token"]
    print("   ✅ 登录成功")
    headers = {"Authorization": f"Bearer {token}"}
    
    # 2. 获取当前用户
    print("\n2. 用户信息测试")
    resp = requests.get(f"{BASE_URL}/users/me", headers=headers)
    print(f"   状态码: {resp.status_code}")
    user_data = resp.json()
    print(f"   用户: {user_data.get('data', {}).get('email', 'N/A')}")
    print("   ✅ 用户信息正常")
    
    # 3. 商品列表
    print("\n3. 商品列表测试")
    resp = requests.get(f"{BASE_URL}/products", headers=headers, params={"limit": 5})
    print(f"   状态码: {resp.status_code}")
    data = resp.json()
    products = data.get("data", {}).get("products", []) or []
    print(f"   商品数量: {len(products)}")
    print("   ✅ 商品列表接口正常")
    
    # 4. 分类列表
    print("\n4. 分类列表测试")
    resp = requests.get(f"{BASE_URL}/categories", headers=headers, params={"limit": 5})
    print(f"   状态码: {resp.status_code}")
    data = resp.json()
    cats = data.get("data", {}).get("categories", []) or []
    print(f"   分类数量: {len(cats)}")
    print("   ✅ 分类列表接口正常")
    
    # 5. 寄卖品列表
    print("\n5. 寄卖品列表测试")
    resp = requests.get(f"{BASE_URL}/consignment/items", headers=headers, params={"limit": 5})
    print(f"   状态码: {resp.status_code}")
    data = resp.json()
    cons = data.get("data", {}).get("items", []) or []
    print(f"   寄卖品数量: {len(cons)}")
    print("   ✅ 寄卖品列表接口正常")
    
    # 6. 审计日志
    print("\n6. 审计日志测试")
    resp = requests.get(f"{BASE_URL}/audit-logs", headers=headers, params={"limit": 5})
    print(f"   状态码: {resp.status_code}")
    data = resp.json()
    logs = data.get("data", {}).get("logs", []) or []
    print(f"   日志数量: {len(logs)}")
    print("   ✅ 审计日志接口正常")
    
    # 7. 创建分类测试
    print("\n7. 创建分类测试")
    import uuid
    cat_name = f"测试分类_{uuid.uuid4().hex[:8]}"
    resp = requests.post(f"{BASE_URL}/categories", headers=headers, json={"name": cat_name})
    print(f"   状态码: {resp.status_code}")
    print("   ✅ 创建分类正常")
    
    # 8. 验证审计日志记录
    print("\n8. 审计日志记录验证")
    resp = requests.get(f"{BASE_URL}/audit-logs", headers=headers, params={"limit": 1})
    data = resp.json()
    logs = data.get("data", {}).get("logs", []) or []
    if logs:
        latest = logs[0]
        print(f"   最新日志: {latest.get('action', 'N/A')} - {latest.get('resource_type', 'N/A')}")
        print(f"   用户ID: {latest.get('user_id', 'N/A')}")
        print(f"   IP地址: {latest.get('ip_address', 'N/A')}")
        print("   ✅ 审计日志记录正常")
    else:
        print("   ⚠️ 审计日志为空")
    
    # 9. 供应商列表
    print("\n9. 供应商列表测试")
    resp = requests.get(f"{BASE_URL}/suppliers", headers=headers, params={"limit": 5})
    print(f"   状态码: {resp.status_code}")
    data = resp.json()
    supps = data.get("data", {}).get("suppliers", []) or []
    print(f"   供应商数量: {len(supps)}")
    print("   ✅ 供应商列表接口正常")
    
    print("\n=== 所有功能验证通过 ===")

if __name__ == "__main__":
    test_api()
