// Package admin contains the implementation of the Transiter admin service.
package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jamespfennell/transiter/config"
	"github.com/jamespfennell/transiter/internal/gen/api"
	"github.com/jamespfennell/transiter/internal/gen/db"
	"github.com/jamespfennell/transiter/internal/public/errors"
	"github.com/jamespfennell/transiter/internal/scheduler"
)

// Service implements the Transiter admin service.
type Service struct {
	database  *sql.DB
	scheduler *scheduler.Scheduler
	api.UnimplementedTransiterAdminServer
}

func New(database *sql.DB, scheduler *scheduler.Scheduler) *Service {
	return &Service{
		database:  database,
		scheduler: scheduler,
	}
}

func (s *Service) GetSystemConfig(ctx context.Context, req *api.GetSystemConfigRequest) (*api.SystemConfig, error) {
	tx, err := s.database.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	querier := db.New(tx)

	system, err := querier.GetSystem(ctx, req.SystemId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.NewNotFoundError(fmt.Sprintf("system %q not found", req.SystemId))
		}
		return nil, err
	}
	feeds, err := querier.ListFeedsInSystem(ctx, system.Pk)
	if err != nil {
		return nil, err
	}

	reply := &api.SystemConfig{
		Name: system.Name,
	}
	for _, feed := range feeds {
		feed := feed
		var feedConfig config.FeedConfig
		if err := json.Unmarshal([]byte(feed.Config), &feedConfig); err != nil {
			log.Panicln("Failed to unmarshal feed config from the DB:", err)
			return nil, err
		}
		reply.Feeds = append(reply.Feeds, config.ConvertFeedConfig(&feedConfig))
	}
	return reply, tx.Commit()
}

func (s *Service) InstallOrUpdateSystem(ctx context.Context, req *api.InstallOrUpdateSystemRequest) (*api.InstallOrUpdateSystemReply, error) {
	log.Printf("Recieved install or update request for system %q", req.SystemId)
	tx, err := s.database.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	querier := db.New(tx)

	var systemConfig *config.SystemConfig
	switch c := req.GetConfig().(type) {
	case *api.InstallOrUpdateSystemRequest_SystemConfig:
		systemConfig = config.ConvertApiSystemConfig(c.SystemConfig)
	case *api.InstallOrUpdateSystemRequest_YamlConfigUrl:
		systemConfig, err = getSystemConfigFromUrl(c.YamlConfigUrl)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("no system configuration provided")
	}
	// log.Printf("Config for %s: %+v\n", req.SystemId, systemConfig)

	{
		system, err := querier.GetSystem(ctx, req.SystemId)
		if err == sql.ErrNoRows {
			if err = querier.InsertSystem(ctx, db.InsertSystemParams{
				ID:     req.SystemId,
				Name:   systemConfig.Name,
				Status: "ACTIVE",
			}); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		} else {
			if err = querier.UpdateSystem(ctx, db.UpdateSystemParams{
				Pk:   system.Pk,
				Name: systemConfig.Name,
			}); err != nil {
				return nil, err
			}
		}
	}
	system, err := querier.GetSystem(ctx, req.SystemId)
	if err != nil {
		return nil, err
	}

	feeds, err := querier.ListFeedsInSystem(ctx, system.Pk)
	if err != nil {
		return nil, err
	}
	feedIdToPk := map[string]int32{}
	for _, feed := range feeds {
		feedIdToPk[feed.ID] = feed.Pk
	}

	for _, newFeed := range systemConfig.Feeds {
		if pk, ok := feedIdToPk[newFeed.Id]; ok {
			querier.UpdateFeed(ctx, db.UpdateFeedParams{
				FeedPk:            pk,
				AutoUpdateEnabled: newFeed.AutoUpdateEnabled,
				AutoUpdatePeriod:  convertNullDuration(newFeed.AutoUpdatePeriod),
				Config:            string(newFeed.MarshalToJson()),
			})
			delete(feedIdToPk, newFeed.Id)
		} else {
			querier.InsertFeed(ctx, db.InsertFeedParams{
				ID:                newFeed.Id,
				SystemPk:          system.Pk,
				AutoUpdateEnabled: newFeed.AutoUpdateEnabled,
				AutoUpdatePeriod:  convertNullDuration(newFeed.AutoUpdatePeriod),
				Config:            string(newFeed.MarshalToJson()),
			})
		}
	}
	for _, pk := range feedIdToPk {
		querier.DeleteFeed(ctx, pk)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	log.Printf("Installed system %+v\n", system.ID)
	s.scheduler.Refresh(ctx, req.SystemId)
	return &api.InstallOrUpdateSystemReply{}, nil
}

func (s *Service) DeleteSystem(ctx context.Context, req *api.DeleteSystemRequest) (*api.DeleteSystemReply, error) {
	log.Printf("Recieved delete request for system %q", req.SystemId)
	tx, err := s.database.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	querier := db.New(tx)

	system, err := querier.GetSystem(ctx, req.SystemId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.NewNotFoundError(fmt.Sprintf("system %q not found", req.SystemId))
		}
		return nil, err
	}

	if err := querier.DeleteSystem(ctx, system.Pk); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	log.Printf("Deleted system %q", req.SystemId)
	s.scheduler.Refresh(ctx, req.SystemId)
	return &api.DeleteSystemReply{}, nil
}

func getSystemConfigFromUrl(url string) (*config.SystemConfig, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return config.UnmarshalFromYaml(body)
}

func convertNullDuration(d *time.Duration) sql.NullInt32 {
	if d == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Valid: true,
		Int32: int32(d.Milliseconds()),
	}
}
