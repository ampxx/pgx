package pgtype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// Int2 represents a PostgreSQL smallint (int2).
type Int2 struct {
	Int16 int16
	Valid bool
}

// Scan implements the database/sql Scanner interface.
func (dst *Int2) Scan(src any) error {
	if src == nil {
		*dst = Int2{}
		return nil
	}

	var n int64

	switch src := src.(type) {
	case int64:
		n = src
	case string:
		var err error
		n, err = strconv.ParseInt(src, 10, 16)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("cannot scan %T into Int2", src)
	}

	if n < math.MinInt16 || n > math.MaxInt16 {
		return fmt.Errorf("%d is out of range for int16", n)
	}

	*dst = Int2{Int16: int16(n), Valid: true}
	return nil
}

// Value implements the database/sql/driver Valuer interface.
func (src Int2) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return int64(src.Int16), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (src Int2) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(src.Int16)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (dst *Int2) UnmarshalJSON(b []byte) error {
	var n *int16
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	if n == nil {
		*dst = Int2{}
	} else {
		*dst = Int2{Int16: *n, Valid: true}
	}
	return nil
}
