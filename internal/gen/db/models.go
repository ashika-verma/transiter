// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Agency struct {
	Pk       int64
	ID       string
	SystemPk int64
	Name     string
	Url      string
	Timezone string
	Language pgtype.Text
	Phone    pgtype.Text
	FareUrl  pgtype.Text
	Email    pgtype.Text
	FeedPk   int64
}

type Alert struct {
	Pk          int64
	ID          string
	SystemPk    int64
	Cause       string
	Effect      string
	Header      string
	Description string
	Url         string
	Hash        string
	FeedPk      int64
}

type AlertActivePeriod struct {
	Pk       int64
	AlertPk  int64
	StartsAt pgtype.Timestamptz
	EndsAt   pgtype.Timestamptz
}

type AlertAgency struct {
	AlertPk  int64
	AgencyPk int64
}

type AlertRoute struct {
	AlertPk int64
	RoutePk int64
}

type AlertStop struct {
	AlertPk int64
	StopPk  int64
}

type AlertTrip struct {
	AlertPk int64
	TripPk  int64
}

type Feed struct {
	Pk                   int64
	ID                   string
	SystemPk             int64
	Config               string
	LastContentHash      pgtype.Text
	LastUpdate           pgtype.Timestamptz
	LastSuccessfulUpdate pgtype.Timestamptz
	LastSkippedUpdate    pgtype.Timestamptz
	LastFailedUpdate     pgtype.Timestamptz
}

type Route struct {
	Pk                int64
	ID                string
	SystemPk          int64
	Color             string
	TextColor         string
	ShortName         pgtype.Text
	LongName          pgtype.Text
	Description       pgtype.Text
	Url               pgtype.Text
	SortOrder         pgtype.Int4
	Type              string
	AgencyPk          int64
	ContinuousDropOff string
	ContinuousPickup  string
	FeedPk            int64
}

type ScheduledService struct {
	Pk        int64
	ID        string
	SystemPk  int64
	Monday    bool
	Tuesday   bool
	Wednesday bool
	Thursday  bool
	Friday    bool
	Saturday  bool
	Sunday    bool
	EndDate   pgtype.Date
	StartDate pgtype.Date
	FeedPk    int64
}

type ScheduledServiceAddition struct {
	Pk        int64
	ServicePk int64
	Date      pgtype.Date
}

type ScheduledServiceRemoval struct {
	Pk        int64
	ServicePk int64
	Date      pgtype.Date
}

type ScheduledTrip struct {
	Pk                   int64
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

type ScheduledTripFrequency struct {
	Pk             int64
	TripPk         int64
	StartTime      time.Time
	EndTime        time.Time
	Headway        int32
	FrequencyBased bool
}

type ScheduledTripStopTime struct {
	Pk                    int64
	TripPk                int64
	StopPk                int64
	StopSequence          int32
	ExactTimes            bool
	Headsign              pgtype.Text
	ShapeDistanceTraveled pgtype.Float8
	ArrivalTime           pgtype.Int4
	DepartureTime         pgtype.Int4
	ContinuousDropOff     int16
	ContinuousPickup      int16
	DropOffType           int16
	PickupType            int16
}

type ServiceMap struct {
	Pk       int64
	RoutePk  int64
	ConfigPk int64
}

type ServiceMapConfig struct {
	Pk       int64
	ID       string
	SystemPk int64
	Config   []byte
}

type ServiceMapVertex struct {
	Pk       int64
	StopPk   int64
	MapPk    int64
	Position int32
}

type Stop struct {
	Pk                 int64
	ID                 string
	SystemPk           int64
	ParentStopPk       pgtype.Int8
	Name               pgtype.Text
	Longitude          pgtype.Numeric
	Latitude           pgtype.Numeric
	Url                pgtype.Text
	Code               pgtype.Text
	Description        pgtype.Text
	PlatformCode       pgtype.Text
	Timezone           pgtype.Text
	Type               string
	WheelchairBoarding pgtype.Bool
	ZoneID             pgtype.Text
	FeedPk             int64
}

type StopHeadsignRule struct {
	Pk       int64
	Priority int32
	StopPk   int64
	Track    pgtype.Text
	Headsign string
	FeedPk   int64
}

type System struct {
	Pk       int64
	ID       string
	Name     string
	Timezone pgtype.Text
	Status   string
}

type Transfer struct {
	Pk              int64
	SystemPk        pgtype.Int8
	FromStopPk      int64
	ToStopPk        int64
	Type            string
	MinTransferTime pgtype.Int4
	Distance        pgtype.Int4
	FeedPk          int64
}

type Trip struct {
	Pk          int64
	ID          string
	RoutePk     int64
	DirectionID pgtype.Bool
	StartedAt   pgtype.Timestamptz
	GtfsHash    string
	FeedPk      int64
}

type TripStopTime struct {
	Pk                   int64
	StopPk               int64
	TripPk               int64
	ArrivalTime          pgtype.Timestamptz
	ArrivalDelay         pgtype.Int4
	ArrivalUncertainty   pgtype.Int4
	DepartureTime        pgtype.Timestamptz
	DepartureDelay       pgtype.Int4
	DepartureUncertainty pgtype.Int4
	StopSequence         int32
	Track                pgtype.Text
	Headsign             pgtype.Text
	Past                 bool
}

type Vehicle struct {
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
}
