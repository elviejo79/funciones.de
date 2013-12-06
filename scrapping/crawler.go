package main

import (
	"fmt"
	"./cinepolis"
)

func main() {
	fmt.Printf("%#v\n", cinepolis.ExtractCities("http://cinepolis.com"))
}
