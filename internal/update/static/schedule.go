package static

import (
	"context"
	"fmt"

	"github.com/jamespfennell/gtfs"
	"github.com/jamespfennell/transiter/internal/convert"
	"github.com/jamespfennell/transiter/internal/gen/db"
	"github.com/jamespfennell/transiter/internal/update/common"
)

func updateSchedule(ctx context.Context, updateCtx common.UpdateContext, data *gtfs.Static, routeIDToPk map[string]int64) error {
	// TODO: delete existing schedules
	serviceIDToPk := map[string]int64{}
	for _, service := range data.Services {
		pk, err := updateCtx.Querier.InsertScheduledService(ctx, db.InsertScheduledServiceParams{
			ID:        service.Id,
			SystemPk:  updateCtx.SystemPk,
			FeedPk:    updateCtx.FeedPk,
			Monday:    service.Monday,
			Tuesday:   service.Tuesday,
			Wednesday: service.Wednesday,
			Thursday:  service.Thursday,
			Friday:    service.Friday,
			Saturday:  service.Saturday,
			Sunday:    service.Sunday,
			// EndDate   pgtype.Date TODO
			//StartDate pgtype.Date
		})
		if err != nil {
			return err
		}
		serviceIDToPk[service.Id] = pk
	}
	tripIDToPk := map[string]int64{}
	for _, trip := range data.Trips {
		routePk, ok := routeIDToPk[trip.Route.Id]
		if !ok {
			updateCtx.Logger.Debug(fmt.Sprintf("invalid route ID %q for scheduled trip", trip.Route.Id))
			continue
		}
		servicePk, ok := serviceIDToPk[trip.Service.Id]
		if !ok {
			updateCtx.Logger.Debug(fmt.Sprintf("invalid service ID %q for scheduled trip", trip.Service.Id))
			continue
		}
		pk, err := updateCtx.Querier.InsertScheduledTrip(ctx, db.InsertScheduledTripParams{
			ID:        trip.ID,
			RoutePk:   routePk,
			ServicePk: servicePk,
			// DirectionID:          convert.DirectionID(trip.DirectionId), TODO: need to update the GTFS package
			// BikesAllowed:         trip.BikesAllowed,
			// BlockID:   convert.NullString(trip.BlockID), TODO
			Headsign:  convert.NullString(trip.Headsign),
			ShortName: convert.NullString(trip.ShortName),
			// WheelchairAccessible: trip.WheelchairAccessible,
		})
		if err != nil {
			return err
		}
		tripIDToPk[trip.ID] = pk
	}
	return nil
}
