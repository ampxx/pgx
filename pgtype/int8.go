package pgtype

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

// Int8 represents a PostgreSQL bigint (int8).
// Note: despite the name, this maps to Go's int64, not int8.
// The name "Int8" refers to the PostgreSQL type name (8-byte integer), not the Go type int8.
type Int8 struct {
	Int64 int64
	Valid bool
}

// ScanInt64 implements the Int64Scanner interface.
func (i *Int8) ScanInt64(v Int8) error {
	*i = v
	return nil
}

// Scan implements the database/sql Scanner interface.
func (i *Int8) Scan(src any) error {
	if src == nil {
		*i = Int8{}
		return nil
	}

	switch v := src.(type) {
	case int64:
		*i = Int8{Int64: v, Valid: true}
		return nil
	case string:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("pgtype: cannot scan %T into Int8: %w", src, err)
		}
		*i = Int8{Int64: n, Valid: true}
		return nil
	}

	return fmt.Errorf("pgtype: cannot scan %T into Int8", src)
}

// Value implements the database/sql/driver Valuer interface.
func (i Int8) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return i.Int64, nil
}

// MarshalJSON implements the json.Marshaler interface.
func (i Int8) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Accepts both numeric JSON values and quoted string representations.
// Note: JSON numbers are floating-point by spec, but integers up to 2^53 are
// represented exactly. For full int64 range safety, prefer the quoted string form.
func (i *Int8) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		*i = Int8{}
		return nil
	}
	// Strip surrounding quotes if the value was encoded as a JSON string.
	s := string(b)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("pgtype: cannot unmarshal JSON into Int8: %w", err)
	}
	*i = Int8{Int64: n, Valid: true}
	return nil
}
