

-- name: InsertScheduledService :one
INSERT INTO scheduled_service 
    (id, system_pk, feed_pk,
    monday, tuesday, wednesday, thursday, friday, saturday, sunday, end_date, start_date)
VALUES
    (sqlc.arg(id),
    sqlc.arg(system_pk),
    sqlc.arg(feed_pk),
    sqlc.arg(monday),
    sqlc.arg(tuesday),
    sqlc.arg(wednesday),
    sqlc.arg(thursday),
    sqlc.arg(friday),
    sqlc.arg(saturday),
    sqlc.arg(sunday),
    sqlc.arg(end_date),
    sqlc.arg(start_date))
RETURNING pk;

-- name: DeleteScheduledServices :exec
DELETE FROM scheduled_service WHERE feed_pk = sqlc.arg(feed_pk);

-- name: InsertScheduledTrip :one
INSERT INTO scheduled_trip 
    (id, route_pk, service_pk, direction_id, bikes_allowed, block_id, headsign,
    short_name, wheelchair_accessible)
VALUES
    (sqlc.arg(id),
    sqlc.arg(route_pk),
    sqlc.arg(service_pk),
    sqlc.arg(direction_id),
    sqlc.arg(bikes_allowed),
    sqlc.arg(block_id),
    sqlc.arg(headsign),
    sqlc.arg(short_name),
    sqlc.arg(wheelchair_accessible))
RETURNING pk;

-- name: InsertScheduledTripStopTime :copyfrom
INSERT INTO scheduled_trip_stop_time
    (trip_pk ,
    stop_pk ,
    arrival_time ,
    departure_time ,
    stop_sequence ,
    continuous_drop_off ,
    continuous_pickup,
    drop_off_type,
    exact_times ,
    headsign ,
    pickup_type ,
    shape_distance_traveled)
VALUES
    (sqlc.arg(trip_pk),
    sqlc.arg(stop_pk),
    sqlc.arg(arrival_time),
    sqlc.arg(departure_time),
    sqlc.arg(stop_sequence),
    sqlc.arg(continuous_drop_off),
    sqlc.arg(continuous_pickup),
    sqlc.arg(drop_off_type),
    sqlc.arg(exact_times),
    sqlc.arg(headsign),
    sqlc.arg(pickup_type) ,
    sqlc.arg(shape_distance_traveled));


-- name: DropScheduledTripStopTimeIndexes :one
DROP INDEX ix_scheduled_trip_stop_time__stop_pk, ix_scheduled_trip_stop_time__trip_pk;

-- name: DropScheduledTripStopTimeConstraints :exec
ALTER TABLE scheduled_trip_stop_time
    DROP CONSTRAINT fk_scheduled_trip_stop_time_stop_pk,
    DROP CONSTRAINT fk_scheduled_trip_stop_time_trip_pk,
    DROP CONSTRAINT scheduled_trip_stop_time_trip_pk_stop_sequence_key;

-- CREATE INDEX ix_scheduled_trip_stop_time__stop_pk ON scheduled_trip_stop_time USING btree (stop_pk);