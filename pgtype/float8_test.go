package pgtype_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestFloat8Scan(t *testing.T) {
	tests := []struct {
		src    any
		want   pgtype.Float8
		wantErr bool
	}{
		{src: float64(3.14), want: pgtype.Float8{Float64: 3.14, Valid: true}},
		{src: float64(0), want: pgtype.Float8{Float64: 0, Valid: true}},
		{src: nil, want: pgtype.Float8{}},
		{src: "invalid", wantErr: true},
	}

	for _, tt := range tests {
		var f pgtype.Float8
		err := f.Scan(tt.src)
		if tt.wantErr {
			if err == nil {
				t.Errorf("Scan(%v): expected error, got nil", tt.src)
			}
			continue
		}
		if err != nil {
			t.Errorf("Scan(%v): unexpected error: %v", tt.src, err)
			continue
		}
		if f != tt.want {
			t.Errorf("Scan(%v): got %v, want %v", tt.src, f, tt.want)
		}
	}
}

func TestFloat8Value(t *testing.T) {
	tests := []struct {
		input pgtype.Float8
		want  any
	}{
		{input: pgtype.Float8{Float64: 1.5, Valid: true}, want: float64(1.5)},
		{input: pgtype.Float8{}, want: nil},
		{input: pgtype.Float8{Float64: math.NaN(), Valid: true}, want: "NaN"},
		{input: pgtype.Float8{Float64: math.Inf(1), Valid: true}, want: "Infinity"},
		{input: pgtype.Float8{Float64: math.Inf(-1), Valid: true}, want: "-Infinity"},
	}

	for _, tt := range tests {
		got, err := tt.input.Value()
		if err != nil {
			t.Errorf("Value() error: %v", err)
			continue
		}
		if got != tt.want {
			t.Errorf("Value(): got %v, want %v", got, tt.want)
		}
	}
}

func TestFloat8JSON(t *testing.T) {
	tests := []struct {
		input pgtype.Float8
		want  string
	}{
		{input: pgtype.Float8{Float64: 2.71, Valid: true}, want: "2.71"},
		{input: pgtype.Float8{}, want: "null"},
	}

	for _, tt := range tests {
		b, err := json.Marshal(tt.input)
		if err != nil {
			t.Errorf("MarshalJSON() error: %v", err)
			continue
		}
		if string(b) != tt.want {
			t.Errorf("MarshalJSON(): got %s, want %s", b, tt.want)
		}

		var f pgtype.Float8
		if err := json.Unmarshal(b, &f); err != nil {
			t.Errorf("UnmarshalJSON(%s) error: %v", b, err)
			continue
		}
		if f != tt.input {
			t.Errorf("UnmarshalJSON(%s): got %v, want %v", b, f, tt.input)
		}
	}
}
