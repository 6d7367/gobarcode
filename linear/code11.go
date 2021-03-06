// implements http://en.wikipedia.org/wiki/Code_11
// 	msg := "1917"
// 	f, _ := os.Create("code11.png")
// 	b := linear.NewCode11(msg)
// 	img := b.GetImage()
// 	png.Encode(f, img)
// 	f.Close()
package linear

import "strconv"
import "image"
import "image/color"

type Code11 struct {
	raw, withChecksum   string
	BarWidth, BarHeight int
	DebugPrint          bool
}

var code11encodeMap map[string]string = map[string]string{"0": "101011",
	"1":  "1101011",
	"2":  "1001011",
	"3":  "1100101",
	"4":  "1011011",
	"5":  "1101101",
	"6":  "1001101",
	"7":  "1010011",
	"8":  "1101001",
	"9":  "110101",
	"-":  "101101",
	"ss": "1011001",
}

func NewCode11(msg string) *Code11 {
	b := new(Code11)
	b.raw = msg
	b.BarWidth = 2
	b.BarHeight = 25
	b.DebugPrint = false

	return b
}

func (this *Code11) GetImage() image.Image {
	this.checksum()

	encoded := this.getEncodedForPrint()

	pos := 0
	barH := this.BarHeight
	barW := this.BarWidth
	imgH := barH
	imgW := barW * len(encoded)

	size := image.Rect(0, 0, imgW, imgH)
	img := image.NewRGBA(size)

	for _, c := range encoded {
		switch string(c) {
		case "1":
			for x := 0; x <= barW; x++ {
				for y := 0; y <= barH; y++ {
					img.Set(x+pos, y, color.Black)
				}
			}
			pos += barW
		case "-":
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

	return img

}

func (this *Code11) getEncodedForPrint() string {
	var interCharSymb string

	if this.DebugPrint {
		interCharSymb = "-"
	} else {
		interCharSymb = "0"
	}

	encoded := code11encodeMap["ss"]
	encoded += interCharSymb
	for _, ch := range this.withChecksum {
		encoded += code11encodeMap[string(ch)]
		encoded += interCharSymb
	}
	encoded += code11encodeMap["ss"]
	encoded += interCharSymb

	return encoded
}

func (this *Code11) checksum() {
	msg := this.raw

	weighted := this.stringWeight(msg)
	checksumC := strconv.Itoa(weighted % 11)

	msg += checksumC

	if len(this.raw) >= 10 {
		weighted = this.stringWeight(msg)
		checksumK := strconv.Itoa(weighted % 9)

		msg += checksumK
	}

	this.withChecksum = msg
}

func (this *Code11) stringWeight(s string) int {
	weight := 1
	var weighted int = 0
	for i := len(s) - 1; i >= 0; i-- {
		var mult int
		ch := string(s[i])
		if ch == "-" {
			mult = 10
		} else {
			multTmp, err := strconv.ParseInt(ch, 10, 32)
			if err != nil {
				panic("shit happens")
			}
			mult = int(multTmp)
		}
		weighted += weight * mult
		weight++
	}

	return weighted
}
