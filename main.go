package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"
)

type ImagePoint struct {
	image.Point
	color color.RGBA
}

var (
	imgHeight = flag.Int("h", 640, "image height")
	imgWidth  = flag.Int("w", 480, "image width")
	points    = flag.Int("p", 8, "amount of points & colors")
)

var (
	colorPoints []ImagePoint
	img         *image.RGBA
)

// initialize global variables, rand seed & flags
func init() {
	flag.Parse()
	rand.Seed(time.Now().Unix())

	img = image.NewRGBA(image.Rectangle{
		Min: image.Pt(0, 0),
		Max: image.Pt(*imgWidth, *imgHeight),
	})

	colorPoints = make([]ImagePoint, *points)
}

func main() {
	// setup color palette
	for i := range colorPoints {
		x := rand.Intn(*imgWidth)
		y := rand.Intn(*imgHeight)

		colorPoints[i] = ImagePoint{
			Point: image.Pt(x, y),
			color: color.RGBA{
				R: uint8(rand.Intn(255)),
				G: uint8(rand.Intn(255)),
				B: uint8(rand.Intn(255)),
				A: 255,
			},
		}
	}

	var wg sync.WaitGroup

	for y := 0; y < *imgHeight; y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			processLine(y)
		}(y)
	}

	wg.Wait()

	f, _ := os.Create("image.png")
	png.Encode(f, img)
}

// processLine takes a line number "y" and calculate every pixel color on that line
func processLine(y int) {
	type distanceColor struct {
		color    color.RGBA
		distance float64
	}
	for x := 0; x < *imgWidth; x++ {
		distances := make([]distanceColor, len(colorPoints))

		// loop over every entry and calculate the distance from (x,y) to each point on the colorPoints array
		for ip := 0; ip < len(colorPoints); ip++ {
			// formula: d=√((x2 – x1)² + (y2 – y1)²)
			distance := math.Sqrt(math.Pow(float64(colorPoints[ip].X-x), 2) + math.Pow(float64(colorPoints[ip].Y-y), 2))
			distances[ip] = distanceColor{
				distance: distance,
				color:    colorPoints[ip].color,
			}
		}

		// sort the by distance so the lowest value will be index 0
		sort.SliceStable(distances[:], func(i, j int) bool {
			return distances[i].distance < distances[j].distance
		})

		// now set the pixel at (x,y) with the closest point's color
		img.Set(x, y, distances[0].color)
	}
}
