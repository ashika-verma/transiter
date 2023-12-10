// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CountAgenciesInSystem(ctx context.Context, systemPk int64) (int64, error)
	CountFeedsInSystem(ctx context.Context, systemPk int64) (int64, error)
	CountRoutesInSystem(ctx context.Context, systemPk int64) (int64, error)
	CountStopsInSystem(ctx context.Context, systemPk int64) (int64, error)
	CountTransfersInSystem(ctx context.Context, systemPk pgtype.Int8) (int64, error)
	DeleteAlerts(ctx context.Context, alertPks []int64) error
	DeleteFeed(ctx context.Context, pk int64) error
	DeleteScheduledServices(ctx context.Context, arg DeleteScheduledServicesParams) error
	DeleteServiceMap(ctx context.Context, arg DeleteServiceMapParams) error
	DeleteServiceMapConfig(ctx context.Context, pk int64) error
	DeleteShapes(ctx context.Context, arg DeleteShapesParams) error
	DeleteStaleAgencies(ctx context.Context, arg DeleteStaleAgenciesParams) error
	DeleteStaleAlerts(ctx context.Context, arg DeleteStaleAlertsParams) error
	DeleteStaleRoutes(ctx context.Context, arg DeleteStaleRoutesParams) error
	DeleteStaleStops(ctx context.Context, arg DeleteStaleStopsParams) error
	DeleteStaleTrips(ctx context.Context, arg DeleteStaleTripsParams) ([]int64, error)
	DeleteStopHeadsignRules(ctx context.Context, feedPk int64) error
	DeleteSystem(ctx context.Context, pk int64) error
	DeleteTransfers(ctx context.Context, feedPk int64) error
	DeleteTripStopTimes(ctx context.Context, pks []int64) error
	DeleteVehicles(ctx context.Context, arg DeleteVehiclesParams) error
	EstimateHeadwaysForRoutes(ctx context.Context, arg EstimateHeadwaysForRoutesParams) ([]EstimateHeadwaysForRoutesRow, error)
	GetAgency(ctx context.Context, arg GetAgencyParams) (Agency, error)
	GetAlertInSystem(ctx context.Context, arg GetAlertInSystemParams) (Alert, error)
	GetDestinationsForTrips(ctx context.Context, tripPks []int64) ([]GetDestinationsForTripsRow, error)
	GetFeed(ctx context.Context, arg GetFeedParams) (Feed, error)
	GetRoute(ctx context.Context, arg GetRouteParams) (Route, error)
	GetScheduledService(ctx context.Context, arg GetScheduledServiceParams) (ScheduledService, error)
	GetScheduledTrip(ctx context.Context, arg GetScheduledTripParams) (ScheduledTrip, error)
	GetShape(ctx context.Context, arg GetShapeParams) (Shape, error)
	GetStop(ctx context.Context, arg GetStopParams) (Stop, error)
	GetSystem(ctx context.Context, id string) (System, error)
	GetTrip(ctx context.Context, arg GetTripParams) (GetTripRow, error)
	GetVehicle(ctx context.Context, arg GetVehicleParams) (GetVehicleRow, error)
	InsertAgency(ctx context.Context, arg InsertAgencyParams) (int64, error)
	InsertAlert(ctx context.Context, arg InsertAlertParams) (int64, error)
	InsertAlertActivePeriod(ctx context.Context, arg InsertAlertActivePeriodParams) error
	InsertAlertAgency(ctx context.Context, arg InsertAlertAgencyParams) error
	InsertAlertRoute(ctx context.Context, arg InsertAlertRouteParams) error
	InsertAlertRouteType(ctx context.Context, arg InsertAlertRouteTypeParams) error
	InsertAlertStop(ctx context.Context, arg InsertAlertStopParams) error
	InsertAlertTrip(ctx context.Context, arg InsertAlertTripParams) error
	InsertFeed(ctx context.Context, arg InsertFeedParams) error
	InsertRoute(ctx context.Context, arg InsertRouteParams) (int64, error)
	InsertScheduledService(ctx context.Context, arg InsertScheduledServiceParams) (int64, error)
	InsertScheduledServiceAddition(ctx context.Context, arg InsertScheduledServiceAdditionParams) error
	InsertScheduledServiceRemoval(ctx context.Context, arg InsertScheduledServiceRemovalParams) error
	InsertScheduledTrip(ctx context.Context, arg []InsertScheduledTripParams) (int64, error)
	InsertScheduledTripFrequency(ctx context.Context, arg InsertScheduledTripFrequencyParams) error
	InsertScheduledTripStopTime(ctx context.Context, arg []InsertScheduledTripStopTimeParams) (int64, error)
	InsertServiceMap(ctx context.Context, arg InsertServiceMapParams) (int64, error)
	InsertServiceMapConfig(ctx context.Context, arg InsertServiceMapConfigParams) error
	InsertServiceMapStop(ctx context.Context, arg []InsertServiceMapStopParams) (int64, error)
	InsertShape(ctx context.Context, arg InsertShapeParams) (int64, error)
	InsertStop(ctx context.Context, arg InsertStopParams) (int64, error)
	InsertStopHeadSignRule(ctx context.Context, arg InsertStopHeadSignRuleParams) error
	InsertSystem(ctx context.Context, arg InsertSystemParams) (int64, error)
	InsertTransfer(ctx context.Context, arg InsertTransferParams) error
	InsertTrip(ctx context.Context, arg InsertTripParams) (int64, error)
	InsertTripStopTime(ctx context.Context, arg []InsertTripStopTimeParams) (int64, error)
	InsertVehicle(ctx context.Context, arg InsertVehicleParams) error
	ListActiveAlertsForAgencies(ctx context.Context, arg ListActiveAlertsForAgenciesParams) ([]ListActiveAlertsForAgenciesRow, error)
	// ListActiveAlertsForRoutes returns preview information about active alerts for the provided routes.
	ListActiveAlertsForRoutes(ctx context.Context, arg ListActiveAlertsForRoutesParams) ([]ListActiveAlertsForRoutesRow, error)
	ListActiveAlertsForStops(ctx context.Context, arg ListActiveAlertsForStopsParams) ([]ListActiveAlertsForStopsRow, error)
	ListActivePeriodsForAlerts(ctx context.Context, pks []int64) ([]ListActivePeriodsForAlertsRow, error)
	ListAgencies(ctx context.Context, systemPk int64) ([]Agency, error)
	ListAgenciesByPk(ctx context.Context, pk []int64) ([]Agency, error)
	ListAlertPksAndHashes(ctx context.Context, arg ListAlertPksAndHashesParams) ([]ListAlertPksAndHashesRow, error)
	ListAlertsInSystem(ctx context.Context, systemPk int64) ([]Alert, error)
	ListAlertsInSystemAndByIDs(ctx context.Context, arg ListAlertsInSystemAndByIDsParams) ([]Alert, error)
	ListAlertsWithActivePeriodsAndAllInformedEntities(ctx context.Context, systemPk int64) ([]ListAlertsWithActivePeriodsAndAllInformedEntitiesRow, error)
	ListFeeds(ctx context.Context, systemPk int64) ([]Feed, error)
	ListRoutes(ctx context.Context, systemPk int64) ([]Route, error)
	ListRoutesByPk(ctx context.Context, routePks []int64) ([]ListRoutesByPkRow, error)
	ListRoutesInAgency(ctx context.Context, agencyPk int64) ([]ListRoutesInAgencyRow, error)
	ListScheduledServices(ctx context.Context, systemPk int64) ([]ListScheduledServicesRow, error)
	ListScheduledTripFrequencies(ctx context.Context, systemPk int64) ([]ListScheduledTripFrequenciesRow, error)
	ListScheduledTripStopTimes(ctx context.Context, systemPk int64) ([]ListScheduledTripStopTimesRow, error)
	ListScheduledTrips(ctx context.Context, systemPk int64) ([]ListScheduledTripsRow, error)
	ListServiceMapConfigsInSystem(ctx context.Context, systemPk int64) ([]ServiceMapConfig, error)
	ListServiceMapsConfigIDsForStops(ctx context.Context, stopPks []int64) ([]ListServiceMapsConfigIDsForStopsRow, error)
	// TODO: make this better?
	ListServiceMapsForRoutes(ctx context.Context, routePks []int64) ([]ListServiceMapsForRoutesRow, error)
	ListServiceMapsForStops(ctx context.Context, stopPks []int64) ([]ListServiceMapsForStopsRow, error)
	ListShapes(ctx context.Context, arg ListShapesParams) ([]Shape, error)
	ListShapesAndTrips(ctx context.Context, systemPk int64) ([]ListShapesAndTripsRow, error)
	ListStopHeadsignRulesForStops(ctx context.Context, stopPks []int64) ([]StopHeadsignRule, error)
	ListStopPksForRealtimeMap(ctx context.Context, routePk int64) ([]ListStopPksForRealtimeMapRow, error)
	ListStops(ctx context.Context, arg ListStopsParams) ([]Stop, error)
	ListStopsByPk(ctx context.Context, stopPks []int64) ([]ListStopsByPkRow, error)
	ListStopsTimesForTrip(ctx context.Context, tripPk int64) ([]ListStopsTimesForTripRow, error)
	ListStops_Geographic(ctx context.Context, arg ListStops_GeographicParams) ([]Stop, error)
	ListSystems(ctx context.Context) ([]System, error)
	ListTransfersFromStops(ctx context.Context, fromStopPks []int64) ([]Transfer, error)
	ListTransfersInSystem(ctx context.Context, systemPk pgtype.Int8) ([]ListTransfersInSystemRow, error)
	ListTripPksInSystem(ctx context.Context, arg ListTripPksInSystemParams) ([]ListTripPksInSystemRow, error)
	ListTripStopTimesByStops(ctx context.Context, stopPks []int64) ([]ListTripStopTimesByStopsRow, error)
	ListTripStopTimesForUpdate(ctx context.Context, tripPks []int64) ([]ListTripStopTimesForUpdateRow, error)
	ListTrips(ctx context.Context, arg ListTripsParams) ([]ListTripsRow, error)
	ListVehicles(ctx context.Context, arg ListVehiclesParams) ([]ListVehiclesRow, error)
	ListVehicles_Geographic(ctx context.Context, arg ListVehicles_GeographicParams) ([]ListVehicles_GeographicRow, error)
	MapAgencyPkToId(ctx context.Context, systemPk int64) ([]MapAgencyPkToIdRow, error)
	MapRouteIDToPkInSystem(ctx context.Context, arg MapRouteIDToPkInSystemParams) ([]MapRouteIDToPkInSystemRow, error)
	MapScheduledTripIDToPkInSystem(ctx context.Context, arg MapScheduledTripIDToPkInSystemParams) ([]MapScheduledTripIDToPkInSystemRow, error)
	MapScheduledTripIDToRoutePkInSystem(ctx context.Context, arg MapScheduledTripIDToRoutePkInSystemParams) ([]MapScheduledTripIDToRoutePkInSystemRow, error)
	MapStopIDAndPkToStationPk(ctx context.Context, arg MapStopIDAndPkToStationPkParams) ([]MapStopIDAndPkToStationPkRow, error)
	MapStopIDToPk(ctx context.Context, arg MapStopIDToPkParams) ([]MapStopIDToPkRow, error)
	MapStopPkToChildPks(ctx context.Context, stopPks []int64) ([]MapStopPkToChildPksRow, error)
	MapStopPkToDescendentPks(ctx context.Context, stopPks []int64) ([]MapStopPkToDescendentPksRow, error)
	MapTripIDToPkInSystem(ctx context.Context, arg MapTripIDToPkInSystemParams) ([]MapTripIDToPkInSystemRow, error)
	MapTripPkToVehicleID(ctx context.Context, arg MapTripPkToVehicleIDParams) ([]MapTripPkToVehicleIDRow, error)
	MarkFailedUpdate(ctx context.Context, arg MarkFailedUpdateParams) error
	MarkSkippedUpdate(ctx context.Context, arg MarkSkippedUpdateParams) error
	MarkSuccessfulUpdate(ctx context.Context, arg MarkSuccessfulUpdateParams) error
	MarkTripStopTimesPast(ctx context.Context, arg []MarkTripStopTimesPastParams) *MarkTripStopTimesPastBatchResults
	UpdateAgency(ctx context.Context, arg UpdateAgencyParams) error
	UpdateFeed(ctx context.Context, arg UpdateFeedParams) error
	UpdateRoute(ctx context.Context, arg UpdateRouteParams) error
	UpdateServiceMapConfig(ctx context.Context, arg UpdateServiceMapConfigParams) error
	UpdateStop(ctx context.Context, arg UpdateStopParams) error
	UpdateStop_Parent(ctx context.Context, arg UpdateStop_ParentParams) error
	UpdateSystem(ctx context.Context, arg UpdateSystemParams) error
	UpdateSystemStatus(ctx context.Context, arg UpdateSystemStatusParams) error
	UpdateTrip(ctx context.Context, arg []UpdateTripParams) *UpdateTripBatchResults
	UpdateTripStopTime(ctx context.Context, arg UpdateTripStopTimeParams) error
}

var _ Querier = (*Queries)(nil)
