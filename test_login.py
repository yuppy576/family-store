import requests

r = requests.post('http://localhost:8090/v1/users/login', json={'email': 'admin@test.com', 'password': 'test123456'})
print("登录响应:", r.text)
