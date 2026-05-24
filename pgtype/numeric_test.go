package pgtype_test

import (
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestNumericScan(t *testing.T) {
	tests := []struct {
		input    string
		wantStr  string
		wantNaN  bool
		wantInf  pgtype.InfinityModifier
		wantErr  bool
	}{
		{input: "123", wantStr: "123"},
		{input: "123.456", wantStr: "123.456"},
		{input: "-42.5", wantStr: "-42.5"},
		{input: "NaN", wantNaN: true},
		{input: "Infinity", wantInf: pgtype.Infinity},
		{input: "-Infinity", wantInf: pgtype.NegativeInfinity},
		{input: "0", wantStr: "0"},
		// Edge cases for large numbers
		{input: "99999999999999999999.99", wantStr: "99999999999999999999.99"},
		{input: "-99999999999999999999.99", wantStr: "-99999999999999999999.99"},
		{input: "not-a-number", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var n pgtype.Numeric
			err := n.Scan(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !n.Valid {
				t.Errorf("expected Valid=true")
			}
			if tt.wantNaN && !n.NaN {
				t.Errorf("expected NaN")
			}
			if tt.wantInf != 0 && n.InfinityModifier != tt.wantInf {
				t.Errorf("expected InfinityModifier=%v, got %v", tt.wantInf, n.InfinityModifier)
			}
			if tt.wantStr != "" && n.String() != tt.wantStr {
				t.Errorf("expected String()=%q, got %q", tt.wantStr, n.String())
			}
		})
	}
}

func TestNumericFloat64(t *testing.T) {
	var n pgtype.Numeric
	err := n.Scan("3.14")
	if err != nil {
		t.Fatal(err)
	}
	f, err := n.Float64()
	if err != nil {
		t.Fatal(err)
	}
	if f < 3.13 || f > 3.15 {
		t.Errorf("expected ~3.14, got %f", f)
	}
}

func TestNumericValue(t *testing.T) {
	n := pgtype.Numeric{
		Int:   big.NewInt(12345),
		Exp:   -2,
		Valid: true,
	}
	v, err := n.Value()
	if err != nil {
		t.Fatal(err)
	}
	if v != "123.45" {
		t.Errorf("expected \"123.45\", got %q", v)
	}
}

func TestNumericInvalidValue(t *testing.T) {
	var n pgtype.Numeric
	v, err := n.Value()
	if err != nil {
		t.Fatal(err)
	}
	if v != nil {
		t.Errorf("expected nil for invalid numeric, got %v", v)
	}
}
