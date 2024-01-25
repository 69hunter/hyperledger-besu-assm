package rlp

import (
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}
func TestEncodeFromAddresses(t *testing.T) {
	extraData := "0xf87ea00000000000000000000000000000000000000000000000000000000000000000f854942c398325fcc8718404a73ee87eb10d839f933eb1940bdedbb2effd48c27311be688cedc8715b225a3494796e1012f2a24425b28e0c947adb0c8b1a6df4dd9403562b0f7fa4d65303b853e8f21d7ba1962ca3c8808400000000c0"
	addresses := []string{"2c398325fcc8718404a73ee87eb10d839f933eb1", "0bdedbb2effd48c27311be688cedc8715b225a34", "796e1012f2a24425b28e0c947adb0c8b1a6df4dd", "03562b0f7fa4d65303b853e8f21d7ba1962ca3c8"}
	a := NewAdapter()
	result, _ := a.EncodeFromAddresses(addresses)
	assert(t, result, extraData)
}
