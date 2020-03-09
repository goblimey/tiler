# Tiler

The tiler program takes an ESRI grid format file
and renders it as a monochrome png
with shades of grey representing the height values.

The tiler draws a picture with one pixel per grid cell
using 256 shades of grey to represent the height data in the grid.
The pixels at the lowest level (the floor) are drawn in white
and the ones at the highest level (the ceiling)
are drawn in black.
Pixels in between are drawn in a shade of grey.

## Buiding and Running

To build the tiler in a command window:

   go install github.com/goblimey/tiler

For a list of options:

    tiler -h

To process a file called in and produce a picture called out.png:

    tiler -i in -o out.png

By default the floor is set to the lowest point in the file and
the ceiling is set to the highest point,
but you can override that.
If you want to compare two tiles,
you should use the same values for the floor and ceiling.

The -f options sets the height of the floor,
-c sets the ceiling,
for example:

    tiler -i in -f 100 -c 1000 -o out.png

## Example data

tilt/tilt.txt is an ESRI grid that can be used for testing.
The highest point is at the top left corner and the lowest is at the bottom right.
The height reduces evenly in between.
The file is produced by tilt/tilt.go.

UK Environment Agency publish Lidar data covering large parts of Britain
as tiles in ESRI grid format.
The file tq1652_DTM_1M.asc is their Digital Terrain Model (DTM)
covering the one kilometre map square TQ1652
with a cell size of one metre.
(A DTM has as much vegetation removed as possible,
so shows the ground surface.)

Produce a picture from that file like so:

    tiler -i tq1652_DTM_1M.asc -o tq1652.png

The result covers a piece of ground to the South of Dorking in Surrey.
The River Mole, the railway line and the highway are clearly visible.
The lower slopes of Box Hill can be seen in the top right corner of the picture.