package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
)

type histogram struct {
	r, g, b, a uint32
}

type quad struct {
	x, y, width, height int
	children            []*quad
	img                 *image.Image
	colorDelta          float64
	color               color.RGBA
	threshold           float64
	maxDepth            int32
	currDepth           int32
}

func newQuad(img *image.Image, x, y, width, height int, t float64, maxDepth int32, currDepth int32) *quad {
	q := quad{
		img:       img,
		x:         x,
		y:         y,
		width:     width,
		height:    height,
		threshold: t,
		currDepth: currDepth,
		maxDepth:  maxDepth,
	}
	// spew.Dump(q)
	q.calcAvgColor()
	q.calcAvgSimpleColorDistance()
	// q.calcAvgEuclidianColorDistance()
	q.subdivide()
	return &q
}

func (q quad) String() string {
	return fmt.Sprintf("<%3d,%3d> %3dx%3d D:%d Δ:%9f T:%d %x", q.x, q.y, q.width, q.height, q.currDepth, q.colorDelta, int(q.threshold), q.color)
}

func (q *quad) calcAvgSimpleColorDistance() {

	var colorSum float64

	for y := q.y; y < q.y+q.height; y++ {
		for x := q.x; x < q.x+q.width; x++ {
			r, g, b, _ := (*q.img).At(x, y).RGBA()
			colorSum += math.Abs(float64(int32(q.color.R) - int32(r>>8)))
			colorSum += math.Abs(float64(int32(q.color.G) - int32(g>>8)))
			colorSum += math.Abs(float64(int32(q.color.B) - int32(b>>8)))
		}
	}

	q.colorDelta = colorSum / float64(3*q.width*q.height)
}

func (q *quad) calcAvgEuclidianColorDistance() {
	var runningDistance float64
	for y := q.y; y < q.y+q.height; y++ {
		for x := q.x; x < q.x+q.width; x++ {
			r, g, b, _ := (*q.img).At(x, y).RGBA()
			runningDistance += math.Sqrt(math.Pow(float64(int32(q.color.R)-int32(r)), 2) +
				math.Pow(float64(int32(q.color.G)-int32(g)), 2) +
				math.Pow(float64(int32(q.color.B)-int32(b)), 2))
		}
	}
	q.colorDelta = runningDistance / float64(q.width*q.height)
}

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

func (q *quad) subdivide() {
	// spew.Dump(q)
	if q.colorDelta > q.threshold && (q.maxDepth == 0 || q.currDepth < q.maxDepth) && q.width > PixelMin && q.height > PixelMin {
		q.children = []*quad{
			newQuad(q.img, q.x, q.y, q.width/2, q.height/2, q.threshold, q.maxDepth, q.currDepth+1),
			newQuad(q.img, q.x+q.width/2, q.y, q.width/2, q.height/2, q.threshold, q.maxDepth, q.currDepth+1),
			newQuad(q.img, q.x, q.y+q.height/2, q.width/2, q.height/2, q.threshold, q.maxDepth, q.currDepth+1),
			newQuad(q.img, q.x+q.width/2, q.y+q.height/2, q.width/2, q.height/2, q.threshold, q.maxDepth, q.currDepth+1),
		}

		for i := range q.children {
			q.children[i].subdivide()
		}
	}
}

func (q *quad) draw(c *image.RGBA) {
	if q.children != nil {
		for i := range q.children {
			q.children[i].draw(c)
		}
	} else {
		draw.Draw(c, image.Rect(q.x, q.y, q.x+q.width, q.y+q.height), &image.Uniform{q.color}, image.Point{q.x, q.y}, draw.Src)
		if Outlines {
			// top
			draw.Draw(c, image.Rect(q.x, q.y, q.x+q.width, q.y+1), &image.Uniform{LineColor}, image.Point{q.x, q.y}, draw.Src)
			// left
			draw.Draw(c, image.Rect(q.x, q.y, q.x+1, q.y+q.height), &image.Uniform{LineColor}, image.Point{q.x, q.y}, draw.Src)
		}
	}

	if q.currDepth == 1 {
		draw.Draw(c, image.Rect(q.x, q.y, q.x+q.width, q.y+1), &image.Uniform{BorderColor}, image.Point{q.x, q.y}, draw.Src)
		draw.Draw(c, image.Rect(q.x, q.y, q.x+1, q.y+q.height), &image.Uniform{BorderColor}, image.Point{q.x, q.y}, draw.Src)
		draw.Draw(c, image.Rect(q.x, q.y+q.height-1, q.x+q.width, q.y+q.height), &image.Uniform{BorderColor}, image.Point{q.x, q.y}, draw.Src)
		draw.Draw(c, image.Rect(q.x+q.width-1, q.y, q.x+q.width, q.y+q.height), &image.Uniform{BorderColor}, image.Point{q.x, q.y}, draw.Src)
	}

}
