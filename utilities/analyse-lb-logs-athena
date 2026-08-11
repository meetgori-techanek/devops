# ALB Logs Athena Setup

## 1. Drop existing table if exist

```sql
DROP TABLE IF EXISTS alb_logs;
```

## 2. Create new table

```sql
CREATE EXTERNAL TABLE IF NOT EXISTS alb_logs_v2 (
    type string,
    time string,
    elb string,
    client_ip string,
    client_port int,
    target_ip string,
    target_port int,
    request_processing_time double,
    target_processing_time double,
    response_processing_time double,
    elb_status_code string,
    target_status_code string,
    received_bytes bigint,
    sent_bytes bigint,
    request_verb string,
    request_url string,
    request_proto string,
    user_agent string,
    ssl_cipher string,
    ssl_protocol string,
    target_group_arn string,
    trace_id string,
    domain_name string,
    chosen_cert_arn string,
    matched_rule_priority string,
    request_creation_time string,
    actions_executed string,
    redirect_url string,
    lambda_error_reason string,
    target_port_list string,
    target_status_code_list string,
    classification string,
    classification_reason string,
    conn_trace_id string
)
ROW FORMAT SERDE 'org.apache.hadoop.hive.serde2.RegexSerDe'
WITH SERDEPROPERTIES (
  'serialization.format' = '1',
  'input.regex'='([^ ]*) ([^ ]*) ([^ ]*) ([^ ]*):([0-9]*) ([^ ]*)[:-]([0-9]*) ([-.0-9]*) ([-.0-9]*) ([-.0-9]*) (|[-0-9]*) (|[-0-9]*) ([-0-9]*) ([-0-9]*) "([^ ]*) (.*) (- |[^ ]*)" "([^"]*)" ([A-Z0-9-_]+) ([A-Za-z0-9.-]*) ([^ ]*) "([^"]*)" "([^"]*)" "([^"]*)" ([-.0-9]*) ([^ ]*) "([^"]*)" "([^"]*)" "([^ ]*)" "([^\\s]+?)" "([^\\s]+)" "([^ ]*)" "([^ ]*)" ?([^ ]*)? ?( .*)?'
)
LOCATION 's3://pxo-prod-alb-logs/prod-live-lb/access-log/AWSLogs/590649111187/elasticloadbalancing/us-east-1/';
```

## 3. Verify data

```sql
SELECT *
FROM alb_logs_v2
LIMIT 5;
```

## 4. Get data for a specific month

```sql
SELECT 
    domain_name,
    SPLIT_PART(request_url, '?', 1) AS url_path,
    COUNT(*) AS request_count,
    SUM(sent_bytes) AS total_sent_bytes,
    ROUND(SUM(sent_bytes) / 1073741824.0, 2) AS total_sent_gb
FROM alb_logs_v2
WHERE time >= '2026-07-01' AND time < '2026-08-01'
GROUP BY domain_name, SPLIT_PART(request_url, '?', 1)
ORDER BY total_sent_bytes DESC
LIMIT 50;
```

## 5. Daily transfer on a specific URL

```sql
SELECT
    DATE(from_iso8601_timestamp(time)) AS day,
    COUNT(*) AS requests,
    ROUND(SUM(sent_bytes) / 1073741824.0, 3) AS transfer_gb,
    ROUND(AVG(sent_bytes) / 1024.0 / 1024.0, 3) AS avg_file_size_mb,
    ROUND(MAX(sent_bytes) / 1024.0 / 1024.0, 3) AS max_file_size_mb
FROM alb_logs_v2
WHERE time >= '2026-07-01'
  AND time < '2026-08-01'
  AND domain_name = 'pxo.rockwelltrading.com'
  AND SPLIT_PART(request_url, '?', 1) = 'https://api-pxo.rockwelltrading.com:443/v2/wheel/scannerV2'
GROUP BY DATE(from_iso8601_timestamp(time))
ORDER BY day;
```
