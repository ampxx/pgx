package pgtype_test

import (
	"math"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestFloat8CodecEncodeBinary(t *testing.T) {
	m := pgtype.NewMap()
	codec := pgtype.Float8Codec{}

	if !codec.FormatSupported(pgtype.BinaryFormatCode) {
		t.Fatal("BinaryFormatCode should be supported")
	}
	if !codec.FormatSupported(pgtype.TextFormatCode) {
		t.Fatal("TextFormatCode should be supported")
	}
	if codec.PreferredFormat() != pgtype.BinaryFormatCode {
		t.Fatal("preferred format should be binary")
	}

	tests := []struct {
		value  any
		wantOK bool
	}{
		{value: float64(3.14), wantOK: true},
		{value: pgtype.Float8{Float64: 3.14, Valid: true}, wantOK: true},
		{value: pgtype.Float8{}, wantOK: true},
		// NaN and infinity are valid float64 values that should also be encodable
		{value: float64(math.NaN()), wantOK: true},
		{value: float64(math.Inf(1)), wantOK: true},
		{value: float64(math.Inf(-1)), wantOK: true}, // also test negative infinity
		// float32 is not a supported type for Float8Codec
		{value: float32(3.14), wantOK: false},
		{value: "unsupported", wantOK: false},
	}

	for _, tt := range tests {
		plan := codec.PlanEncode(m, 701, pgtype.BinaryFormatCode, tt.value)
		if tt.wantOK && plan == nil {
			t.Errorf("PlanEncode(%T): expected plan, got nil", tt.value)
		}
		if !tt.wantOK && plan != nil {
			t.Errorf("PlanEncode(%T): expected nil plan", tt.value)
		}
	}
}

func TestFloat8CodecScanBinary(t *testing.T) {
	m := pgtype.NewMap()
	codec := pgtype.Float8Codec{}

	tests := []struct {
		src    []byte
		target any
		wantOK bool
	}{
		{
			// Encode 1.23 as IEEE 754 big-endian bytes for the scan source
			src: func() []byte {
				bits := math.Float64bits(1.23)
				b := make([]byte, 8)
				b[0] = byte(bits >> 56)
				b[1] = byte(bits >> 48)
				b[2] = byte(bits >> 40)
				b[3] = byte(bits >> 32)
				b[4] = byte(bits >> 24)
				b[5] = byte(bits >> 16)
				b[6] = byte(bits >> 8)
				b[7] = byte(bits)
				return b
			}(),
			target: new(float64),
			wantOK: true,
		},
		{
			// nil src represents a SQL NULL value; scanning into Float8 should succeed
			src:    nil,
			target: new(pgtype.Float8),
			wantOK: true,
		},
	}

	for _, tt := range tests {
		plan := codec.PlanScan(m, 701, pgtype.BinaryFormatCode, tt.target)
		if tt.wantOK && plan == nil {
			t.Errorf("PlanScan(%T): expected plan, got nil", tt.target)
			continue
		}
		if plan == nil {
			continue
		}
		if err := plan.Scan(tt.src, tt.target); err != nil && tt.src != nil {
			t.Errorf("Scan error: %v", err)
		}
	}
}
