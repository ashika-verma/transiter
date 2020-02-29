from transiter import models
from transiter.data.dams import systemdam


def test_create(db_session):
    new_system = systemdam.create()
    new_system.id = "new_id"
    new_system.status = models.System.SystemStatus.ACTIVE
    db_session.flush()

    assert new_system == systemdam.get_by_id("new_id")


def test_delete_by_id(system_1):
    system_1_id = system_1.id

    result = systemdam.delete_by_id(system_1_id)

    assert True is result
    assert None is systemdam.get_by_id(system_1_id)


def test_delete_by_id__invalid_id(system_1):
    result = systemdam.delete_by_id("unknown_id")

    assert False is result


def test_list_all(system_1, system_2):
    assert [system_1, system_2] == systemdam.list_all()


def test_list_all__no_systems(db_session):
    assert [] == systemdam.list_all()


def test_get_by_id(system_1):
    assert system_1 == systemdam.get_by_id(system_1.id)


def test_get_by_id__unknown(system_1):
    assert None is systemdam.get_by_id("unknown_id")


def test_get_by_id__not_active(installing_system):
    assert None is systemdam.get_by_id(installing_system.id, only_return_active=True)


def test_count_stops_in_system(system_1, stop_1_1, stop_1_2, stop_2_1):
    count = systemdam.count_stops_in_system(system_1.id)

    assert 2 == count


def test_count_stops_in_system__no_stops(system_1, stop_2_1):
    count = systemdam.count_stops_in_system(system_1.id)

    assert 0 == count


def test_count_routes_in_system(system_1, route_1_1, route_1_2, route_1_3, route_2_1):
    count = systemdam.count_routes_in_system(system_1.id)

    assert 3 == count


def test_count_routes_in_system__no_routes(system_1, route_2_1):
    count = systemdam.count_routes_in_system(system_1.id)

    assert 0 == count


def test_count_feeds_in_system(system_1, feed_1_1, feed_1_2, feed_2_1):
    count = systemdam.count_feeds_in_system(system_1.id)

    assert 2 == count


def test_count_feeds_in_system__no_feeds(system_1, feed_2_1):
    count = systemdam.count_routes_in_system(system_1.id)

    assert 0 == count


def test_list_all_alerts_associated_to_system(
    db_session, add_model, system_1, system_2
):
    alert_1 = add_model(models.Alert(pk=703))
    alert_2 = add_model(models.Alert(pk=704))
    alert_1.system_pk = system_1.pk
    alert_2.system_pk = system_2.pk
    db_session.flush()

    assert [alert_1] == systemdam.list_all_alerts_associated_to_system(system_1.pk)


def test_list_all_alerts_associated_to_system__no_alerts(system_1, system_2):
    assert [] == systemdam.list_all_alerts_associated_to_system(system_1.pk)
