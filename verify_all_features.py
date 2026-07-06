import urllib.request
import json
import sys

base_url = "http://127.0.0.1:8081"

def test_login():
    print("=== 1. 测试登录接口 ===")
    try:
        r = urllib.request.urlopen(urllib.request.Request(
            f"{base_url}/v1/users/login",
            data=json.dumps({'email': 'admin@test.com', 'password': 'test123456'}).encode(),
            headers={'Content-Type': 'application/json'}
        ))
        result = json.loads(r.read())
        if result['success']:
            print("✅ 登录成功")
            return result['data']['token']
        else:
            print(f"❌ 登录失败: {result}")
            return None
    except Exception as e:
        print(f"❌ 登录异常: {e}")
        return None

def test_me(token):
    print("\n=== 2. 测试用户信息接口 (/me) ===")
    try:
        r = urllib.request.urlopen(urllib.request.Request(
            f"{base_url}/v1/users/me",
            headers={'Authorization': f'Bearer {token}'}
        ))
        result = json.loads(r.read())
        if result['success'] and result['data']['id'] == 2:
            print(f"✅ 用户信息正确: {result['data']['name']}")
        else:
            print(f"❌ 用户信息错误: {result}")
    except Exception as e:
        print(f"❌ 用户信息接口异常: {e}")

def test_user_by_id(token):
    print("\n=== 3. 测试用户信息接口 (/users/{id}) ===")
    try:
        r = urllib.request.urlopen(urllib.request.Request(
            f"{base_url}/v1/users/2",
            headers={'Authorization': f'Bearer {token}'}
        ))
        result = json.loads(r.read())
        if result['success'] and result['data']['id'] == 2:
            print(f"✅ 用户查询正确: {result['data']['name']}")
        else:
            print(f"❌ 用户查询错误: {result}")
    except Exception as e:
        print(f"❌ 用户查询接口异常: {e}")

def test_audit_logs(token):
    print("\n=== 4. 测试审计日志接口 ===")
    try:
        r = urllib.request.urlopen(urllib.request.Request(
            f"{base_url}/v1/audit-logs?skip=0&limit=5",
            headers={'Authorization': f'Bearer {token}'}
        ))
        result = json.loads(r.read())
        if result['success']:
            logs = result['data'].get('audit_logs', [])
            print(f"✅ 审计日志查询成功，共 {len(logs)} 条记录")
            if logs:
                log = logs[0]
                print(f"   最新日志: user_id={log['user_id']}, action={log['action']}, resource={log['resource_type']}")
                if log['ip_address']:
                    print(f"   ✅ IP地址正确记录: {log['ip_address']}")
                else:
                    print(f"   ⚠️ IP地址为空")
        else:
            print(f"❌ 审计日志查询错误: {result}")
    except Exception as e:
        print(f"❌ 审计日志接口异常: {e}")

def test_categories(token):
    print("\n=== 5. 测试商品分类列表接口 ===")
    try:
        r = urllib.request.urlopen(urllib.request.Request(
            f"{base_url}/v1/categories",
            headers={'Authorization': f'Bearer {token}'}
        ))
        result = json.loads(r.read())
        if result['success']:
            print(f"✅ 分类列表查询成功")
            return True
        else:
            print(f"❌ 分类列表查询错误: {result}")
            return False
    except Exception as e:
        print(f"❌ 分类列表接口异常: {e}")
        return False

def test_create_category(token):
    print("\n=== 6. 测试创建分类（验证审计日志记录） ===")
    import time
    category_name = f"验证分类_{int(time.time())}"
    try:
        r = urllib.request.urlopen(urllib.request.Request(
            f"{base_url}/v1/categories",
            data=json.dumps({'name': category_name, 'description': '验证分类描述'}).encode(),
            headers={'Authorization': f'Bearer {token}', 'Content-Type': 'application/json'}
        ))
        result = json.loads(r.read())
        if result['success']:
            print(f"✅ 创建分类成功: {result['data']['name']}")
            return result['data']['id']
        else:
            print(f"❌ 创建分类错误: {result}")
            return None
    except Exception as e:
        print(f"❌ 创建分类异常: {e}")
        return None

def test_verify_audit_log_after_create(token, category_id):
    print("\n=== 7. 验证创建分类后审计日志记录 ===")
    try:
        r = urllib.request.urlopen(urllib.request.Request(
            f"{base_url}/v1/audit-logs?skip=0&limit=1",
            headers={'Authorization': f'Bearer {token}'}
        ))
        result = json.loads(r.read())
        if result['success']:
            logs = result['data'].get('audit_logs', [])
            if logs and logs[0]['action'] == 'CREATE' and logs[0]['resource_type'] == 'categories':
                print(f"✅ 审计日志正确记录了创建操作")
                print(f"   user_id={logs[0]['user_id']}, ip_address={logs[0]['ip_address']}")
                print(f"   new_data={logs[0]['new_data'][:50]}...")
            else:
                print(f"❌ 审计日志未记录创建操作")
        else:
            print(f"❌ 查询审计日志错误: {result}")
    except Exception as e:
        print(f"❌ 查询审计日志异常: {e}")

def test_products(token):
    print("\n=== 8. 测试商品列表接口 ===")
    try:
        r = urllib.request.urlopen(urllib.request.Request(
            f"{base_url}/v1/products",
            headers={'Authorization': f'Bearer {token}'}
        ))
        result = json.loads(r.read())
        if result['success']:
            print(f"✅ 商品列表查询成功")
        else:
            print(f"❌ 商品列表查询错误: {result}")
    except Exception as e:
        print(f"❌ 商品列表接口异常: {e}")

def test_consignment(token):
    print("\n=== 9. 测试寄卖品列表接口 ===")
    try:
        r = urllib.request.urlopen(urllib.request.Request(
            f"{base_url}/v1/consignment/items",
            headers={'Authorization': f'Bearer {token}'}
        ))
        result = json.loads(r.read())
        if result['success']:
            items = result['data'].get('consignments', [])
            print(f"✅ 寄卖品列表查询成功，共 {len(items)} 条")
        else:
            print(f"❌ 寄卖品列表查询错误: {result}")
    except Exception as e:
        print(f"❌ 寄卖品列表接口异常: {e}")

def test_suppliers(token):
    print("\n=== 10. 测试供应商列表接口 ===")
    try:
        r = urllib.request.urlopen(urllib.request.Request(
            f"{base_url}/v1/suppliers",
            headers={'Authorization': f'Bearer {token}'}
        ))
        result = json.loads(r.read())
        if result['success']:
            suppliers = result['data'].get('suppliers', [])
            print(f"✅ 供应商列表查询成功，共 {len(suppliers)} 条")
        else:
            print(f"❌ 供应商列表查询错误: {result}")
    except Exception as e:
        print(f"❌ 供应商列表接口异常: {e}")

def main():
    print("=" * 60)
    print("家族门店管理系统 - 功能验证测试")
    print("=" * 60)
    
    token = test_login()
    if not token:
        print("\n❌ 登录失败，无法继续测试")
        sys.exit(1)
    
    print(f"\nToken: {token[:30]}...")
    
    test_me(token)
    test_user_by_id(token)
    test_audit_logs(token)
    test_categories(token)
    category_id = test_create_category(token)
    test_verify_audit_log_after_create(token, category_id)
    test_products(token)
    test_consignment(token)
    test_suppliers(token)
    
    print("\n" + "=" * 60)
    print("测试完成！")
    print("=" * 60)

if __name__ == "__main__":
    main()
