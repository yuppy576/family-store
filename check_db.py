import subprocess
import os

cmd = [
    'ssh', '-i', f'{os.getenv("USERPROFILE")}\\.ssh\\aliyunid_rsa',
    'root@8.162.7.31',
    "sudo -u postgres psql -d family_store -c 'SELECT id, email, store_id FROM users WHERE email='\\''admin@test.com'\\'';'"
]

result = subprocess.run(cmd, capture_output=True, text=True)
print("管理员用户信息:")
print(result.stdout)
print(result.stderr)
