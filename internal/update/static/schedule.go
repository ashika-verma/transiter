package static

import (
	"context"

	"github.com/jamespfennell/gtfs"
	"github.com/jamespfennell/transiter/internal/gen/db"
	"github.com/jamespfennell/transiter/internal/update/common"
)

func updateSchedule(ctx context.Context, updateCtx common.UpdateContext, services []gtfs.Service, stopIDToPk map[string]int64) error {
	// TODO: delete existing schedules
	serviceIDToPk := map[string]int64{}
	for _, service := range services {
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
	return nil
}
