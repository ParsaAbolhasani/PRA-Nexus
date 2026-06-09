# ==================== Input از Filebeat و GELF ====================
input {
  beats {
    port => 5044
    add_field => { "data_source" => "filebeat" }
  }
  
  gelf {
    port => 5000
    type => gelf
    add_field => { "data_source" => "gelf" }
  }
  
  # برای تست مستقیم (اختیاری)
  tcp {
    port => 5000
    type => syslog
    add_field => { "data_source" => "tcp" }
  }
  
  udp {
    port => 5000
    type => syslog
    add_field => { "data_source" => "udp" }
  }
}

# ==================== Filter ====================
filter {
  # جداسازی لاگ‌های PRA Exchange
  if [container][name] =~ /pra/ or [service] =~ /pra/ {
    grok {
      match => { 
        "message" => [
          "(?<timestamp>%{TIMESTAMP_ISO8601})\s+(?<log_level>%{LOGLEVEL})\s+(?<log_message>.*)",
          "%{TIMESTAMP_ISO8601}\s+\[%{DATA:module}\]\s+%{LOGLEVEL:log_level}\s+-\s+%{GREEDYDATA:log_message}",
          "%{TIMESTAMP_ISO8601}\s+%{LOGLEVEL:log_level}\s+\[%{DATA:module}\]\s+-\s+%{GREEDYDATA:log_message}",
          "\[%{TIMESTAMP_ISO8601:timestamp}\]\s+%{LOGLEVEL:log_level}\s+%{GREEDYDATA:log_message}"
        ]
      }
    }
    
    # استخراج فیلدهای سفارشی PRA
    if [log_message] =~ /Trade/ {
      grok {
        match => { "log_message" => "Trade\s+(?<trade_id>%{NUMBER})" }
      }
    }
    
    if [log_message] =~ /Transaction/ {
      grok {
        match => { "log_message" => "Transaction\s+(?<tx_hash>0x[0-9a-fA-F]+)" }
      }
    }
    
    # اضافه کردن تگ‌ها
    mutate {
      add_tag => ["pra-exchange"]
    }
  }
  
  # لاگ‌های Docker
  if [container][name] {
    mutate {
      add_field => { "service_name" => "%{[container][name]}" }
    }
  }
  
  # لاگ‌های PostgreSQL
  if [container][name] =~ /postgres/ {
    grok {
      match => { "message" => "(?<db_timestamp>%{TIMESTAMP_ISO8601})\s+\[%{DATA:db_pid}\]\s+%{LOGLEVEL:db_level}:\s+%{GREEDYDATA:db_message}" }
    }
    mutate {
      add_tag => ["postgres", "database"]
    }
  }
  
  # لاگ‌های Redis
  if [container][name] =~ /redis/ {
    mutate {
      add_tag => ["redis", "cache"]
    }
  }
  
  # حذف فیلدهای اضافی
  mutate {
    remove_field => ["agent", "ecs", "input", "log", "tags"]
  }
  
  # تنظیم timestamp
  date {
    match => ["timestamp", "ISO8601"]
    target => "@timestamp"
  }
}

# ==================== Output ====================
output {
  # Elasticsearch
  elasticsearch {
    hosts => ["${ELASTIC_HOSTS:-http://elasticsearch:9200}"]
    user => "${ELASTIC_USER:-elastic}"
    password => "${ELASTIC_PASSWORD}"
    index => "pra-logs-%{+YYYY.MM.dd}"
    manage_template => true
    template_name => "pra-logs"
    template_overwrite => true
    template => '[{"index_patterns":["pra-logs-*"],"settings":{"number_of_shards":1,"number_of_replicas":0},"mappings":{"properties":{"message":{"type":"text"},"log_level":{"type":"keyword"},"service_name":{"type":"keyword"},"container_name":{"type":"keyword"},"data_source":{"type":"keyword"},"@timestamp":{"type":"date"}}}}]'
  }
  
  # خروجی برای اشکال‌زدایی
  stdout {
    codec => rubydebug
    enabled => false
  }
}