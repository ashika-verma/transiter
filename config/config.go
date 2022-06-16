// Package config contains the definition of the system config Yaml format used by Transiter
package config

// TODO: move into the /systems directory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jamespfennell/gtfs/extensions/nyctalerts"
	"github.com/jamespfennell/transiter/internal/gen/api"
	"gopkg.in/yaml.v3"
)

type SystemConfig struct {
	Name                 string
	Feeds                []FeedConfig
	ServiceMaps          []ServiceMapConfig `yaml:"serviceMaps"`
	NoDefaultServiceMaps bool               `yaml:"noDefaultServiceMaps"`
}

type FeedConfig struct {
	ID string

	RequiredForInstall bool `yaml:"requiredForInstall"`

	PeriodicUpdateEnabled bool           `yaml:"periodicUpdateEnabled"`
	PeriodicUpdatePeriod  *time.Duration `yaml:"periodicUpdatePeriod"`

	URL         string
	HTTPTimeout *time.Duration    `yaml:"httpTimeout"`
	HTTPHeaders map[string]string `yaml:"httpHeaders"`

	Parser Parser
	// TODO: these next options should be pointers
	GtfsStaticOptions   GtfsStaticOptions   `yaml:"gtfsStaticOptions"`
	GtfsRealtimeOptions GtfsRealtimeOptions `yaml:"gtfsRealtimeOptions"`
}

type Parser string

const (
	GtfsStatic    Parser = "GTFS_STATIC"
	GtfsRealtime  Parser = "GTFS_REALTIME"
	NyctSubwayCsv Parser = "NYCT_SUBWAY_CSV"
)

type GtfsStaticOptions struct{}

type GtfsRealtimeExtension string

const (
	NoExtension GtfsRealtimeExtension = ""
	NyctTrips   GtfsRealtimeExtension = "NYCT_TRIPS"
	NyctAlerts  GtfsRealtimeExtension = "NYCT_ALERTS"
)

type GtfsRealtimeOptions struct {
	Extension         GtfsRealtimeExtension
	NyctAlertsOptions *nyctalerts.ExtensionOpts `yaml:"nyctAlertsOptions"`
}

type ServiceMapSource string

const (
	ServiceMapSourceStatic   ServiceMapSource = "STATIC"
	ServiceMapSourceRealtime ServiceMapSource = "REALTIME"
)

type ServiceMapConfig struct {
	ID        string
	Source    ServiceMapSource
	Threshold float64

	Days              []string
	StartsEarlierThan *time.Duration `yaml:"startsEarlierThan"`
	StartsLaterThan   *time.Duration `yaml:"startsLaterThan"`
	EndsEarlierThan   *time.Duration `yaml:"endsEarlierThan"`
	EndsLaterThan     *time.Duration `yaml:"endsLaterThan"`

	DefaultForRoutesAtStop bool `yaml:"defaultForRoutesAtStop"`
	DefaultForStopsInRoute bool `yaml:"defaultForStopsInRoute"`
}

func ConvertAPISystemConfig(sc *api.SystemConfig) *SystemConfig {
	result := &SystemConfig{
		Name: sc.Name,
	}
	for _, feed := range sc.Feeds {
		result.Feeds = append(result.Feeds, *ConvertAPIFeedConfig(feed))
	}
	return result
}

func ConvertAPIFeedConfig(fc *api.FeedConfig) *FeedConfig {
	result := &FeedConfig{
		RequiredForInstall:    fc.RequiredForInstall,
		PeriodicUpdateEnabled: fc.PeriodicUpdateEnabled,
		PeriodicUpdatePeriod:  convertMilliseconds(fc.PeriodicUpdatePeriod),
		URL:                   fc.Url,
		HTTPTimeout:           convertMilliseconds(fc.HttpTimeout),
		HTTPHeaders:           fc.HttpHeaders,
	}

	switch parser := fc.Parser.(type) {
	case *api.FeedConfig_GtfsStaticParser_:
		result.Parser = GtfsStatic
	case *api.FeedConfig_GtfsRealtimeParser_:
		result.Parser = GtfsRealtime
		switch extension := parser.GtfsRealtimeParser.Extension.(type) {
		case *api.FeedConfig_GtfsRealtimeParser_NoExtension_:
			result.GtfsRealtimeOptions.Extension = NoExtension
		case *api.FeedConfig_GtfsRealtimeParser_NyctTripsExtension_:
			result.GtfsRealtimeOptions.Extension = NyctTrips
		case *api.FeedConfig_GtfsRealtimeParser_NyctAlertsExtension_:
			result.GtfsRealtimeOptions.Extension = NyctAlerts
			deduplicationPolicy := nyctalerts.NoDeduplication
			switch extension.NyctAlertsExtension.ElevatorAlertsDeduplicationPolicy {
			case api.FeedConfig_GtfsRealtimeParser_NyctAlertsExtension_DEDUPLICATE_IN_STATION:
				deduplicationPolicy = nyctalerts.DeduplicateInStation
			case api.FeedConfig_GtfsRealtimeParser_NyctAlertsExtension_DEDUPLICATE_IN_COMPLEX:
				deduplicationPolicy = nyctalerts.DeduplicateInComplex
			}
			result.GtfsRealtimeOptions.NyctAlertsOptions = &nyctalerts.ExtensionOpts{
				ElevatorAlertsDeduplicationPolicy:   deduplicationPolicy,
				ElevatorAlertsInformUsingStationIDs: extension.NyctAlertsExtension.ElevatorAlertsInformUsingStationIds,
				SkipTimetabledNoServiceAlerts:       extension.NyctAlertsExtension.SkipTimetabledNoServiceAlerts,
				AddNyctMetadata:                     extension.NyctAlertsExtension.AddNyctMetadata,
			}
		}
	case *api.FeedConfig_NyctSubwayCsvParser_:
		result.Parser = NyctSubwayCsv
	}
	return result
}

func ConvertAPIServiceMapConfig(in *api.ServiceMapConfig) *ServiceMapConfig {
	out := &ServiceMapConfig{
		ID:                     in.Id,
		Threshold:              in.Threshold,
		DefaultForRoutesAtStop: in.DefaultForRoutesAtStop,
		DefaultForStopsInRoute: in.DefaultForStopsInRoute,
	}
	switch source := in.Source.(type) {
	case *api.ServiceMapConfig_RealtimeSource:
		out.Source = ServiceMapSourceRealtime
	case *api.ServiceMapConfig_StaticSource:
		out.Source = ServiceMapSourceStatic
		out.StartsEarlierThan = convertAPITimeInDay(source.StaticSource.StartsEarlierThan)
		out.StartsLaterThan = convertAPITimeInDay(source.StaticSource.StartsLaterThan)
		out.EndsEarlierThan = convertAPITimeInDay(source.StaticSource.EndsEarlierThan)
		out.EndsLaterThan = convertAPITimeInDay(source.StaticSource.EndsLaterThan)
		out.Days = source.StaticSource.Days
	}
	return out
}

func convertAPITimeInDay(in *int64) *time.Duration {
	if in == nil {
		return nil
	}
	out := time.Duration(*in) * time.Minute
	return &out
}

func ConvertSystemConfig(sc *SystemConfig) *api.SystemConfig {
	result := &api.SystemConfig{
		Name: sc.Name,
	}
	for _, feed := range sc.Feeds {
		result.Feeds = append(result.Feeds, ConvertFeedConfig(&feed))
	}
	return result
}

func ConvertFeedConfig(fc *FeedConfig) *api.FeedConfig {
	result := &api.FeedConfig{
		RequiredForInstall:    fc.RequiredForInstall,
		PeriodicUpdateEnabled: fc.PeriodicUpdateEnabled,
		PeriodicUpdatePeriod:  convertDuration(fc.PeriodicUpdatePeriod),
		Url:                   fc.URL,
		HttpTimeout:           convertDuration(fc.HTTPTimeout),
		HttpHeaders:           fc.HTTPHeaders,
	}
	switch fc.Parser {
	case GtfsStatic:
		result.Parser = &api.FeedConfig_GtfsStaticParser_{}
	case GtfsRealtime:
		p := &api.FeedConfig_GtfsRealtimeParser{
			Extension: &api.FeedConfig_GtfsRealtimeParser_NoExtension_{},
		}
		switch fc.GtfsRealtimeOptions.Extension {
		case NyctTrips:
			p.Extension = &api.FeedConfig_GtfsRealtimeParser_NyctTripsExtension_{}
		case NyctAlerts:
			var inOpts nyctalerts.ExtensionOpts
			if fc.GtfsRealtimeOptions.NyctAlertsOptions != nil {
				inOpts = *fc.GtfsRealtimeOptions.NyctAlertsOptions
			}
			deduplicationPolicy := api.FeedConfig_GtfsRealtimeParser_NyctAlertsExtension_NO_DEDUPLICATION
			switch inOpts.ElevatorAlertsDeduplicationPolicy {
			case nyctalerts.DeduplicateInStation:
				deduplicationPolicy = api.FeedConfig_GtfsRealtimeParser_NyctAlertsExtension_DEDUPLICATE_IN_STATION
			case nyctalerts.DeduplicateInComplex:
				deduplicationPolicy = api.FeedConfig_GtfsRealtimeParser_NyctAlertsExtension_DEDUPLICATE_IN_COMPLEX
			}
			p.Extension = &api.FeedConfig_GtfsRealtimeParser_NyctAlertsExtension_{
				NyctAlertsExtension: &api.FeedConfig_GtfsRealtimeParser_NyctAlertsExtension{
					ElevatorAlertsDeduplicationPolicy:   deduplicationPolicy,
					ElevatorAlertsInformUsingStationIds: inOpts.ElevatorAlertsInformUsingStationIDs,
					SkipTimetabledNoServiceAlerts:       inOpts.SkipTimetabledNoServiceAlerts,
					AddNyctMetadata:                     inOpts.AddNyctMetadata,
				},
			}
		}
		result.Parser = &api.FeedConfig_GtfsRealtimeParser_{
			GtfsRealtimeParser: p,
		}
	case NyctSubwayCsv:
		result.Parser = &api.FeedConfig_NyctSubwayCsvParser_{}
	}
	return result
}

func ConvertServiceMapConfig(in *ServiceMapConfig) *api.ServiceMapConfig {
	out := &api.ServiceMapConfig{
		Id:                     in.ID,
		Threshold:              in.Threshold,
		DefaultForRoutesAtStop: in.DefaultForRoutesAtStop,
		DefaultForStopsInRoute: in.DefaultForStopsInRoute,
	}
	switch in.Source {
	case ServiceMapSourceRealtime:
		out.Source = &api.ServiceMapConfig_RealtimeSource{}
	case ServiceMapSourceStatic:
		out.Source = &api.ServiceMapConfig_StaticSource{
			StaticSource: &api.ServiceMapConfig_Static{
				StartsEarlierThan: convertTimeInDay(in.StartsEarlierThan),
				StartsLaterThan:   convertTimeInDay(in.StartsLaterThan),
				EndsEarlierThan:   convertTimeInDay(in.EndsEarlierThan),
				EndsLaterThan:     convertTimeInDay(in.EndsLaterThan),
				Days:              in.Days,
			},
		}
	}
	return out
}

func convertTimeInDay(in *time.Duration) *int64 {
	if in == nil {
		return nil
	}
	out := in.Milliseconds() / (1000 * 60)
	return &out
}

func convertMilliseconds(t *int64) *time.Duration {
	if t == nil {
		return nil
	}
	r := time.Millisecond * time.Duration(*t)
	return &r
}

func convertDuration(t *time.Duration) *int64 {
	if t == nil {
		return nil
	}
	r := t.Milliseconds()
	return &r
}

func UnmarshalFromYaml(b []byte) (*SystemConfig, error) {
	var config SystemConfig
	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)
	if err := dec.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse Transiter system config Yaml: %w", err)
	}
	if len(config.ServiceMaps) == 0 && !config.NoDefaultServiceMaps {
		sevenAm := 7 * time.Hour
		sevenPm := 19 * time.Hour
		config.ServiceMaps = []ServiceMapConfig{
			{
				ID:                     "alltimes",
				Source:                 ServiceMapSourceStatic,
				Threshold:              0.1,
				DefaultForRoutesAtStop: false,
				DefaultForStopsInRoute: true,
			},
			{
				ID:                     "weekday",
				Source:                 ServiceMapSourceStatic,
				Threshold:              0.1,
				StartsLaterThan:        &sevenAm,
				EndsEarlierThan:        &sevenPm,
				Days:                   []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
				DefaultForRoutesAtStop: true,
				DefaultForStopsInRoute: false,
			},
			{
				ID:                     "realtime",
				Source:                 ServiceMapSourceRealtime,
				DefaultForRoutesAtStop: true,
				DefaultForStopsInRoute: true,
			},
		}
	}
	return &config, nil
}

func (fc *FeedConfig) MarshalToJSON() []byte {
	b, err := json.Marshal(fc)
	if err != nil {
		panic(fmt.Sprintf("unexpected error when marhsalling feed config: %s", err))
	}
	return b
}

func UnmarshalFromJSON(b []byte) (*FeedConfig, error) {
	var config FeedConfig
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (smc *ServiceMapConfig) MarshalToJSON() []byte {
	b, err := json.Marshal(smc)
	if err != nil {
		panic(fmt.Sprintf("unexpected error when marhsalling service map config: %s", err))
	}
	return b
}
