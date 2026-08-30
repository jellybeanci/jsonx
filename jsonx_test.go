package jsonx

import (
	"bytes"
	jsonv2 "encoding/json/v2"
	"testing"
)

type marshalRouteModel struct {
	Items []string `json:"items"`
}

func TestMarshalRoutesTopLevelDynamicContainersToJSONV2(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
	}{
		{name: "map", value: map[string]any(nil), want: "{}"},
		{name: "slice", value: []any(nil), want: "[]"},
		{name: "map behind any", value: any(map[string]any(nil)), want: "{}"},
		{name: "slice behind any", value: any([]any(nil)), want: "[]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.value)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tt.want {
				t.Fatalf("Marshal() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMarshalKeepsTypedValuesOnJSONIterator(t *testing.T) {
	got, err := Marshal(marshalRouteModel{})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"items":null}` {
		t.Fatalf("Marshal() = %s, want jsoniter-compatible nil slice", got)
	}
}

func TestMarshalAndMarshalStringMatch(t *testing.T) {
	values := []any{
		map[string]any{"items": []any{"one", 2.0, true}},
		[]any{"one", 2.0, true},
		marshalRouteModel{Items: []string{"one"}},
		nil,
	}

	for _, value := range values {
		gotBytes, bytesErr := Marshal(value)
		gotString, stringErr := MarshalString(value)
		if bytesErr != nil || stringErr != nil {
			t.Fatalf("Marshal errors: bytes=%v string=%v", bytesErr, stringErr)
		}
		if !bytes.Equal(gotBytes, []byte(gotString)) {
			t.Fatalf("Marshal outputs differ: %q != %q", gotBytes, gotString)
		}
	}
}

func TestMarshalPropagatesErrors(t *testing.T) {
	tests := []any{
		make(chan int),
		map[string]any{"invalid": make(chan int)},
	}
	for _, value := range tests {
		if _, err := Marshal(value); err == nil {
			t.Fatalf("Marshal(%T) unexpectedly succeeded", value)
		}
	}
}

func TestMarshalNil(t *testing.T) {
	got, err := Marshal(nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "null" {
		t.Fatalf("Marshal(nil) = %s, want null", got)
	}
}

var marshalRouteSink []byte

func BenchmarkMarshalRouteTyped(b *testing.B) {
	value := marshalRouteModel{Items: []string{"one", "two", "three"}}
	b.ReportAllocs()
	for b.Loop() {
		out, err := Marshal(value)
		if err != nil {
			b.Fatal(err)
		}
		marshalRouteSink = out
	}
}

func BenchmarkMarshalRouteDynamic(b *testing.B) {
	value := map[string]any{
		"items": []any{"one", 2.0, true},
		"meta":  map[string]any{"active": true},
	}
	b.Run("jsonx", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			out, err := Marshal(value)
			if err != nil {
				b.Fatal(err)
			}
			marshalRouteSink = out
		}
	})
	b.Run("json-v2-direct", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			out, err := jsonv2.Marshal(value)
			if err != nil {
				b.Fatal(err)
			}
			marshalRouteSink = out
		}
	})
	b.Run("jsoniter-direct", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			out, err := jiter.Marshal(value)
			if err != nil {
				b.Fatal(err)
			}
			marshalRouteSink = out
		}
	})
}
