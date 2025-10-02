package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/png"
	"os"
)

func displayImage(img image.Image, maxWidth int) {
	bounds := img.Bounds()
	width := bounds.Dx()

	// Автоматически вычисляем scale для подгонки под maxWidth
	scale := 1
	if width > maxWidth {
		scale = width / maxWidth
		if scale < 1 {
			scale = 1
		}
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y += scale * 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x += scale {
			if y+scale >= bounds.Max.Y {
				break
			}

			r1, g1, b1, a1 := img.At(x, y).RGBA()
			r2, g2, b2, a2 := img.At(x, y+scale).RGBA()

			if a1 == 0 && a2 == 0 {
				fmt.Print(" ")
				continue
			}

			fmt.Printf("\033[38;2;%d;%d;%d;48;2;%d;%d;%dm▀\033[0m",
				r1>>8, g1>>8, b1>>8,
				r2>>8, g2>>8, b2>>8)
		}
		fmt.Println()
	}
}

func main() {
	maxWidth := flag.Int("width", 40, "Максимальная ширина в символах")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: go run main.go [-width N] image.png")
		fmt.Println("Например: go run main.go -width 20 sprite.png")
		return
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	bounds := img.Bounds()
	fmt.Printf("Изображение: %dx%d (%s)\n", bounds.Dx(), bounds.Dy(), format)

	displayImage(img, *maxWidth)
}

// func displayImage(img image.Image) {
// 	bounds := img.Bounds()

// 	// Масштабируем до разумного размера
// 	scale := 2

// 	for y := bounds.Min.Y; y < bounds.Max.Y; y += scale * 2 {
// 		for x := bounds.Min.X; x < bounds.Max.X; x += scale {
// 			if y+scale >= bounds.Max.Y {
// 				break
// 			}

// 			// Верхний пиксель
// 			r1, g1, b1, a1 := img.At(x, y).RGBA()
// 			// Нижний пиксель
// 			r2, g2, b2, a2 := img.At(x, y+scale).RGBA()

// 			// Пропускаем прозрачные
// 			if a1 == 0 && a2 == 0 {
// 				fmt.Print(" ")
// 				continue
// 			}

// 			// True color (24-bit)
// 			fmt.Printf("\033[38;2;%d;%d;%d;48;2;%d;%d;%dm▀\033[0m",
// 				r1>>8, g1>>8, b1>>8,
// 				r2>>8, g2>>8, b2>>8)
// 		}
// 		fmt.Println()
// 	}
// }

// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("Usage: go run main.go image.png")
// 		return
// 	}

// 	file, err := os.Open(os.Args[1])
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	defer file.Close()

// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}

// 	displayImage(img)
// }
