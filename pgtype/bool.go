package pgtype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Bool struct {
	Bool  bool
	Valid bool
}

func (b *Bool) ScanBool(v bool) error {
	b.Bool = v
	b.Valid = true
	return nil
}

func (b Bool) BoolValue() (Bool, error) {
	return b, nil
}

func (b *Bool) Scan(src any) error {
	if src == nil {
		*b = Bool{}
		return nil
	}

	switch src := src.(type) {
	case bool:
		*b = Bool{Bool: src, Valid: true}
		return nil
	case string:
		v, err := parseBoolString(src)
		if err != nil {
			return err
		}
		*b = Bool{Bool: v, Valid: true}
		return nil
	}

	return fmt.Errorf("cannot scan type %T into Bool", src)
}

func (b Bool) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bool, nil
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(b.Bool)
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	var v *bool
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v == nil {
		*b = Bool{}
		return nil
	}
	*b = Bool{Bool: *v, Valid: true}
	return nil
}

func parseBoolString(s string) (bool, error) {
	switch s {
	case "t", "true", "TRUE", "1", "yes", "YES", "on", "ON":
		return true, nil
	case "f", "false", "FALSE", "0", "no", "NO", "off", "OFF":
		return false, nil
	}
	return false, fmt.Errorf("cannot parse %q as bool", s)
}
