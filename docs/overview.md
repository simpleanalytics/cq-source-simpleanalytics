A [Simple Analytics](https://simpleanalytics.com/) source plugin for CloudQuery that loads raw page view and event data from Simple Analytics to any database, data warehouse or data lake supported by CloudQuery, such as PostgreSQL, BigQuery, Athena, and many more.

## Configuration

The following source configuration file will sync all data points for `mywebsite.com` to a PostgreSQL database. See [Quickstart](https://cloudquery.io/docs/quickstart) for more information on how to configure the source and destination.

```yaml
kind: source
spec:
  name: "simple-analytics"
  path: "simple-analytics/simple-analytics"
  registry: cloudquery
  version: "VERSION_SOURCE_SIMPLE_ANALYTICS"
  # use this to enable incremental syncing
  # backend_options:
  #   table_name: "cq_state_simpleanalytics
  #   connection: "@@plugins.DESTINATION_NAME.connection"
  tables: 
    ["*"]
  destinations: 
    - "postgresql"
  spec:
    # plugin spec section
    user_id: "${SA_USER_ID}"
    api_key: "${SA_API_KEY}"
    websites:
      - hostname: mywebsite.com
        metadata_fields: 
          - fieldname_text
          - fieldname_int
          # - ... 
```

### Plugin Spec

- `user_id` (`string`) (required)

  A user ID from Simple Analytics, obtained from the [account settings](https://simpleanalytics.com/account) page. It should start with `sa_user_id...`

- `api_key` (`string`) (required)

  An API Key from Simple Analytics, obtained from the [account settings](https://simpleanalytics.com/account) page. It should start with `sa_api_key...`

- `websites` (`list`) (required)

  A list of websites to sync data for. Each website should have the following fields:

    - `hostname` (`string`) (required)
    
      The hostname of the website to sync data for. This should be the same as the hostname in Simple Analytics.
  
    - `metadata_fields` (`[]string`) (optional)

      A list of metadata fields to sync, e.g. `["path_text", "created_at_time"]`. If not specified, no metadata fields will be synced.

- `start_date` (`string`) (optional)

  The date to start syncing data from. If not specified, the plugin will sync data from the beginning of time (or use a start time defined by `period`, if set).

- `end_date` (`string`) (optional) 

  The date to stop syncing data at. If not specified, the plugin will sync data until the current date.

- `period` (`string`) (optional)
  
  The duration of the time window to fetch historical data for, in days, months or years. It is used to calculate `start_date` if it is not specified. If `start_date` is specified, duration is ignored. Examples:
    - `7d`: last 7 days
    - `3m`: last 3 months
    - `1y`: last year

- `concurrency` (`integer`) (optional`) (default: `1000`)

  Best effort maximum number of Go routines to use. Lower this number to reduce memory usage.

## Incremental Syncing

The Simple Analytics plugin supports incremental syncing. This means that only new data points will be fetched from Simple Analytics and loaded into your destination. This is done by keeping track of the last date a sync was done, and only fetching new data from that date onwards.
To enable this, `backend_options` must be set in the spec (as shown in the example config above). By default, incremental syncing is turned off. Also note that this will introduce duplicates, unless the destination is using `overwrite-delete-stale` mode. Care should be taken to remove these duplicates after loading them. For more information, see [Managing Incremental Tables](https://cloudquery.io/docs/advanced-topics/managing-incremental-tables).

## Example Queries

### List the top 10 pages by views for a given period, excluding robots

```sql
select 
  path, 
  count(*) 
from 
  simple_analytics_page_views 
where 
  hostname = 'mywebsite.com'
  and is_robot is false 
  and added_iso between '2023-01-01' 
  and '2023-02-01'
group by 
  path 
order by
  count desc 
limit 
  10
```

```text
+----------------------------------+---------+
| path                             | count   |
|----------------------------------+---------|
| /                                | 100333  |
| /page                            | 91234   |
| /another-page                    | 84567   |
| /blog/introduction               | 74342   |
| /index                           | 69333   |
| /another/page                    | 64935   |
| /deeply/nested/page              | 50404   |
| /yet/another                     | 42309   |
| /some/page                       | 34433   |
| /about-us                        | 20334   |
+----------------------------------+---------+
```


### List events

```sql
select 
  added_iso, 
  datapoint, 
  path, 
  browser_name 
from 
  simple_analytics_events 
order by 
  added_iso desc 
limit 
  5
```

```text
+-------------------------+-----------+-----------------------------------------------+---------------+
| added_iso               | datapoint | path                                          | browser_name  |
|-------------------------+-----------+-----------------------------------------------+---------------|
| 2023-01-23 19:32:25.68  | 404       | /security                                     | Google Chrome |
| 2023-01-22 20:23:23.379 | 404       | /blog/running-cloudquery-in-gcp               | Google Chrome |
| 2023-01-19 12:04:57.095 | 404       | /docs/plugins/sources/vercel/configuration.md | Brave         |
| 2023-01-19 12:04:36.567 | 404       | /docsss                                       | Firefox       |
| 2023-01-19 01:50:19.259 | 404       | /imgs/gcp-cross-project-service-account       | Google Chrome |
+-------------------------+-----------+-----------------------------------------------+---------------+
```

## Development

### Run tests

```bash
make test
```

### Run linter

```bash
make lint
```

### Generate docs

```bash
make gen-docs
```

### Release a new version

1. Run `git tag v1.0.0` to create a new tag for the release (replace `v1.0.0` with the new version number)
2. Run `git push origin v1.0.0` to push the tag to GitHub

Once the tag is pushed, a new GitHub Actions workflow will be triggered to build the release binaries and create the new release on GitHub.
To customize the release notes, see the Go releaser [changelog configuration docs](https://goreleaser.com/customization/changelog/#changelog).

### Publish a new version to the Cloudquery Hub

After tagging a release, you can build and publish a new version to the [Cloudquery Hub](https://hub.cloudquery.io/) by running the following commands.
Replace `v1.0.0` with the new version number.

```bash
# -m parameter adds release notes message, output is created in dist/ directory
go run main.go package -m "Release v1.0.0" v1.0.0 .

# Login to cloudquery hub and publish the new version
cloudquery login -t simple-analytics
cloudquery plugin publish --finalize
```

After publishing the new version, it will [show up](https://hub.cloudquery.io/plugins/source/simple-analytics/simple-analytics) in the [hub](https://hub.cloudquery.io/).

For more information please refer to the official [Publishing a Plugin to the Hub](https://cloudquery.io/docs/developers/publishing-a-plugin-to-the-hub) guide.
