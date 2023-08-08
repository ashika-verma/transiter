

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
