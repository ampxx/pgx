package pgtype_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestInt8Scan(t *testing.T) {
	tests := []struct {
		name    string
		src     any
		want    pgtype.Int8
		wantErr bool
	}{
		{name: "nil", src: nil, want: pgtype.Int8{}},
		{name: "int64", src: int64(42), want: pgtype.Int8{Int64: 42, Valid: true}},
		{name: "string", src: "9223372036854775807", want: pgtype.Int8{Int64: 9223372036854775807, Valid: true}},
		{name: "negative", src: int64(-1), want: pgtype.Int8{Int64: -1, Valid: true}},
		{name: "invalid string", src: "not-a-number", wantErr: true},
		{name: "unsupported type", src: 3.14, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i pgtype.Int8
			err := i.Scan(tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && i != tt.want {
				t.Errorf("Scan() = %v, want %v", i, tt.want)
			}
		})
	}
}

func TestInt8Value(t *testing.T) {
	valid := pgtype.Int8{Int64: 100, Valid: true}
	v, err := valid.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != int64(100) {
		t.Errorf("Value() = %v, want %v", v, int64(100))
	}

	null := pgtype.Int8{}
	v, err = null.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != nil {
		t.Errorf("Value() = %v, want nil", v)
	}
}

func TestInt8JSON(t *testing.T) {
	i := pgtype.Int8{Int64: 42, Valid: true}
	b, err := i.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error: %v", err)
	}
	if string(b) != "42" {
		t.Errorf("MarshalJSON() = %s, want 42", b)
	}

	var i2 pgtype.Int8
	if err := i2.UnmarshalJSON(b); err != nil {
		t.Fatalf("UnmarshalJSON() error: %v", err)
	}
	if i2 != i {
		t.Errorf("UnmarshalJSON() = %v, want %v", i2, i)
	}

	nullI := pgtype.Int8{}
	nullB, _ := nullI.MarshalJSON()
	if string(nullB) != "null" {
		t.Errorf("MarshalJSON() null = %s, want null", nullB)
	}
}
