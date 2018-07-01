package fspress

import (
	"errors"
	"reflect"
	"testing"
)

func eq(t *testing.T, expecting, returned interface{}) {
	if !reflect.DeepEqual(expecting, returned) {
		t.Logf("expecting: %+v\n", expecting)
		t.Logf("returned: %+v\n", returned)
		t.Fatalf("assertion error")
	}
}

func TestCleanURL(t *testing.T) {
	eq(t, "one", cleanURL("12345567890-one.md"))
	eq(t, "one-two", cleanURL("12345567890-one-two.md"))
	eq(t, "one-two-three", cleanURL("12345567890-one-two-three.md"))
}

func TestMustPanics(t *testing.T) {
	err := errors.New("panic")
	defer func() { eq(t, err, recover()) }()
	Must(nil, err)
}

func TestMustReturns(t *testing.T) {
	eq(t, &Blog{}, Must(&Blog{}, nil))
}
