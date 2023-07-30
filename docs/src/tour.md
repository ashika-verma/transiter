# The Transiter Tour

The Tour is a grounds-up introduction to Transiter.
You'll learn how to get Transiter running,
    add some transit systems,
    and read data from the API.

There are some prerequisites for doing the tour:

- Go is installed.

- You've checked out the Transiter Git repo.

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
  and have your connection details handy.

Transiter is written in Go.
With Go installed and Postgres running, the Transiter server is launched using:

```
$ go run . server --log-level debug
```

If you're using a non-default Postgres configuration, launch Transiter using:

```
$ go run . server --log-level debug --$postgres://${USERNAME}:${PASSWORD}@${HOST}:${PORT}/${DATABASE_NAME}
```

We've passed `--log-level debug` to start the server with debug logging.
This will allow us to have more insight into what's happening inside the server later.

During the Tour we're going to interact with the Transiter server in two ways:

*HTTP API*: after launching the Transiter server, 
  Transiter exports a public HTTP API on port 8080.
  Let's see what it says:

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

GOT TO HERE

To begin, we're going to launch Transiter.
The easiest way to do this is with Docker compose and the
    [standard Transiter compose config file](https://github.com/jamespfennell/transiter/blob/master/docker/docker-compose.yml).
Simply run,

    docker-compose up -f path/to/docker-compose.yml

It will take a minute for the images to download from Docker Hub and for the containers to be launched successfully.

When everything is launched, Transiter will be listening on port 8000.
If you navigate to `localhost:8000` in your web browser (or use `curl`), you will find the Transiter landing page,

```json
{
  "transiter": {
    "version": "0.4.5",
    "href": "https://github.com/jamespfennell/transiter",
    "docs": {
      "href": "https://demo.transiter.io/docs/"
    }
  },
  "systems": {
    "count": 0,
    "href": "https://demo.transiter.io/systems"
  }
}
```

As you can see, there are no (transit) systems installed, and the next step is to install one!

??? info "Running Transiter without Docker"
    It's possible to run Transiter on "bare metal" without Docker; 
    the [running Transiter](deployment/running-transiter.md) page details how.
    It's quite a bit more work though, so for getting started we recommend the Docker approach.

??? info "Building the Docker images locally"
    If you want to build the Docker images locally that's easy, too:
    just check out the [Transiter Git repository](https://github.com/jamespfennell/transiter)
    and in the root of repository run `make docker`.
   
## Install a system

Each deployment of Transiter can have multiple transit systems installed side-by-side.
A transit system is installed using a YAML configuration file that 
    contains basic metadata about the system (like its name),
    the URLs of the data feeds,
    and how to parse those data feeds (GTFS Static, GTFS Realtime, or a custom format).

For this tour, we're going to start by installing the BART system in San Francisco.
The YAML configuration file is stored in Github, you can [inspect it here](https://github.com/jamespfennell/transiter-sfbart/blob/master/Transiter-SF-BART-config.yaml).
The system in installed by sending a `PUT` HTTP request to the desired system ID.
In this case we'll install the system using ID `bart`,

```
curl -X PUT "localhost:8000/systems/bart?sync=true" \
     -F 'config_file=https://raw.githubusercontent.com/jamespfennell/transiter-sfbart/master/Transiter-SF-BART-config.yaml'
```

As you can see, we've provided a `config_file` form parameter that contains the URL of the config file.
It's also possible to provide the config as a file upload using the same `config_file` form parameter.

The request will take a few seconds to complete;
    most of the time is spent loading the BART's schedule into the database.
After it finishes, hit the Transiter landing page again to get,

```json
{
  "transiter": {
    "version": "0.4.5",
    "href": "https://github.com/jamespfennell/transiter",
    "docs": {
      "href": "https://demo.transiter.io/docs/"
    }
  },
  "systems": {
    "count": 1,
    "href": "http://localhost:8000/systems"
  }
}
```

It's installed! 
Next, navigate to the list systems endpoint.
The URL `http://localhost:8000/systems` is helpfully given in the JSON response.
We get,

```json
[
  {
    "id": "bart",
    "status": "ACTIVE",
    "name": "San Francisco BART",
    "href": "http://localhost:8000/systems/bart"
  }
]
```

Now navigating to the system itself, we get,


```json
{
  "id": "bart",
  "status": "ACTIVE",
  "name": "San Francisco BART",
  "agencies": {
    "count": 1,
    "href": "http://localhost:8000/systems/bart/agencies"
  },
  "feeds": {
    "count": 3,
    "href": "http://localhost:8000/systems/bart/feeds"
  },
  "routes": {
    "count": 14,
    "href": "http://localhost:8000/systems/bart/routes"
  },
  "stops": {
    "count": 177,
    "href": "http://localhost:8000/systems/bart/stops"
  },
  "transfers": {
    "count": 0,
    "href": "http://localhost:8000/systems/bart/transfers"
  }
}
```

This is an overview of the system, showing the number of various things like stops and routes,
as well as URLs for those.

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

- Create a transit system config for a system you're interested in.

- Consult the API reference to find other endpoints and data that Transiter exposes.
