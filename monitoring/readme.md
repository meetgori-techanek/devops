# Monitoring

## Types of monitoring
 1. server monitoring
 2. network monitoring
 3. cost monitoring
 4. infrastructure monitoring
 5. database monitoring


## observablity 
key components
1. monitoring
2. logging
3. tracing
4. alerting
5. visuliszation

- promithius is used for monitoring, logging and tracing
- grafana is used for visuliszation


### flow
```
     +-------------------------+
     |      Kubernetes         |
     |   (Apps & Services)     |
     +-----------+-------------+
                 |
                 v
        +--------+--------+
        |     Alloy       |  <-- Unified collector agent
        |  (on each node) |
        +--------+--------+
                 |
     +-----------+-----------+-------------+
     |                       |             |
     v                       v             v
+----------+         +-------------+  +--------------+
|   Loki   |  <---   |   Tempo     |  |   Metrics DB |
| (Logs)   |         | (Traces)    |  | (Optional: Prometheus, Mimir) |
+----------+         +-------------+  +--------------+
       \                     |                /
        \____________________|_______________/
                             |
                             v
                        +----------+
                        | Grafana  |
                        | (UI/UX)  |
                        +----------+

```
