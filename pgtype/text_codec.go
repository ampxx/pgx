package pgtype

import (
	"database/sql/driver"
	"fmt"
)

// TextCodec handles encoding and decoding of PostgreSQL text and varchar types.
// Both text format and binary format are supported (binary format for text types
// is just the raw bytes, same as text format).
type TextCodec struct{}

func (TextCodec) FormatSupported(format int16) bool {
	return format == TextFormatCode || format == BinaryFormatCode
}

func (TextCodec) PreferredFormat() int16 {
	return TextFormatCode
}

func (TextCodec) PlanEncode(m *Map, oid uint32, format int16, value any) EncodePlan {
	switch format {
	case TextFormatCode, BinaryFormatCode:
		switch value.(type) {
		case string:
			return encodePlanTextCodecEitherFormatString{}
		case Text:
			return encodePlanTextCodecEitherFormatText{}
		}
	}
	return nil
}

type encodePlanTextCodecEitherFormatString struct{}

func (encodePlanTextCodecEitherFormatString) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v := value.(string)
	return append(buf, v...), nil
}

type encodePlanTextCodecEitherFormatText struct{}

func (encodePlanTextCodecEitherFormatText) Encode(value any, buf []byte) (newBuf []byte, err error) {
	v := value.(Text)
	if !v.Valid {
		return nil, nil
	}
	return append(buf, v.String...), nil
}

func (TextCodec) PlanScan(m *Map, oid uint32, format int16, target any) ScanPlan {
	switch format {
	case TextFormatCode, BinaryFormatCode:
		switch target.(type) {
		case *string:
			return scanPlanTextCodecEitherFormatToString{}
		case *Text:
			return scanPlanTextCodecEitherFormatToText{}
		}
	}
	return nil
}

type scanPlanTextCodecEitherFormatToString struct{}

func (scanPlanTextCodecEitherFormatToString) Scan(src []byte, dst any) error {
	p := dst.(*string)
	if src == nil {
		// NULL values cannot be scanned into a plain *string; use *Text (pgtype.Text) instead
		// if you need to handle NULLs gracefully.
		return fmt.Errorf("cannot scan NULL into *string; consider using *pgtype.Text to handle NULL values")
	}
	*p = string(src)
	return nil
}

type scanPlanTextCodecEitherFormatToText struct{}

func (scanPlanTextCodecEitherFormatToText) Scan(src []byte, dst any) error {
	p := dst.(*Text)
	if src == nil {
		*p = Text{}
		return nil
	}
	*p = Text{String: string(src), Valid: true}
	return nil
}

func (c TextCodec) DecodeDatabaseSQLValue(m *Map, oid uint32, format int16, src []byte) (driver.Value, error) {
	if src == nil {
		return nil, nil
	}
	return string(src), nil
}

func (c TextCodec) DecodeValue(m *Map, oid uint32, format int16, src []byte) (any, error) {
	if src == nil {
		return nil, nil
	}
	return string(src), nil
}
