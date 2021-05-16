package main

import (
  "flag"
  "fmt"
)

var n = flag.Uint("n", 10, "Number order in a sequence")

//fibonacci function computes a value of Fibonacci
//sequence member of order n
func fibonacci(n uint) uint {
  if n == 0 {
    return 0
  }
  if n == 1 {
    return 1
  }
  return fibonacci(n-2) + fibonacci(n-1)
}

func main() {
  flag.Parse()
  fmt.Println(fibonacci(*n))
}