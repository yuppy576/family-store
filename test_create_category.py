import urllib.request
import json

login_data = json.loads(urllib.request.urlopen(urllib.request.Request(
    'http://127.0.0.1:8081/v1/users/login',
    data=json.dumps({'email': 'admin@test.com', 'password': 'test123456'}).encode(),
    headers={'Content-Type': 'application/json'}
)).read())

token = login_data['data']['token']
print(f"Token: {token[:30]}...")

import time
category_name = f"测试创建_{int(time.time())}"

try:
    req = urllib.request.Request(
        'http://127.0.0.1:8081/v1/categories/',
        data=json.dumps({'name': category_name, 'description': '测试创建分类'}).encode(),
        headers={'Authorization': f'Bearer {token}', 'Content-Type': 'application/json'}
    )
    r = urllib.request.urlopen(req)
    result = json.loads(r.read())
    print(f"创建分类响应: {result}")
except urllib.error.HTTPError as e:
    print(f"HTTP错误: {e.code} - {e.reason}")
    print(f"响应内容: {e.read().decode()}")
except Exception as e:
    print(f"异常: {e}")
    import traceback
    traceback.print_exc()
