# AWS WAF Logs with S3 and Athena

## Purpose

Enable AWS WAF logging to Amazon S3 and use Amazon Athena to query requests, investigate blocked traffic, and troubleshoot API errors such as `502` and `504`.

## 1. Enable AWS WAF logging and configure the S3 destination

1. Open **AWS WAF & Shield** in the AWS Console.
2. Select the relevant Web ACL, for example, `waf-prod`.
3. Open **Logging and metrics**.
4. Enable logging and select the S3 bucket as the logging destination.
5. Use or create a bucket with a name that meets AWS WAF's logging destination requirements.
6. Save the configuration and confirm that log objects are being delivered to S3.

Example S3 log path:

```text
s3://aws-waf-logs-prod/AWSLogs/590649111187/WAFLogs/us-east-1/waf-prod/
```

Replace the example account ID, region, Web ACL name, bucket, and prefix with the values for your environment.

## 2. Create the Athena table

1. Open **Amazon Athena**.
2. Select the intended workgroup and database. The example below creates the table in the `default` database.
3. Replace the `LOCATION` value with the S3 prefix containing the WAF log objects.
4. Run the statement.

```sql
CREATE EXTERNAL TABLE IF NOT EXISTS default.waf_logs (
    `timestamp` BIGINT,
    formatversion INT,
    webaclid STRING,
    terminatingruleid STRING,
    terminatingruletype STRING,
    action STRING,
    httprequest STRUCT<
        clientip:STRING,
        country:STRING,
        uri:STRING,
        args:STRING,
        httpmethod:STRING,
        httpversion:STRING,
        requestid:STRING,
        headers:ARRAY<STRUCT<name:STRING,value:STRING>>,
        host:STRING
    >
)
ROW FORMAT SERDE 'org.openx.data.jsonserde.JsonSerDe'
LOCATION 's3://YOUR_WAF_BUCKET/PREFIX/'
TBLPROPERTIES ('ignore.malformed.json'='true');
```

**Important:** Replace `s3://YOUR_WAF_BUCKET/PREFIX/` with the actual S3 log prefix. The table schema must match the JSON fields in your WAF logs. If the table already exists, `IF NOT EXISTS` will not update its schema or location.

## 3. Validate the table

Run a small query to confirm that Athena can read the logs. Update the dates to the required time range.

```sql
SELECT
    from_unixtime("timestamp" / 1000) AS request_time,
    action,
    httprequest.clientip AS source_ip,
    httprequest.uri AS uri,
    terminatingruleid
FROM default.waf_logs
WHERE "timestamp" >= to_unixtime(
          TIMESTAMP '2026-10-12 00:00:00'
      ) * 1000
  AND "timestamp" < to_unixtime(
          TIMESTAMP '2026-10-13 00:00:00'
      ) * 1000
ORDER BY "timestamp" DESC
LIMIT 1000;
```

The query assumes the WAF `timestamp` field is Unix time in milliseconds. It uses a half-open interval, including the start time and excluding the end time.

## 4. Filter requests for affected API endpoints

Use this query to narrow the results to endpoints related to stock prices, options, historical data, expiration dates, trading logs, income statements, and authorization.

```sql
SELECT
    from_unixtime("timestamp" / 1000) AS request_time,
    action,
    httprequest.clientip AS source_ip,
    httprequest.uri AS uri,
    terminatingruleid
FROM default.waf_logs
WHERE "timestamp" >= to_unixtime(
          TIMESTAMP '2026-10-12 00:00:00'
      ) * 1000
  AND "timestamp" < to_unixtime(
          TIMESTAMP '2026-10-13 00:00:00'
      ) * 1000
  AND (
      httprequest.uri LIKE '%price%'
      OR httprequest.uri LIKE '%options%'
      OR httprequest.uri LIKE '%history%'
      OR httprequest.uri LIKE '%expiration%'
      OR httprequest.uri LIKE '%tradingLog%'
      OR httprequest.uri LIKE '%income%'
      OR httprequest.uri LIKE '%users/self%'
  )
ORDER BY "timestamp" DESC
LIMIT 500;
```

Change the start and end timestamps to the incident window you are investigating. URL path matching is case-sensitive with `LIKE`, so adjust the patterns if your API paths use different capitalization.

## 5. Interpret the results

- `ALLOW`: WAF allowed the request. Continue investigating the load balancer, EC2 Node.js backend, and database logs for the same timestamp.
- `BLOCK`: WAF blocked the request. Inspect `terminatingruleid` to identify the rule responsible.
- `terminatingruleid = Default_Action`: No earlier terminating rule matched; the Web ACL's default action was applied.

WAF logs show the request and WAF decision, but do not show how long the Node.js application or database took to process an allowed request.

## Troubleshooting

- **`TABLE_NOT_FOUND`:** Confirm the database and table name. Run the `CREATE EXTERNAL TABLE` statement before querying.
- **No rows returned:** Confirm the S3 `LOCATION`, date range, region, Web ACL prefix, and that log objects exist for the selected date.
- **Schema or SerDe errors:** Inspect a raw JSON log object in S3 and update the Athena table schema to match it.
- **Access denied:** Check Athena permissions, S3 read permissions, and the Athena query-results output location permissions.
- **High scan cost:** Narrow the S3 prefix to the required date where practical and consider partitioning the table for recurring analysis.


