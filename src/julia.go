// Stefan Nilsson 2013-02-27

// This program creates pictures of Julia sets (en.wikipedia.org/wiki/Julia_set).
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type ComplexFunc func(complex128) complex128

var Funcs []ComplexFunc = []ComplexFunc{
	func(z complex128) complex128 { return z*z - 0.61803398875 },
	func(z complex128) complex128 { return z*z + complex(0, 1) },
	func(z complex128) complex128 { return z*z + complex(-0.835, -0.2321) },
	func(z complex128) complex128 { return z*z + complex(0.45, 0.1428) },
	func(z complex128) complex128 { return z*z*z + 0.400 },
	func(z complex128) complex128 { return cmplx.Exp(z*z*z) - 0.621 },
	func(z complex128) complex128 { return (z*z+z)/cmplx.Log(z) + complex(0.268, 0.060) },
	func(z complex128) complex128 { return cmplx.Sqrt(cmplx.Sinh(z*z)) + complex(0.065, 0.122) },
}

func main() {

	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu) // Try to use all available CPUs.
	var wg sync.WaitGroup
	start := time.Now()

	for n, fn := range Funcs {
		wg.Add(1)
		go func(n int, fn ComplexFunc) {
			err := CreatePng("picture-"+strconv.Itoa(n)+".png", fn, 1024)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(n, fn)
	}
	wg.Wait()
	fmt.Println("Took: ", time.Since(start))
}

// CreatePng creates a PNG picture file with a Julia image of size n x n.
func CreatePng(filename string, f ComplexFunc, n int) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	err = png.Encode(file, Julia(f, n))
	return
}

// Julia returns an image of size n x n of the Julia set for f.
func Julia(f ComplexFunc, n int) image.Image {

	bounds := image.Rect(-n/2, -n/2, n/2, n/2)
	img := image.NewRGBA(bounds)
	s := float64(n / 4)

	size := 300 //max(1, 2^8/n)

	var wg sync.WaitGroup

	for i, j := 0, size; i < n; i, j = j, j+size {
		if j > n {
			j = n
		}
		wg.Add(1)
		go func(i, j int) {
			for x := bounds.Min.X + i; x < bounds.Min.X+j; x++ {
				wg.Add(1)
				go func(x int) {
					for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
						n := Iterate(f, complex(float64(x)/s, float64(y)/s), 256)
						r := uint8(0)
						g := uint8(0)
						b := uint8(n % 32 * 8)
						img.Set(x, y, color.RGBA{r, g, b, 255})
					}
					wg.Done()
				}(x)
			}
			wg.Done()
		}(i, j)
	}
	wg.Wait()
	return img
}

// Iterate sets z_0 = z, and repeatedly computes z_n = f(z_{n-1}), n â‰¥ 1,
// until |z_n| > 2  or n = max and returns this n.
func Iterate(f ComplexFunc, z complex128, max int) (n int) {

	for ; n < max; n++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			break
		}
		z = f(z)
	}

	return
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
