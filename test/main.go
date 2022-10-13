package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(runtime.NumCPU())

	a := []int{2, 3}
	fmt.Println(a[3:])
}
