#!/bin/bash
# Test script for role-based access control

BASE_URL="http://localhost:8080"

echo "=== Testing Role-Based Access Control ==="
echo ""

# 1. Register an admin user
echo "1. Registering admin user..."
ADMIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","name":"Admin User","password":"admin123","role":"admin"}')

ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])" 2>/dev/null)
echo "Admin user created. Role: $(echo $ADMIN_RESPONSE | python3 -c "import sys,json; print(json.load(sys.stdin)['user']['role'])" 2>/dev/null)"
echo ""

# 2. Register a regular user
echo "2. Registering regular user..."
USER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","name":"Regular User","password":"user123"}')

USER_TOKEN=$(echo $USER_RESPONSE | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])" 2>/dev/null)
echo "Regular user created. Role: $(echo $USER_RESPONSE | python3 -c "import sys,json; print(json.load(sys.stdin)['user']['role'])" 2>/dev/null)"
echo ""

# 3. Both users can GET users
echo "3. Testing GET /api/v1/users (both admin and user should succeed)..."
echo "  Admin GET users:"
curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$BASE_URL/api/v1/users" | python3 -c "import sys,json; data=json.load(sys.stdin); print(f'  ✓ Success: Found {len(data)} users')" 2>/dev/null || echo "  ✗ Failed"

echo "  User GET users:"
curl -s -H "Authorization: Bearer $USER_TOKEN" "$BASE_URL/api/v1/users" | python3 -c "import sys,json; data=json.load(sys.stdin); print(f'  ✓ Success: Found {len(data)} users')" 2>/dev/null || echo "  ✗ Failed"
echo ""

# 4. Admin can DELETE user, regular user cannot
echo "4. Testing DELETE /api/v1/users/{id} (only admin should succeed)..."
TEST_USER_ID=$(echo $USER_RESPONSE | python3 -c "import sys,json; print(json.load(sys.stdin)['user']['id'])" 2>/dev/null)

echo "  User trying to DELETE (should fail with 403):"
DELETE_USER_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X DELETE \
  -H "Authorization: Bearer $USER_TOKEN" \
  "$BASE_URL/api/v1/users/$TEST_USER_ID")
HTTP_CODE=$(echo "$DELETE_USER_RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
if [ "$HTTP_CODE" == "403" ]; then
  echo "  ✓ Correctly denied (403 Forbidden)"
else
  echo "  ✗ Unexpected result: HTTP $HTTP_CODE"
fi

echo "  Admin trying to DELETE (should succeed):"
DELETE_ADMIN_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X DELETE \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  "$BASE_URL/api/v1/users/$TEST_USER_ID")
HTTP_CODE=$(echo "$DELETE_ADMIN_RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
if [ "$HTTP_CODE" == "200" ]; then
  echo "  ✓ Successfully deleted user (200 OK)"
else
  echo "  ✗ Failed: HTTP $HTTP_CODE"
fi
echo ""

# 5. Admin can CREATE user, regular user cannot
echo "5. Testing POST /api/v1/users (only admin should succeed)..."

echo "  User trying to CREATE (should fail with 403):"
CREATE_USER_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"newuser@test.com","name":"New User","password":"pass123"}' \
  "$BASE_URL/api/v1/users")
HTTP_CODE=$(echo "$CREATE_USER_RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
if [ "$HTTP_CODE" == "403" ]; then
  echo "  ✓ Correctly denied (403 Forbidden)"
else
  echo "  ✗ Unexpected result: HTTP $HTTP_CODE"
fi

echo "  Admin trying to CREATE (should succeed):"
CREATE_ADMIN_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin.created@test.com","name":"Admin Created User","password":"pass123"}' \
  "$BASE_URL/api/v1/users")
HTTP_CODE=$(echo "$CREATE_ADMIN_RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
if [ "$HTTP_CODE" == "201" ]; then
  echo "  ✓ Successfully created user (201 Created)"
else
  echo "  ✗ Failed: HTTP $HTTP_CODE"
fi
echo ""

echo "=== Role-Based Access Control Test Complete ==="
