

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
