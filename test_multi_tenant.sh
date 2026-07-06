#!/bin/bash

TOKEN=$(curl -s -X POST http://127.0.0.1:8081/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"test123456"}' | jq -r '.data.token')

echo "Token obtained: ${TOKEN:0:30}..."

echo -e "\n=== Testing Categories ==="
curl -s -X GET http://127.0.0.1:8081/v1/categories/?skip=0\&limit=20 \
  -H "Authorization: Bearer $TOKEN"

echo -e "\n\n=== Testing Products ==="
curl -s -X GET http://127.0.0.1:8081/v1/products/?skip=0\&limit=10 \
  -H "Authorization: Bearer $TOKEN"

echo -e "\n\n=== Testing Users ==="
curl -s -X GET http://127.0.0.1:8081/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
