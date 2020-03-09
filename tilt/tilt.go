package main

import (
	"fmt"
)
func main() {
	nrows := 1000
	ncols := 1000
	fmt.Printf("ncols %d\n", nrows)
	fmt.Printf("nrows %d\n", ncols)
		fmt.Printf("xllcorner %d\n",    513000)
		fmt.Printf("yllcorner %d\n",    152000)
	fmt.Printf("cellsize 1\n")
	fmt.Printf("NODATA_value -9999\n")
	for i := 1; i <= nrows; i++ {
		for j := 1; j <= ncols; j++ {
		
			start := float32(i) / 2.0
			number := start + (float32(j) / 2.0)
			fmt.Printf("%f ", number)
		}
		fmt.Printf("\n")
	}
}