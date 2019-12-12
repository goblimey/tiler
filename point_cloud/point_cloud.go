package pointCloud

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// The PointCloud interface defines a data structure that holds a 3D point cloud map read from a 
// data file. A Point Cloud map gives a rectangular grid of height values describing a surface 
// within some map.  Point Cloud format files are used in mapping.  Equipment such as Lidar 
// mapping sensors produce point cloud format data and there is a host of software available to 
// process and visualise it.
//
// The interface provides a method that reads a data file and produces
//
// Point Cloud data is held in plain text files.  This is a very simple example:
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
// the bottom line.  In the real world the point cloud would be much bigger, for example 
// 1,000 by 1,000 points across 1m cells, giving a 1Km square.
//
// Some of the values only make sense in the context of a local mapping system, for example
// UK point clouds use Ordnance Survey map references for xllcorner and yllcorner, and the 
// cell sizes are in metres.  

type PointCloud interface {
	// NCols returns the number of columns
	Ncols() int
	// Nrows returns the 
	Nrows() int
	Xllcorner() int
	Yllcorner() int
	CellSize() int
	NoDataValue() int
	MaxHeight() float32
	MinHeight() float32
	SetNCols(ncols int)
	SetNRows(nrows int)
	SetXllcorner(xllcorner int)
	SetYllcorner(yllcorner int)
	SetCellSize(cellsize int)
	SetNoDataValue(noDataValue int)
	Height(i, col int) float32
	HeightAtCoordinate(x, y int) float32
	SetHeight(i, col int, height float32)
	ReadPointCloudFromFile(filename string, verbose bool) error 
}

type ConcretePointCloud struct {
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

func (pc ConcretePointCloud) Ncols() int {
	return pc.ncols
}

func (pc ConcretePointCloud) Nrows() int {
	return pc.nrows
}

func (pc ConcretePointCloud) Xllcorner() float32 {
	return pc.xllcorner
}

func (pc ConcretePointCloud) Yllcorner() float32 {
	return pc.yllcorner
}

func (pc ConcretePointCloud) CellSize() float32 {
	return pc.cellsize
}

func (pc ConcretePointCloud) NoDataValue() int {
	return pc.noDataValue
}

func (pc ConcretePointCloud) MaxHeight() float32 {
	return pc.maxHeight
}

func (pc ConcretePointCloud) MinHeight() float32 {
	return pc.minHeight
}

func (pc *ConcretePointCloud) SetNCols(ncols int) {
	pc.ncols = ncols
}

func (pc *ConcretePointCloud) SetNRows(nrows int) {
	pc.nrows = nrows
}

func (pc *ConcretePointCloud) SetXllcorner(xllcorner float32) {
	pc.xllcorner = xllcorner
}

func (pc *ConcretePointCloud) SetYllcorner(yllcorner float32) {
	pc.yllcorner = yllcorner
}

func (pc *ConcretePointCloud) SetCellSize(cellsize float32) {
	pc.cellsize = cellsize
}

func (pc *ConcretePointCloud) SetNoDataValue(noDataValue int) {
	pc.noDataValue = noDataValue
}

func (pc ConcretePointCloud) Height(i, col int) float32 {
	return pc.height[i][col]
}

func (pc *ConcretePointCloud) SetHeight(row, col int, height float32) {

	if row >= pc.nrows || col >= pc.ncols {
		log.Printf("SetHeight(%d,%d) - out of range", row, col)
		return
	} 
	pc.height[row][col] = height

	if pc.maxHeightSet {
		if height > pc.maxHeight {
			pc.maxHeight = height
		}
	} else {
		pc.maxHeight = height
		pc.maxHeightSet = true
	}

	if pc.minHeightSet {
		if height < pc.minHeight {
			pc.minHeight = height
		}
	} else {
		pc.minHeight = height
		pc.minHeightSet = true
	}
}

func (pc *ConcretePointCloud) ReadPointCloudFromFile(filename string, verbose bool) error {
	m := "ReadPointCloudFromFile"
	if verbose {
		log.Printf("%s: %s", m, filename)
	}

	in, err := os.Open(filename)
	if err != nil {
		log.Printf(filename + err.Error())
		return err
	}

	r := bufio.NewReader(in)

	lineNum := 0
	fieldName := "ncols"
	pc.ncols, err = readIntFromHeader(r, fieldName, verbose)
	if err != nil {
		return err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %d", m, fieldName, pc.ncols)
	}

	fieldName = "nrows"
	pc.nrows, err = readIntFromHeader(r, fieldName, verbose)
	if err != nil {
		return err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %d", m, fieldName, pc.nrows)
	}

	pc.height = make([][]float32, pc.nrows)

	for i := 0; i < pc.nrows; i++ {
		pc.height[i] = make([]float32, pc.ncols)
	}

	fieldName = "xllcorner"
	pc.xllcorner, err = readFloat32FromHeader(r, fieldName, verbose)
	if err != nil {
		return err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %f", m, fieldName, pc.xllcorner)
	}

	fieldName = "yllcorner"
	pc.yllcorner, err = readFloat32FromHeader(r, fieldName, verbose)
	if err != nil {
		return err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %f", m, fieldName, pc.yllcorner)
	}

	fieldName = "cellsize"
	pc.cellsize, err = readFloat32FromHeader(r, fieldName, verbose)
	if err != nil {
		return err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %f", m, fieldName, pc.cellsize)
	}

	fieldName = "NODATA_value"
	pc.noDataValue, err = readIntFromHeader(r, fieldName, verbose)
	if err != nil {
		return err
	}
	lineNum++

	log.Printf("NODATA_value %d", pc.noDataValue)

	// Read nrows of lines each containing ncols floats, space separated.
	log.Printf("%s: reading %d data lines", m, pc.nrows)

	linesExpected := pc.nrows + 6

	for row := 0;; row++ {
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
			return err
		}
		if verbose {
			log.Println(line)
		}

		numbers := strings.Split(line, " ")
		if len(numbers) > pc.ncols {
			log.Printf("warning: line %d has too many columns - got %d expected %d\n",
				lineNum, len(numbers), pc.ncols)
			continue
		}
		if len(numbers) < pc.ncols {
			log.Printf("warning: line %d has too few columns - got %d expected %d\n",
				lineNum, len(numbers), pc.ncols)
			continue
		}
		for col := range numbers {
			var f float32
			_, err := fmt.Sscanf(numbers[col], "%f", &f)
			if err != nil {
				log.Printf("%d %d %s", row, col, err.Error())
				return err
			}

			// Set height, maxheight and minHeight
			pc.SetHeight(row, col, f)

			if verbose {
				log.Printf("height[%d][%d] %f", row, col, pc.height[row][col])
			}
		}
	}

	if lineNum < linesExpected {
		log.Printf("warning: file %s has too few lines - got %d expected %d\n",
			filename, lineNum, linesExpected)
	}

	if verbose {
		log.Printf("maxHeight %f minheight %f", pc.maxHeight, pc.minHeight)
	}

	return nil
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
	s = strings.TrimSpace(s)
	re, err := regexp.Compile("  +")
	if err != nil {
		return s, err
	}

	return re.ReplaceAllLiteralString(s, " "), nil
}
