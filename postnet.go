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
	"0": "11000",
	"1": "00011",
	"2": "00101",
	"3": "00110",
	"4": "01001",
	"5": "01010",
	"6": "01100",
	"7": "10001",
	"8": "10010",
	"9": "10100",
}

type Postnet struct {
	msg                 string
	BarHeight, BarWidth int
	DebugPrint          bool
}

func NewPostnet(msg string) *Postnet {
	b := new(Postnet)
	b.msg = msg

	b.BarWidth = 2
	b.BarHeight = 25
	b.DebugPrint = false

	return b
}

func (this *Postnet) EncodeToPNG(w io.Writer) {
	encoded := this.getEncodedForPrint()

	pos := 0
	barH := this.BarHeight
	barW := this.BarWidth

	imgH := barH * 4
	imgW := len(encoded) * barW * 2

	size := image.Rect(0, 0, imgW, imgH)
	img := image.NewRGBA(size)

	for _, c := range encoded {

		switch string(c) {
		case "1":
			for x := 0; x <= barW; x++ {
				for y := barH * 2; y > 0; y-- {
					img.Set(x+pos, y, color.Black)
				}
			}
			pos += barW
		case "0":
			for x := 0; x <= barW; x++ {
				for y := barH; y <= barH*2; y++ {
					img.Set(x+pos, y, color.Black)
				}
			}
			pos += barW
		case "_":
			pos -= barW
			for x := 0; x <= barW; x++ {
				for y := 0; y <= barH; y++ {
					img.Set(x+pos, y, color.RGBA{255, 0, 0, 255})
				}
			}
		}
		pos += barW
	}

	png.Encode(w, img)
}

func (this *Postnet) getEncodedForPrint() string {
	var interCharSymb string

	if this.DebugPrint {
		interCharSymb = "_"
	} else {
		interCharSymb = ""
	}

	checkDigit := this.checksum()

	encoded := "1"
	encoded += interCharSymb
	for _, c := range this.msg {
		ch := string(c)
		encoded += encodeMap[ch]
		encoded += interCharSymb
	}
	encoded += encodeMap[checkDigit]
	encoded += interCharSymb
	encoded += "1"
	encoded += interCharSymb

	return encoded
}

func (this *Postnet) checksum() string {
	var sum int64 = 0
	var r int = 0

	for _, v := range this.msg {
		rV, _ := strconv.ParseInt(string(v), 10, 32)

		sum += rV
	}

	for (sum % 10) != 0 {
		sum += 1
		r += 1
	}

	return fmt.Sprint(r)
}

func main() {
	msg := "0123456789"
	f, _ := os.Create(msg + ".png")

	postnet := NewPostnet(msg)
	postnet.DebugPrint = true
	postnet.EncodeToPNG(f)

	f.Close()
}
