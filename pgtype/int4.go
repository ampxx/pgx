package pgtype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

type Int4 struct {
	Int32 int32
	Valid bool
}

func (i *Int4) ScanInt64(v Int8) error {
	if !v.Valid {
		*i = Int4{}
		return nil
	}
	if v.Int64 < math.MinInt32 || v.Int64 > math.MaxInt32 {
		return fmt.Errorf("%d is out of range for int4", v.Int64)
	}
	*i = Int4{Int32: int32(v.Int64), Valid: true}
	return nil
}

func (i Int4) Int64Value() (Int8, error) {
	return Int8{Int64: int64(i.Int32), Valid: i.Valid}, nil
}

func (i *Int4) Scan(src any) error {
	if src == nil {
		*i = Int4{}
		return nil
	}
	switch v := src.(type) {
	case int64:
		if v < math.MinInt32 || v > math.MaxInt32 {
			return fmt.Errorf("%d is out of range for int4", v)
		}
		*i = Int4{Int32: int32(v), Valid: true}
		return nil
	case string:
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return err
		}
		*i = Int4{Int32: int32(n), Valid: true}
		return nil
	}
	return fmt.Errorf("cannot scan %T into Int4", src)
}

func (i Int4) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}
	return int64(i.Int32), nil
}

func (i Int4) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(i.Int32)
}

func (i *Int4) UnmarshalJSON(b []byte) error {
	var n *int32
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	if n == nil {
		*i = Int4{}
		return nil
	}
	*i = Int4{Int32: *n, Valid: true}
	return nil
}
