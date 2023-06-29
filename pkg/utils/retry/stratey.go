package retry

import "time"

type BackoffStrategy int

const (
	StrategyConstant BackoffStrategy = iota
	StrategyLinear
	StrategyFibonacci
)

type Strategy interface {
	Sleep(times int) time.Duration
}
type Constant struct {
	startInterval time.Duration
}

func NewConstant(d time.Duration) *Constant {
	return &Constant{startInterval: d}
}

type Linear struct {
	startInterval time.Duration
}

func NewLinear(d time.Duration) *Linear {
	return &Linear{startInterval: d}
}

type Fibonacci struct {
	startInterval time.Duration
}

func NewFibonacci(d time.Duration) *Fibonacci {
	return &Fibonacci{startInterval: d}
}

func (c *Constant) Sleep(_ int) time.Duration {
	return c.startInterval
}
func (l *Linear) Sleep(times int) time.Duration {
	return l.startInterval * time.Duration(times)

}
func (f *Fibonacci) Sleep(times int) time.Duration {
	return f.startInterval * time.Duration(fibonacciNumber(times))

}
func fibonacciNumber(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return fibonacciNumber(n-1) + fibonacciNumber(n-2)
}
