#!/bin/bash
BACKUP_DIR="/backup"
DB_NAME="family_store"
DB_USER="family_store"
DB_PASSWORD="store2024"
DB_HOST="127.0.0.1"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/family_store_${TIMESTAMP}.sql.gz"

export PGPASSWORD="${DB_PASSWORD}"
pg_dump -h "${DB_HOST}" -U "${DB_USER}" "${DB_NAME}" | gzip > "${BACKUP_FILE}"

if [ $? -eq 0 ]; then
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Backup successful: ${BACKUP_FILE}"
else
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Backup FAILED!"
    exit 1
fi

find "${BACKUP_DIR}" -name "family_store_*.sql.gz" -type f -mtime +7 -exec rm {} \;
echo "[$(date '+%Y-%m-%d %H:%M:%S')] Cleanup completed"
