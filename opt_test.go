package opt_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"testing/quick"
	"time"

	"github.com/lukasngl/opt"
)

func ExampleT_OrZero() {
	nothing := opt.None[string]()

	fmt.Printf("%q", nothing.OrZero())
	// Output: ""
}

func ExampleT_OrElse() {
	something := opt.Some("hello")
	nothing := opt.None[string]()

	fmt.Printf("%s %s",
		something.OrElse("bye"),
		nothing.OrElse("world!"),
	)
	// Output: hello world!
}

func ExampleT_Unwrap() {
	something := opt.Some("hello")

	if value, present := something.Unwrap(); present {
		fmt.Printf("%s unwrapped to %s", something, value)
	}
	// Output: Some[string](hello) unwrapped to hello
}

func ExampleFromNillable_nil() {
	fmt.Printf("%s\n", opt.FromNillable((*time.Time)(nil)))
	// Output: None[time.Time]()
}

func ExampleFromNillable_zeroValue() {
	// Time that is zero and the zero value.
	value := time.Time{}

	fmt.Printf("%s\n", opt.FromNillable(&value))
	// Output: Some[time.Time](0001-01-01 00:00:00 +0000 UTC)
}

func ExampleFromNillable_zero() {
	// Time that is zero, but not the zero value.
	value := time.Time{}.In(time.FixedZone("XTC", 42))

	fmt.Printf("%s\n", opt.FromNillable(&value))
	// Output: Some[time.Time](0001-01-01 00:00:42 +0000 XTC)
}

func ExampleFromNillable_notZero() {
	// Time that is not zero
	value := time.Time{}.Add(time.Hour + time.Minute + time.Second)

	fmt.Printf("%s\n", opt.FromNillable(&value))

	// Output: Some[time.Time](0001-01-01 01:01:01 +0000 UTC)
}

func ExampleFromZeroable_nil() {
	fmt.Printf("%s\n", opt.FromZeroable((*time.Time)(nil)))
	// Output: None[*time.Time]()
}

func ExampleFromZeroable_zeroValue() {
	// Time that is zero and the zero value.
	value := time.Time{}

	fmt.Printf("%s\n", opt.FromZeroable(value))
	fmt.Printf("%s\n", opt.FromZeroable(&value))
	// Output: None[time.Time]()
	// None[*time.Time]()
}

func ExampleFromZeroable_zero() {
	// Time that is zero, but not the zero value.
	value := time.Time{}.In(time.FixedZone("XTC", 42))

	fmt.Printf("%s\n", opt.FromZeroable(value))
	fmt.Printf("%s\n", opt.FromZeroable(&value))
	// Output: None[time.Time]()
	// None[*time.Time]()
}

func ExampleFromZeroable_notZero() {
	// Time that is not zero
	value := time.Time{}.Add(time.Hour + time.Minute + time.Second)

	fmt.Printf("%s\n", opt.FromZeroable(value))
	fmt.Printf("%s\n", opt.FromZeroable(&value))
	// Output: Some[time.Time](0001-01-01 01:01:01 +0000 UTC)
	// Some[*time.Time](0001-01-01 01:01:01 +0000 UTC)
}

type Thing struct {
	Bool    opt.Bool    `json:"bool,"`
	Byte    opt.Byte    `json:"byte"`
	Float32 opt.Float32 `json:"float32"`
	Float64 opt.Float64 `json:"float64"`
	Int8    opt.Int8    `json:"int8"`
	Int16   opt.Int16   `json:"int16"`
	Int32   opt.Int32   `json:"int32"`
	Int64   opt.Int64   `json:"int64"`
	Rune    opt.Rune    `json:"rune"`
	String  opt.String  `json:"string"`
	Uint8   opt.Uint8   `json:"uint8"`
	Uint16  opt.Uint16  `json:"uint16"`
	Uint32  opt.Uint32  `json:"uint32"`
	Uint64  opt.Uint64  `json:"uint64"`
	Struct  opt.T[struct {
		Test  string
		Test2 int
	}] `json:"struct"`
}

func TestMarshalIdentity(t *testing.T) {
	err := quick.Check(func(ser Thing) bool {
		var de Thing

		data, err := json.Marshal(ser)
		if err != nil {
			t.Log(err.Error())
			return false
		}

		err = json.Unmarshal(data, &de)
		if err != nil {
			t.Log(err.Error())
			return false
		}

		if de != ser {
			t.Logf("ser: %#v", ser)
			t.Logf("de: %#v", de)
			t.Log(string(data))
		}

		return de == ser
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFromNillableIdentity(t *testing.T) {
	err := quick.Check(func(input opt.T[string]) bool {
		to := input.ToNillable()
		from := opt.FromNillable(to)

		if input.String() != from.String() {
			t.Logf("input: %#v", to)
			t.Logf("from:  %#v", from)

			return false
		}

		return true
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
