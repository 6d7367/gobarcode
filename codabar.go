// implements http://en.wikipedia.org/wiki/Codabar
package gobarcode

import "io"
import "image"
import "image/color"
import "image/png"

var encodeMap map[string]string = map[string]string{
	"0": "||| ▮",
	"2": "|| |▮",
	"6": "| ||▮",
	"1": "||▮ |",
	"-": "|| ▮|",
	"7": "| |▮|",
	"4": "|▮| |",
	"$": "|▮ ||",
	"8": "| ▮||",
	"5": "▮|| |",
	"9": "▮| ||",
	"3": "▮ |||",
	"C": "|| | ▮",
	"B": "| | |▮",
	"D": "|| ▮ |",
	"A": "|▮ | |",
	".": "▮▮▮|",
	"/": "▮▮|▮",
	":": "▮|▮▮",
	"+": "|▮▮▮",
}

type Codabar struct {
	msg, Start, Stop    string
	BarWidth, BarHeight int
	DebugPrint          bool
}

// Example
// 	msg := "1234567890"
// 	f, _ := os.Create(msg)
// 	codabar := gobarcode.NewCodabar(msg)
// 	codabar.BarHeight = 150
// 	codabar.BarWidth = 3
// 	codabar.EncodeToPNG(f)
// 	f.Close()
func NewCodabar(msg string) *Codabar {
	b := new(Codabar)
	b.msg = msg
	b.Start = "A"
	b.Stop = "B"
	b.BarWidth = 2
	b.BarHeight = 25
	b.DebugPrint = false

	return b
}

func (this *Codabar) EncodeToPNG(w io.Writer) {
	encoded := this.getEncodedForPrint()

	pos := 0
	barH := this.BarHeight
	barWNarrow := this.BarWidth
	barWWide := barWNarrow * 3

	imgH := barH
	imgW := computeImgW(encoded, barWWide, barWNarrow)

	size := image.Rect(0, 0, imgW, imgH)
	img := image.NewRGBA(size)

	for _, c := range encoded {
		var curW int = 0

		switch string(c) {
		case "|":
			curW = barWNarrow
		case "▮":
			curW = barWWide
		case " ":
			pos += barWWide
		case "_":
			pos -= barWNarrow
			for x := 0; x <= barWNarrow; x++ {
				for y := 0; y <= barH; y++ {
					img.Set(x+pos, y, color.RGBA{255, 0, 0, 255})
				}
			}
			pos += barWNarrow
		}

		if curW > 0 {
			for x := 0; x <= curW; x++ {
				for y := 0; y <= barH; y++ {
					img.Set(x+pos, y, color.Black)
				}
			}
			pos += curW
			pos += barWNarrow
		}
	}

	png.Encode(w, img)
}

func (this *Codabar) getEncodedForPrint() string {
	var interCharSymb string

	if this.DebugPrint {
		interCharSymb = "_"
	} else {
		interCharSymb = "+"
	}

	encoded := encodeMap[this.Start]
	encoded += interCharSymb
	for _, c := range this.msg {
		ch := string(c)
		encoded += encodeMap[ch]
		encoded += interCharSymb
	}
	encoded += encodeMap[this.Stop]
	encoded += interCharSymb

	return encoded
}
