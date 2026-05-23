package pgtype

import "fmt"

// InfinityModifier represents PostgreSQL infinity values.
type InfinityModifier int8

const (
	// Infinity represents positive infinity.
	Infinity InfinityModifier = 1
	// Finite represents a finite value (not infinity).
	Finite InfinityModifier = 0
	// NegativeInfinity represents negative infinity.
	NegativeInfinity InfinityModifier = -1
)

// String returns the string representation of InfinityModifier.
func (im InfinityModifier) String() string {
	switch im {
	case Infinity:
		return "Infinity"
	case Finite:
		return "Finite"
	case NegativeInfinity:
		return "-Infinity"
	default:
		return fmt.Sprintf("invalid InfinityModifier %d", im)
	}
}

// MarshalText implements encoding.TextMarshaler.
func (im InfinityModifier) MarshalText() ([]byte, error) {
	return []byte(im.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (im *InfinityModifier) UnmarshalText(src []byte) error {
	switch string(src) {
	case "Infinity":
		*im = Infinity
	case "Finite":
		*im = Finite
	case "-Infinity":
		*im = NegativeInfinity
	default:
		return fmt.Errorf("cannot unmarshal %q into InfinityModifier", src)
	}
	return nil
}
