package pgtype

import (
	"fmt"
	"math"
)

type Int4Codec struct{}

func (Int4Codec) FormatSupported(format int16) bool {
	return format == TextFormatCode || format == BinaryFormatCode
}

func (Int4Codec) PreferredFormat() int16 {
	return BinaryFormatCode
}

func (Int4Codec) PlanEncode(m *Map, oid uint32, format int16, value any) EncodePlan {
	switch format {
	case BinaryFormatCode:
		switch value.(type) {
		case int32:
			return encodePlanInt4CodecBinaryInt32{}
		case Int64Valuer:
			return encodePlanInt4CodecBinaryInt64Valuer{}
		}
	case TextFormatCode:
		switch value.(type) {
		case int32:
			return encodePlanInt4CodecTextInt32{}
		case Int64Valuer:
			return encodePlanInt4CodecTextInt64Valuer{}
		}
	}
	return nil
}

type encodePlanInt4CodecBinaryInt32 struct{}

func (encodePlanInt4CodecBinaryInt32) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v := value.(int32)
	return pgio.AppendInt32(buf, v), nil
}

type encodePlanInt4CodecBinaryInt64Valuer struct{}

func (encodePlanInt4CodecBinaryInt64Valuer) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v, err := value.(Int64Valuer).Int64Value()
	if err != nil {
		return nil, err
	}
	if !v.Valid {
		return nil, nil
	}
	if v.Int64 < math.MinInt32 || v.Int64 > math.MaxInt32 {
		return nil, fmt.Errorf("%d is out of range for int4", v.Int64)
	}
	return pgio.AppendInt32(buf, int32(v.Int64)), nil
}

type encodePlanInt4CodecTextInt32 struct{}

func (encodePlanInt4CodecTextInt32) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v := value.(int32)
	return append(buf, fmt.Sprintf("%d", v)...), nil
}

type encodePlanInt4CodecTextInt64Valuer struct{}

func (encodePlanInt4CodecTextInt64Valuer) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v, err := value.(Int64Valuer).Int64Value()
	if err != nil {
		return nil, err
	}
	if !v.Valid {
		return nil, nil
	}
	if v.Int64 < math.MinInt32 || v.Int64 > math.MaxInt32 {
		return nil, fmt.Errorf("%d is out of range for int4", v.Int64)
	}
	return append(buf, fmt.Sprintf("%d", v.Int64)...), nil
}

func (Int4Codec) PlanScan(m *Map, oid uint32, format int16, target any) ScanPlan {
	switch format {
	case BinaryFormatCode:
		switch target.(type) {
		case *int32:
			return scanPlanBinaryInt4ToInt32{}
		case Int64Scanner:
			return scanPlanBinaryInt4ToInt64Scanner{}
		}
	case TextFormatCode:
		switch target.(type) {
		case *int32:
			return scanPlanTextAnyToInt32{}
		}
	}
	return nil
}

type scanPlanBinaryInt4ToInt32 struct{}

func (scanPlanBinaryInt4ToInt32) Scan(src []byte, dst any) error {
	if src == nil {
		return fmt.Errorf("cannot scan NULL into *int32")
	}
	if len(src) != 4 {
		return fmt.Errorf("invalid length for int4: %v", len(src))
	}
	p := (dst).(*int32)
	*p = int32(binary.BigEndian.Uint32(src))
	return nil
}

type scanPlanBinaryInt4ToInt64Scanner struct{}

func (scanPlanBinaryInt4ToInt64Scanner) Scan(src []byte, dst any) error {
	s := dst.(Int64Scanner)
	if src == nil {
		return s.ScanInt64(Int8{})
	}
	if len(src) != 4 {
		return fmt.Errorf("invalid length for int4: %v", len(src))
	}
	n := int64(int32(binary.BigEndian.Uint32(src)))
	return s.ScanInt64(Int8{Int64: n, Valid: true})
}

type scanPlanTextAnyToInt32 struct{}

func (scanPlanTextAnyToInt32) Scan(src []byte, dst any) error {
	if src == nil {
		return fmt.Errorf("cannot scan NULL into *int32")
	}
	n, err := strconv.ParseInt(string(src), 10, 32)
	if err != nil {
		return err
	}
	p := dst.(*int32)
	*p = int32(n)
	return nil
}
