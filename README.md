![bonsai](https://user-images.githubusercontent.com/7358/42061154-c7a36824-7aee-11e8-822a-486bf40dc065.png)

Ever since [Michael Fogelman](https://www.michaelfogleman.com/) posted his [Top 10 Projects from 2014](https://medium.com/@fogleman/my-top-10-side-projects-from-2014-713a78d6fc9d), I've been following him just to remind myself how prolific one person can be.  One of his projects was a [really cool quadtree art generator](https://www.michaelfogleman.com/static/quads/) and [I just had to make my own](https://github.com/bwiggs/quadtree-art)!

## Quadtrees

[Quadtrees](https://en.wikipedia.org/wiki/Quadtree) aren't something I was familiar with, so it was fun getting to work with a new data structure.

> A **quadtree** is a tree data structure in which each internal node has exactly four children.

## Art Generation

The basic idea is that you recursively divide an image into 4 quadrants (hence the name *quad*-tree) either filling that quadrant with the average color of that quad or subdivide it.

1. Choose some threshold between 0 and 1. This threshold represents the amount of color difference required in a quadrant to trigger a recursive call.
2. Loop through each quadrant
  1. calculate the amount of color difference of all the pixels in that quad
  3. if the amount of color difference is more than your threshold
    1. recursively process this quadrant
  4. otherwise fill that quadrant with the average color of the quadrant.

### Calculating Color Difference

Essentially, loop through all the pixels in this quad adding up all the RGP color values.
Divide the total amount of color (`colorSum`) by the number of pixels to get your average color distance.

```go
func (q *quad) calcAvgSimpleColorDistance() {

	var colorSum float64

	for y := q.y; y < q.y+q.height; y++ {
		for x := q.x; x < q.x+q.width; x++ {
			// use _ to ignore the alpha channel since we don't care about that
			r, g, b, _ := (*q.img).At(x, y).RGBA()
			colorSum += math.Abs(float64(int32(q.color.R) - int32(r>>8)))
			colorSum += math.Abs(float64(int32(q.color.G) - int32(g>>8)))
			colorSum += math.Abs(float64(int32(q.color.B) - int32(b>>8)))
		}
	}

	q.colorDelta = colorSum / float64(3*q.width*q.height)
}
```

### Calculating Average Color

To get the average color for the quadrant, add up all the channels for each pixel individually.
Then divide the total almount for each channel by the size of the area.

```go
func (q *quad) calcAvgColor() {
	h := histogram{}

	for y := q.y; y < q.y+q.height; y++ {
		for x := q.x; x < q.x+q.width; x++ {
			r, g, b, a := (*q.img).At(x, y).RGBA()
			h.r += r >> 8
			h.g += g >> 8
			h.b += b >> 8
			h.a += a >> 8
		}
	}

	area := uint32(q.width * q.height)
	h.r = h.r / area
	h.g = h.g / area
	h.b = h.b / area
	h.a = h.a / area

	q.color = color.RGBA{
		uint8(h.r),
		uint8(h.g),
		uint8(h.b),
		uint8(h.a),
	}
}
```

### CLI Tool

I've built a few interesting features into this tool to generate different output from the same image.
Some options are:

- adding a border color to each quadrant
- changing the max recursion level to generate more detailed or blocky images
- setting a line color that is different from the image border
- setting the minimum block size to also generate more blocky images


```shell
$ quadtree-art -h
Usage of quadtree-art:
  -bc string
    	border color (hex) (default "333333")
  -l int
    	max recursive levels (default 7)
  -lc string
    	line color (hex)
  -m int
    	minimum size a block can be (default 1)
  -n	do not render block outlines
  -o string
    	output file name with extension (default "art.png")
  -t float
    	minimum size a block can be (default 25)
```

![sushi](https://user-images.githubusercontent.com/7358/42061162-d21afce0-7aee-11e8-83e2-ffaab1da5870.png)
