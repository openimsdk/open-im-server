// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retry

import "time"

// renamed int to backofstrategy
type BackoffStrategy int

// const a constant
const (
	StrategyConstant BackoffStrategy = iota
	StrategyLinear
	StrategyFibonacci
)

// define a interface of strategy
type Strategy interface {
	Sleep(times int) time.Duration
}

// define constant struct
type Constant struct {
	startInterval time.Duration
}

// define a new constant
func NewConstant(d time.Duration) *Constant {
	return &Constant{startInterval: d}
}

// define a linear struct
type Linear struct {
	startInterval time.Duration
}

// define a new linear function
func NewLinear(d time.Duration) *Linear {
	return &Linear{startInterval: d}
}

// define fibonacci struct
type Fibonacci struct {
	startInterval time.Duration
}

// new fibonacci function
func NewFibonacci(d time.Duration) *Fibonacci {
	return &Fibonacci{startInterval: d}
}

// sleep reload
func (c *Constant) Sleep(_ int) time.Duration {
	return c.startInterval
}

// sleep reload
func (l *Linear) Sleep(times int) time.Duration {
	return l.startInterval * time.Duration(times)

}

// sleep
func (f *Fibonacci) Sleep(times int) time.Duration {
	return f.startInterval * time.Duration(fibonacciNumber(times))
}
func fibonacciNumber(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return fibonacciNumber(n-1) + fibonacciNumber(n-2)
}
