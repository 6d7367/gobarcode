package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"strconv"
)

var encodeMap map[string]string = map[string]string{
	"0":     "100100100100",
	"1":     "110100100100",
	"2":     "100110100100",
	"3":     "110110100100",
	"4":     "100100110100",
	"5":     "110100110100",
	"6":     "100110110100",
	"7":     "110110110100",
	"8":     "100100100110",
	"9":     "110100100110",
	"A":     "100110100110",
	"B":     "110110100110",
	"C":     "100100110110",
	"D":     "110100110110",
	"E":     "100110110110",
	"F":     "110110110110",
	"TERM":  "1111",
	"START": "110110100110",
	"STOP":  "100100110110",
}

type Plessey struct {
	msg                 string
	BarHeight, BarWidth int
	DebugPrint          bool
}

func NewPlessey(msg string) *Plessey {
	b := new(Plessey)
	b.msg = msg
	b.BarWidth = 3
	b.BarHeight = 50
	b.DebugPrint = false

	return b
}

func (this *Plessey) EncodeToPNG(w io.Writer) {
	encoded := this.getEncodedForPrint()

	barH := this.BarHeight

	barW := this.BarWidth
	imgH := barH
	imgW := len(encoded) * barW

	size := image.Rect(0, 0, imgW, imgH)
	img := image.NewRGBA(size)

	pos := 0

	for _, c := range encoded {
		switch string(c) {
		case "1":
			for x := 0; x <= barW; x++ {
				for y := 0; y <= barH; y++ {
					img.Set(x+pos, y, color.Black)
				}
			}
			pos += barW
		case "_":
			for x := 0; x <= barW; x++ {
				for y := 0; y <= barH; y++ {
					img.Set(x+pos, y, color.RGBA{255, 0, 0, 255})
				}
			}
			pos += barW
		case "0":
			pos += barW
		}

	}

	png.Encode(w, img)
}

func (this *Plessey) getEncodedForPrint() string {
	var interCharSymb string
	if this.DebugPrint {
		interCharSymb = "_"
	} else {
		interCharSymb = "0"
	}

	encoded := encodeMap["START"]
	encoded += interCharSymb

	for _, c := range this.msg {
		ch := string(c)
		encoded += encodeMap[ch]
		encoded += interCharSymb
	}

	encoded += encodeMap["TERM"]
	encoded += interCharSymb
	encoded += encodeMap["STOP"]
	encoded += interCharSymb

	return encoded
}

func main() {
	msg := "0F"
	f, _ := os.Create("plessey.png")
	plessey := NewPlessey(msg)
	plessey.DebugPrint = true
	plessey.EncodeToPNG(f)
	f.Close()

	x, _ := strconv.ParseInt("10000001000000010000", 2, 32)
	fmt.Printf("%x\n", x)

}
