package main

import "image"
import "image/color"
import _ "image/jpeg"
import "image/png"
import "os"
import "log"
import "math"
import "fmt"

func tofloat(p color.Color) ([]float64) {
	r, g, b, a := p.RGBA()
	return []float64{float64(r) / 65535.0,
			float64(g) / 65535.0,
			float64(b) / 65535.0,
			float64(a) / 65535.0}
}

func main() {
	f1, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("Old image could not be opened")
	}
	f2, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal("New image could not be opened")
	}
	img1, _, err := image.Decode(f1)
	if err != nil {
		log.Fatal("Old image could not be decoded")
	}
	img2, _, err := image.Decode(f2)
	if err != nil {
		log.Fatal("New image could not be decoded")
	}
	
	if img1.Bounds() != img2.Bounds() {
		log.Fatal("Images don't have the same size")
	}
	
	totaldiff := 0.0
	mapimg := image.NewNRGBA(img1.Bounds())

	for y := img1.Bounds().Min.Y; y < img1.Bounds().Max.Y; y += 1 {
		for x := img1.Bounds().Min.X; x < img1.Bounds().Max.X; x += 1 {
			p1 := tofloat(img1.At(x, y))
			p2 := tofloat(img2.At(x, y))

			var totplus float64 = 0.0
			absdiff := math.Abs(p2[3] - p1[3])
			diffpixel := []float64{1.0, 1.0, 1.0}

			for i := 0; i < 3; i += 1 {
				diff := p2[i] - p1[i]
				absdiff += math.Abs(diff)
				totplus += math.Max(0, diff)
				diffpixel[i] += diff
			}

			totaldiff += absdiff

			for i := 0; i < 3; i += 1 {
				diffpixel[i] -= totplus
				if absdiff > 0 && absdiff < (5.0 / 255.0) {
					diffpixel[i] -= (5.0 / 255.0)
				}
				diffpixel[i] = math.Max(0.0, diffpixel[i])
			}

			p := color.NRGBA{R: uint8(diffpixel[0] * 255.0),
					G: uint8(diffpixel[1] * 255.0), 
					B: uint8(diffpixel[2] * 255.0),
					A: 255}
			mapimg.Set(x, y, p)
		}
	}

	mapfile, _ := os.Create(os.Args[3])
	defer mapfile.Close()

    	png.Encode(mapfile, mapimg)
	fmt.Printf("Difference: %v\n", int(0.5 + totaldiff * 255.0))
}
