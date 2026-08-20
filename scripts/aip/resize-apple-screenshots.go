package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	if len(os.Args) == 6 && os.Args[1] == "--pad" {
		width, _ := strconv.Atoi(os.Args[2])
		height, _ := strconv.Atoi(os.Args[3])
		paths, err := filepath.Glob(os.Args[4])
		if err != nil {
			panic(err)
		}
		if err := os.MkdirAll(os.Args[5], 0755); err != nil {
			panic(err)
		}
		for _, path := range paths {
			output := filepath.Join(os.Args[5], filepath.Base(path))
			if err := pad(path, output, width, height); err != nil {
				panic(err)
			}
			fmt.Printf("created %s at %dx%d\n", output, width, height)
		}
		return
	}
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: resize-apple-screenshots <width> <height> <glob>")
		fmt.Fprintln(os.Stderr, "   or: resize-apple-screenshots --pad <width> <height> <glob> <output-dir>")
		os.Exit(2)
	}
	width, err := strconv.Atoi(os.Args[1])
	if err != nil || width <= 0 {
		panic("invalid width")
	}
	height, err := strconv.Atoi(os.Args[2])
	if err != nil || height <= 0 {
		panic("invalid height")
	}
	paths, err := filepath.Glob(os.Args[3])
	if err != nil {
		panic(err)
	}
	for _, path := range paths {
		if err := resize(path, width, height); err != nil {
			panic(err)
		}
		fmt.Printf("resized %s to %dx%d\n", path, width, height)
	}
}

func pad(inputPath, outputPath string, width, height int) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	source, _, err := image.Decode(input)
	input.Close()
	if err != nil {
		return err
	}

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	background := color.RGBA{R: 248, G: 244, B: 252, A: 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			canvas.Set(x, y, background)
		}
	}
	scale := float64(height) / float64(source.Bounds().Dy())
	scaledWidth := int(float64(source.Bounds().Dx())*scale + 0.5)
	left := (width - scaledWidth) / 2
	for y := 0; y < height; y++ {
		for x := 0; x < scaledWidth; x++ {
			canvas.Set(left+x, y, bilinear(source, float64(x)/scale, float64(y)/scale))
		}
	}
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	if err := writeRGBPNG(file, canvas); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

func resize(path string, width int, height int) error {
	input, err := os.Open(path)
	if err != nil {
		return err
	}
	source, _, err := image.Decode(input)
	input.Close()
	if err != nil {
		return err
	}

	output := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			output.Set(x, y, bilinear(source, float64(x)*float64(source.Bounds().Dx())/float64(width), float64(y)*float64(source.Bounds().Dy())/float64(height)))
		}
	}

	temporary := path + ".resized"
	file, err := os.Create(temporary)
	if err != nil {
		return err
	}
	if err := writeRGBPNG(file, output); err != nil {
		file.Close()
		os.Remove(temporary)
		return err
	}
	if err := file.Close(); err != nil {
		os.Remove(temporary)
		return err
	}
	return os.Rename(temporary, path)
}

func writeRGBPNG(file *os.File, source image.Image) error {
	var encoded bytes.Buffer
	encoded.Write([]byte{137, 80, 78, 71, 13, 10, 26, 10})
	ihdr := make([]byte, 13)
	binary.BigEndian.PutUint32(ihdr[0:4], uint32(source.Bounds().Dx()))
	binary.BigEndian.PutUint32(ihdr[4:8], uint32(source.Bounds().Dy()))
	ihdr[8] = 8
	ihdr[9] = 2 // truecolor RGB; no alpha channel
	writePNGChunk(&encoded, "IHDR", ihdr)

	var raw bytes.Buffer
	for y := source.Bounds().Min.Y; y < source.Bounds().Max.Y; y++ {
		raw.WriteByte(0)
		for x := source.Bounds().Min.X; x < source.Bounds().Max.X; x++ {
			red, green, blue, alpha := source.At(x, y).RGBA()
			blend := float64(alpha) / 65535
			raw.WriteByte(uint8(float64(red>>8)*blend + 248*(1-blend)))
			raw.WriteByte(uint8(float64(green>>8)*blend + 244*(1-blend)))
			raw.WriteByte(uint8(float64(blue>>8)*blend + 252*(1-blend)))
		}
	}
	var compressed bytes.Buffer
	writer := zlib.NewWriter(&compressed)
	if _, err := writer.Write(raw.Bytes()); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	writePNGChunk(&encoded, "IDAT", compressed.Bytes())
	writePNGChunk(&encoded, "IEND", nil)
	_, err := file.Write(encoded.Bytes())
	return err
}

func writePNGChunk(output *bytes.Buffer, kind string, data []byte) {
	binary.Write(output, binary.BigEndian, uint32(len(data)))
	output.WriteString(kind)
	output.Write(data)
	checksum := crc32.NewIEEE()
	checksum.Write([]byte(kind))
	checksum.Write(data)
	binary.Write(output, binary.BigEndian, checksum.Sum32())
}

func bilinear(source image.Image, x float64, y float64) color.RGBA {
	bounds := source.Bounds()
	x0 := clamp(int(x), bounds.Min.X, bounds.Max.X-1)
	y0 := clamp(int(y), bounds.Min.Y, bounds.Max.Y-1)
	x1 := clamp(x0+1, bounds.Min.X, bounds.Max.X-1)
	y1 := clamp(y0+1, bounds.Min.Y, bounds.Max.Y-1)
	xWeight := x - float64(x0)
	yWeight := y - float64(y0)
	return color.RGBA{
		R: interpolate(channel(source.At(x0, y0), 0), channel(source.At(x1, y0), 0), channel(source.At(x0, y1), 0), channel(source.At(x1, y1), 0), xWeight, yWeight),
		G: interpolate(channel(source.At(x0, y0), 1), channel(source.At(x1, y0), 1), channel(source.At(x0, y1), 1), channel(source.At(x1, y1), 1), xWeight, yWeight),
		B: interpolate(channel(source.At(x0, y0), 2), channel(source.At(x1, y0), 2), channel(source.At(x0, y1), 2), channel(source.At(x1, y1), 2), xWeight, yWeight),
		A: interpolate(channel(source.At(x0, y0), 3), channel(source.At(x1, y0), 3), channel(source.At(x0, y1), 3), channel(source.At(x1, y1), 3), xWeight, yWeight),
	}
}

func channel(value color.Color, index int) uint8 {
	r, g, b, a := value.RGBA()
	values := []uint32{r, g, b, a}
	return uint8(values[index] >> 8)
}

func interpolate(topLeft, topRight, bottomLeft, bottomRight uint8, x, y float64) uint8 {
	top := float64(topLeft) + (float64(topRight)-float64(topLeft))*x
	bottom := float64(bottomLeft) + (float64(bottomRight)-float64(bottomLeft))*x
	value := top + (bottom-top)*y
	if value < 0 {
		value = 0
	}
	if value > 255 {
		value = 255
	}
	return uint8(value + 0.5)
}

func clamp(value, minimum, maximum int) int {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}
