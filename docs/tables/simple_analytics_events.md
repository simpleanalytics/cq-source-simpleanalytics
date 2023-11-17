# Table: simple_analytics_events

This table shows data for Simple Analytics Events.

https://docs.simpleanalytics.com/api/export-data-points

The primary key for this table is **_cq_id**.
It supports incremental syncs.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|metadata|`json`|
|added_iso|`timestamp[us, tz=UTC]`|
|added_unix|`int64`|
|browser_name|`utf8`|
|browser_version|`utf8`|
|country_code|`utf8`|
|datapoint|`utf8`|
|device_type|`utf8`|
|document_referrer|`utf8`|
|hostname|`utf8`|
|hostname_original|`utf8`|
|is_robot|`bool`|
|lang_language|`utf8`|
|lang_region|`utf8`|
|os_name|`utf8`|
|os_version|`utf8`|
|path|`utf8`|
|path_and_query|`utf8`|
|query|`utf8`|
|screen_height|`int64`|
|screen_width|`int64`|
|session_id|`utf8`|
|utm_campaign|`utf8`|
|utm_content|`utf8`|
|utm_medium|`utf8`|
|utm_source|`utf8`|
|utm_term|`utf8`|
|uuid|`utf8`|
|user_agent|`utf8`|
|viewport_height|`int64`|
|viewport_width|`int64`|