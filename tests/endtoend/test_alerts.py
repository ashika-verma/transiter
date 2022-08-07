import datetime
import time

import pytest
import requests
from . import gtfs_realtime_pb2 as gtfs

ONE_DAY_IN_SECONDS = 60 * 60 * 24
TIME_1 = datetime.datetime.utcfromtimestamp(time.time() - ONE_DAY_IN_SECONDS)
TIME_2 = datetime.datetime.utcfromtimestamp(time.time() + ONE_DAY_IN_SECONDS)

ALERT_SMALL_JSON = [{"id": "alert_id", "cause": "STRIKE", "effect": "MODIFIED_SERVICE"}]

ALERT_LARGE_JSON = [
    dict(
        **ALERT_SMALL_JSON[0],
        **{
            "currentActivePeriod": {
                "startsAt": str(int(TIME_1.timestamp())),
                "endsAt": str(int(TIME_2.timestamp())),
            },
            "allActivePeriods": [{
                "startsAt": str(int(TIME_1.timestamp())),
                "endsAt": str(int(TIME_2.timestamp())),
            }],
            "header": [
                {
                    "text": "Advertencia",
                    "language": "es",
                }
            ], 
            "description": [
                {
                    "text": "Description",
                    "language": "en",
                }
            ],
            "url": [
                {
                    "text": "URL",
                    "language": "en",
                }
            ],
        }
    )
]


def setup_test(
    system_id, install_system_1, transiter_host, source_server
):

    __, realtime_feed_url = install_system_1(system_id)

    message = gtfs.FeedMessage(
        header=gtfs.FeedHeader(gtfs_realtime_version="2.0", timestamp=int(time.time())),
        entity=[
            gtfs.FeedEntity(
                id="trip_id",
                trip_update=gtfs.TripUpdate(
                    trip=gtfs.TripDescriptor(trip_id="trip_id", route_id="A")
                ),
            ),
            gtfs.FeedEntity(
                id="alert_id",
                alert=gtfs.Alert(
                    active_period=[
                        gtfs.TimeRange(
                            start=int(TIME_1.timestamp()), end=int(TIME_2.timestamp()),
                        )
                    ],
                    header_text=gtfs.TranslatedString(
                        translation=[
                            gtfs.TranslatedString.Translation(
                                text="Advertencia", language="es"
                            )
                        ],
                    ),
                    description_text=gtfs.TranslatedString(
                        translation=[
                            gtfs.TranslatedString.Translation(
                                text="Description", language="en"
                            )
                        ],
                    ),
                    url=gtfs.TranslatedString(
                        translation=[
                            gtfs.TranslatedString.Translation(
                                text="URL", language="en"
                            )
                        ],
                    ),
                    informed_entity=[
                        gtfs.EntitySelector(agency_id="AgencyId"),
                        gtfs.EntitySelector(route_id="A"),
                        gtfs.EntitySelector(stop_id="1A"),
                        gtfs.EntitySelector(trip=gtfs.TripDescriptor(trip_id="trip_id")),
                    ],
                    cause=gtfs.Alert.Cause.STRIKE,
                    effect=gtfs.Alert.Effect.MODIFIED_SERVICE,
                ),
            ),
        ],
    )

    source_server.put(realtime_feed_url, message.SerializeToString())
    requests.post(
        "{}/admin/systems/{}/feeds/GtfsRealtimeFeed?sync=true".format(
            transiter_host, system_id
        )
    ).raise_for_status()

@pytest.mark.parametrize(
    "path,entity_id",
    [
        ["routes", "A"],
        ["stops", "1A"],
        # TODO: renable
        # [
        #     "routes/A/trips",
        #     "trip_id",
        #     ,
        #     ALERT_SMALL_JSON,
        # ],
        [
            "agencies",
            "AgencyId",
        ],
    ],
)
def test_alerts_list_informed_entities(
    install_system_1,
    transiter_host,
    source_server,
    path,
    entity_id,
):
    system_id = "test_alerts__list_informed_entities_" + str(hash(path))
    setup_test(
        system_id=system_id,
        install_system_1=install_system_1,
        transiter_host=transiter_host,
        source_server=source_server,
    )

    url = "{}/systems/{}/{}".format(transiter_host, system_id, path)

    actual_data = requests.get(url).json()
    actual_data = actual_data[path]
    actual_data = {response["id"]: response for response in actual_data}[entity_id]

    assert ALERT_SMALL_JSON == actual_data["alerts"]


@pytest.mark.parametrize(
    "path",
    [
        "routes/A",
        "stops/1A",
        # TODO: renable
        # [
        #     "routes/A/trips/trip_id",
        #     None,
        # ],
        "stops/1A",
        "agencies/AgencyId",
    ],
)
def test_alerts_get_informed_entity(
    install_system_1,
    transiter_host,
    source_server,
    path,
):
    system_id = "test_alerts__get_informed_entity_" + str(hash(path))
    setup_test(
        system_id=system_id,
        install_system_1=install_system_1,
        transiter_host=transiter_host,
        source_server=source_server,
    )

    url = "{}/systems/{}/{}".format(transiter_host, system_id, path)

    actual_data = requests.get(url).json()

    expected_json = ALERT_SMALL_JSON
    assert expected_json == actual_data["alerts"]


@pytest.mark.parametrize(
    "path,entity_id",
    [
        [
            "alerts",
            "alert_id",
        ],
        [
            "alerts/alert_id",
            None,
        ],
    ],
)
def test_alerts_get_or_list_alert(
    install_system_1,
    transiter_host,
    source_server,
    path,
    entity_id,
):
    system_id = "test_alerts__get_or_list_alert_" + str(hash(path))
    setup_test(
        system_id=system_id,
        install_system_1=install_system_1,
        transiter_host=transiter_host,
        source_server=source_server,
    )

    url = "{}/systems/{}/{}".format(transiter_host, system_id, path)

    actual_data = requests.get(url).json()

    if entity_id is not None:
        actual_data = actual_data[path]
        actual_data = {response["id"]: response for response in actual_data}[entity_id]

    expected_json = ALERT_LARGE_JSON[0]
    assert expected_json == actual_data
