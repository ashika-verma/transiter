// Package types contains Go types corresponding to database types
package types

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

type GeographyType uint32

const (
	Point GeographyType = 0x20000001
)

type Geography struct {
	Valid     bool
	Type      GeographyType
	Longitude float64
	Latitude  float64
}

func (g *Geography) NullableLongitude() *float64 {
	if !g.Valid {
		return nil
	}
	return &g.Longitude
}

func (g *Geography) NullableLatitude() *float64 {
	if !g.Valid {
		return nil
	}
	return &g.Latitude
}

func (g *Geography) Scan(src any) error {
	b, err := hex.DecodeString(src.(string))
	if err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch b[0] {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("invalid byte order 0x%02x, require 0x00 (big endian) or 0x01 (little endian)", b[0])
	}

	geographyType := GeographyType(byteOrder.Uint32(b[1:5]))
	switch geographyType {
	case Point:
		g.Valid = true
		g.Type = Point
		g.Longitude = math.Float64frombits(byteOrder.Uint64(b[9:17]))
		g.Latitude = math.Float64frombits(byteOrder.Uint64(b[17:25]))
		return nil
	default:
		return fmt.Errorf("unsupported PostGIS type code 0x%x", geographyType)
	}
}

func (g *Geography) encode() []byte {
	if !g.Valid {
		return nil
	}
	b := make([]byte, 25) // TODO: can this be [25]byte instead?
	b[0] = 1
	binary.LittleEndian.PutUint32(b[1:5], uint32(Point))
	binary.LittleEndian.PutUint32(b[5:9], uint32(4326))
	binary.LittleEndian.PutUint64(b[9:17], math.Float64bits(g.Longitude))
	binary.LittleEndian.PutUint64(b[17:25], math.Float64bits(g.Latitude))
	return b
}

func (g Geography) Value() (driver.Value, error) {
	if !g.Valid {
		return nil, nil
	}
	return hex.EncodeToString(g.encode()), nil
}

func (g *Geography) ScanText(v pgtype.Text) error {
	if !v.Valid {
		return nil
	}
	return g.Scan(v.String)
}

func (g *Geography) TextValue() (pgtype.Text, error) {
	return pgtype.Text{
		String: fmt.Sprintf("%x", g.encode()),
		Valid:  g.Valid,
	}, nil
}
func (g *Geography) PlanEncode(m *pgtype.Map, oid uint32, format int16, value any) pgtype.EncodePlan {
	panic("PlanEncode called!")
}
