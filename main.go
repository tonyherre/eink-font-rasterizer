package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func LoadTtf() *truetype.Font {
	f, err := os.Open("Roboto-Black.ttf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var ttf []byte
	for {
		buffer := make([]byte, 1024)
		_, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		ttf = append(ttf, buffer...)
	}

	font, err := truetype.Parse(ttf)
	if err != nil {
		panic(err)
	}
	return font
}

func RenderToBytes(s string, face font.Face, height int) (bytes []byte, byteWidth int, advance int) {
	point := fixed.Point26_6{X: fixed.I(0), Y: fixed.I(30)}
	d := &font.Drawer{
		Dst:  nil,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  point,
	}

	box, advanceDec := d.BoundString(s)

	img := image.NewRGBA(image.Rect(0, 0, box.Max.X.Ceil(), height))
	d.Dst = img
	d.DrawString(s)
	byteWidth = (box.Max.X.Ceil() + 7) / 8

	for y := 0; y < height; y++ {
		for x := 0; x < byteWidth; x++ {
			bt := byte(0)
			for bit := 0; bit < 8; bit++ {
				colour := img.RGBAAt(x*8+bit, y)
				if colour.A > 200 {
					bt |= 1 << (7 - bit)
					colour.A = 255
				} else {
					colour.A = 0
				}
				img.SetRGBA(x*8+bit, y, colour)
			}
			bytes = append(bytes, bt)
		}
	}

	f, err := os.Create(fmt.Sprintf("%s.png", s))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}

	return bytes, byteWidth, advanceDec.Ceil()
}

func PrintByteArr(name string, bytes []byte) {
	fmt.Printf("const unsigned char %s[%d] PROGMEM = {\n", name, len(bytes))
	for i, b := range bytes {
		fmt.Printf("0x%02x", b)
		if i != len(bytes)-1 {
			fmt.Print(",")
		}
		if i%16 == 15 {
			fmt.Println()
		}
	}
	fmt.Println("};")

}

func main() {
	fmt.Println("Loading TTF")

	myFont := LoadTtf()
	face := truetype.NewFace(myFont, &truetype.Options{Size: 40})
	fmt.Println("Got font")

	height := 31
	fmt.Printf("const int DIGIT_HEIGHT = %d;\n", height)
	fmt.Printf("struct Element {int byte_width; int advance; const unsigned char* data;};\n\n")

	for i := 0; i < 10; i++ {
		bytes, width, advance := RenderToBytes(fmt.Sprintf("%d", i), face, height)

		//writeToArduino(bytes)
		fmt.Printf("const int DIGIT_%d_BYTE_WIDTH = %d;\n", i, width)
		fmt.Printf("const int DIGIT_%d_ADVANCE = %d;\n", i, advance)
		PrintByteArr(fmt.Sprintf("DIGIT_%d_DATA", i), bytes)
		fmt.Printf("Element DIGIT_%d{DIGIT_%d_BYTE_WIDTH, DIGIT_%d_ADVANCE, DIGIT_%d_DATA};\n\n", i, i, i, i)
	}

	bytes, width, advance := RenderToBytes("mins", face, height)
	fmt.Printf("const int MINS_BYTE_WIDTH = %d;\n", width)
	fmt.Printf("const int MINS_ADVANCE = %d;\n", advance)
	PrintByteArr("MINS_DATA", bytes)
	fmt.Printf("Element MINS{MINS_BYTE_WIDTH, MINS_ADVANCE, MINS_DATA};\n\n")

	bytes, width, advance = RenderToBytes(" - ", face, height)
	fmt.Printf("const int SEP_BYTE_WIDTH = %d;\n", width)
	fmt.Printf("const int SEP_ADVANCE = %d;\n", advance)
	PrintByteArr("SEP_DATA", bytes)
	fmt.Printf("Element SEP{SEP_BYTE_WIDTH, SEP_ADVANCE, SEP_DATA};\n\n")

	bytes, width, advance = RenderToBytes(":", face, height)
	fmt.Printf("const int COLON_BYTE_WIDTH = %d;\n", width)
	fmt.Printf("const int COLON_ADVANCE = %d;\n", advance)
	PrintByteArr("COLON_DATA", bytes)
	fmt.Printf("Element COLON{COLON_BYTE_WIDTH, COLON_ADVANCE, COLON_DATA};\n\n")

	bytes, width, advance = RenderToBytes("Status:", face, height)
	fmt.Printf("const int STATUS_BYTE_WIDTH = %d;\n", width)
	fmt.Printf("const int STATUS_ADVANCE = %d;\n", advance)
	PrintByteArr("STATUS_DATA", bytes)
	fmt.Printf("Element STATUS{STATUS_BYTE_WIDTH, STATUS_ADVANCE, STATUS_DATA};\n\n")

	bytes, width, advance = RenderToBytes("Batt ", face, height)
	fmt.Printf("const int BATT_BYTE_WIDTH = %d;\n", width)
	fmt.Printf("const int BATT_ADVANCE = %d;\n", advance)
	PrintByteArr("BATT_DATA", bytes)
	fmt.Printf("Element BATT{BATT_BYTE_WIDTH, BATT_ADVANCE, BATT_DATA};\n\n")

	bytes, width, advance = RenderToBytes("%", face, height)
	fmt.Printf("const int PERCENT_BYTE_WIDTH = %d;\n", width)
	fmt.Printf("const int PERCENT_ADVANCE = %d;\n", advance)
	PrintByteArr("PERCENT_DATA", bytes)
	fmt.Printf("Element PERCENT{PERCENT_BYTE_WIDTH, PERCENT_ADVANCE, PERCENT_DATA};\n\n")

	// Helper arrays
	fmt.Print(
		`const Element DIGITS[] = {
  DIGIT_0, DIGIT_1, DIGIT_2, DIGIT_3, DIGIT_4, DIGIT_5, DIGIT_6, DIGIT_7, DIGIT_8, DIGIT_9
};`)
}
