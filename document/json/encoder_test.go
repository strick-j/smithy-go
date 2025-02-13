package json_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/strick-j/smithy-go/document"
	"github.com/strick-j/smithy-go/document/internal/serde"
	"github.com/strick-j/smithy-go/document/json"
)

func TestEncoder_Encode(t *testing.T) {
	t.Run("Object", func(t *testing.T) {
		for name, tt := range sharedObjectTests {
			t.Run(name, func(t *testing.T) {
				testEncode(t, tt)
			})
		}
	})
	t.Run("Array", func(t *testing.T) {
		for name, tt := range sharedArrayTestCases {
			t.Run(name, func(t *testing.T) {
				testEncode(t, tt)
			})
		}
	})
	t.Run("Number", func(t *testing.T) {
		for name, tt := range sharedNumberTestCases {
			t.Run(name, func(t *testing.T) {
				testEncode(t, tt)
			})
		}
	})
	t.Run("String", func(t *testing.T) {
		for name, tt := range sharedStringTests {
			t.Run(name, func(t *testing.T) {
				testEncode(t, tt)
			})
		}
	})
}

func TestNewEncoderUnsupportedTypes(t *testing.T) {
	type customTime time.Time
	type noSerde = document.NoSerde
	type NestedThing struct {
		SomeThing string
		noSerde
	}
	type Thing struct {
		OtherThing  string
		NestedThing NestedThing
	}

	cases := []interface{}{
		time.Now().UTC(),
		customTime(time.Now().UTC()),
		Thing{OtherThing: "foo", NestedThing: NestedThing{SomeThing: "bar"}},
	}

	encoder := json.NewEncoder()
	for _, tt := range cases {
		_, err := encoder.Encode(tt)
		if err == nil {
			t.Errorf("expect error, got nil")
		}
	}
}

func testEncode(t *testing.T, tt testCase) {
	t.Helper()

	e := json.NewEncoder(func(options *json.EncoderOptions) {
		*options = tt.encoderOptions
	})

	encodeBytes, err := e.Encode(tt.actual)
	if (err != nil) != tt.wantErr {
		t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
	}

	got := MustJSONUnmarshal(encodeBytes, !tt.disableJSONNumber)

	if diff := cmp.Diff(
		serde.PtrToValue(MustJSONUnmarshal(tt.json, !tt.disableJSONNumber)),
		serde.PtrToValue(got),
		cmp.AllowUnexported(StructA{}, StructB{}),
		cmp.Comparer(cmpBigFloat()),
		cmp.Comparer(cmpBigInt()),
	); len(diff) > 0 {
		t.Error(diff)
	}
}
