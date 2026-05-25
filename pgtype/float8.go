package pgtype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

type Float8 struct {
	Float64 float64
	Valid   bool
}

func (f *Float8) ScanFloat64(v Float8) error {
	*f = v
	return nil
}

func (f Float8) Float64Value() (Float8, error) {
	return f, nil
}

func (f *Float8) Scan(src any) error {
	if src == nil {
		*f = Float8{}
		return nil
	}
	switch v := src.(type) {
	case float64:
		*f = Float8{Float64: v, Valid: true}
		return nil
	case string:
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("cannot scan %T into Float8: %w", src, err)
		}
		*f = Float8{Float64: n, Valid: true}
		return nil
	}
	return fmt.Errorf("cannot scan %T into Float8", src)
}

// Value implements the driver.Valuer interface.
// Special float values (NaN, Infinity, -Infinity) are returned as their
// PostgreSQL string representations since the wire protocol requires it.
func (f Float8) Value() (driver.Value, error) {
	if !f.Valid {
		return nil, nil
	}
	if math.IsNaN(f.Float64) {
		return "NaN", nil
	}
	if math.IsInf(f.Float64, 1) {
		return "Infinity", nil
	}
	if math.IsInf(f.Float64, -1) {
		return "-Infinity", nil
	}
	return f.Float64, nil
}

func (f Float8) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	// Note: json.Marshal does not support NaN or Infinity for float64,
	// so callers should handle special values before marshaling.
	return json.Marshal(f.Float64)
}

func (f *Float8) UnmarshalJSON(b []byte) error {
	var v *float64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	if v == nil {
		*f = Float8{}
		return nil
	}
	*f = Float8{Float64: *v, Valid: true}
	return nil
}
