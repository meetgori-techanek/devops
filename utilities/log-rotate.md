# Enable Log Rotation for pxo-processing (nohup.log)

## 1. Update logrotate Configuration

Edit or create the logrotate config file:

```bash
sudo nano /etc/logrotate.d/pxo-processing
```

Use the following configuration:

```conf
/home/ubuntu/pxo-2.0-processing/nohup.log {
    daily
    rotate 10
    dateext
    dateformat -%Y-%m-%d-%H%M%S
    compress
    delaycompress
    missingok
    notifempty
    copytruncate
}
```

Save and exit.

---

## 2. Resulting Filenames

Instead of numeric rotations like:

```text
nohup.log.1.gz
nohup.log.2.gz
```

Logs will be rotated using date-based filenames:

```text
nohup.log-2025-01-14-073000.gz
nohup.log-2025-01-15-073000.gz
nohup.log-2025-01-16-073000.gz
```

One compressed file is created per day.

---

## 3. Test Log Rotation

### Dry Run (No Changes Applied)
```
sudo logrotate -d /etc/logrotate.d/pxo-processing
```

### Force Rotation
```
sudo logrotate -f /etc/logrotate.d/pxo-processing
```

---

## 4. Configuration Behavior Summary

| Setting          | Description |
|------------------|-------------|
| `daily`          | Rotate logs once per day |
| `rotate 10`      | Keep last 10 rotated logs |
| `dateext`        | Append date to rotated files |
| `dateformat`     | Use `YYYY-MM-DD-HHMMSS` format |
| `compress`       | Compress rotated logs (`.gz`) |
| `delaycompress` | Compress from the next cycle (safer for running processes) |
| `missingok`      | Do not error if log file is missing |
| `notifempty`     | Skip rotation if log is empty |
| `copytruncate`  | Truncate log without stopping the Node process |

---

## 5. Important Notes

- `nohup.log` is never deleted.
- Logrotate truncates the file while the Node.js process continues running.
- This setup is safe for long-running production Node services started via `nohup`.
```
