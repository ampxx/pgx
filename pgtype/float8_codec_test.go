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
			src:    func() []byte { b := make([]byte, 8); _ = math.Float64bits(1.23); return b }(),
			target: new(float64),
			wantOK: true,
		},
		{
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
