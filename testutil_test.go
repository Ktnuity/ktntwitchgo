package ktntwitchgo

import "testing"

type TestUtil struct {
	t *testing.T
	msg string
}

func (tu *TestUtil) expect(expect, value any) {
	if value != expect {
		tu.t.Errorf("Failed to %s. Expected '%v', got '%v'", tu.msg, expect, value)
	}
}

func formTest(t *testing.T, msg string) TestUtil {
	return TestUtil{
		t: t,
		msg: msg,
	}
}

func asRef[T any](value T) *T {
	return &value
}
