package pgtype

import (
	"fmt"
	"math"
	"encoding/binary"
)

type Float8Codec struct{}

func (Float8Codec) FormatSupported(format int16) bool {
	return format == TextFormatCode || format == BinaryFormatCode
}

func (Float8Codec) PreferredFormat() int16 {
	return BinaryFormatCode
}

func (Float8Codec) PlanEncode(m *Map, oid uint32, format int16, value any) EncodePlan {
	switch format {
	case BinaryFormatCode:
		switch value.(type) {
		case float64:
			return encodePlanFloat8CodecBinaryFloat64{}
		case Float8:
			return encodePlanFloat8CodecBinaryFloat8{}
		}
	case TextFormatCode:
		switch value.(type) {
		case float64:
			return encodePlanFloat8CodecTextFloat64{}
		case Float8:
			return encodePlanFloat8CodecTextFloat8{}
		}
	}
	return nil
}

type encodePlanFloat8CodecBinaryFloat64 struct{}

func (encodePlanFloat8CodecBinaryFloat64) Encode(value any, buf []byte) ([]byte, error) {
	v := value.(float64)
	buf = binary.BigEndian.AppendUint64(buf, math.Float64bits(v))
	return buf, nil
}

type encodePlanFloat8CodecBinaryFloat8 struct{}

func (encodePlanFloat8CodecBinaryFloat8) Encode(value any, buf []byte) ([]byte, error) {
	v := value.(Float8)
	if !v.Valid {
		return nil, nil
	}
	buf = binary.BigEndian.AppendUint64(buf, math.Float64bits(v.Float64))
	return buf, nil
}

type encodePlanFloat8CodecTextFloat64 struct{}

func (encodePlanFloat8CodecTextFloat64) Encode(value any, buf []byte) ([]byte, error) {
	v := value.(float64)
	buf = append(buf, fmt.Sprintf("%g", v)...)
	return buf, nil
}

type encodePlanFloat8CodecTextFloat8 struct{}

func (encodePlanFloat8CodecTextFloat8) Encode(value any, buf []byte) ([]byte, error) {
	v := value.(Float8)
	if !v.Valid {
		return nil, nil
	}
	buf = append(buf, fmt.Sprintf("%g", v.Float64)...)
	return buf, nil
}

func (Float8Codec) PlanScan(m *Map, oid uint32, format int16, target any) ScanPlan {
	switch format {
	case BinaryFormatCode:
		switch target.(type) {
		case *float64:
			return scanPlanFloat8BinaryFloat64{}
		case *Float8:
			return scanPlanFloat8BinaryFloat8{}
		}
	}
	return nil
}

type scanPlanFloat8BinaryFloat64 struct{}

func (scanPlanFloat8BinaryFloat64) Scan(src []byte, dst any) error {
	if src == nil {
		return fmt.Errorf("cannot scan NULL into *float64")
	}
	if len(src) != 8 {
		return fmt.Errorf("invalid length for float8: %d", len(src))
	}
	p := dst.(*float64)
	*p = math.Float64frombits(binary.BigEndian.Uint64(src))
	return nil
}

type scanPlanFloat8BinaryFloat8 struct{}

func (scanPlanFloat8BinaryFloat8) Scan(src []byte, dst any) error {
	p := dst.(*Float8)
	if src == nil {
		*p = Float8{}
		return nil
	}
	if len(src) != 8 {
		return fmt.Errorf("invalid length for float8: %d", len(src))
	}
	*p = Float8{Float64: math.Float64frombits(binary.BigEndian.Uint64(src)), Valid: true}
	return nil
}
