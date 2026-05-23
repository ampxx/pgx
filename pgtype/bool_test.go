package pgtype_test

import (
	"encoding/json"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestBoolScan(t *testing.T) {
	tests := []struct {
		src      any
		want     pgtype.Bool
		wantErr  bool
	}{
		{src: nil, want: pgtype.Bool{}},
		{src: true, want: pgtype.Bool{Bool: true, Valid: true}},
		{src: false, want: pgtype.Bool{Bool: false, Valid: true}},
		{src: "true", want: pgtype.Bool{Bool: true, Valid: true}},
		{src: "false", want: pgtype.Bool{Bool: false, Valid: true}},
		{src: "1", want: pgtype.Bool{Bool: true, Valid: true}},
		{src: "0", want: pgtype.Bool{Bool: false, Valid: true}},
		{src: "invalid", wantErr: true},
		{src: 42, wantErr: true},
	}

	for _, tt := range tests {
		var b pgtype.Bool
		err := b.Scan(tt.src)
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
		if b != tt.want {
			t.Errorf("Scan(%v): got %v, want %v", tt.src, b, tt.want)
		}
	}
}

func TestBoolValue(t *testing.T) {
	tests := []struct {
		input pgtype.Bool
		want  any
	}{
		{input: pgtype.Bool{}, want: nil},
		{input: pgtype.Bool{Bool: true, Valid: true}, want: true},
		{input: pgtype.Bool{Bool: false, Valid: true}, want: false},
	}

	for _, tt := range tests {
		got, err := tt.input.Value()
		if err != nil {
			t.Errorf("Value(): unexpected error: %v", err)
			continue
		}
		if got != tt.want {
			t.Errorf("Value(): got %v, want %v", got, tt.want)
		}
	}
}

func TestBoolJSON(t *testing.T) {
	tests := []struct {
		input pgtype.Bool
		want  string
	}{
		{input: pgtype.Bool{}, want: "null"},
		{input: pgtype.Bool{Bool: true, Valid: true}, want: "true"},
		{input: pgtype.Bool{Bool: false, Valid: true}, want: "false"},
	}

	for _, tt := range tests {
		data, err := json.Marshal(tt.input)
		if err != nil {
			t.Errorf("MarshalJSON(): unexpected error: %v", err)
			continue
		}
		if string(data) != tt.want {
			t.Errorf("MarshalJSON(): got %s, want %s", data, tt.want)
		}

		var b pgtype.Bool
		if err := json.Unmarshal(data, &b); err != nil {
			t.Errorf("UnmarshalJSON(%s): unexpected error: %v", data, err)
			continue
		}
		if b != tt.input {
			t.Errorf("UnmarshalJSON(%s): got %v, want %v", data, b, tt.input)
		}
	}
}
