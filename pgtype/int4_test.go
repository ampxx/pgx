package pgtype_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestInt4Scan(t *testing.T) {
	tests := []struct {
		name    string
		src     any
		want    pgtype.Int4
		wantErr bool
	}{
		{name: "nil", src: nil, want: pgtype.Int4{}},
		{name: "int64 valid", src: int64(42), want: pgtype.Int4{Int32: 42, Valid: true}},
		{name: "int64 min", src: int64(math.MinInt32), want: pgtype.Int4{Int32: math.MinInt32, Valid: true}},
		{name: "int64 max", src: int64(math.MaxInt32), want: pgtype.Int4{Int32: math.MaxInt32, Valid: true}},
		{name: "int64 overflow", src: int64(math.MaxInt32 + 1), wantErr: true},
		{name: "string valid", src: "100", want: pgtype.Int4{Int32: 100, Valid: true}},
		{name: "string invalid", src: "abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got pgtype.Int4
			err := got.Scan(tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt4Value(t *testing.T) {
	tests := []struct {
		name string
		input pgtype.Int4
		want any
	}{
		{name: "null", input: pgtype.Int4{}, want: nil},
		{name: "valid", input: pgtype.Int4{Int32: 7, Valid: true}, want: int64(7)},
		{name: "negative", input: pgtype.Int4{Int32: -1, Valid: true}, want: int64(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Value()
			if err != nil {
				t.Fatalf("Value() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt4JSON(t *testing.T) {
	tests := []struct {
		name  string
		input pgtype.Int4
		want  string
	}{
		{name: "null", input: pgtype.Int4{}, want: "null"},
		{name: "valid", input: pgtype.Int4{Int32: 42, Valid: true}, want: "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("MarshalJSON() unexpected error: %v", err)
			}
			if string(b) != tt.want {
				t.Errorf("MarshalJSON() = %s, want %s", b, tt.want)
			}

			var roundtrip pgtype.Int4
			if err := json.Unmarshal(b, &roundtrip); err != nil {
				t.Fatalf("UnmarshalJSON() unexpected error: %v", err)
			}
			if roundtrip != tt.input {
				t.Errorf("UnmarshalJSON() = %v, want %v", roundtrip, tt.input)
			}
		})
	}
}
