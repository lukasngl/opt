package omitzero_test

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/lukasngl/opt"
)

func TestOmitZero_(t *testing.T) {
	err := quick.Check(func(seed int64) bool {
		rand := rand.New(rand.NewSource(seed))

		return typedTests[rand.Intn(len(typedTests))](t, rand)
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func typedTest[V any](t *testing.T, rand *rand.Rand) bool {
	value, ok := quick.Value(reflect.TypeFor[opt.T[V]](), rand)
	if !ok {
		panic("failed to generate type")
	}

	optValue := value.Interface().(opt.T[V])

	ser := struct {
		Value opt.T[V] `json:"value,omitzero"`
	}{optValue}

	data, err := json.Marshal(ser)
	if err != nil {
		t.Log(err.Error())
		return false
	}

	if (string(data) == "{}") != optValue.IsZero() {
		t.Logf("%s Marshaled to %s", optValue, data)
		return false
	}

	return true
}

var typedTests = []func(*testing.T, *rand.Rand) bool{
	typedTest[bool],
	typedTest[byte],
	typedTest[float32],
	typedTest[float64],
	typedTest[int8],
	typedTest[int16],
	typedTest[int32],
	typedTest[int64],
	typedTest[rune],
	typedTest[string],
	typedTest[uint8],
	typedTest[uint16],
	typedTest[uint32],
	typedTest[uint64],
	typedTest[struct {
		Test  string
		Test2 int
	}],
}
