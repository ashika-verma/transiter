// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: schedule_queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteScheduledServices = `-- name: DeleteScheduledServices :exec
DELETE FROM scheduled_service WHERE feed_pk = $1
`

func (q *Queries) DeleteScheduledServices(ctx context.Context, feedPk int64) error {
	_, err := q.db.Exec(ctx, deleteScheduledServices, feedPk)
	return err
}

const insertScheduledService = `-- name: InsertScheduledService :one
INSERT INTO scheduled_service 
    (id, system_pk, feed_pk,
    monday, tuesday, wednesday, thursday, friday, saturday, sunday, end_date, start_date)
VALUES
    ($1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12)
RETURNING pk
`

type InsertScheduledServiceParams struct {
	ID        string
	SystemPk  int64
	FeedPk    int64
	Monday    bool
	Tuesday   bool
	Wednesday bool
	Thursday  bool
	Friday    bool
	Saturday  bool
	Sunday    bool
	EndDate   pgtype.Date
	StartDate pgtype.Date
}

func (q *Queries) InsertScheduledService(ctx context.Context, arg InsertScheduledServiceParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertScheduledService,
		arg.ID,
		arg.SystemPk,
		arg.FeedPk,
		arg.Monday,
		arg.Tuesday,
		arg.Wednesday,
		arg.Thursday,
		arg.Friday,
		arg.Saturday,
		arg.Sunday,
		arg.EndDate,
		arg.StartDate,
	)
	var pk int64
	err := row.Scan(&pk)
	return pk, err
}

const insertScheduledTrip = `-- name: InsertScheduledTrip :one
INSERT INTO scheduled_trip 
    (id, route_pk, service_pk, direction_id, bikes_allowed, block_id, headsign,
    short_name, wheelchair_accessible)
VALUES
    ($1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9)
RETURNING pk
`

type InsertScheduledTripParams struct {
	ID                   string
	RoutePk              int64
	ServicePk            int64
	DirectionID          pgtype.Bool
	BikesAllowed         string
	BlockID              pgtype.Text
	Headsign             pgtype.Text
	ShortName            pgtype.Text
	WheelchairAccessible string
}

func (q *Queries) InsertScheduledTrip(ctx context.Context, arg InsertScheduledTripParams) (int64, error) {
	row := q.db.QueryRow(ctx, insertScheduledTrip,
		arg.ID,
		arg.RoutePk,
		arg.ServicePk,
		arg.DirectionID,
		arg.BikesAllowed,
		arg.BlockID,
		arg.Headsign,
		arg.ShortName,
		arg.WheelchairAccessible,
	)
	var pk int64
	err := row.Scan(&pk)
	return pk, err
}

type InsertScheduledTripStopTimeParams struct {
	TripPk                int64
	StopPk                int64
	ArrivalTime           pgtype.Int4
	DepartureTime         pgtype.Int4
	StopSequence          int32
	ContinuousDropOff     int16
	ContinuousPickup      int16
	DropOffType           int16
	ExactTimes            bool
	Headsign              pgtype.Text
	PickupType            int16
	ShapeDistanceTraveled pgtype.Float8
}
