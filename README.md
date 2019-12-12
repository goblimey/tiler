# tiler
Tiler takes an ESRI grid format file
and renders it as a monochrome png
with shades of grey representing the height values.

UK Environment Agency Lidar tiles are in ESRI grid format.
tq1652_DTM_1M.asc is their Digital Terrain Model (DTM)
covering the one kilometre map square TQ1652
with a cell size of one metre.
(A DTM has as much vegetation removed as possible,
so shows the ground surface.)

For a list of options:

    $ tiler -h

To process the DTM file and produce tq1652.png:

    $ tiler -i tq1652_DTM_1M.asc -o tq1652.png

TQ1652 covers a piece of ground between Dorking and Leatherhead.
The River Mole, the railway line and the highway are clearly visible.
In the top right corner the land rises steadily.
That's the lower slopes of Box Hill

The tiler draws a picture with one pixel per cell.
It uses 256 shades of grey to represent the height data in the point cloud.
The floor is drawn in white
and the ceiling is drawn in black.
Intermediate points are drawn in a shade of grey.

By default the floor is set to the lowest point in the file and
the ceiling is set to the highest point,
but you can override that.
If you want to compare two tiles,
you should use the same values for the floor and ceiling.

The -f options sets the height of the floor,
-c sets the ceiling.

