package pgtype_test

import (
	"encoding/json"
	"testing"

	"github.com/your-org/pgx/pgtype"
)

func TestTextScan(t *testing.T) {
	tests := []struct {
		src    any
		want   pgtype.Text
		wantErr bool
	}{
		{src: "hello", want: pgtype.Text{String: "hello", Valid: true}},
		{src: nil, want: pgtype.Text{}},
		{src: 42, wantErr: true},
	}

	for _, tt := range tests {
		var got pgtype.Text
		err := got.Scan(tt.src)
		if tt.wantErr {
			if err == nil {
				t.Errorf("Scan(%v) expected error, got nil", tt.src)
			}
			continue
		}
		if err != nil {
			t.Errorf("Scan(%v) unexpected error: %v", tt.src, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Scan(%v) = %v, want %v", tt.src, got, tt.want)
		}
	}
}

func TestTextValue(t *testing.T) {
	tests := []struct {
		input pgtype.Text
		want  any
	}{
		{input: pgtype.Text{String: "world", Valid: true}, want: "world"},
		{input: pgtype.Text{}, want: nil},
	}

	for _, tt := range tests {
		got, err := tt.input.Value()
		if err != nil {
			t.Errorf("Value() unexpected error: %v", err)
			continue
		}
		if got != tt.want {
			t.Errorf("Value() = %v, want %v", got, tt.want)
		}
	}
}

func TestTextJSON(t *testing.T) {
	valid := pgtype.Text{String: "pgx", Valid: true}
	b, err := json.Marshal(valid)
	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}
	if string(b) != `"pgx"` {
		t.Errorf("MarshalJSON = %s, want \"pgx\"", b)
	}

	var got pgtype.Text
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("UnmarshalJSON error: %v", err)
	}
	if got != valid {
		t.Errorf("UnmarshalJSON = %v, want %v", got, valid)
	}

	null := pgtype.Text{}
	b, _ = json.Marshal(null)
	if string(b) != "null" {
		t.Errorf("MarshalJSON null = %s, want null", b)
	}
}
