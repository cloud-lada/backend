# backend

This repository contains the entire backend implementation for the cloud-lada project. It operates as a monorepo
of go services.

## Services

This section outlines the individual applications that go in to the cloud-lada platform:

### Ingestor

The ingestor is an HTTP API that handles inbound streams of reading data and writes them to an event bus for downstream
services to consume. Streams of JSON-encoded objects are received via the ingestor, here is a sample format:

```json lines
{"sensor": "speed", "value": 55.5, "timestamp": "2012-04-23T18:25:43.511Z"}
{"sensor": "speed", "value": 60.2, "timestamp": "2012-04-23T18:26:43.511Z"}
{"sensor": "speed", "value": 61.9, "timestamp": "2012-04-23T18:27:43.511Z"}
```

The reason for accepting JSON streams is that in the event that there is no internet connection for the on-board
instruments (which is likely for some portions of the trip), they can be stored on disk. Once the connection is
reestablished, these local files can be streamed to the ingestor. For the same reason, each reading contains a timestamp
of when it was taken.

#### Configuration

The ingestor accepts a small number of command-line flags to modify its behaviour:

* `--port` - The port to serve HTTP traffic on
* `--api-key` - The API key to check for in basic authentication for inbound HTTP requests
* `--event-writer-url` - A URL that describes the event bus to write events to, see the [gocloud](https://gocloud.dev/howto/pubsub/publish/) documentation for more information

#### Endpoints

* `/ingest` (POST) - Handles inbound sensor data.

### Persistor

The persistor is a service that handles readings from an event bus and persists them to a PostgreSQL compatible database.
For more information on the expected structure of the database schema, see the [migrations](./pkg/postgres/migrations)
directory.

#### Configuration

The persistor accepts a small number of command-line flags to modify its behaviour:

* `--event-writer-url` - A URL that describes the event bus to read events from, see the [gocloud](https://gocloud.dev/howto/pubsub/subscribe/) documentation for more information
* `--database-url` - A URL that describes the database to persist reading data to, see the [gocloud](https://gocloud.dev/howto/sql/) documentation for more information

### Dumper

The dumper is a job that writes an entire day's worth of reading data into a blob storage provider, with the intention of
making the daily dumps publicly accessible.

#### Configuration

The dumper accepts a small number of command-line flags to modify its behaviour:

* `--database-url` - A URL that describes the database to query reading data from, see the [gocloud](https://gocloud.dev/howto/sql/) documentation for more information
* `--blob-store-url` - A URL that describes the blob storage provider to write data dumps to, [gocloud](https://gocloud.dev/howto/blob/) documentation for more information
* `--dump-date` - A `YYYY-MM-DD` formatted string that specifies the date to produce a dump for.

### API

The API is an HTTP API that serves statistical data based on the readings stored in the database.

#### Configuration

The api accepts a small number of command-line flags to modify its behaviour:

* `--port` - The port to serve HTTP traffic on
* `--database-url` - A URL that describes the database to query reading data from, see the [gocloud](https://gocloud.dev/howto/sql/) documentation for more information

#### Endpoints

* `/api/statistics/latest` (GET) - Returns the latest sensor data.
* `/api/location/latest` (GET) - Returns the latest location data.

## CI

When opening a pull request, go code will be vetted and tests will be run. The same will happen when merging into the
`master` branch. When a tag is created, all binaries, docker images and manifests will be bundled and uploaded to the
release and GitHub docker image registry. Image tags will be managed by the CI when bundling the static manifests.

## Deployment

This repository contains a [Kustomize](https://kustomize.io/) manifest in the `manifests` directory which can be used
to run one or more ingestor instances in a Kubernetes cluster. You can use `make kustomize` to produce a static manifest
that contains everything you need. You will need to override the `ConfigMap` and `Secret` resources for your specific
setup. On each release, a static manifest is produced that can deploy all the services within this repository.
