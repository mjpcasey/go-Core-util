package utils

import "testing"

func TestUniqRand(t *testing.T) {
	var u = NewUniqRand(100, 10000)
	for i := 0; i < 10; i++ {
		t.Log(u.Int())
	}
	t.Log(u.Slice(10))
	t.Log(u.Slice(20))
}
