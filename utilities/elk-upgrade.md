1. first need to identify what version we are using
Central server (172-30-0-157): Elasticsearch 8.18.3, Kibana 8.18.3, Filebeat 8.18.3, APM Server 8.18.3

App servers:

| Server | Filebeat | Metricbeat |
|---|---|---|
| PROD-A | 8.18.1 | 8.18.1 |
| DEV-A | 8.18.3 | 8.18.3 |
| STAGE | 8.18.1 | 8.18.1 |
| PROD-B | 8.18.1 | 8.18.1 |
| PROD-C | 8.18.1 | 8.18.1 |
| PROD-D | 8.18.1 | 8.18.1 |
| dev-process-new | 7.17.29 | 7.17.29 |
| prod-process-new | 7.17.29 | 7.17.29 |

dev-process-new and prod-process-new are on 7.17.29, a major version behind the 8.18.3 cluster. These two need attention first in your upgrade plan.
