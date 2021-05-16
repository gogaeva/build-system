package main

import (
  "testing"
)

func TestFibonacci(t *testing.T) {
  cases := map[uint]uint{
    0:  0,
    1:  1,
    2:  1,
    3:  2,
    10: 55,
  }
  for input, wanted := range cases {
    if res := fibonacci(input); res != wanted {
      t.Errorf("Expected: %d, got: %d", wanted, res)
    }
  }
}