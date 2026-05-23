package pgtype

// FormatCode represents the wire protocol format.
type FormatCode = int16

const (
	// TextFormatCode is the text wire protocol format.
	TextFormatCode FormatCode = 0
	// BinaryFormatCode is the binary wire protocol format.
	BinaryFormatCode FormatCode = 1
)

// EncodePlan is an interface for encoding values.
type EncodePlan interface {
	Encode(value any, buf []byte) (newBuf []byte, err error)
}

// ScanPlan is an interface for scanning values.
type ScanPlan interface {
	Scan(src []byte, dst any) error
}

// Codec is the interface for PostgreSQL type codecs.
type Codec interface {
	// FormatSupported returns true if the codec supports the given format.
	FormatSupported(format int16) bool
	// PreferredFormat returns the preferred format for the codec.
	PreferredFormat() int16
	// PlanEncode returns an EncodePlan for encoding the given value.
	PlanEncode(m *Map, oid uint32, format int16, value any) EncodePlan
	// PlanScan returns a ScanPlan for scanning into the given target.
	PlanScan(m *Map, oid uint32, format int16, target any) ScanPlan
	// DecodeDatabaseSQLValue decodes a value for database/sql compatibility.
	DecodeDatabaseSQLValue(m *Map, oid uint32, format int16, src []byte) (any, error)
	// DecodeValue decodes a value to a native Go type.
	DecodeValue(m *Map, oid uint32, format int16, src []byte) (any, error)
}

// Map holds the mapping of OIDs to codecs.
type Map struct {
	oidToCodec map[uint32]Codec
}

// NewMap creates a new Map with default type mappings.
func NewMap() *Map {
	m := &Map{
		oidToCodec: make(map[uint32]Codec),
	}
	// Register built-in codecs
	const numericOID = 1700
	m.oidToCodec[numericOID] = NumericCodec{}
	return m
}

// RegisterCodec registers a codec for the given OID.
func (m *Map) RegisterCodec(oid uint32, c Codec) {
	m.oidToCodec[oid] = c
}

// CodecForOID returns the codec registered for the given OID.
func (m *Map) CodecForOID(oid uint32) (Codec, bool) {
	c, ok := m.oidToCodec[oid]
	return c, ok
}
