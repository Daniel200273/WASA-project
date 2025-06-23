#!/bin/bash

echo "🚀 WASA Photo Upload Quick Test"
echo "==============================="

# Check if test image exists
echo "📸 Checking for test image..."
if [ ! -f "test_user.png" ]; then
    echo "❌ test_user.png not found in current directory"
    echo "💡 Please make sure test_user.png exists before running this script"
    exit 1
fi
echo "✅ Found test_user.png ($(wc -c < test_user.png) bytes)"
echo ""
echo "⏸️  Press Enter to continue to build step..."
read -r

# Build and start server
echo ""
echo "📦 Building application..."
if ! go build ./cmd/webapi; then
    echo "❌ Build failed!"
    exit 1
fi
echo "✅ Build successful!"
echo ""
echo "⏸️  Press Enter to start the server..."
read -r

echo ""
echo "🔄 Starting server..."
./webapi &
SERVER_PID=$!
sleep 3

# Check server is running
echo "🔍 Checking server health..."
echo "🔧 Curl command:"
echo "curl -s http://localhost:3000/liveness"
echo ""
if ! curl -s http://localhost:3000/liveness >/dev/null; then
    echo "❌ Server not responding"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi
echo "✅ Server running"
echo ""
echo "⏸️  Press Enter to proceed with authentication..."
read -r

# Login to get authentication token
echo ""
echo "🔐 Authenticating..."
echo "🔧 Curl command:"
echo "curl -s -X POST -H 'Content-Type: application/json' -d '{\"name\": \"testuser\"}' http://localhost:3000/session"
echo ""

TOKEN=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"name": "testuser"}' \
    http://localhost:3000/session | \
    grep -o '"identifier":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Authentication failed"
    kill $SERVER_PID 2>/dev/null
    exit 1
fi
echo "✅ Token: $TOKEN"
echo ""
echo "⏸️  Press Enter to upload the photo..."
read -r

# Upload photo
echo ""
echo "📤 Uploading photo..."
echo "🔧 Curl command:"
echo "curl -s -w \"%{http_code}\" -X PUT -H \"Authorization: Bearer $TOKEN\" -F \"photo=@test_user.png\" http://localhost:3000/users/me/photo -o /dev/null"
echo ""

UPLOAD_RESPONSE=$(curl -s -w "%{http_code}" \
    -X PUT \
    -H "Authorization: Bearer $TOKEN" \
    -F "photo=@test_user.png" \
    http://localhost:3000/users/me/photo \
    -o /dev/null)

echo "📊 Upload status: $UPLOAD_RESPONSE"

if [ "$UPLOAD_RESPONSE" = "204" ]; then
    echo "✅ Upload successful!"
    echo ""
    echo "⏸️  Press Enter to check uploaded files..."
    read -r
    
    # Check uploaded file
    echo ""
    echo "📁 Uploaded files:"
    find tmp/uploads -name "*.png" 2>/dev/null || echo "No files found"
    echo ""
    echo "⏸️  Press Enter to test static file access..."
    read -r
    
    # Test static file access
    UPLOADED_FILE=$(find tmp/uploads -name "*.png" | head -1)
    if [ -f "$UPLOADED_FILE" ]; then
        FILENAME=$(basename "$UPLOADED_FILE")
        echo ""
        echo "🌐 Testing static access..."
        echo "🔧 Curl command:"
        echo "curl -s -w \"%{http_code}\" http://localhost:3000/uploads/profiles/$FILENAME -o /dev/null"
        echo ""
        STATIC_RESPONSE=$(curl -s -w "%{http_code}" \
            "http://localhost:3000/uploads/profiles/$FILENAME" \
            -o /dev/null)
        echo "📊 Static file status: $STATIC_RESPONSE"
        
        if [ "$STATIC_RESPONSE" = "200" ]; then
            echo "✅ Static serving works!"
        else
            echo "⚠️  Static serving issue"
        fi
    fi
else
    echo "❌ Upload failed"
fi

echo ""
echo "⏸️  Press Enter to stop the server and cleanup..."
read -r

# Cleanup
echo ""
echo "🛑 Stopping server..."
kill $SERVER_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null

echo "🧹 Cleanup complete"
echo "📁 Your test_user.png file is preserved"

echo ""
echo "✅ Test completed!"
echo ""
echo "📋 Manual commands for testing:"
echo "1. Start server: ./webapi &"
echo "2. Login: TOKEN=\$(curl -s -X POST -H 'Content-Type: application/json' -d '{\"name\": \"user\"}' http://localhost:3000/session | grep -o '\"identifier\":\"[^\"]*\"' | cut -d'\"' -f4)"
echo "3. Upload: curl -X PUT -H \"Authorization: Bearer \$TOKEN\" -F \"photo=@test_user.png\" http://localhost:3000/users/me/photo"
echo "4. Stop: pkill webapi"
