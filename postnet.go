package gobarcode

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"strconv"
)

var postnetEncodeMap map[string]string = map[string]string{
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

	b.BarWidth = 4
	b.BarHeight = 50
	b.DebugPrint = false

	return b
}

// Example
// 	msg := "555551237"
// 	f, _ := os.Create(msg + ".png")
// 	postnet := NewPostnet(msg)
// 	postnet.DebugPrint = true
// 	postnet.EncodeToPNG(f)
// 	f.Close()
func (this *Postnet) EncodeToPNG(w io.Writer) {
	encoded := this.getEncodedForPrint()

	pos := 0
	barH := this.BarHeight
	barW := this.BarWidth
	barWWide := barW * 2

	imgH := barH * 2
	imgW := len(encoded) * barWWide * 2

	size := image.Rect(0, 0, imgW, imgH)
	img := image.NewRGBA(size)

	for _, c := range encoded {

		switch string(c) {
		case "1":
			for x := 0; x <= barW; x++ {
				for y := barH * 2; y >= 0; y-- {
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
			pos -= barWWide
			for x := 0; x <= barW; x++ {
				for y := barH * 2; y >= 0; y-- {
					img.Set(x+pos, y, color.RGBA{255, 0, 0, 255})
				}
			}
		}
		pos += barWWide
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
		encoded += postnetEncodeMap[ch]
		encoded += interCharSymb
	}
	encoded += postnetEncodeMap[checkDigit]
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

	r = 10 - (int(sum) % 10)

	return fmt.Sprint(r)
}
