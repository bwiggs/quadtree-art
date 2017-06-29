package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
)

const PixelMin = 1
const Threshold = 25
const Levels = 7
const Outlines = true

func main() {
	fmt.Println("Quadtree Art generation with GoLang")

	// open the image to quadify

	reader, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := img.Bounds()

	// process the image

	fmt.Println("processing image")
	q := newQuad(&img, bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y, Threshold, Levels, 1)

	// save art to the filesystem
	fmt.Println("rendering artwork")
	canvas := image.NewRGBA(image.Rect(0, 0, q.width, q.height))
	q.draw(canvas)
	o, _ := os.Create("art.jpg")
	defer o.Close()
	jpeg.Encode(o, canvas, &jpeg.Options{100})
}
