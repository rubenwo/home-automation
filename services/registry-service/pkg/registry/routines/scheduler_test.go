package routines

import (
	"testing"
)

func TestNewScheduler(t *testing.T) {
	//s := NewScheduler(nil, runtime.NumCPU())
	//if err := s.UpdateRoutines(); err != nil{
	//	t.Errorf("error updating routines: %v", err)
	//}
}

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

// from fib_test.go
func BenchmarkFib10(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		Fib(10)
	}
}