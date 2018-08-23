package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"

	"golang.org/x/image/bmp"
)

var (
	PixelMin         int
	Threshold        float64
	Levels           int
	ShowGrid         bool
	oFile            string
	InputBorderColor string
	InputGridColor   string
	BorderColor      color.RGBA
	GridColor        color.RGBA
	AverageColor     color.RGBA
)

func init() {
	flag.IntVar(&PixelMin, "m", 1, "minimum size a block can be")
	flag.Float64Var(&Threshold, "t", 25, "color difference threshold")
	flag.IntVar(&Levels, "l", 7, "max recursive levels")
	flag.BoolVar(&ShowGrid, "g", false, "render grid lines")
	flag.StringVar(&InputBorderColor, "bc", "333333", "border color (hex)")
	flag.StringVar(&InputGridColor, "gc", "", "grid color (hex)")
	flag.StringVar(&oFile, "o", "quad.png", "output file name with extension")
	flag.Parse()
}

func main() {
	fmt.Println("Quadtree Art generation with GoLang")

	// open the image to quadify
	reader, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := img.Bounds()

	fmt.Println("processing image")
	q := newQuad(&img, bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y, Threshold, int32(Levels), 1)

	// use the input color if available or the default color for borders.
	if len(InputBorderColor) == 6 {
		colorHex, _ := hex.DecodeString(InputBorderColor)
		BorderColor = color.RGBA{uint8(colorHex[0]), uint8(colorHex[1]), uint8(colorHex[2]), 0xff}
	} else {
		// use the average color of the image
		BorderColor = q.color
	}

	// use the input color if available or the default color for borders.
	if len(InputGridColor) == 6 {
		colorHex, _ := hex.DecodeString(InputGridColor)
		GridColor = color.RGBA{uint8(colorHex[0]), uint8(colorHex[1]), uint8(colorHex[2]), 0xff}
	} else {
		// use the average color of the image
		GridColor = q.color
	}

	fmt.Println("rendering artwork")
	canvas := image.NewRGBA(image.Rect(0, 0, q.width, q.height))
	q.draw(canvas)

	// save art to the filesystem
	o, _ := os.Create(oFile)
	defer o.Close()
	ext := strings.Split(oFile, ".")[1]
	switch ext {
	case "jpg":
		jpeg.Encode(o, canvas, &jpeg.Options{100})
	case "png":
		png.Encode(o, canvas)
	case "bmp":
		bmp.Encode(o, canvas)
	}
}
