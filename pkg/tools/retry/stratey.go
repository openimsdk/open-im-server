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
