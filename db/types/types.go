// Package types contains Go types corresponding to database types
package types

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
)

type GeographyType uint32

const (
	Point GeographyType = 0x20000001
)

type Geography struct {
	Longitude float64
	Latitude  float64
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
		g.Latitude = math.Float64frombits(byteOrder.Uint64(b[9:17]))
		g.Longitude = math.Float64frombits(byteOrder.Uint64(b[17:25]))
		return nil
	default:
		return fmt.Errorf("unsupported PostGIS type code 0x%x", geographyType)
	}
}
