package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"github.com/goblimey/tiler/esri"
)

var filename string // The file to display.
var output string   // The .png results file.
var ceiling64 float64 // parameter - the maximum height expected.
var ceiling float32	// ceiling as a float32
var floor64 float64   // parameter - the minimum height expected.
var floor float32	// floor as a float32
var verbose bool    // verbose mode

var maxHeight float64 = 0
var maxHeightSet = false
var minHeight float64 = 0
var minHeightSet = false
var maxShade uint8 = 0
var maxShadeSet = false
var minShade uint8 = 0
var minShadeSet = false

func init() {
	flag.StringVar(&filename, "input", "", "data file")
	flag.StringVar(&filename, "i", "", "data file")
	flag.StringVar(&output, "output", "", ".png results file")
	flag.StringVar(&output, "o", "", ".png results file")
	flag.Float64Var(&ceiling64, "ceiling", 0.0, "maximum height expected")
	flag.Float64Var(&ceiling64, "c", 0.0, "maximum height expected")
	flag.Float64Var(&floor64, "floor", 0.0, "mimimum height expected")
	flag.Float64Var(&floor64, "f", 0.0, "minimum height expected")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")
	flag.BoolVar(&verbose, "v", false, "verbose mode")
}

func main() {
	flag.Parse()

	// filename = "TT"
	// output := "tile.png"

	flagset := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { flagset[f.Name] = true })

	if flagset["floor"] {
		floor = float32(floor64)
		minHeightSet = true
	}
	if flagset["f"] {
		floor = float32(floor64)
		minHeightSet = true
	}
	if flagset["ceiling"] {
		ceiling = float32(ceiling64)
		maxHeightSet = true
	}
	if flagset["c"] {
		ceiling = float32(ceiling64)
		maxHeightSet = true
	}

	out, err := os.Create(output)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	grid, err := esri.ReadGridFromFile(filename, verbose)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	// If floor or ceiling not already set, set them from the data.
	if !minHeightSet {
		floor = grid.MinHeight() - 0.1
	}

	if !maxHeightSet {
		ceiling = grid.MaxHeight() + 0.1
	}

	log.Printf("creating image - floor %f ceiling %f\n", floor, ceiling)
	img := image.NewRGBA(image.Rect(0, 0, grid.Nrows(), grid.Ncols()))
	maxRow := grid.Nrows() - 1
	for row := maxRow; row >= 0; row-- {
		for col := 0; col < grid.Ncols(); col++ {
			c := shade(floor, ceiling, grid.Height(row, col))
			if verbose {
				log.Printf("colouring cell[%d[%d] %d\n", row, col, c)
			}
			img.Set(col, row, c)
		}
	}

	log.Printf("encoding image")
	err = png.Encode(out, img)

	log.Printf("%d %d %f %f %d %d", grid.Nrows(), grid.Ncols(), grid.MinHeight(), grid.MaxHeight(), minShade, maxShade)
}

func shade(floor, ceiling, height float32) color.Color {
	// Get height and ceiling relative to the floor.
	height = height - floor
	ceiling = ceiling - floor
	shade := uint8(255 - uint8(height*256.0/ceiling))
	if verbose {
		log.Printf("shade %d", shade)
	}
	if maxShadeSet {
		if shade > maxShade {
			maxShade = shade
		}
	} else {
		maxShade = shade
		maxShadeSet = true
	}
	if minShadeSet {
		if shade < minShade {
			minShade = shade
		}
	} else {
		minShade = shade
		minShadeSet = true
	}
	return color.Gray{shade}
}
