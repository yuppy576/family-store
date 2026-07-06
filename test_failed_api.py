import requests

base_url = 'http://localhost:8090/v1'

r = requests.post(base_url + '/users/login', json={'email': 'admin@test.com', 'password': 'test123456'})
token = r.json()['data']['token']
headers = {'Authorization': 'Bearer ' + token}

print("测试失败接口的详细错误：")

r = requests.get(base_url + '/payments?skip=0', headers=headers)
print('支付方式列表:', r.status_code, '-', r.text)

r = requests.get(base_url + '/orders?skip=0', headers=headers)
print('订单列表:', r.status_code, '-', r.text)

r = requests.get(base_url + '/consignment/items', headers=headers)
print('寄卖品列表:', r.status_code, '-', r.text)

r = requests.get(base_url + '/suppliers?skip=0', headers=headers)
print('供应商列表:', r.status_code, '-', r.text)

r = requests.get(base_url + '/consignment/consignors', headers=headers)
print('寄卖人列表:', r.status_code, '-', r.text)
