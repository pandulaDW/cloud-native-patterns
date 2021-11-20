package testing

import "testing"

func TestAdd(t *testing.T) {
	a := 1

	t.Run("+1", func(t *testing.T) {
		if a+1 != 2 {
			t.Fatal("fail!")
		}
	})

	t.Run("+2", func(t *testing.T) {
		if a+2 != 3 {
			t.Fatal("fail!")
		}
	})
}
