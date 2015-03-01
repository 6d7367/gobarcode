// implements http://en.wikipedia.org/wiki/Code_39
// 	msg := "GOBARCODE"
// 	f, _ := os.Create("code39.png")
// 	b := linear.NewCode39(msg)
// 	img := b.GetImage()
// 	png.Encode(f, img)
// 	f.Close()
package linear

import "image"
import "image/color"

type Code39 struct {
	msg                 string
	BarHeight, BarWidth int
	DebugPrint          bool
}

var code39encodeMap map[string]string = map[string]string{
	"1": "▮| ||▮",
	"2": "|▮ ||▮",
	"3": "▮▮ |||",
	"4": "|| ▮|▮",
	"5": "▮| ▮||",
	"6": "|▮ ▮||",
	"7": "|| |▮▮",
	"8": "▮| |▮|",
	"9": "|▮ |▮|",
	"0": "|▮ |▮|",
	"A": "▮|| |▮",
	"B": "|▮| |▮",
	"C": "▮▮| ||",
	"D": "||▮ |▮",
	"E": "▮|▮ ||",
	"F": "|▮▮ ||",
	"G": "||| ▮▮",
	"H": "▮|| ▮|",
	"I": "|▮| ▮|",
	"J": "||▮ ▮|",
	"K": "▮||| ▮",
	"L": "|▮|| ▮",
	"M": "▮▮|| |",
	"N": "||▮| ▮",
	"O": "▮|▮| |",
	"P": "|▮▮| |",
	"Q": "|||▮ ▮",
	"R": "▮||▮ |",
	"S": "|▮|▮ |",
	"T": "||▮▮ |",
	"U": "▮ |||▮",
	"V": "| ▮||▮",
	"W": "▮ ▮|||",
	"X": "| |▮|▮",
	"Y": "▮ |▮||",
	"Z": "| ▮▮||",
	"-": "| ||▮▮",
	".": "▮ ||▮|",
	" ": "| ▮|▮|",
	"*": "| |▮▮|",
}

func NewCode39(msg string) *Code39 {
	b := new(Code39)
	b.msg = msg
	b.BarWidth = 2
	b.BarHeight = 25
	b.DebugPrint = false

	return b
}

func (this *Code39) GetImage() image.Image {
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

	return img
}

func (this *Code39) getEncodedForPrint() string {
	var interCharSymb string

	if this.DebugPrint {
		interCharSymb = "_"
	} else {
		interCharSymb = ""
	}

	encoded := code39encodeMap["*"]
	encoded += interCharSymb
	for _, c := range this.msg {
		ch := string(c)
		encoded += code39encodeMap[ch]
		encoded += interCharSymb
	}
	encoded += code39encodeMap["*"]
	encoded += interCharSymb

	return encoded
}
