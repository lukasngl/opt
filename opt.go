package opt

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"testing/quick"
)

// T represents an option, i.e. a value that may be empty or present.
type T[V any] struct {
	v       V
	present bool
}

// IsZero returns whether the option is empty.
// From go1.24 this can be used with omitzero struct tag.
func (t T[V]) IsZero() bool {
	return !t.present
}

// IsEmpty returns whether the option is empty.
func (t T[V]) IsEmpty() bool {
	return !t.present
}

// IsPresent returns whether the option is present.
func (t T[V]) IsPresent() bool {
	return t.present
}

// String implements [fmt.Stringer].
func (t T[V]) String() string {
	value, present := t.Unwrap()
	if !present {
		return fmt.Sprintf("None[%T]()", value)
	}

	return fmt.Sprintf("Some[%T](%s)", value, coerceString(value))
}

func coerceString[V any](v V) string {
	rv := reflect.ValueOf(v)

	if reflect.TypeFor[V]().Implements(reflect.TypeFor[fmt.Stringer]()) {
		if rv.Kind() == reflect.Pointer && rv.IsNil() {
			return "<nil>"
		}

		return rv.Interface().(fmt.Stringer).String()
	}

	return fmt.Sprintf("%v", v)
}

// None creates a new empty option.
func None[V any]() T[V] {
	//nolint:exhaustruct
	return T[V]{}
}

// Some creates a new present option, that contains the given value.
func Some[V any](value V) T[V] {
	return T[V]{
		v:       value,
		present: true,
	}
}

// FromNillable creates a new option from a pointer.
//
// If the pointer is nil an empty option is returned,
// otherwise the value referenced by the pointer will be used as the value
// for a new present option.
//
// Inverse of [T.ToNillable].
func FromNillable[V any](value *V) T[V] {
	if value == nil {
		return None[V]()
	}

	return Some(*value)
}

// ToNillable returns a pointer to the wrapped value if present,
// otherwise a nil pointer.
//
// Inverse of [FromNillable].
func (t T[V]) ToNillable() *V {
	value, present := t.Unwrap()
	if !present {
		return nil
	}

	return &value
}

// FromZeroable creates a new option from a value.
//
// If the value is zero, an empty option is returned,
// otherwise a present option containing the option is returned.
//
// Zeroness is determined as follows:
//
//  1. If the value has an "IsZero() bool" method, it is used to determine zeroness,
//  2. otherwise [reflect.Value#IsZero] is used.
func FromZeroable[V any](value V) T[V] {
	if isZero(value) {
		return None[V]()
	}

	return Some(value)
}

type isZeroer interface {
	IsZero() bool
}

func isZero[V any](value V) bool {
	rv := reflect.ValueOf(value)

	if reflect.TypeFor[V]().Implements(reflect.TypeFor[isZeroer]()) {
		if rv.Kind() == reflect.Pointer && rv.IsNil() {
			return true
		}

		return rv.Interface().(isZeroer).IsZero()
	}

	return rv.IsZero()
}

// Unwrap returns the wrapped value and whether it is empty.
func (t T[V]) Unwrap() (V, bool) {
	return t.v, t.present
}

// Must returns the wrapped value and panics if the [github.com/lukasngl/opt.T] is empty.
func (t T[V]) Must() V {
	value, present := t.Unwrap()
	if !present {
		panic("called Must() on an empty opt.T")
	}

	return value
}

// OrElse returns the wrapped value if not empty and the given default value otherwise.
func (t T[V]) OrElse(defaultValue V) V {
	value, present := t.Unwrap()
	if !present {
		return defaultValue
	}

	return value
}

// OrElse returns the wrapped value if not empty and the zero value otherwise.
func (t T[V]) OrZero() V {
	value, _ := t.Unwrap()

	return value
}

// Alias for the builtin type.
type (
	Bool = T[bool]

	Complex128 = T[complex64]
	Complex64  = T[complex128]

	Byte = T[byte]

	Error = T[error]

	Float32 = T[float32]
	Float64 = T[float64]

	Int8  = T[uint8]
	Int16 = T[uint16]
	Int32 = T[uint32]
	Int64 = T[uint64]

	Rune = T[rune]

	String = T[string]

	Uint8  = T[uint8]
	Uint16 = T[uint16]
	Uint32 = T[uint32]
	Uint64 = T[uint64]
)

// JSON Marshalling und Unmarshalling.
var (
	_ json.Unmarshaler = &T[any]{}
	_ json.Marshaler   = T[any]{}
)

// MarshalJSON implements [json.Marshaler].
func (t T[V]) MarshalJSON() ([]byte, error) {
	if !t.present {
		return []byte("null"), nil
	}

	return json.Marshal(t.v)
}

// UnmarshalJSON implements [json.Unmarshaler].
func (t *T[V]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		t.present = false

		return nil
	}

	err := json.Unmarshal(data, &t.v)
	if err != nil {
		return err
	}

	t.present = true

	return nil
}

// Database Value and Scanner.
var (
	_ driver.Valuer = T[any]{}
	_ sql.Scanner   = &T[any]{}
)

// Scan implements [sql.Scanner].
func (t *T[V]) Scan(src any) error {
	null := sql.Null[V]{}
	err := null.Scan(src)

	t.v = null.V
	t.present = null.Valid

	return err
}

// Value implements [driver.Valuer].
func (t T[V]) Value() (driver.Value, error) {
	return sql.Null[V]{V: t.v, Valid: t.present}.Value()
}

// Generator for quick testing
var _ quick.Generator = T[any]{}

// Generate implements [quick.Generator].
func (t T[V]) Generate(rand *rand.Rand, _ int) reflect.Value {
	if rand.Intn(2) == 0 {
		return reflect.ValueOf(None[V]())
	}

	value, ok := quick.Value(reflect.TypeFor[V](), rand)
	if !ok {
		panic(fmt.Sprintf("failed to generate value for type %s", reflect.TypeFor[V]().Name()))
	}

	concrete, ok := value.Interface().(V)
	if !ok {
		panic(fmt.Sprintf("failed to cast value to type %s", reflect.TypeFor[V]().Name()))
	}

	return reflect.ValueOf(Some(concrete))
}
