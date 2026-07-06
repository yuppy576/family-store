import subprocess
import json

ssh_cmd = [
    'ssh', '-i', f'{subprocess.getenv("USERPROFILE")}\\.ssh\\aliyunid_rsa',
    'root@8.162.7.31',
    'sudo -u postgres psql -d family_store -c "SELECT COUNT(*) FROM audit_logs;"'
]

result = subprocess.run(ssh_cmd, capture_output=True, text=True)
print("审计日志总数查询:")
print(result.stdout)
print(result.stderr)
