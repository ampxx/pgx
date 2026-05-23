package pgtype

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// Numeric represents a PostgreSQL numeric type with arbitrary precision.
type Numeric struct {
	Int    *big.Int
	Exp    int32
	NaN    bool
	InfinityModifier InfinityModifier
	Valid  bool
}

func (n *Numeric) ScanNumeric(v Numeric) error {
	*n = v
	return nil
}

// Scan implements the database/sql Scanner interface.
func (n *Numeric) Scan(src any) error {
	if src == nil {
		*n = Numeric{}
		return nil
	}
	switch v := src.(type) {
	case string:
		return n.scanString(v)
	case []byte:
		return n.scanString(string(v))
	}
	return fmt.Errorf("cannot scan %T into Numeric", src)
}

func (n *Numeric) scanString(s string) error {
	if s == "NaN" {
		*n = Numeric{NaN: true, Valid: true}
		return nil
	}
	if s == "Infinity" || s == "+Infinity" {
		*n = Numeric{InfinityModifier: Infinity, Valid: true}
		return nil
	}
	if s == "-Infinity" {
		*n = Numeric{InfinityModifier: NegativeInfinity, Valid: true}
		return nil
	}

	parts := strings.SplitN(s, ".", 2)
	intStr := parts[0]
	var exp int32
	if len(parts) == 2 {
		exp = -int32(len(parts[1]))
		intStr += parts[1]
	}

	i := new(big.Int)
	if _, ok := i.SetString(intStr, 10); !ok {
		return fmt.Errorf("cannot parse %q as numeric", s)
	}
	*n = Numeric{Int: i, Exp: exp, Valid: true}
	return nil
}

// Value implements the database/sql/driver Valuer interface.
func (n Numeric) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.String(), nil
}

// String returns a string representation of the Numeric.
func (n Numeric) String() string {
	if !n.Valid {
		return ""
	}
	if n.NaN {
		return "NaN"
	}
	if n.InfinityModifier == Infinity {
		return "Infinity"
	}
	if n.InfinityModifier == NegativeInfinity {
		return "-Infinity"
	}
	if n.Int == nil {
		return "0"
	}
	s := n.Int.String()
	if n.Exp == 0 {
		return s
	}
	if n.Exp < 0 {
		pos := len(s) + int(n.Exp)
		if pos <= 0 {
			return "0." + strings.Repeat("0", -pos) + s
		}
		return s[:pos] + "." + s[pos:]
	}
	return s + strings.Repeat("0", int(n.Exp))
}

// Float64 returns the float64 representation.
func (n Numeric) Float64() (float64, error) {
	if !n.Valid {
		return 0, fmt.Errorf("numeric is not valid")
	}
	return strconv.ParseFloat(n.String(), 64)
}
