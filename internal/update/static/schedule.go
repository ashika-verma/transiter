package static

import (
	"context"
	"fmt"

	"github.com/jamespfennell/gtfs"
	"github.com/jamespfennell/transiter/internal/convert"
	"github.com/jamespfennell/transiter/internal/gen/db"
	"github.com/jamespfennell/transiter/internal/update/common"
)

func updateSchedule(ctx context.Context, updateCtx common.UpdateContext, data *gtfs.Static, routeIDToPk map[string]int64, stopIDToPk map[string]int64) error {
	if err := updateCtx.Querier.DeleteScheduledServices(ctx, updateCtx.FeedPk); err != nil {
		return err
	}
	// TODOs
	// (1) Delete existing schedules from other feeds
	// (2) change the GTFS package to return all of the fields here
	// (3) change the inserters to use the fields
	// (4) add an integration test...or maybe cannot without adding an API
	// (4) then try to performance optimize
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
			// TODO StartDate: convert.NullTime(&service.StartDate),
			// TODO EndDate:   convert.NullTime(&service.EndDate),
		})
		if err != nil {
			return err
		}
		serviceIDToPk[service.Id] = pk
	}
	var stopTimeParams []db.InsertScheduledTripStopTimeParams
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
		tripPk, err := updateCtx.Querier.InsertScheduledTrip(ctx, db.InsertScheduledTripParams{
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
		for _, stopTime := range trip.StopTimes {
			stopPk, ok := stopIDToPk[stopTime.Stop.Id]
			if !ok {
				updateCtx.Logger.Debug(fmt.Sprintf("invalid stop ID %q for scheduled trip stop time", stopTime.Stop.Id))
				continue
			}
			stopTimeParams = append(stopTimeParams, db.InsertScheduledTripStopTimeParams{
				TripPk:       tripPk,
				StopPk:       stopPk,
				StopSequence: int32(stopTime.StopSequence),
				Headsign:     convert.NullString(stopTime.Headsign),
				// ArrivalTime: stopTime.ArrivalTime,
				/*
					DepartureTime         pgtype.Time
					StopSequence          int32
					ContinuousDropOff     string
					ContinuousPickup      string
					DropOffType           string
					ExactTimes            bool
					Headsign              pgtype.Text
					PickupType            string
					ShapeDistanceTraveled pgtype.Float8*/
			})
		}
	}
	if _, err := updateCtx.Querier.InsertScheduledTripStopTime(ctx, stopTimeParams); err != nil {
		return err
	}
	return nil
}
