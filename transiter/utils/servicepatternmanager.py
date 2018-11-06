"""
ALGO

USE GTFS STATIC! It does everything!

Step 3: for each route, pass each RawTrip through the matchers for that route
    to get a set of RawTrips for that matchers. By defining a hash function
    that only takes the list of stops, we eliminate duplicates automatically!
"""

from transiter.utils import graphutils
from transiter.database import models
"""
SELECT stops.* 
FROM stops
INNER JOIN service_pattern_vertices sp_v ON sp_v.stop_pri_key = stops.id
INNER JOIN service_patterns sp ON sp_v.service_pattern_pri_key = sp.id
INNER JOIN routes ON sp.id = routes.regular_service_pattern_pri_key
WHERE  routes.route_id = 'A'
ORDER BY sp_v.position;
"""

def construct_sps_from_gtfs_static_data(
        gtfs_static_parser,
        route_sp_settings=[],
        general_sp_settings=None):
    route_id_to_route = gtfs_static_parser.route_id_to_route
    stop_id_to_stop = gtfs_static_parser.stop_id_to_stop
    trips = gtfs_static_parser.trip_id_to_trip.values()

    route_id_to_trips = {
        route_id: set() for route_id in route_id_to_route.keys()}
    for trip in trips:
        # TODO: what happens if this is not set?
        if trip.direction_id:
            trip.reverse()

        if trip.route_id in route_id_to_trips:
            route_id_to_trips[trip.route_id].add(trip)

    # TODO: invert the for loops here for an easy optimization
    # Better: split off a new method construct from trips and settings
    for route_id, trips in route_id_to_trips.items():
        route = route_id_to_route[route_id]
        for sp_setting in route_sp_settings:
            name = sp_setting.get('name', None)
            default = sp_setting.get('default', False)
            regular = sp_setting.get('regular', False)
            threshold = sp_setting.get('threshold', 0)
            conditions = sp_setting.get('conditions', None)

            if conditions is not None:
                sp_trips = _filter_trips_by_conditions(trips, threshold, conditions)
            else:
                sp_trips = trips

            service_pattern = _construct_for_static_trips(sp_trips, stop_id_to_stop)
            service_pattern.name = name
            service_pattern.route = route
            if default:
                route.default_service_pattern = service_pattern
            if regular:
                route.regular_service_pattern = service_pattern


# TODO rename this method
def _construct_for_static_trips(trips, stop_id_to_stop):
    path_lists = set()
    for trip in trips:
        if len(trip.stop_ids) == 0:
            continue
        path_lists.add(tuple(trip.stop_ids))
    sorted_graph = _path_lists_to_sorted_graph(path_lists)
    service_pattern = _sorted_graph_to_service_pattern(sorted_graph, stop_id_to_stop)
    return service_pattern


def _sorted_graph_to_service_pattern(sorted_graph, stop_id_to_stop):
    service_pattern = models.ServicePattern()
    for index, vertex in enumerate(sorted_graph.vertices()):
        stop_id = vertex.label
        stop = stop_id_to_stop[stop_id]
        sp_vertex = models.ServicePatternVertex()
        sp_vertex.stop = stop
        sp_vertex.service_pattern = service_pattern
        sp_vertex.position = index
    return service_pattern


def _path_lists_to_sorted_graph(path_lists):
    if len(path_lists) == 0:
        return graphutils.graphdatastructs.DirectedPath([])
    if len(path_lists) == 1:
        unique_element = next(iter(path_lists))
        return graphutils.graphdatastructs.DirectedPath(unique_element)
    paths = [
        graphutils.graphdatastructs.DirectedPath(path_list) for
        path_list in path_lists
    ]
    graph = graphutils.pathstitcher.stitch(paths)
    # short circuit if the route_graph is actually a path
    if graph.is_path():
        return graph.cast_to_path()
    return graphutils.topologicalsort.sort(graph)


def _filter_trips_by_conditions(trips, threshold, matching_conditions):
    trip_matcher = _TripMatcher(matching_conditions)
    stop_ids_to_trips = {}
    total_count = 0
    for trip in trips:
        if not trip_matcher.match(trip):
            continue
        total_count += 1
        stop_ids = tuple(trip.stop_ids)
        stop_ids_to_trips.setdefault(stop_ids, [])
        stop_ids_to_trips[stop_ids].append(trip)

    filtered_trips = []
    for stop_ids, grouped_trips in stop_ids_to_trips.items():
        if len(grouped_trips) >= threshold * total_count:
            filtered_trips += grouped_trips
    return filtered_trips


class _TripMatcher:

    def __init__(self, raw_conds):
        self._primary_conditions = self._convert_raw_conditions(raw_conds)

    def match(self, trip):
        for condition in self._primary_conditions:
            if not condition(trip):
                return False
        return True

    @classmethod
    def _convert_raw_conditions(cls, dictionary):
        conditions = []
        for key, value in dictionary.items():
            conditions.append(cls._convert_raw_condition(key, value))
        return conditions

    _logical_operators = None
    _key_to_function = None
    _key_to_extra_args = None

    @classmethod
    def _convert_raw_condition(cls, key, value):
        # Having these variables as class variables that are populated on the
        # first run is basically a form of caching
        if cls._logical_operators is None:
            cls._logical_operators = {'one_of', 'all_of', 'none_of'}
            cls._key_to_function = {
                'one_of': cls.one_of_factory,
                'all_of': cls.all_of_factory,
                'none_of': cls.none_of_factory,
                'starts_earlier_than': cls.order_factory,
                'starts_later_than': cls.order_factory,
                'ends_earlier_than': cls.order_factory,
                'ends_later_than': cls.order_factory,
                'route_id': cls.equality_factory,
                'weekend': cls.weekend_factory,
                'weekday': cls.weekday_factory,
            }
            cls._key_to_extra_args = {
                'starts_earlier_than': ('start_time', True),
                'starts_later_than': ('start_time', False),
                'ends_earlier_than': ('end_time', True),
                'ends_later_than': ('end_time', False),
                'route_id': ('route_id',),
            }
        if key in cls._logical_operators:
            value = cls._convert_raw_conditions(value)
        try:
            func = cls._key_to_function[key]
        except KeyError:
            raise NotImplementedError
        extra_args = cls._key_to_extra_args.get(key, ())
        return func(value, *extra_args)

    @staticmethod
    def one_of_factory(conditions):
        def one_of(trip):
            for condition in conditions:
                if condition(trip):
                    return True
            return False
        return one_of

    @staticmethod
    def none_of_factory(conditions):
        def none_of(trip):
            for condition in conditions:
                if condition(trip):
                    return False
            return True
        return none_of

    @staticmethod
    def all_of_factory(conditions):
        def all_of(trip):
            for condition in conditions:
                if not condition(trip):
                    return False
            return True
        return all_of

    @staticmethod
    def order_factory(value, trip_attr, less_than=True):
        def order(trip):
            return (getattr(trip, trip_attr) < value) == less_than
        return order

    @staticmethod
    def equality_factory(value, trip_attr):
        def equality(trip):
            return getattr(trip, trip_attr) == value

        def contains(trip):
            return getattr(trip, trip_attr) in value

        if isinstance(value, list):
            return contains
        return equality

    @staticmethod
    def weekday_factory(value):
        def weekday(trip):
            weekday_cond = trip.monday or trip.tuesday or trip.wednesday or trip.thursday or trip.friday
            weekend_cond = not (trip.saturday or trip.sunday)
            return (weekday_cond and weekend_cond) == value
        return weekday

    @classmethod
    def weekend_factory(cls, value):
        return cls.weekday_factory(not value)





