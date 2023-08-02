# The Transiter Tour

The Transiter Tour is a grounds-up introduction to Transiter.
You'll learn how to get Transiter running,
    add some transit systems,
    and read data from the API.

There are some prerequisites for the Tour:

- You have Go installed.

- You've cloned the Transiter Git repo and are in the root of the repo.

- You have a Postgres instance to use, or Docker is available.

## Launch Transiter

Transiter requires a Postgres database.
By default it tries to connect to Postgres on `localhost:5432`
  with the username/password/database combination `transiter`/`transiter`/`transiter`.
If you don't already have Postgres running and Docker is available,
  you can easily spin up Postgres with these credentials by running the following command in the root of the repo:

```
$ docker-compose up -d postgres
```

Otherwise, just create new empty database in your preexisting Postgres instance
  and have your connection details ready.

Transiter is written in Go.
With Go installed and Postgres running, the Transiter server is launched using:

```
$ go run . server --log-level debug
```

If you're using a non-default Postgres configuration, launch Transiter using:

```
$ go run . server --log-level debug \
    -p postgres://${USERNAME}:${PASSWORD}@${HOST}:${PORT}/${DATABASE_NAME}
```

We've passed `--log-level debug` to start the server with debug logging.
This will give us more insight into what's happening inside the server later.

During the Tour we're going to interact with the Transiter server in two ways.

### HTTP API

After launching the Transiter server, 
  Transiter exports a public HTTP API on port 8080.
  Let's see what it returns:

```
$ curl localhost:8080
{
  "transiter":  {
    "version":  "1.0.0-dev",
    "href":  "https://github.com/jamespfennell/transiter"
  },
  "systems":  {
    "count":  "0"
  }
}
```

The main thing to note here is that the number of systems is 0.
This means we have no transit system installed.
Because of this there's not much else to see in the API yet.

Transiter also exports an admin HTTP API on port 8082.
This is a superset of the public API,
  so endpoints accessible through the public API are also accessible through the admin API:

```
$ curl localhost:8082
# same response as before
```

However admin API has additional functionality which is used to manage Transiter.
For example, we can discover the current log level:

```
curl localhost:8082/loglevel 
{
  "logLevel":  "DEBUG"
}
```

By sending a `POST` request to the same endpoint, it's possible to change log level.

In the Tour we will mainly avoid the admin HTTP API,
  and will instead manage Transiter using the Transiter CLI.

??? info "gRPC APIs"
    In addition to the two HTTP APIs,
    Transiter also runs a public gRPC API on port 8081 and an admin gRPC API on 8083.
    These are identical to the HTTP APIs, but may easier to use programmatically.

??? info "Changing the default port numbers"
    The default port numbers (8080-8083) can be changed using the flags 
    `--public-http-addr`,
    `--public-grpc-addr`,
    `--admin-http-addr`,
    `--admin-grpc-addr`.
    If you set any of these flags to `-`, Transiter won't run the associated API.


### Transiter CLI

We've already used the Transiter CLI to launch the Transiter server above (`go run . server`).
The CLI has many other commands, all of which are client commands.
These commands are used to interact with a running Transiter server.
For example, we can list the currently installed Transit systems:

```
$ go run . list
No transit systems installed.
```

As expected, there is nothing to return
  because there are no transit systems installed!
Let's do that!

## Installing transit systems

Each deployment of Transiter can have multiple transit systems side-by-side.
A transit system is installed by providing Transiter with a YAML configuration file for the system.
This config contains basic metadata about the system like its name,
    as well as URLs for the system's GTFS static and realtime feeds.

We're going to start by installing the [PATH train](https://en.wikipedia.org/wiki/PATH_(rail_system)) in New York and New Jersey.
The YAML configuration file for this system is included in
  [Transiter's library of systems on GitHub](https://github.com/jamespfennell/transiter/blob/master/systems).
We install it by provided the URL of the config to Transiter:

```
go run . install us-ny-path https://raw.githubusercontent.com/jamespfennell/transiter/master/systems/us-ny-path.yaml
```

The install will take a few seconds to complete, with
    most of the time spent loading the PATH's schedule into the database.
After it finishes, hit the Transiter root endpoint again:

```json
$ curl localhost:8080
{
  "transiter": {
    "version": "1.0.0-dev",
    "href": "https://github.com/jamespfennell/transiter",
  },
  "systems": {
    "count": 1,
    "href": "http://localhost:8000/systems"
  }
}
```

It's installed! 
Next, let's navigate to the list systems endpoint.
The URL `http://localhost:8000/systems` was given in the response above.
In general the Transiter API is designed with discoverability in mind.
We get:

```
$ curl localhost:8000/systems
{
  "systems":  [
    {
      "id":  "us-ny-path",
      "resource":  null,
      "name":  "Port Authority Trans-Hudson (PATH)",
      "status":  "ACTIVE",
      "agencies":  {
        "count":  "1",
        "href":  "https://demo.transiter.dev/systems/us-ny-path/agencies"
      },
      "feeds":  {
        "count":  "2",
        "href":  "https://demo.transiter.dev/systems/us-ny-path/feeds"
      },
      "routes":  {
        "count":  "6",
        "href":  "https://demo.transiter.dev/systems/us-ny-path/routes"
      },
      "stops":  {
        "count":  "64",
        "href":  "https://demo.transiter.dev/systems/us-ny-path/stops"
      },
      "transfers":  {
        "count":  "0",
        "href":  "https://demo.transiter.dev/systems/us-ny-path/transfers"
      }
    },
  ]
}
```

This is an overview of the system, showing the number of various entities like stops and routes.
All of these entities correspond to entities in the GTFS static and realtime specs.


??? info "Installing systems using the HTTP admin API"
    It's possible to install transit systems using the HTTP admin API by sending a `PUT` request
    to the system endpoint:

    ```
    curl -X PUT localhost:8082/systems/us-ny-path -d '{"yaml_config": {"url": "https://raw.githubusercontent.com/jamespfennell/transiter/master/systems/us-ny-path.yaml"}}'
    ```

GOT TO HERE

## Explore route data

Let's dive into the routes data that Transiter exposes.
Navigating to the list routes endpoint (given above; `http://localhost:8000/systems/bart/routes`)
lists all the routes.
We'll focus on a specific route, the *Berryessa/North San José–Richmond line* or *Orange Line* 
([Wikipedia page](https://en.wikipedia.org/wiki/Berryessa/North_San_Jos%C3%A9%E2%80%93Richmond_line)).
The route ID for this system is `3`, so we can find it by navigating to,

```text
http://localhost:8000/systems/bart/routes/3
```

The start of response will look like this,

```json
{
  "id": "3",
  "color": "FF9933",
  "short_name": "OR-N",
  "long_name": "Berryessa/North San Jose to Richmond",
  "description": "",
  "url": "http://www.bart.gov/schedules/bylineresults?route=3",
  "type": "SUBWAY",
  "estimated_headway": 300,
  "agency": null,
  "alerts": [],
  "service_maps": // public map definitions in here
}
```

Most of the basic data here, such as [the color FF9933](https://www.color-hex.com/color/ff9933),
is taken from the GTFS Static feed.
The *alerts* are taken from the GTFS Realtime feed.
Depending on the current state of the system when you take the tour, this may be empty.
The *estiamted headway* is calculated by Transiter.
It is the current average time between realtime trips on this route.
If there is insufficient data to estimate the headway, it will be `null`.

Arguably the most useful data here, though, are the *service maps*.

### Service maps

When transit consumers think of a route, they often think of the list of stops the route usually calls at.
In our example above, "the Orange Line goes from Richmond to San José."
Even though it's so central to how people think about routes, 
    GTFS does not directly give this kind of information.
However, Transiter has a system for automatically 
    generating such lists of stops using the timetable and realtime data in the GTFS feeds.
They're called *service maps* in Transiter.

Each route can have multiple service maps.
In the BART example there are two maps presented in the routes endpoint: `all-times` and `realtime`:

```json
  "service_maps": [
    {
      "group_id": "all-times",
      "stops": [
        {
          "id": "place_RICH",
          "name": "Richmond",
          "href": "http://localhost:8000/systems/bart/stops/place_RICH"
        },
        {
          "id": "place_DELN",
          "name": "El Cerrito Del Norte",
          "href": "http://localhost:8000/systems/bart/stops/place_DELN"
        },
        // More stops ...
        {
          "id": "place_BERY",
          "name": "Berryessa",
          "href": "http://localhost:8000/systems/bart/stops/place_BERY"
        }
      ]
    },
    {
      "group_id": "realtime",
      "stops": [
        // List of stops ...
      ]
    }
  ]
```

Transiter enables generating service maps from two sources:

1. The realtime trips in the GTFS Realtime feeds.
    The `realtime` service map was generated in this way.

1. The timetable in the GTFS Static feeds.
    Transiter can calculate the service maps using every trip in the timetable, like the `all-times` service map.
    Transiter can also calculate service maps for a subset of the timetable - for example, just using
        the weekend trips, or just using the weekday trips.
    (The weekday service map will make an appearance below.)

More information on configuring service maps can be found in the 
    [service maps section](systems.md#service-maps)
    of the transit system config documentation page.
    
## Explore stop data

Having looked at routes, let's look at some stops data.
The endpoint for the BART transit system (`http://localhost:8000/systems/bart`)
tells us where the find the list of all stops (`http://localhost:8000/systems/bart/stops`).
Let's go straight to the Downtown Berkeley station, which has stop ID `place_DBRK`
    and URL `http://localhost:8000/systems/bart/stops/place_DBRK` 
    (both the ID and the URL can be found in the Orange Line's service map).
This is the start of what comes back:

```json
{
  "id": "place_DBRK",
  "name": "Downtown Berkeley",
  "longitude": "-122.268109",
  "latitude": "37.870110",
  "url": "",
  "service_maps": [
    {
      "group_id": "weekday",
      "routes": [
        {
          "id": "3",
          "color": "FF9933",
          "href": "http://localhost:8000/systems/bart/routes/3"
        },
        {
          "id": "4",
          "color": "FF9933",
          "href": "http://localhost:8000/systems/bart/routes/4"
        },
        {
          "id": "7",
          "color": "FF0000",
          "href": "http://localhost:8000/systems/bart/routes/7"
        },
        {
          "id": "8",
          "color": "FF0000",
          "href": "http://localhost:8000/systems/bart/routes/8"
        }
      ]
    }
  // more data
  ]
}
```

The basic data at the start again comes directly from GTFS Static.

Next, we again see service maps!
In the route endpoint, the service maps returned the list of stops a route called at.
At the stop endpoint we see the inverse of this: the list of routes that a stop is included in.
Like for routes, this is important data that consumers associate with a station
    ("at Downtown Berkeley, I can get on the Orange Line"),
    which is not given explicitly in the GTFS Static feed,
    and which Transiter calculates automatically.
The two service maps shown are `weekday`, which is built using the weekday schedule, and `realtime`.

Going further down we see *alerts* (as mentioned above) and *transfers* (will be mentioned below).
Probably the most important data at the stop however is the *stop times*:
    these show the realtime arrivals for the station.
```json
{
  // ...
  "stop_times": [
    {
      "arrival": {
        "time": 1593271377.0,
        "delay": 0,
        "uncertainty": 30
      },
      "departure": {
        "time": 1593271401.0,
        "delay": 0,
        "uncertainty": 30
      },
      "track": null,
      "future": true,
      "stop_sequence": 12,
      "direction": null,
      "trip": {
        "id": "2570751",
        "direction_id": false,
        "started_at": null,
        "updated_at": null,
        "delay": null,
        "vehicle": null,
        "route": {
          "id": "3",
          "color": "FF9933",
          "href": "http://localhost:8000/systems/bart/routes/3"
        },
        "last_stop": {
          "id": "DELN",
          "name": "El Cerrito Del Norte",
          "href": "http://localhost:8000/systems/bart/stops/DELN"
        },
        "href": "http://localhost:8000/systems/bart/routes/3/trips/2570751"
      }
    }
  ]
}
```

The response shows the arrival and departure times, as well as details
    on the associated trip, the associated vehicle (if defined in the feed),
    the route of trip,
    and the last stop the trip will make.


## Search for stops by location

So far we have navigated though Transiter's data using the links in each endpoint,
    starting from the system, to a route, to a stop.
Apps built on top of Transiter (like [realtimerail.nyc](https://www.realtimerail.nyc)) often follow this pattern.

Another way to find stops is to search for them by location,
    and Transiter supports this too.
Suppose we're at [Union Square in San Francisco](https://en.wikipedia.org/wiki/Union_Square%2C_San_Francisco)
    whose coordinates are latitude 37.788056 and longitude -122.4075.
    and we want to see what BART stations are nearby.
We can perform a geographical search in Transiter by sending a `POST` requqsts to the stops endpoint.
Using `curl`,

```
curl -X POST "http://localhost:8000/systems/bart/stops?latitude=37.788056&longitude=-122.4075"
```

This returns two stations, starting with Powell Street, which is closest to Union Square.
Its distance is returned as 383 meters.
The next closest station is Montgomery Street, which is 534 meters away.

```json

[
  {
    "id": "place_POWL",
    "name": "Powell Street",
    "distance": 383,
    "service_maps": [
      // ...
    ],
    "href": "http://localhost:8000/systems/bart/stops/place_POWL"
  },
  {
    "id": "place_MONT",
    "name": "Montgomery Street",
    "distance": 534,
    "service_maps": [
      // ...
    ],
    "href": "http://localhost:8000/systems/bart/stops/place_MONT"
  }
]
```

It's possible to search for more stops by passing a distance URL parameter to the search.

## Where to go next?

- Learn more about [Transiter's system configuration YAML format](systems.md).

- Read advice about [deploying and monitoring Transiter](deployment.md).

- Consult the [API reference](api/index.md) to discover more endpoints and data that Transiter exposes.
