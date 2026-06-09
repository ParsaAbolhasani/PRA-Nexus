#!/bin/bash

# تنظیمات
KIBANA_URL="http://localhost:5601"
KIBANA_USER="elastic"
KIBANA_PASS="${ELASTIC_PASSWORD:-changeme}"

# مسیر داشبوردها
DASHBOARD_DIR="../monitoring/elk/kibana/dashboards"

echo "🚀 Importing Kibana dashboards for PRA Exchange..."

for dashboard in "$DASHBOARD_DIR"/*.ndjson; do
  filename=$(basename "$dashboard")
  echo "📊 Importing $filename..."
  
  curl -X POST "${KIBANA_URL}/api/saved_objects/_import" \
    -H "kbn-xsrf: true" \
    -u "${KIBANA_USER}:${KIBANA_PASS}" \
    --form file=@"$dashboard"
  
  echo ""
done

echo "✅ Import completed!"