package esri

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// Grid defines a data structure that holds a 3D ESRI Grid read from a
// data file. The file contains a rectangular Grid of height values describing a surface
// within some map.  ESRI Grid format files are used in mapping.
//
// ReadGridFromFile() reads a data file and fills in the grid object.
//
// The grid files contain plain text.  This is a very simple example:
//
// ncols 4
// nrows 4
// xllcorner    513000
// yllcorner    152000
// cellsize     1
// NODATA_value -9999
// 500 500 500 500
// 500 500 500 500
// 1000 1000 1000 1000
// 1000 1000 1000 1000
//
// The file starts with six header lines defining the rest of the data.  ncols is the number
// of columns, nrows the number of rows.  xllcorner gives the x map reference of the bottom
// left corner of the grid, yllcorner the y map reference.  cellsize is the size of the cells
// (the grid).  The NODATA value is used for points on the grid where the sensor couldn't
// figure out the height.
//
// The header lines are followed by the rows and columns of height data.  The values can be
// floating point numbers, here they are integers.  This example defines a four by four grid.
// The first row defines the top (most northern) line of the grid and the last row describes
// the bottom line.  In the real world the grid file would be much bigger, for example
// 1,000 by 1,000 points across 1m cells, giving a 1Km square.
//
// Some of the values only make sense in the context of a local mapping system, for example
// UK point clouds use Ordnance Survey map references for xllcorner and yllcorner, and the
// cell sizes are in metres.
//
type Grid struct {
	ncols        int
	nrows        int
	xllcorner    float32
	yllcorner    float32
	cellsize     float32
	noDataValue  int
	maxHeightSet bool
	maxHeight    float32
	minHeightSet bool
	minHeight    float32
	height       [][]float32
	verbose      bool
}

//ReadGridFromFile is a factory method that reads data from an ESRI Grid
// format file and returns a Grid object.
//
func ReadGridFromFile(filename string, verbose bool) (*Grid, error) {
	m := "ReadGridFromFile"
	if verbose {
		log.Printf("%s: %s", m, filename)
	}

	in, err := os.Open(filename)
	if err != nil {
		log.Printf(filename + err.Error())
		return nil, err
	}

	grid := new(Grid)

	r := bufio.NewReader(in)

	lineNum := 0
	fieldName := "ncols"
	grid.ncols, err = readIntFromHeader(r, fieldName, verbose)
	if err != nil {
		return nil, err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %d", m, fieldName, grid.ncols)
	}

	fieldName = "nrows"
	grid.nrows, err = readIntFromHeader(r, fieldName, verbose)
	if err != nil {
		return nil, err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %d", m, fieldName, grid.nrows)
	}

	grid.height = make([][]float32, grid.nrows)

	for i := 0; i < grid.nrows; i++ {
		grid.height[i] = make([]float32, grid.ncols)
	}

	fieldName = "xllcorner"
	grid.xllcorner, err = readFloat32FromHeader(r, fieldName, verbose)
	if err != nil {
		return nil, err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %f", m, fieldName, grid.xllcorner)
	}

	fieldName = "yllcorner"
	grid.yllcorner, err = readFloat32FromHeader(r, fieldName, verbose)
	if err != nil {
		return nil, err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %f", m, fieldName, grid.yllcorner)
	}

	fieldName = "cellsize"
	grid.cellsize, err = readFloat32FromHeader(r, fieldName, verbose)
	if err != nil {
		return nil, err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %f", m, fieldName, grid.cellsize)
	}

	fieldName = "NODATA_value"
	grid.noDataValue, err = readIntFromHeader(r, fieldName, verbose)
	if err != nil {
		return nil, err
	}
	lineNum++

	log.Printf("NODATA_value %d", grid.noDataValue)

	// Read nrows of lines each containing ncols floats, space separated.
	log.Printf("%s: reading %d data lines", m, grid.nrows)

	linesExpected := grid.nrows + 6

	for row := 0; ; row++ {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		lineNum++
		if lineNum > linesExpected {
			log.Printf("%s: warning: file %s has too many lines - expected %d\n", m, filename, linesExpected)
			break
		}
		line, err = stripSpaces(line)
		if err != nil {
			log.Printf("%s: stripSpaces failed - %s", m, err.Error())
			return nil, err
		}
		if verbose {
			log.Println(line)
		}

		numbers := strings.Split(line, " ")
		if len(numbers) > grid.ncols {
			log.Printf("warning: line %d has too many columns - got %d expected %d\n",
				lineNum, len(numbers), grid.ncols)
			continue
		}
		if len(numbers) < grid.ncols {
			log.Printf("warning: line %d has too few columns - got %d expected %d\n",
				lineNum, len(numbers), grid.ncols)
			continue
		}
		for col := range numbers {
			var f float32
			_, err := fmt.Sscanf(numbers[col], "%f", &f)
			if err != nil {
				log.Printf("%d %d %s", row, col, err.Error())
				return nil, err
			}

			// Set height, maxheight and minHeight
			grid.SetHeight(row, col, f)

			if verbose {
				log.Printf("height[%d][%d] %f", row, col, grid.height[row][col])
			}
		}
	}

	if lineNum < linesExpected {
		log.Printf("warning: file %s has too few lines - got %d expected %d\n",
			filename, lineNum, linesExpected)
	}

	if verbose {
		log.Printf("maxHeight %f minheight %f", grid.maxHeight, grid.minHeight)
	}

	return grid, nil
}

// Ncols returns the number of columns in the Grid.
func (g Grid) Ncols() int {
	return g.ncols
}

// Nrows returns the number of rows in the Grid.
func (g Grid) Nrows() int {
	return g.nrows
}

// Xllcorner returns the x coordinate of the lower left corner of the Grid.
func (g Grid) Xllcorner() float32 {
	return g.xllcorner
}

// Yllcorner returns the y coordinate of the lower left corner of the Grid.
func (g Grid) Yllcorner() float32 {
	return g.yllcorner
}

// CellSize returns the size of the Grid cells in metres.
func (g Grid) CellSize() float32 {
	return g.cellsize
}

// NoDataValue returns the No Data value.
func (g Grid) NoDataValue() int {
	return g.noDataValue
}

// MaxHeight returns the largest height reading in the Grid.
func (g Grid) MaxHeight() float32 {
	return g.maxHeight
}

// MinHeight returns the smallest height reading in the Grid.
func (g Grid) MinHeight() float32 {
	return g.minHeight
}

// SetNCols sets the number of columns in the Grid.
func (g *Grid) SetNCols(ncols int) {
	g.ncols = ncols
}

// SetNRows sets the number of rows in the Grid.
func (g *Grid) SetNRows(nrows int) {
	g.nrows = nrows
}

// SetXllcorner sets the x coordinate of the lower left corner of the Grid.
func (g *Grid) SetXllcorner(xllcorner float32) {
	g.xllcorner = xllcorner
}

// SetYllcorner sets the y coordinate of the lower left corner of the Grid.
func (g *Grid) SetYllcorner(yllcorner float32) {
	g.yllcorner = yllcorner
}

// SetCellSize sets the size of the grid cells in metres.
func (g *Grid) SetCellSize(cellsize float32) {
	g.cellsize = cellsize
}

// SetNoData sets the No Data value.
func (g *Grid) SetNoDataValue(noDataValue int) {
	g.noDataValue = noDataValue
}

// Height gets the height of cell (row, col).
func (g Grid) Height(row, col int) float32 {
	return g.height[row][col]
}

// SetHeight sets the height of cell (row, col).
func (g *Grid) SetHeight(row, col int, height float32) {

	if row >= g.nrows || col >= g.ncols {
		log.Printf("SetHeight(%d,%d) - out of range", row, col)
		return
	}
	g.height[row][col] = height

	if g.maxHeightSet {
		if height > g.maxHeight {
			g.maxHeight = height
		}
	} else {
		g.maxHeight = height
		g.maxHeightSet = true
	}

	if g.minHeightSet {
		if height < g.minHeight {
			g.minHeight = height
		}
	} else {
		g.minHeight = height
		g.minHeightSet = true
	}
}

func readIntFromHeader(r *bufio.Reader, fieldName string, verbose bool) (int, error) {
	m := "readIntHeader"
	line, err := r.ReadString('\n')
	if err != nil {
		return 0, err
	}
	if verbose {
		log.Printf("%s: line %s", m, line)
	}
	line, err = stripSpaces(line)
	field := strings.Split(line, " ")
	if field[0] != fieldName {
		log.Printf("%s: expected %s, got %s", m, fieldName, line)
	}
	var result int
	_, err = fmt.Sscanf(field[1], "%d", &result)
	if err != nil {
		return 0, err
	}
	if verbose {
		log.Printf("%s: %s %d", m, fieldName, result)
	}

	return result, nil
}

func readFloat32FromHeader(r *bufio.Reader, fieldName string, verbose bool) (float32, error) {
	m := "readFloat32FromHeader"
	line, err := r.ReadString('\n')
	if err != nil {
		return 0, err
	}
	if verbose {
		log.Printf("%s: line %s", m, line)
	}
	line, err = stripSpaces(line)
	field := strings.Split(line, " ")
	if field[0] != fieldName {
		log.Printf("%s: expected %s, got %s", m, fieldName, line)
	}
	var result float32
	_, err = fmt.Sscanf(field[1], "%f", &result)
	if err != nil {
		return 0, err
	}
	if verbose {
		log.Printf("%s: %s %f", m, fieldName, result)
	}

	return result, nil
}

func stripSpaces(s string) (string, error) {
	// Remove spaces from the beginning and the end of the staring.
	s = strings.TrimSpace(s)
	// Reduce multiple adjacent spaces within the string to a single space.
	re, err := regexp.Compile("  +")
	if err != nil {
		return s, err
	}
	return re.ReplaceAllLiteralString(s, " "), nil
}
