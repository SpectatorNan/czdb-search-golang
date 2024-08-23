package bytex

import "testing"

func TestGetIntLong(t *testing.T) {
	b := []byte{0x01, 0x02, 0x03, 0x04}
	offset := 0
	expected := int64(0x04030201)
	if actual := getIntLong1(b, offset); actual != expected {
		t.Errorf("Expected %d, but got %d", expected, actual)
	}
	if actual := getIntLong2(b); actual != expected {
		t.Errorf("Expected %d, but got %d", expected, actual)
	}
}
