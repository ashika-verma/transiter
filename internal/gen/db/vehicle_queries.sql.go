// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: vehicle_queries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deleteStaleVehicles = `-- name: DeleteStaleVehicles :exec
DELETE FROM vehicle
WHERE
  feed_pk = $1
  AND NOT id = ANY($2::text[])
`

type DeleteStaleVehiclesParams struct {
	FeedPk           int64
	ActiveVehicleIds []string
}

func (q *Queries) DeleteStaleVehicles(ctx context.Context, arg DeleteStaleVehiclesParams) error {
	_, err := q.db.Exec(ctx, deleteStaleVehicles, arg.FeedPk, arg.ActiveVehicleIds)
	return err
}

const getVehicle = `-- name: GetVehicle :one
SELECT vehicle.pk, vehicle.id, vehicle.system_pk, vehicle.trip_pk, vehicle.label, vehicle.license_plate, vehicle.current_status, vehicle.latitude, vehicle.longitude, vehicle.bearing, vehicle.odometer, vehicle.speed, vehicle.congestion_level, vehicle.updated_at, vehicle.current_stop_pk, vehicle.current_stop_sequence, vehicle.occupancy_status, vehicle.feed_pk, vehicle.occupancy_percentage,
       stop.id as stop_id,
       stop.name as stop_name,
       trip.id as trip_id,
       trip.direction_id as trip_direction_id,
       route.id as route_id,
       route.color as route_color
FROM vehicle
LEFT JOIN stop ON vehicle.current_stop_pk = stop.pk
LEFT JOIN trip ON vehicle.trip_pk = trip.pk
LEFT JOIN route ON trip.route_pk = route.pk
WHERE vehicle.system_pk = $1 AND vehicle.id = $2
`

type GetVehicleParams struct {
	SystemPk  int64
	VehicleID pgtype.Text
}

type GetVehicleRow struct {
	Pk                  int64
	ID                  pgtype.Text
	SystemPk            int64
	TripPk              pgtype.Int8
	Label               pgtype.Text
	LicensePlate        pgtype.Text
	CurrentStatus       pgtype.Text
	Latitude            pgtype.Numeric
	Longitude           pgtype.Numeric
	Bearing             pgtype.Float4
	Odometer            pgtype.Float8
	Speed               pgtype.Float4
	CongestionLevel     string
	UpdatedAt           pgtype.Timestamptz
	CurrentStopPk       pgtype.Int8
	CurrentStopSequence pgtype.Int4
	OccupancyStatus     pgtype.Text
	FeedPk              int64
	OccupancyPercentage pgtype.Int4
	StopID              pgtype.Text
	StopName            pgtype.Text
	TripID              pgtype.Text
	TripDirectionID     pgtype.Bool
	RouteID             pgtype.Text
	RouteColor          pgtype.Text
}

func (q *Queries) GetVehicle(ctx context.Context, arg GetVehicleParams) (GetVehicleRow, error) {
	row := q.db.QueryRow(ctx, getVehicle, arg.SystemPk, arg.VehicleID)
	var i GetVehicleRow
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.SystemPk,
		&i.TripPk,
		&i.Label,
		&i.LicensePlate,
		&i.CurrentStatus,
		&i.Latitude,
		&i.Longitude,
		&i.Bearing,
		&i.Odometer,
		&i.Speed,
		&i.CongestionLevel,
		&i.UpdatedAt,
		&i.CurrentStopPk,
		&i.CurrentStopSequence,
		&i.OccupancyStatus,
		&i.FeedPk,
		&i.OccupancyPercentage,
		&i.StopID,
		&i.StopName,
		&i.TripID,
		&i.TripDirectionID,
		&i.RouteID,
		&i.RouteColor,
	)
	return i, err
}

const insertVehicle = `-- name: InsertVehicle :exec
INSERT INTO vehicle
    (id, system_pk, trip_pk, label, license_plate, current_status, latitude, longitude, bearing, odometer, speed, congestion_level, updated_at, current_stop_pk, current_stop_sequence, occupancy_status, feed_pk, occupancy_percentage)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
`

type InsertVehicleParams struct {
	ID                  pgtype.Text
	SystemPk            int64
	TripPk              pgtype.Int8
	Label               pgtype.Text
	LicensePlate        pgtype.Text
	CurrentStatus       pgtype.Text
	Latitude            pgtype.Numeric
	Longitude           pgtype.Numeric
	Bearing             pgtype.Float4
	Odometer            pgtype.Float8
	Speed               pgtype.Float4
	CongestionLevel     string
	UpdatedAt           pgtype.Timestamptz
	CurrentStopPk       pgtype.Int8
	CurrentStopSequence pgtype.Int4
	OccupancyStatus     pgtype.Text
	FeedPk              int64
	OccupancyPercentage pgtype.Int4
}

func (q *Queries) InsertVehicle(ctx context.Context, arg InsertVehicleParams) error {
	_, err := q.db.Exec(ctx, insertVehicle,
		arg.ID,
		arg.SystemPk,
		arg.TripPk,
		arg.Label,
		arg.LicensePlate,
		arg.CurrentStatus,
		arg.Latitude,
		arg.Longitude,
		arg.Bearing,
		arg.Odometer,
		arg.Speed,
		arg.CongestionLevel,
		arg.UpdatedAt,
		arg.CurrentStopPk,
		arg.CurrentStopSequence,
		arg.OccupancyStatus,
		arg.FeedPk,
		arg.OccupancyPercentage,
	)
	return err
}

const listVehicleUniqueColumns = `-- name: ListVehicleUniqueColumns :many
SELECT id, pk, trip_pk FROM vehicle
WHERE id = ANY($1::text[])
AND system_pk = $2
`

type ListVehicleUniqueColumnsParams struct {
	VehicleIds []string
	SystemPk   int64
}

type ListVehicleUniqueColumnsRow struct {
	ID     pgtype.Text
	Pk     int64
	TripPk pgtype.Int8
}

func (q *Queries) ListVehicleUniqueColumns(ctx context.Context, arg ListVehicleUniqueColumnsParams) ([]ListVehicleUniqueColumnsRow, error) {
	rows, err := q.db.Query(ctx, listVehicleUniqueColumns, arg.VehicleIds, arg.SystemPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListVehicleUniqueColumnsRow
	for rows.Next() {
		var i ListVehicleUniqueColumnsRow
		if err := rows.Scan(&i.ID, &i.Pk, &i.TripPk); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listVehicles = `-- name: ListVehicles :many
SELECT vehicle.pk, vehicle.id, vehicle.system_pk, vehicle.trip_pk, vehicle.label, vehicle.license_plate, vehicle.current_status, vehicle.latitude, vehicle.longitude, vehicle.bearing, vehicle.odometer, vehicle.speed, vehicle.congestion_level, vehicle.updated_at, vehicle.current_stop_pk, vehicle.current_stop_sequence, vehicle.occupancy_status, vehicle.feed_pk, vehicle.occupancy_percentage,
       stop.id as stop_id,
       stop.name as stop_name,
       trip.id as trip_id,
       trip.direction_id as trip_direction_id,
       route.id as route_id,
       route.color as route_color
FROM vehicle
LEFT JOIN stop ON vehicle.current_stop_pk = stop.pk
LEFT JOIN trip ON vehicle.trip_pk = trip.pk
LEFT JOIN route ON trip.route_pk = route.pk
WHERE vehicle.system_pk = $1
  AND vehicle.id >= $2
  AND (
    NOT $3::bool OR
    vehicle.id = ANY($4::text[])
  )
ORDER BY vehicle.id
LIMIT $5
`

type ListVehiclesParams struct {
	SystemPk               int64
	FirstVehicleID         pgtype.Text
	OnlyReturnSpecifiedIds bool
	VehicleIds             []string
	NumVehicles            int32
}

type ListVehiclesRow struct {
	Pk                  int64
	ID                  pgtype.Text
	SystemPk            int64
	TripPk              pgtype.Int8
	Label               pgtype.Text
	LicensePlate        pgtype.Text
	CurrentStatus       pgtype.Text
	Latitude            pgtype.Numeric
	Longitude           pgtype.Numeric
	Bearing             pgtype.Float4
	Odometer            pgtype.Float8
	Speed               pgtype.Float4
	CongestionLevel     string
	UpdatedAt           pgtype.Timestamptz
	CurrentStopPk       pgtype.Int8
	CurrentStopSequence pgtype.Int4
	OccupancyStatus     pgtype.Text
	FeedPk              int64
	OccupancyPercentage pgtype.Int4
	StopID              pgtype.Text
	StopName            pgtype.Text
	TripID              pgtype.Text
	TripDirectionID     pgtype.Bool
	RouteID             pgtype.Text
	RouteColor          pgtype.Text
}

func (q *Queries) ListVehicles(ctx context.Context, arg ListVehiclesParams) ([]ListVehiclesRow, error) {
	rows, err := q.db.Query(ctx, listVehicles,
		arg.SystemPk,
		arg.FirstVehicleID,
		arg.OnlyReturnSpecifiedIds,
		arg.VehicleIds,
		arg.NumVehicles,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListVehiclesRow
	for rows.Next() {
		var i ListVehiclesRow
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.TripPk,
			&i.Label,
			&i.LicensePlate,
			&i.CurrentStatus,
			&i.Latitude,
			&i.Longitude,
			&i.Bearing,
			&i.Odometer,
			&i.Speed,
			&i.CongestionLevel,
			&i.UpdatedAt,
			&i.CurrentStopPk,
			&i.CurrentStopSequence,
			&i.OccupancyStatus,
			&i.FeedPk,
			&i.OccupancyPercentage,
			&i.StopID,
			&i.StopName,
			&i.TripID,
			&i.TripDirectionID,
			&i.RouteID,
			&i.RouteColor,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listVehicles_Geographic = `-- name: ListVehicles_Geographic :many
WITH distance AS (
  SELECT
  pk vehicle_pk,
  (6371 * acos(cos(radians(latitude)) * cos(radians($3::numeric)) * cos(radians($4::numeric) - radians(longitude)) + sin(radians(latitude)) * sin(radians($3::numeric)))) val
  FROM vehicle
  WHERE vehicle.system_pk = $5 AND latitude IS NOT NULL AND longitude IS NOT NULL
)
SELECT vehicle.pk, vehicle.id, vehicle.system_pk, vehicle.trip_pk, vehicle.label, vehicle.license_plate, vehicle.current_status, vehicle.latitude, vehicle.longitude, vehicle.bearing, vehicle.odometer, vehicle.speed, vehicle.congestion_level, vehicle.updated_at, vehicle.current_stop_pk, vehicle.current_stop_sequence, vehicle.occupancy_status, vehicle.feed_pk, vehicle.occupancy_percentage,
       stop.id as stop_id,
       stop.name as stop_name,
       trip.id as trip_id,
       trip.direction_id as trip_direction_id,
       route.id as route_id,
       route.color as route_color
FROM vehicle
INNER JOIN distance ON vehicle.pk = distance.vehicle_pk
AND distance.val <= $1::numeric
LEFT JOIN stop ON vehicle.current_stop_pk = stop.pk
LEFT JOIN trip ON vehicle.trip_pk = trip.pk
LEFT JOIN route ON trip.route_pk = route.pk
ORDER BY distance.val
LIMIT $2
`

type ListVehicles_GeographicParams struct {
	MaxDistance pgtype.Numeric
	NumVehicles int32
	Latitude    pgtype.Numeric
	Longitude   pgtype.Numeric
	SystemPk    int64
}

type ListVehicles_GeographicRow struct {
	Pk                  int64
	ID                  pgtype.Text
	SystemPk            int64
	TripPk              pgtype.Int8
	Label               pgtype.Text
	LicensePlate        pgtype.Text
	CurrentStatus       pgtype.Text
	Latitude            pgtype.Numeric
	Longitude           pgtype.Numeric
	Bearing             pgtype.Float4
	Odometer            pgtype.Float8
	Speed               pgtype.Float4
	CongestionLevel     string
	UpdatedAt           pgtype.Timestamptz
	CurrentStopPk       pgtype.Int8
	CurrentStopSequence pgtype.Int4
	OccupancyStatus     pgtype.Text
	FeedPk              int64
	OccupancyPercentage pgtype.Int4
	StopID              pgtype.Text
	StopName            pgtype.Text
	TripID              pgtype.Text
	TripDirectionID     pgtype.Bool
	RouteID             pgtype.Text
	RouteColor          pgtype.Text
}

func (q *Queries) ListVehicles_Geographic(ctx context.Context, arg ListVehicles_GeographicParams) ([]ListVehicles_GeographicRow, error) {
	rows, err := q.db.Query(ctx, listVehicles_Geographic,
		arg.MaxDistance,
		arg.NumVehicles,
		arg.Latitude,
		arg.Longitude,
		arg.SystemPk,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListVehicles_GeographicRow
	for rows.Next() {
		var i ListVehicles_GeographicRow
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.TripPk,
			&i.Label,
			&i.LicensePlate,
			&i.CurrentStatus,
			&i.Latitude,
			&i.Longitude,
			&i.Bearing,
			&i.Odometer,
			&i.Speed,
			&i.CongestionLevel,
			&i.UpdatedAt,
			&i.CurrentStopPk,
			&i.CurrentStopSequence,
			&i.OccupancyStatus,
			&i.FeedPk,
			&i.OccupancyPercentage,
			&i.StopID,
			&i.StopName,
			&i.TripID,
			&i.TripDirectionID,
			&i.RouteID,
			&i.RouteColor,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateVehicle = `-- name: UpdateVehicle :exec
UPDATE vehicle
SET trip_pk = $1,
    label = $2,
    license_plate = $3,
    current_status = $4,
    latitude = $5,
    longitude = $6,
    bearing = $7,
    odometer = $8,
    speed = $9,
    congestion_level = $10,
    updated_at = $11,
    current_stop_pk = $12,
    current_stop_sequence = $13,
    occupancy_status = $14,
    feed_pk = $15,
    occupancy_percentage = $16
WHERE vehicle.pk = $17
`

type UpdateVehicleParams struct {
	TripPk              pgtype.Int8
	Label               pgtype.Text
	LicensePlate        pgtype.Text
	CurrentStatus       pgtype.Text
	Latitude            pgtype.Numeric
	Longitude           pgtype.Numeric
	Bearing             pgtype.Float4
	Odometer            pgtype.Float8
	Speed               pgtype.Float4
	CongestionLevel     string
	UpdatedAt           pgtype.Timestamptz
	CurrentStopPk       pgtype.Int8
	CurrentStopSequence pgtype.Int4
	OccupancyStatus     pgtype.Text
	FeedPk              int64
	OccupancyPercentage pgtype.Int4
	Pk                  int64
}

func (q *Queries) UpdateVehicle(ctx context.Context, arg UpdateVehicleParams) error {
	_, err := q.db.Exec(ctx, updateVehicle,
		arg.TripPk,
		arg.Label,
		arg.LicensePlate,
		arg.CurrentStatus,
		arg.Latitude,
		arg.Longitude,
		arg.Bearing,
		arg.Odometer,
		arg.Speed,
		arg.CongestionLevel,
		arg.UpdatedAt,
		arg.CurrentStopPk,
		arg.CurrentStopSequence,
		arg.OccupancyStatus,
		arg.FeedPk,
		arg.OccupancyPercentage,
		arg.Pk,
	)
	return err
}
