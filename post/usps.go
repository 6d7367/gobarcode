// implements http://en.wikipedia.org/wiki/POSTNET and
// http://en.wikipedia.org/wiki/Postal_Alpha_Numeric_Encoding_Technique
// 	msg := "555551237"
// 	f, _ := os.Create("postnet.png")
// 	b := post.NewPostnet(msg) // or post.NewPlanet
// 	img := b.GetImage()
// 	png.Encode(f, img)
// 	f.Close()
package post

import (
	"fmt"
	"image"
	"image/color"
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

var planetEncodeMap map[string]string = map[string]string{
	"0": "00111",
	"1": "11100",
	"2": "11010",
	"3": "11001",
	"4": "10110",
	"5": "10101",
	"6": "10011",
	"7": "01110",
	"8": "01101",
	"9": "01011",
}

type USPS struct {
	msg                 string
	BarHeight, BarWidth int
	DebugPrint          bool
	encodeMap           map[string]string
}

func newUSPS(msg string) *USPS {
	b := new(USPS)
	b.msg = msg

	b.BarWidth = 2
	b.BarHeight = 25
	b.DebugPrint = false

	return b
}

func NewPostnet(msg string) *USPS {
	b := newUSPS(msg)
	b.encodeMap = postnetEncodeMap

	return b
}

func NewPlanet(msg string) *USPS {
	b := newUSPS(msg)
	b.encodeMap = planetEncodeMap

	return b
}

func (this *USPS) GetImage() image.Image {
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

	return img
}

func (this *USPS) getEncodedForPrint() string {
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
		encoded += this.encodeMap[ch]
		encoded += interCharSymb
	}
	encoded += this.encodeMap[checkDigit]
	encoded += interCharSymb
	encoded += "1"
	encoded += interCharSymb

	return encoded
}

func (this *USPS) checksum() string {
	var sum int64 = 0
	var r int = 0

	for _, v := range this.msg {
		rV, _ := strconv.ParseInt(string(v), 10, 32)

		sum += rV
	}

	r = 10 - (int(sum) % 10)

	return fmt.Sprint(r)
}
