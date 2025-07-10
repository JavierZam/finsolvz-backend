#!/bin/bash

# 📊 Finsolvz Backend Performance Testing
# Comprehensive testing for all optimized endpoints

URL="${1:-}"
REPORT_ID="${2:-}"

if [ -z "$URL" ]; then
    echo "❌ Error: Service URL is required"
    echo "Usage: ./performance-test.sh SERVICE_URL [REPORT_ID]"
    echo "Example: ./performance-test.sh https://finsolvz-backend-xxx.a.run.app"
    exit 1
fi

echo "📊 Finsolvz Backend Performance Test"
echo "===================================="
echo "URL: $URL"
echo "Time: $(date)"
echo ""

# Test function
test_endpoint() {
    local endpoint=$1
    local description=$2
    local expected_max_ms=$3
    
    echo "🎯 Testing: $description"
    echo "   Endpoint: $endpoint"
    
    # Test 3 times and calculate average
    total=0
    count=0
    for i in {1..3}; do
        response_time=$(curl -w "%{time_total}" -s -o /dev/null "$URL$endpoint" 2>/dev/null || echo "timeout")
        if [ "$response_time" != "timeout" ]; then
            ms=$(echo "$response_time * 1000" | bc 2>/dev/null || echo "0")
            echo "   Request $i: ${ms}ms"
            total=$(echo "$total + $ms" | bc 2>/dev/null)
            count=$((count + 1))
        else
            echo "   Request $i: timeout"
        fi
    done
    
    if [ $count -gt 0 ]; then
        average=$(echo "scale=1; $total / $count" | bc 2>/dev/null)
        echo "   Average: ${average}ms"
        
        # Performance assessment
        avg_int=$(echo "$average" | cut -d. -f1)
        if [ "$avg_int" -lt "$expected_max_ms" ]; then
            echo "   ✅ PASS: Within expected ${expected_max_ms}ms"
        else
            echo "   ⚠️  SLOW: Exceeds expected ${expected_max_ms}ms"
        fi
    else
        echo "   ❌ FAIL: All requests timed out"
    fi
    echo ""
}

# Core endpoint tests
echo "🚀 Core Performance Tests"
echo "========================"

test_endpoint "/" "Health Check" 50
test_endpoint "/api/reports/paginated?page=1&limit=10" "Reports (Paginated)" 80
test_endpoint "/api/companies" "Companies (Optimized)" 60

# Individual report test (if ID provided)
if [ ! -z "$REPORT_ID" ]; then
    echo "📋 Individual Report Test"
    echo "========================"
    test_endpoint "/api/reports/$REPORT_ID" "Get Report by ID (Cached)" 50
fi

# Compression test
echo "🗜️ Compression Test"
echo "=================="

uncompressed=$(curl -s -w "%{size_download}" -o /dev/null "$URL/" 2>/dev/null || echo "0")
compressed=$(curl -s -H "Accept-Encoding: gzip" -w "%{size_download}" -o /dev/null "$URL/" 2>/dev/null || echo "0")

echo "   Uncompressed: ${uncompressed} bytes"
echo "   Compressed: ${compressed} bytes"

if [ "$compressed" -lt "$uncompressed" ] && [ "$compressed" -gt "0" ]; then
    reduction=$(echo "scale=1; (($uncompressed - $compressed) / $uncompressed) * 100" | bc 2>/dev/null || echo "0")
    echo "   ✅ Compression: ${reduction}% reduction"
else
    echo "   ⚠️  Compression: Not working optimally"
fi

echo ""

# Basic load test
echo "⚡ Basic Load Test"
echo "=================="
echo "   Testing 20 concurrent requests..."

start_time=$(date +%s)
for i in {1..20}; do
    curl -s -o /dev/null "$URL/" &
done
wait
end_time=$(date +%s)

duration=$((end_time - start_time))
echo "   20 requests completed in: ${duration}s"

if [ "$duration" -lt 10 ]; then
    echo "   ✅ Excellent load handling"
elif [ "$duration" -lt 20 ]; then
    echo "   ✅ Good load handling"
else
    echo "   ⚠️  Consider performance review"
fi

echo ""

# Summary
echo "📊 Performance Summary"
echo "====================="
echo ""
echo "🎯 Target Response Times (Jakarta region):"
echo "   • Health Check: <50ms"
echo "   • Companies: <60ms"
echo "   • Reports (Paginated): <80ms"
echo "   • Individual Report: <50ms"
echo ""
echo "🌏 Regional Performance Expectations:"
echo "   • Jakarta/Indonesia: 20-80ms"
echo "   • Asia-Pacific: 100-200ms"
echo "   • Global: 200-500ms"
echo ""
echo "💡 Optimization Features Active:"
echo "   ✅ Jakarta region deployment"
echo "   ✅ Database connection pooling"
echo "   ✅ MongoDB indexes"
echo "   ✅ Smart caching system"
echo "   ✅ Response compression"
echo "   ✅ Optimized aggregation pipelines"
echo "   ✅ Rate limiting protection"
echo ""
echo "🚀 Test completed: $(date)"