import unittest.mock as mock
import unittest
import itertools
from transiter.services import routeservice
import datetime

class TestRouteService(unittest.TestCase):

    SYSTEM_ID = '1'
    ROUTE_ONE_ID = '3'
    ROUTE_ONE_HREF = '30'
    ROUTE_ONE_REPR = {'route_id': ROUTE_ONE_ID}
    ROUTE_ONE_STATUS = 'Bad'
    ROUTE_TWO_STATUS = 'Good'
    ROUTE_TWO_HREF = '40'
    ROUTE_TWO_REPR = {'route_id': '4'}
    STOP_REPR = {'stop_id': '5'}
    STOP_EVENT_REPR = {'track': 'Track Two'}
    GOOD_STATUS = "Good Service"
    OKAY_STATUS = "Okay"
    BAD_STATUS = "Quite Bad"
    REALLY_BAD_STATUS = "Really Bad"

    @classmethod
    def setUp(cls):
        cls.route_one = mock.MagicMock()
        cls.route_one.short_repr.return_value = cls.ROUTE_ONE_REPR.copy()
        cls.route_one.pk = 500
        cls.route_one.short_repr.return_value = cls.ROUTE_ONE_REPR.copy()

        cls.route_two = mock.MagicMock()
        cls.route_two.pk = 501
        cls.route_two.short_repr.return_value = cls.ROUTE_TWO_REPR


        """
        stop = mock.MagicMock()
        stop.repr_for_list.return_value = cls.STOP_REPR

        stop_event = mock.MagicMock()
        stop_event.repr_for_list.return_value = cls.STOP_EVENT_REPR
        stop_event.stop = stop
        cls.trip_one.stop_events = [stop_event]
        """

    def _construct_status_mock(self, route):
        if route == self.route_one:
            return self.ROUTE_ONE_STATUS
        elif route == self.route_two:
            return self.ROUTE_TWO_STATUS
        else:
            self.fail('Unwanted interaction with _construct_status')

    @mock.patch('transiter.services.routeservice.linksutil')
    @mock.patch('transiter.services.routeservice.routedam')
    def test_list_all_in_system(self, route_dao, linksutil):
        """[Route service] Listing all routes in a system"""

        def RouteEntityLink(system):
            if system == self.route_one:
                return self.ROUTE_ONE_HREF
            if system == self.route_two:
                return self.ROUTE_TWO_HREF
        linksutil.RouteEntityLink.side_effect = RouteEntityLink

        expected = [{
            'service_status': self.ROUTE_ONE_STATUS,
            'href': self.ROUTE_ONE_HREF,
            **self.ROUTE_ONE_REPR
        },{
            'service_status': self.ROUTE_TWO_STATUS,
            'href': self.ROUTE_TWO_HREF,
            **self.ROUTE_TWO_REPR
        }]
        route_dao.list_all_in_system.return_value = [self.route_one,
                                                     self.route_two]
        route_dao.list_route_statuses.return_value = {
            500: self.ROUTE_ONE_STATUS,
            501: self.ROUTE_TWO_STATUS,
        }

        actual = routeservice.list_all_in_system(self.SYSTEM_ID)

        self.assertEqual(actual, expected)
        route_dao.list_all_in_system.assert_called_once_with(self.SYSTEM_ID)
        self.route_one.short_repr.assert_called_once()
        self.route_two.short_repr.assert_called_once()

    @mock.patch('transiter.services.routeservice.routedam')
    def test_get_in_system_by_id(self, route_dao):
        """[Route service] Getting a specific route in a system"""
        route_dao.get_in_system_by_id.return_value = self.route_one
        sp_vertex = mock.MagicMock()
        self.route_one.default_service_pattern.vertices = [sp_vertex]

        actual = routeservice.get_in_system_by_id(
            self.SYSTEM_ID,
            self.ROUTE_ONE_ID)

        #self.assertDictEqual(actual, expected)
        route_dao.get_in_system_by_id.assert_called_once_with(
            self.SYSTEM_ID,
            self.ROUTE_ONE_ID)




