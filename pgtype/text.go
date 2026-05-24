package pgtype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Text represents a PostgreSQL text type. The Valid field indicates whether
// the value is non-NULL. A zero-value Text (Valid == false) represents NULL.
//
// Note: Text is also used for PostgreSQL varchar, char, and other string-like
// types. The String field holds the actual text value.
type Text struct {
	String string
	Valid  bool
}

// ScanText implements the TextScanner interface.
func (t *Text) ScanText(v Text) error {
	*t = v
	return nil
}

// Scan implements the database/sql Scanner interface.
// Accepts string, []byte, or nil (NULL) values.
func (t *Text) Scan(src any) error {
	if src == nil {
		*t = Text{}
		return nil
	}
	switch src := src.(type) {
	case string:
		*t = Text{String: src, Valid: true}
		return nil
	case []byte:
		// Also handle []byte since some drivers return text columns as []byte.
		*t = Text{String: string(src), Valid: true}
		return nil
	}
	return fmt.Errorf("cannot scan %T into Text", src)
}

// Value implements the database/sql/driver Valuer interface.
func (t Text) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.String, nil
}

// MarshalJSON implements the json.Marshaler interface.
// A NULL (Valid == false) Text marshals to JSON null.
func (t Text) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.String)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// A JSON null value results in a Text with Valid == false.
func (t *Text) UnmarshalJSON(b []byte) error {
	var s *string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == nil {
		*t = Text{}
		return nil
	}
	*t = Text{String: *s, Valid: true}
	return nil
}
