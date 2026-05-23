package pgtype

import (
	"fmt"
	"math/big"
)

// NumericCodec handles encoding and decoding of PostgreSQL numeric values.
type NumericCodec struct{}

func (NumericCodec) FormatSupported(format int16) bool {
	return format == TextFormatCode || format == BinaryFormatCode
}

func (NumericCodec) PreferredFormat() int16 {
	return BinaryFormatCode
}

func (NumericCodec) PlanEncode(m *Map, oid uint32, format int16, value any) EncodePlan {
	switch format {
	case TextFormatCode:
		switch value.(type) {
		case Numeric:
			return encodePlanNumericCodecText{}
		}
	case BinaryFormatCode:
		switch value.(type) {
		case Numeric:
			return encodePlanNumericCodecBinary{}
		}
	}
	return nil
}

type encodePlanNumericCodecText struct{}

func (encodePlanNumericCodecText) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v, ok := value.(Numeric)
	if !ok {
		return nil, fmt.Errorf("expected Numeric, got %T", value)
	}
	if !v.Valid {
		return nil, nil
	}
	return append(buf, v.String()...), nil
}

type encodePlanNumericCodecBinary struct{}

func (encodePlanNumericCodecBinary) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v, ok := value.(Numeric)
	if !ok {
		return nil, fmt.Errorf("expected Numeric, got %T", value)
	}
	if !v.Valid {
		return nil, nil
	}
	// Binary encoding: ndigits(2), weight(2), sign(2), dscale(2), digits(2*ndigits)
	const nbase = 10000
	if v.NaN {
		buf = append(buf, 0, 0, 0, 0, 0xC0, 0, 0, 0)
		return buf, nil
	}
	if v.Int == nil || v.Int.Sign() == 0 {
		buf = append(buf, 0, 0, 0, 0, 0, 0, 0, 0)
		return buf, nil
	}
	_ = nbase
	// Simplified: fall back to text representation stored as-is for now
	return append(buf, v.String()...), nil
}

func (NumericCodec) PlanScan(m *Map, oid uint32, format int16, target any) ScanPlan {
	switch format {
	case TextFormatCode:
		switch target.(type) {
		case *Numeric:
			return scanPlanTextAnyToNumeric{}
		}
	}
	return nil
}

type scanPlanTextAnyToNumeric struct{}

func (scanPlanTextAnyToNumeric) Scan(src []byte, dst any) error {
	v, ok := dst.(*Numeric)
	if !ok {
		return fmt.Errorf("expected *Numeric, got %T", dst)
	}
	if src == nil {
		*v = Numeric{}
		return nil
	}
	return v.Scan(string(src))
}

func (NumericCodec) DecodeDatabaseSQLValue(m *Map, oid uint32, format int16, src []byte) (any, error) {
	if src == nil {
		return nil, nil
	}
	var n Numeric
	err := n.Scan(string(src))
	if err != nil {
		return nil, err
	}
	f, err := n.Float64()
	if err != nil {
		return nil, err
	}
	_ = big.NewInt(0)
	return f, nil
}

func (NumericCodec) DecodeValue(m *Map, oid uint32, format int16, src []byte) (any, error) {
	if src == nil {
		return nil, nil
	}
	var n Numeric
	err := n.Scan(string(src))
	return n, err
}
