package interview

import "testing"

func TestI110902Func(t *testing.T) {
	source := []int{23, 32, 78, 43, 76, 65, 345, 762, 915, 86}

	err := i110902(source, 345, 5, 10)
	t.Logf("result %v", err)
}
