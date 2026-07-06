import os
import subprocess

cmd1 = [
    'ssh', '-i', f'{os.getenv("USERPROFILE")}\\.ssh\\aliyunid_rsa',
    'root@8.162.7.31',
    'sudo -u postgres psql -d family_store -c "SELECT store_id, COUNT(*) FROM audit_logs GROUP BY store_id;"'
]

result1 = subprocess.run(cmd1, capture_output=True, text=True)
print("审计日志按store_id分组:")
print(result1.stdout)

cmd2 = [
    'ssh', '-i', f'{os.getenv("USERPROFILE")}\\.ssh\\aliyunid_rsa',
    'root@8.162.7.31',
    'sudo -u postgres psql -d family_store -c "SELECT id, name, domain FROM stores;"'
]

result2 = subprocess.run(cmd2, capture_output=True, text=True)
print("\n租户列表:")
print(result2.stdout)

cmd3 = [
    'ssh', '-i', f'{os.getenv("USERPROFILE")}\\.ssh\\aliyunid_rsa',
    'root@8.162.7.31',
    'sudo -u postgres psql -d family_store -c "SELECT id, email, store_id FROM users WHERE email=\'admin@test.com\';"'
]

result3 = subprocess.run(cmd3, capture_output=True, text=True)
print("\n管理员用户信息:")
print(result3.stdout)
