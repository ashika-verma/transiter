// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: query.sql

package db

import (
	"context"
	"database/sql"

	"github.com/jackc/pgtype"
)

const countAgenciesInSystem = `-- name: CountAgenciesInSystem :one
SELECT COUNT(*) FROM agency WHERE system_pk = $1
`

func (q *Queries) CountAgenciesInSystem(ctx context.Context, systemPk int64) (int64, error) {
	row := q.db.QueryRow(ctx, countAgenciesInSystem, systemPk)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countFeedsInSystem = `-- name: CountFeedsInSystem :one
SELECT COUNT(*) FROM feed WHERE system_pk = $1
`

func (q *Queries) CountFeedsInSystem(ctx context.Context, systemPk int64) (int64, error) {
	row := q.db.QueryRow(ctx, countFeedsInSystem, systemPk)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countRoutesInSystem = `-- name: CountRoutesInSystem :one
SELECT COUNT(*) FROM route WHERE system_pk = $1
`

func (q *Queries) CountRoutesInSystem(ctx context.Context, systemPk int64) (int64, error) {
	row := q.db.QueryRow(ctx, countRoutesInSystem, systemPk)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countStopsInSystem = `-- name: CountStopsInSystem :one
SELECT COUNT(*) FROM stop WHERE system_pk = $1
`

func (q *Queries) CountStopsInSystem(ctx context.Context, systemPk int64) (int64, error) {
	row := q.db.QueryRow(ctx, countStopsInSystem, systemPk)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countTransfersInSystem = `-- name: CountTransfersInSystem :one
SELECT COUNT(*) FROM transfer WHERE system_pk = $1
`

func (q *Queries) CountTransfersInSystem(ctx context.Context, systemPk sql.NullInt64) (int64, error) {
	row := q.db.QueryRow(ctx, countTransfersInSystem, systemPk)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getDestinationsForTrips = `-- name: GetDestinationsForTrips :many
WITH last_stop_sequence AS (
  SELECT trip_pk, MAX(stop_sequence) as stop_sequence
    FROM trip_stop_time
    WHERE trip_pk = ANY($1::bigint[])
    GROUP BY trip_pk
)
SELECT lss.trip_pk, stop.pk destination_pk
  FROM last_stop_sequence lss
  INNER JOIN trip_stop_time
    ON lss.trip_pk = trip_stop_time.trip_pk
    AND lss.stop_sequence = trip_stop_time.stop_sequence
  INNER JOIN stop
    ON trip_stop_time.stop_pk = stop.pk
`

type GetDestinationsForTripsRow struct {
	TripPk        int64
	DestinationPk int64
}

func (q *Queries) GetDestinationsForTrips(ctx context.Context, tripPks []int64) ([]GetDestinationsForTripsRow, error) {
	rows, err := q.db.Query(ctx, getDestinationsForTrips, tripPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDestinationsForTripsRow
	for rows.Next() {
		var i GetDestinationsForTripsRow
		if err := rows.Scan(&i.TripPk, &i.DestinationPk); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRouteInSystem = `-- name: GetRouteInSystem :one
SELECT route.pk, route.id, route.system_pk, route.source_pk, route.color, route.text_color, route.short_name, route.long_name, route.description, route.url, route.sort_order, route.type, route.agency_pk, route.continuous_drop_off, route.continuous_pickup FROM route
    INNER JOIN system ON route.system_pk = system.pk
    WHERE system.pk = $1
    AND route.id = $2
`

type GetRouteInSystemParams struct {
	SystemPk int64
	RouteID  string
}

func (q *Queries) GetRouteInSystem(ctx context.Context, arg GetRouteInSystemParams) (Route, error) {
	row := q.db.QueryRow(ctx, getRouteInSystem, arg.SystemPk, arg.RouteID)
	var i Route
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.SystemPk,
		&i.SourcePk,
		&i.Color,
		&i.TextColor,
		&i.ShortName,
		&i.LongName,
		&i.Description,
		&i.Url,
		&i.SortOrder,
		&i.Type,
		&i.AgencyPk,
		&i.ContinuousDropOff,
		&i.ContinuousPickup,
	)
	return i, err
}

const getStopInSystem = `-- name: GetStopInSystem :one
SELECT stop.pk, stop.id, stop.system_pk, stop.source_pk, stop.parent_stop_pk, stop.name, stop.longitude, stop.latitude, stop.url, stop.code, stop.description, stop.platform_code, stop.timezone, stop.type, stop.wheelchair_boarding, stop.zone_id FROM stop
    INNER JOIN system ON stop.system_pk = system.pk
    WHERE system.id = $1
    AND stop.id = $2
`

type GetStopInSystemParams struct {
	SystemID string
	StopID   string
}

func (q *Queries) GetStopInSystem(ctx context.Context, arg GetStopInSystemParams) (Stop, error) {
	row := q.db.QueryRow(ctx, getStopInSystem, arg.SystemID, arg.StopID)
	var i Stop
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.SystemPk,
		&i.SourcePk,
		&i.ParentStopPk,
		&i.Name,
		&i.Longitude,
		&i.Latitude,
		&i.Url,
		&i.Code,
		&i.Description,
		&i.PlatformCode,
		&i.Timezone,
		&i.Type,
		&i.WheelchairBoarding,
		&i.ZoneID,
	)
	return i, err
}

const getSystem = `-- name: GetSystem :one

SELECT pk, id, name, timezone, status FROM system
WHERE id = $1 LIMIT 1
`

// TODO: move all queries from this file in the $x_queries.sql files.
func (q *Queries) GetSystem(ctx context.Context, id string) (System, error) {
	row := q.db.QueryRow(ctx, getSystem, id)
	var i System
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.Name,
		&i.Timezone,
		&i.Status,
	)
	return i, err
}

const getTrip = `-- name: GetTrip :one
SELECT pk, id, route_pk, source_pk, direction_id, started_at FROM trip
WHERE trip.id = $1
    AND trip.route_pk = $2
`

type GetTripParams struct {
	TripID  string
	RoutePk int64
}

func (q *Queries) GetTrip(ctx context.Context, arg GetTripParams) (Trip, error) {
	row := q.db.QueryRow(ctx, getTrip, arg.TripID, arg.RoutePk)
	var i Trip
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.RoutePk,
		&i.SourcePk,
		&i.DirectionID,
		&i.StartedAt,
	)
	return i, err
}

const getTripByPk = `-- name: GetTripByPk :one
SELECT pk, id, route_pk, source_pk, direction_id, started_at FROM trip WHERE pk = $1
`

func (q *Queries) GetTripByPk(ctx context.Context, pk int64) (Trip, error) {
	row := q.db.QueryRow(ctx, getTripByPk, pk)
	var i Trip
	err := row.Scan(
		&i.Pk,
		&i.ID,
		&i.RoutePk,
		&i.SourcePk,
		&i.DirectionID,
		&i.StartedAt,
	)
	return i, err
}

const listActiveAlertsForRoutes = `-- name: ListActiveAlertsForRoutes :many
SELECT route.pk route_pk, alert.id, alert.cause, alert.effect
FROM route
    INNER JOIN alert_route ON route.pk = alert_route.route_pk
    INNER JOIN alert ON alert_route.alert_pk = alert.pk
    INNER JOIN alert_active_period ON alert_active_period.alert_pk = alert.pk
WHERE route.pk = ANY($1::bigint[])
    AND (
        alert_active_period.starts_at < $2
        OR alert_active_period.starts_at IS NULL
    )
    AND (
        alert_active_period.ends_at > $2
        OR alert_active_period.ends_at IS NULL
    )
ORDER BY alert.id ASC
`

type ListActiveAlertsForRoutesParams struct {
	RoutePks    []int64
	PresentTime sql.NullTime
}

type ListActiveAlertsForRoutesRow struct {
	RoutePk int64
	ID      string
	Cause   string
	Effect  string
}

// ListActiveAlertsForRoutes returns preview information about active alerts for the provided routes.
func (q *Queries) ListActiveAlertsForRoutes(ctx context.Context, arg ListActiveAlertsForRoutesParams) ([]ListActiveAlertsForRoutesRow, error) {
	rows, err := q.db.Query(ctx, listActiveAlertsForRoutes, arg.RoutePks, arg.PresentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListActiveAlertsForRoutesRow
	for rows.Next() {
		var i ListActiveAlertsForRoutesRow
		if err := rows.Scan(
			&i.RoutePk,
			&i.ID,
			&i.Cause,
			&i.Effect,
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

const listActiveAlertsForStops = `-- name: ListActiveAlertsForStops :many
SELECT stop.pk stop_pk, alert.pk, alert.id, alert.cause, alert.effect, alert_active_period.starts_at, alert_active_period.ends_at
FROM stop
    INNER JOIN alert_stop ON stop.pk = alert_stop.stop_pk
    INNER JOIN alert ON alert_stop.alert_pk = alert.pk
    INNER JOIN alert_active_period ON alert_active_period.alert_pk = alert.pk
WHERE stop.pk = ANY($1::bigint[])
    AND (
        alert_active_period.starts_at < $2
        OR alert_active_period.starts_at IS NULL
    )
    AND (
        alert_active_period.ends_at > $2
        OR alert_active_period.ends_at IS NULL
    )
ORDER BY alert.id ASC
`

type ListActiveAlertsForStopsParams struct {
	StopPks     []int64
	PresentTime sql.NullTime
}

type ListActiveAlertsForStopsRow struct {
	StopPk   int64
	Pk       int64
	ID       string
	Cause    string
	Effect   string
	StartsAt sql.NullTime
	EndsAt   sql.NullTime
}

func (q *Queries) ListActiveAlertsForStops(ctx context.Context, arg ListActiveAlertsForStopsParams) ([]ListActiveAlertsForStopsRow, error) {
	rows, err := q.db.Query(ctx, listActiveAlertsForStops, arg.StopPks, arg.PresentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListActiveAlertsForStopsRow
	for rows.Next() {
		var i ListActiveAlertsForStopsRow
		if err := rows.Scan(
			&i.StopPk,
			&i.Pk,
			&i.ID,
			&i.Cause,
			&i.Effect,
			&i.StartsAt,
			&i.EndsAt,
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

const listRoutesByPk = `-- name: ListRoutesByPk :many
SELECT pk, id, system_pk, source_pk, color, text_color, short_name, long_name, description, url, sort_order, type, agency_pk, continuous_drop_off, continuous_pickup FROM route WHERE route.pk = ANY($1::bigint[])
`

func (q *Queries) ListRoutesByPk(ctx context.Context, routePks []int64) ([]Route, error) {
	rows, err := q.db.Query(ctx, listRoutesByPk, routePks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Route
	for rows.Next() {
		var i Route
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.SourcePk,
			&i.Color,
			&i.TextColor,
			&i.ShortName,
			&i.LongName,
			&i.Description,
			&i.Url,
			&i.SortOrder,
			&i.Type,
			&i.AgencyPk,
			&i.ContinuousDropOff,
			&i.ContinuousPickup,
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

const listRoutesInAgency = `-- name: ListRoutesInAgency :many
SELECT route.id, route.color FROM route
WHERE route.agency_pk = $1
`

type ListRoutesInAgencyRow struct {
	ID    string
	Color string
}

func (q *Queries) ListRoutesInAgency(ctx context.Context, agencyPk int64) ([]ListRoutesInAgencyRow, error) {
	rows, err := q.db.Query(ctx, listRoutesInAgency, agencyPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRoutesInAgencyRow
	for rows.Next() {
		var i ListRoutesInAgencyRow
		if err := rows.Scan(&i.ID, &i.Color); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRoutesInSystem = `-- name: ListRoutesInSystem :many
SELECT pk, id, system_pk, source_pk, color, text_color, short_name, long_name, description, url, sort_order, type, agency_pk, continuous_drop_off, continuous_pickup FROM route WHERE system_pk = $1 ORDER BY id
`

func (q *Queries) ListRoutesInSystem(ctx context.Context, systemPk int64) ([]Route, error) {
	rows, err := q.db.Query(ctx, listRoutesInSystem, systemPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Route
	for rows.Next() {
		var i Route
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.SourcePk,
			&i.Color,
			&i.TextColor,
			&i.ShortName,
			&i.LongName,
			&i.Description,
			&i.Url,
			&i.SortOrder,
			&i.Type,
			&i.AgencyPk,
			&i.ContinuousDropOff,
			&i.ContinuousPickup,
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

const listServiceMapsConfigIDsForStops = `-- name: ListServiceMapsConfigIDsForStops :many
SELECT stop.pk, service_map_config.id
FROM service_map_config
    INNER JOIN stop ON service_map_config.system_pk = stop.system_pk
WHERE stop.pk = ANY($1::bigint[])
`

type ListServiceMapsConfigIDsForStopsRow struct {
	Pk int64
	ID string
}

func (q *Queries) ListServiceMapsConfigIDsForStops(ctx context.Context, stopPks []int64) ([]ListServiceMapsConfigIDsForStopsRow, error) {
	rows, err := q.db.Query(ctx, listServiceMapsConfigIDsForStops, stopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListServiceMapsConfigIDsForStopsRow
	for rows.Next() {
		var i ListServiceMapsConfigIDsForStopsRow
		if err := rows.Scan(&i.Pk, &i.ID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listServiceMapsForRoutes = `-- name: ListServiceMapsForRoutes :many
SELECT DISTINCT route.pk route_pk, service_map_config.id config_id, service_map_vertex.position, stop.id stop_id, stop.name stop_name
FROM service_map_config
  INNER JOIN system ON service_map_config.system_pk = system.pk
  INNER JOIN route ON route.system_pk = system.pk
  LEFT JOIN service_map ON service_map.config_pk = service_map_config.pk AND service_map.route_pk = route.pk
  LEFT JOIN service_map_vertex ON service_map_vertex.map_pk = service_map.pk
  LEFT JOIN stop ON stop.pk = service_map_vertex.stop_pk
WHERE route.pk = ANY($1::bigint[])
ORDER BY service_map_config.id, service_map_vertex.position
`

type ListServiceMapsForRoutesRow struct {
	RoutePk  int64
	ConfigID string
	Position sql.NullInt32
	StopID   sql.NullString
	StopName sql.NullString
}

// TODO: make this better?
func (q *Queries) ListServiceMapsForRoutes(ctx context.Context, routePks []int64) ([]ListServiceMapsForRoutesRow, error) {
	rows, err := q.db.Query(ctx, listServiceMapsForRoutes, routePks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListServiceMapsForRoutesRow
	for rows.Next() {
		var i ListServiceMapsForRoutesRow
		if err := rows.Scan(
			&i.RoutePk,
			&i.ConfigID,
			&i.Position,
			&i.StopID,
			&i.StopName,
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

const listServiceMapsForStops = `-- name: ListServiceMapsForStops :many
SELECT stop.pk stop_pk, service_map_config.id config_id, service_map.route_pk route_pk
FROM stop
  INNER JOIN service_map_vertex vertex ON vertex.stop_pk = stop.pk
  INNER JOIN service_map ON service_map.pk = vertex.map_pk
  INNER JOIN service_map_config ON service_map_config.pk = service_map.config_pk
WHERE stop.pk = ANY($1::bigint[])
`

type ListServiceMapsForStopsRow struct {
	StopPk   int64
	ConfigID string
	RoutePk  int64
}

func (q *Queries) ListServiceMapsForStops(ctx context.Context, stopPks []int64) ([]ListServiceMapsForStopsRow, error) {
	rows, err := q.db.Query(ctx, listServiceMapsForStops, stopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListServiceMapsForStopsRow
	for rows.Next() {
		var i ListServiceMapsForStopsRow
		if err := rows.Scan(&i.StopPk, &i.ConfigID, &i.RoutePk); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listStopHeadsignRulesForStops = `-- name: ListStopHeadsignRulesForStops :many
SELECT pk, source_pk, priority, stop_pk, track, headsign FROM stop_headsign_rule
WHERE stop_pk = ANY($1::bigint[])
ORDER BY priority ASC
`

func (q *Queries) ListStopHeadsignRulesForStops(ctx context.Context, stopPks []int64) ([]StopHeadsignRule, error) {
	rows, err := q.db.Query(ctx, listStopHeadsignRulesForStops, stopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []StopHeadsignRule
	for rows.Next() {
		var i StopHeadsignRule
		if err := rows.Scan(
			&i.Pk,
			&i.SourcePk,
			&i.Priority,
			&i.StopPk,
			&i.Track,
			&i.Headsign,
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

const listStopTimesAtStops = `-- name: ListStopTimesAtStops :many
SELECT trip_stop_time.pk, trip_stop_time.stop_pk, trip_stop_time.trip_pk, trip_stop_time.arrival_time, trip_stop_time.arrival_delay, trip_stop_time.arrival_uncertainty, trip_stop_time.departure_time, trip_stop_time.departure_delay, trip_stop_time.departure_uncertainty, trip_stop_time.stop_sequence, trip_stop_time.track, trip_stop_time.headsign, trip_stop_time.past, trip.pk, trip.id, trip.route_pk, trip.source_pk, trip.direction_id, trip.started_at, vehicle.id vehicle_id FROM trip_stop_time
    INNER JOIN trip ON trip_stop_time.trip_pk = trip.pk
    LEFT JOIN vehicle ON vehicle.trip_pk = trip.pk
    WHERE trip_stop_time.stop_pk = ANY($1::bigint[])
    AND NOT trip_stop_time.past
    ORDER BY trip_stop_time.departure_time, trip_stop_time.arrival_time
`

type ListStopTimesAtStopsRow struct {
	Pk                   int64
	StopPk               int64
	TripPk               int64
	ArrivalTime          sql.NullTime
	ArrivalDelay         sql.NullInt32
	ArrivalUncertainty   sql.NullInt32
	DepartureTime        sql.NullTime
	DepartureDelay       sql.NullInt32
	DepartureUncertainty sql.NullInt32
	StopSequence         int32
	Track                sql.NullString
	Headsign             sql.NullString
	Past                 bool
	Pk_2                 int64
	ID                   string
	RoutePk              int64
	SourcePk             int64
	DirectionID          sql.NullBool
	StartedAt            sql.NullTime
	VehicleID            sql.NullString
}

func (q *Queries) ListStopTimesAtStops(ctx context.Context, stopPks []int64) ([]ListStopTimesAtStopsRow, error) {
	rows, err := q.db.Query(ctx, listStopTimesAtStops, stopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListStopTimesAtStopsRow
	for rows.Next() {
		var i ListStopTimesAtStopsRow
		if err := rows.Scan(
			&i.Pk,
			&i.StopPk,
			&i.TripPk,
			&i.ArrivalTime,
			&i.ArrivalDelay,
			&i.ArrivalUncertainty,
			&i.DepartureTime,
			&i.DepartureDelay,
			&i.DepartureUncertainty,
			&i.StopSequence,
			&i.Track,
			&i.Headsign,
			&i.Past,
			&i.Pk_2,
			&i.ID,
			&i.RoutePk,
			&i.SourcePk,
			&i.DirectionID,
			&i.StartedAt,
			&i.VehicleID,
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

const listStopsInStopTree = `-- name: ListStopsInStopTree :many
WITH RECURSIVE
ancestor AS (
    SELECT initial.pk, initial.parent_stop_pk
      FROM stop initial
      WHERE	initial.pk = $1
    UNION
    SELECT parent.pk, parent.parent_stop_pk
        FROM stop parent
        INNER JOIN ancestor ON ancestor.parent_stop_pk = parent.pk
),
descendent AS (
    SELECT pk, parent_stop_pk FROM ancestor
    UNION
    SELECT child.pk, child.parent_stop_pk
        FROM stop child
        INNER JOIN descendent ON descendent.pk = child.parent_stop_pk
)
SELECT stop.pk, stop.id, stop.system_pk, stop.source_pk, stop.parent_stop_pk, stop.name, stop.longitude, stop.latitude, stop.url, stop.code, stop.description, stop.platform_code, stop.timezone, stop.type, stop.wheelchair_boarding, stop.zone_id FROM stop
  INNER JOIN descendent
  ON stop.pk = descendent.pk
`

func (q *Queries) ListStopsInStopTree(ctx context.Context, pk int64) ([]Stop, error) {
	rows, err := q.db.Query(ctx, listStopsInStopTree, pk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Stop
	for rows.Next() {
		var i Stop
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.SourcePk,
			&i.ParentStopPk,
			&i.Name,
			&i.Longitude,
			&i.Latitude,
			&i.Url,
			&i.Code,
			&i.Description,
			&i.PlatformCode,
			&i.Timezone,
			&i.Type,
			&i.WheelchairBoarding,
			&i.ZoneID,
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

const listStopsInSystem = `-- name: ListStopsInSystem :many
SELECT pk, id, system_pk, source_pk, parent_stop_pk, name, longitude, latitude, url, code, description, platform_code, timezone, type, wheelchair_boarding, zone_id FROM stop
WHERE system_pk = $1
  AND id >= $2
  AND (
    NOT $3::bool OR
    id = ANY($4::text[])
  )
ORDER BY id
LIMIT $5
`

type ListStopsInSystemParams struct {
	SystemPk               int64
	FirstStopID            string
	OnlyReturnSpecifiedIds bool
	StopIds                []string
	NumStops               int32
}

func (q *Queries) ListStopsInSystem(ctx context.Context, arg ListStopsInSystemParams) ([]Stop, error) {
	rows, err := q.db.Query(ctx, listStopsInSystem,
		arg.SystemPk,
		arg.FirstStopID,
		arg.OnlyReturnSpecifiedIds,
		arg.StopIds,
		arg.NumStops,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Stop
	for rows.Next() {
		var i Stop
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.SourcePk,
			&i.ParentStopPk,
			&i.Name,
			&i.Longitude,
			&i.Latitude,
			&i.Url,
			&i.Code,
			&i.Description,
			&i.PlatformCode,
			&i.Timezone,
			&i.Type,
			&i.WheelchairBoarding,
			&i.ZoneID,
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

const listStopsInSystemGeoFilter = `-- name: ListStopsInSystemGeoFilter :many
SELECT pk, id, system_pk, source_pk, parent_stop_pk, name, longitude, latitude, url, code, description, platform_code, timezone, type, wheelchair_boarding, zone_id FROM stop
WHERE system_pk = $1
  AND (
    NOT $2::bool OR
    id = ANY($3::text[])
  )
  AND (6371 * acos(cos(radians(latitude)) * cos(radians($4::numeric)) * cos(radians($5::numeric) - radians(longitude)) + sin(radians(latitude)) * sin(radians($4::numeric)))) <= $6::numeric
ORDER BY id
`

type ListStopsInSystemGeoFilterParams struct {
	SystemPk               int64
	OnlyReturnSpecifiedIds bool
	StopIds                []string
	Latitude               pgtype.Numeric
	Longitude              pgtype.Numeric
	MaxDistance            pgtype.Numeric
}

func (q *Queries) ListStopsInSystemGeoFilter(ctx context.Context, arg ListStopsInSystemGeoFilterParams) ([]Stop, error) {
	rows, err := q.db.Query(ctx, listStopsInSystemGeoFilter,
		arg.SystemPk,
		arg.OnlyReturnSpecifiedIds,
		arg.StopIds,
		arg.Latitude,
		arg.Longitude,
		arg.MaxDistance,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Stop
	for rows.Next() {
		var i Stop
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.SystemPk,
			&i.SourcePk,
			&i.ParentStopPk,
			&i.Name,
			&i.Longitude,
			&i.Latitude,
			&i.Url,
			&i.Code,
			&i.Description,
			&i.PlatformCode,
			&i.Timezone,
			&i.Type,
			&i.WheelchairBoarding,
			&i.ZoneID,
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

const listStopsTimesForTrip = `-- name: ListStopsTimesForTrip :many
SELECT trip_stop_time.pk, trip_stop_time.stop_pk, trip_stop_time.trip_pk, trip_stop_time.arrival_time, trip_stop_time.arrival_delay, trip_stop_time.arrival_uncertainty, trip_stop_time.departure_time, trip_stop_time.departure_delay, trip_stop_time.departure_uncertainty, trip_stop_time.stop_sequence, trip_stop_time.track, trip_stop_time.headsign, trip_stop_time.past, stop.id stop_id, stop.name stop_name
FROM trip_stop_time
    INNER JOIN stop ON trip_stop_time.stop_pk = stop.pk
WHERE trip_stop_time.trip_pk = $1
ORDER BY trip_stop_time.stop_sequence ASC
`

type ListStopsTimesForTripRow struct {
	Pk                   int64
	StopPk               int64
	TripPk               int64
	ArrivalTime          sql.NullTime
	ArrivalDelay         sql.NullInt32
	ArrivalUncertainty   sql.NullInt32
	DepartureTime        sql.NullTime
	DepartureDelay       sql.NullInt32
	DepartureUncertainty sql.NullInt32
	StopSequence         int32
	Track                sql.NullString
	Headsign             sql.NullString
	Past                 bool
	StopID               string
	StopName             sql.NullString
}

func (q *Queries) ListStopsTimesForTrip(ctx context.Context, tripPk int64) ([]ListStopsTimesForTripRow, error) {
	rows, err := q.db.Query(ctx, listStopsTimesForTrip, tripPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListStopsTimesForTripRow
	for rows.Next() {
		var i ListStopsTimesForTripRow
		if err := rows.Scan(
			&i.Pk,
			&i.StopPk,
			&i.TripPk,
			&i.ArrivalTime,
			&i.ArrivalDelay,
			&i.ArrivalUncertainty,
			&i.DepartureTime,
			&i.DepartureDelay,
			&i.DepartureUncertainty,
			&i.StopSequence,
			&i.Track,
			&i.Headsign,
			&i.Past,
			&i.StopID,
			&i.StopName,
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

const listSystems = `-- name: ListSystems :many
SELECT pk, id, name, timezone, status FROM system ORDER BY id, name
`

func (q *Queries) ListSystems(ctx context.Context) ([]System, error) {
	rows, err := q.db.Query(ctx, listSystems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []System
	for rows.Next() {
		var i System
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.Name,
			&i.Timezone,
			&i.Status,
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

const listTransfersFromStops = `-- name: ListTransfersFromStops :many
  SELECT transfer.pk, transfer.source_pk, transfer.config_source_pk, transfer.system_pk, transfer.from_stop_pk, transfer.to_stop_pk, transfer.type, transfer.min_transfer_time, transfer.distance
  FROM transfer
  WHERE transfer.from_stop_pk = ANY($1::bigint[])
`

func (q *Queries) ListTransfersFromStops(ctx context.Context, fromStopPks []int64) ([]Transfer, error) {
	rows, err := q.db.Query(ctx, listTransfersFromStops, fromStopPks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfer
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.Pk,
			&i.SourcePk,
			&i.ConfigSourcePk,
			&i.SystemPk,
			&i.FromStopPk,
			&i.ToStopPk,
			&i.Type,
			&i.MinTransferTime,
			&i.Distance,
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

const listTransfersInSystem = `-- name: ListTransfersInSystem :many
SELECT
    transfer.pk, transfer.source_pk, transfer.config_source_pk, transfer.system_pk, transfer.from_stop_pk, transfer.to_stop_pk, transfer.type, transfer.min_transfer_time, transfer.distance,
    from_stop.id from_stop_id, from_stop.name from_stop_name, from_system.id from_system_id,
    to_stop.id to_stop_id, to_stop.name to_stop_name, to_system.id to_system_id
FROM transfer
    INNER JOIN stop from_stop ON from_stop.pk = transfer.from_stop_pk
    INNER JOIN system from_system ON from_stop.system_pk = from_system.pk
    INNER JOIN stop to_stop ON to_stop.pk = transfer.to_stop_pk
    INNER JOIN system to_system ON to_stop.system_pk = to_system.pk
WHERE transfer.system_pk = $1
ORDER BY transfer.pk
`

type ListTransfersInSystemRow struct {
	Pk              int64
	SourcePk        sql.NullInt64
	ConfigSourcePk  sql.NullInt64
	SystemPk        sql.NullInt64
	FromStopPk      int64
	ToStopPk        int64
	Type            string
	MinTransferTime sql.NullInt32
	Distance        sql.NullInt32
	FromStopID      string
	FromStopName    sql.NullString
	FromSystemID    string
	ToStopID        string
	ToStopName      sql.NullString
	ToSystemID      string
}

func (q *Queries) ListTransfersInSystem(ctx context.Context, systemPk sql.NullInt64) ([]ListTransfersInSystemRow, error) {
	rows, err := q.db.Query(ctx, listTransfersInSystem, systemPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTransfersInSystemRow
	for rows.Next() {
		var i ListTransfersInSystemRow
		if err := rows.Scan(
			&i.Pk,
			&i.SourcePk,
			&i.ConfigSourcePk,
			&i.SystemPk,
			&i.FromStopPk,
			&i.ToStopPk,
			&i.Type,
			&i.MinTransferTime,
			&i.Distance,
			&i.FromStopID,
			&i.FromStopName,
			&i.FromSystemID,
			&i.ToStopID,
			&i.ToStopName,
			&i.ToSystemID,
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

const listTrips = `-- name: ListTrips :many
SELECT pk, id, route_pk, source_pk, direction_id, started_at FROM trip
WHERE trip.route_pk = $1
ORDER BY trip.id
`

func (q *Queries) ListTrips(ctx context.Context, routePk int64) ([]Trip, error) {
	rows, err := q.db.Query(ctx, listTrips, routePk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Trip
	for rows.Next() {
		var i Trip
		if err := rows.Scan(
			&i.Pk,
			&i.ID,
			&i.RoutePk,
			&i.SourcePk,
			&i.DirectionID,
			&i.StartedAt,
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

const listUpdatesInFeed = `-- name: ListUpdatesInFeed :many
SELECT pk, feed_pk, started_at, finished, finished_at, result, content_length, content_hash, error_message FROM feed_update
WHERE feed_pk = $1
ORDER BY pk DESC
LIMIT 100
`

func (q *Queries) ListUpdatesInFeed(ctx context.Context, feedPk int64) ([]FeedUpdate, error) {
	rows, err := q.db.Query(ctx, listUpdatesInFeed, feedPk)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedUpdate
	for rows.Next() {
		var i FeedUpdate
		if err := rows.Scan(
			&i.Pk,
			&i.FeedPk,
			&i.StartedAt,
			&i.Finished,
			&i.FinishedAt,
			&i.Result,
			&i.ContentLength,
			&i.ContentHash,
			&i.ErrorMessage,
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
