import urllib.request
import json

login_data = json.loads(urllib.request.urlopen(urllib.request.Request(
    'http://127.0.0.1:8081/v1/users/login',
    data=json.dumps({'email': 'admin@test.com', 'password': 'test123456'}).encode()
)).read())

token = login_data['data']['token']
print(f"Token: {token[:30]}...")

headers = {'Authorization': f'Bearer {token}'}

print("\nTest /v1/users/me:")
try:
    r = urllib.request.urlopen(urllib.request.Request('http://127.0.0.1:8081/v1/users/me', headers=headers))
    print(json.loads(r.read()))
except Exception as e:
    print(f"Error: {e}")

print("\nTest /v1/users/2:")
try:
    r = urllib.request.urlopen(urllib.request.Request('http://127.0.0.1:8081/v1/users/2', headers=headers))
    print(json.loads(r.read()))
except Exception as e:
    print(f"Error: {e}")
