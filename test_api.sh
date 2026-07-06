#!/bin/bash
curl -s -X POST http://127.0.0.1:8081/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"test123456"}'
