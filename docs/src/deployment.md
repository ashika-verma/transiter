# Running Transiter


## Container based deployments

The simplest way to deploy Transiter is to use
containers along with the included configurations
for either Kubernetes (using the Transiter Helm Chart) or Docker compose.

During Transiter's continuous integration process, a
Docker image `jamespfennell/transiter:<version>`
 is built and uploaded to Docker hub.
    
Here `<version>` can be any of the following:

- For stable releases, this is the version number of the release; for example `0.3.1`.

- Setting `<version>` to be `latest` will spin up the most recent stable release.

- Setting `<version>` to be `dev-latest` will spin up the most recent green build on master.

For production settings, pinning to a fixed stable release is, of course, recommended.
The `latest` and `dev-latest` versions may introduce breaking API changes at any time.


### Kubernetes: the Transiter Helm Chart

The Transiter Helm Chart enables easy installation of 
Transiter on a Kubernetes cluster.

The Helm chart can be installed by referencing the built `tgz` file 
in GitHub:
```sh
helm install https://github.com/jamespfennell/transiter/raw/master/docker/helmchart-0.1.0.tgz
```
Or, by checking out the Git repository and installing from the file system:
```sh
helm install ./docker/helmchart
```


The Helm chart contains many configuration options.
These are defined and described in the 
[Helm Chart's values file](https://github.com/jamespfennell/transiter/blob/master/docker/helmchart/values.yaml).
Some particularly important options are as follows.

- The `version` is, by default, `latest`. As per the above, in production this should usually
    be overridden with a fixed version.
    
- If you are using custom feed parsers distributed as Python packages,
    these packages can be made available to Transiter using the 
    `pythonPackagesToProvide` option.

- By default the Postgres container does *not* use a persistent volume. 
    In production a persistent volume should be used so that data lives
    beyond the lifecycle of the container.
    This is configured in the `postgres.persistentVolume` configuration.


### Docker compose

There is a [Docker compose configuration in the Github repository](https://raw.githubusercontent.com/jamespfennell/transiter/master/docker/docker-compose.yml)
that may be used a basis for a deployment.
Using Docker compose has its limitations in terms of customization.
The 
configuration file is mostly maintained to enable running a simpler CI pipeline on Travis,
rather than being a recommended deployment mode.

When launched, the Docker compose config starts a Transiter instance listening on port 8000.

In production, you should use a volume for the Postgres data directory to ensure
that data persists if the Postgres container is reset.
By default that Postgres data directory does *not* do this.
To enable a volume, change
```yaml
      - data-volume:/var/lib/postgresql/data-not-used
```
to
```yaml
      - data-volume:/var/lib/postgresql/data
```
and configure the volume `data-volume` appropriately.

## Interlude: Transiter service topology

Transiter uses five Docker containers to run, as
well as an sixth init container:

1. The web service, which is a Gunicorn process 
    running inside of the Transiter Docker image.
    This can be replicated to support more HTTP traffic.
    
1. The scheduler, which is also a Gunicorn process 
    running inside of the Transiter Docker image.
    The scheduler can _not_ be replicated, as this will result in
     duplicated updates being executed.    
     
1. The executor, which is a Celery process
    running inside of the Transiter Docker image.
    This can be replicated to support more feed update processing throughput.
     
1. A Postgres instance, using the vanilla Postgres Docker image.
    This is where all the data lives.

1. A RabbitMQ instance, for the Celery cluster.
    This uses the vanilla rabbitmq Docker image.

1. The init container, which is based on the the Transiter Docker image.
    This container just initializes the database schema and then exits.
    
The topology is identical for non-Docker deployments, described below.

## Running Transiter without containers

Transiter can be run on "bare metal" and in development often is.
You system needs to have the following things installed:
Python 3.6+, Postgres and RabbitMQ.

### Setting up Postgres

Postgres is the only officially supported database for Transiter right now.
Other systems such as SQLite and MySQL can be used with some tinkering and
 may be supported in the future;
start an issue on the issue tracker if you're interested.

By default, Transiter connects to Postgres using the following configuration.

Setting     | Default   | Environment variable to change the default
------------|-----------|---------------------------------
Host        | localhost | `TRANSITER_DB_HOST`
Port        | 5432      | `TRANSITER_DB_PORT`
Database    | transiter | `TRANSITER_DB_DATABASE`
Username    | transiter | `TRANSITER_DB_USERNAME`
Password    | transiter | `TRANSITER_DB_PASSWORD`

In order to change the default, set the environment variable
in any environment that a Transiter process (web service or task server)
will be running.

After setting up Postgres and a database, you need to 
initialize the database with the Transiter schema.
Assuming you're in a (virtual) environment in which the `transiter`
Python package has been installed, this can be done by running:
```sh
transiterclt db reset
```

### Running the web service

The Transiter web service response to HTTP requests.
Transiter was built using the Flask framework and thus 
contains a WSGI app for launching a HTTP service.

Assuming you have the `transiter` Python package installed,
Transiter's WSGI app is located at `transiter.http:wsgi_app`
So, for example, to launch a Transiter web service using Gunicorn:
```sh
gunicorn transiter.http:wsgi_app
```

If you're developing Transiter
it can be helpful to launch the web service in debug mode so that it 
automatically restarts when you've made changes.
Assuming you're in a Python (virtual) environment in which Transiter has 
been installed, run:

```sh
transiterclt launch webservice
```
This launches the Flask debugging server.


### Running the scheduler

The task server is responsible for scheduler periodic feed updates.
It is launched using:
```sh
transiterclt launch scheduler
```
The scheduler contains an HTTP interface and by default listens on localhost's port 5000.
The host and port can be changed using the environment variables
`TRANSITER_SCHEDULER_HOST` and `TRANSITER_SCHEDULER_PORT` respectively.
It is critical that these variables are also set for the web service and executor; otherwise,
the web service will be unable to contact the scheduler server which can result
in feeds not being auto updated after a transit system install.
 
### Running the executor

Simply:
```sh
transiterclt launch executor --logging-level info
```




# Documentation

Transiter can be configured to serve the documentation 
on the `/docs` path.
This requires that the documentation has been built into static
files using `mkdocs`; Transiter then serves files directly from 
the documentation build directory.

## Security and performance notes

Please be aware that there is an inherent security concern here.
If the documentation is misconfigured, Transiter may serve
files from a random directory on the host machine and may thus leak 
data and files on the machine. 
To address this concern,

1. The documentation is disabled by default. A person setting up a Transiter instance
    has to explicitly enable it by setting the environment variable
    `TRANSITER_DOCUMENTATION_ENABLED` to be true.
    If the documentation is disabled any path on `/docs` returns a
    `404 NOT FOUND` response.
    
2. When serving documentation, Transiter will make an effort to ensure that the
    configuration is correct and that it is really serving the documentation
    and not some other directory. 
    To do this, Transiter embeds a 96 character hex string in each documentation page
    and before serving files verifies that the hex string is in the `index.html`
    file at the root of the directory.
    This security mechanism addresses accidental mistakes.
    If the security check fails, a `503 SERVICE UNAVAILABLE` response is sent.
    
In addition to the security concern, there is also a performance concern.
The files are served using Python's Flask and Werkzeug libraries; this 
is a highly non-performant way to serve static content.
In addition, the security check described above entails an additional fixed computational
cost for serving each file.
For these performance reasons, consider not enabling the documentation
in production if malicious users could access it and use it as the basis
for a DoS attack.


## Setting up the documentation

First, to enable the documentation set the environment 
variable
    `TRANSITER_DOCUMENTATION_ENABLED` to be true.
    
Next, you need to compile the documentation.
In the Python environment you're working in, ensure 
the Transiter developer requirements have been installed;
in the root of the Github repo, run
```sh
pip install -r dev-requirements.txt
```
Then `cd` into the `docs` directory and run
```sh
mkdocs build
```
This will result in the documentation static files being built
and placed in the `site` subdirectory.

Finally, configure the 
`TRANSITER_DOCUMENTATION_ROOT` environment variable.
This should point to the `site` subdirectory. 
This can be either:

1. An absolute path.

2. A relative path, relative to the location of the Flask application.

If you're launching the Flask app based on the checked out Github
repo, it is located in `transiter/http` and hence the 
environment variable should be set
to `../../docs/site`. (This is, in fact, the default.)

Verify it's working by visiting the `/docs` path.
If it's not working, consult the console output which will detail
exactly what's happening.

## Documentation in the Docker image

The Docker image contains the built HTML documentation in the directory
`/transiter-docs`.
The internal documentation is disabled by default.


# other

# Monitoring using Prometheus

When running Transiter, it is critical that successful feed updates are occurring regularly.
Otherwise, the data in the Transiter instance will be stale.
In order to monitor the feed update process, Transiter exports metrics in Prometheus format on the `/metrics`
    endpoint.

There are four metrics exported.
For architectural reasons, these metrics are ultimately managed in the scheduler process.
If that process is not running, metrics will be unavailable.

## TRANSITER_NUM_UPDATES

This metric has four labels: `system_id`, `feed_id`, `status` (either `FAILURE` or `SUCCESS`), and `result`.
For each label, the corresponding time series reports the number of updates that have finished for that feed
    and with that status

Generally one will be interested in the `SUCCESS` case.
For this case there are two possible results: `UPDATED` and `NOT_NEEDED`.
The latter indicates that Transiter performed the feed update, 
    but that the data it retrieved from the transit agency was identical to the last successful update.
In this case the rest of the feed update process is skipped.
If `NOT_NEEDED` is seen consecutively for a long time, it indicates that the transit agency is not provided updated data.
This may be a failure case, even though the status is `SUCCESS`.
Or, it may not be a failure case if it's nighttime, and the transit agency not running any services.


## TRANSITER_LAST_UPDATE

This metric has the same four labels as the last.
For each label, the corresponding time series reports the timestamp when the last update for that feed and
with that status occurred.


## TRANSITER_SUCCESSFUL_UPDATE_LATENCY

This metric has two labels: `system_id` and `feed_id`.
The metric reports the number of seconds between the last two feed updates with status `SUCCESS` and result `UPDATED`.
Over time, the metric is an accurate measure of how often data from that feed is being updated.

This metric is a good candidate for alerting because if is large, it directly indicates stale data.
There is a catch, though: if the reason feed updates are not occurring is because of internal problems with Transiter
(for example, RabbitMQ is down), the most recent value of the metric will stay constant and not indicate a problem.
Conversely, the `TRANSITER_NUM_UPDATES` metric will still catch this case: if that metric remains constant for a long
    time, it indicates no updates are happening.


## TRANSITER_NUM_ENTITIES

This metric reports the number of entities (trips, alerts, routes, etc.)
that are in the system, by feed.
It has three labels: `system_id`, `feed_id` and `entity_type`.
