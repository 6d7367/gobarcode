// implements http://en.wikipedia.org/wiki/Code_93
// 	msg := "TEST93"
// 	f, _ := os.Create("code93.png")
// 	b := linear.NewCode93(msg)
// 	img := b.GetImage()
// 	png.Encode(f, img)
// 	f.Close()
package linear

import "image"
import "image/color"

type encodeChar struct {
	value  int
	encode string
}

var code93EncodeMap map[string]encodeChar = map[string]encodeChar{
	"0":   encodeChar{0, "100010100"},
	"1":   encodeChar{1, "101001000"},
	"2":   encodeChar{2, "101000100"},
	"3":   encodeChar{3, "101000010"},
	"4":   encodeChar{4, "100101000"},
	"5":   encodeChar{5, "100100100"},
	"6":   encodeChar{6, "100100010"},
	"7":   encodeChar{7, "101010000"},
	"8":   encodeChar{8, "100010010"},
	"9":   encodeChar{9, "100001010"},
	"A":   encodeChar{10, "110101000"},
	"B":   encodeChar{11, "110100100"},
	"C":   encodeChar{12, "110100010"},
	"D":   encodeChar{13, "110010100"},
	"E":   encodeChar{14, "110010010"},
	"F":   encodeChar{15, "110001010"},
	"G":   encodeChar{16, "101101000"},
	"H":   encodeChar{17, "101100100"},
	"I":   encodeChar{18, "101100010"},
	"J":   encodeChar{19, "100110100"},
	"K":   encodeChar{20, "100011010"},
	"L":   encodeChar{21, "101011000"},
	"M":   encodeChar{22, "101001100"},
	"N":   encodeChar{23, "101000110"},
	"O":   encodeChar{24, "100101100"},
	"P":   encodeChar{25, "100010110"},
	"Q":   encodeChar{26, "110110100"},
	"R":   encodeChar{27, "110110010"},
	"S":   encodeChar{28, "110101100"},
	"T":   encodeChar{29, "110100110"},
	"U":   encodeChar{30, "110010110"},
	"V":   encodeChar{31, "110011010"},
	"W":   encodeChar{32, "101101100"},
	"X":   encodeChar{33, "101100110"},
	"Y":   encodeChar{34, "100110110"},
	"Z":   encodeChar{35, "100111010"},
	"-":   encodeChar{36, "100101110"},
	".":   encodeChar{37, "111010100"},
	" ":   encodeChar{38, "111010010"},
	"$":   encodeChar{39, "111001010"},
	"/":   encodeChar{40, "101101110"},
	"+":   encodeChar{41, "101110110"},
	"%":   encodeChar{42, "110101110"},
	"($)": encodeChar{43, "100100110"}, // unused
	"(%)": encodeChar{44, "111011010"}, // unused
	"(/)": encodeChar{45, "111010110"}, // unused
	"(+)": encodeChar{46, "100110010"}, // unused
	"SS":  encodeChar{0, "101011110"},
}

type Code93 struct {
	msg                 string
	BarHeight, BarWidth int
	DebugPrint          bool
}

func NewCode93(msg string) *Code93 {
	b := new(Code93)

	b.msg = msg
	b.BarHeight = 150
	b.BarWidth = 3
	b.DebugPrint = false

	return b
}

func (this *Code93) GetImage() image.Image {
	encoded := this.getEncodedForPrint()

	pos := 0
	barH := this.BarHeight
	barW := this.BarWidth

	imgH := barH
	imgW := len(encoded) * barW

	size := image.Rect(0, 0, imgW, imgH)
	img := image.NewRGBA(size)

	for _, c := range encoded {
		ch := string(c)
		switch ch {
		case "1":
			for x := 0; x <= barW; x++ {
				for y := 0; y <= barH; y++ {
					img.Set(x+pos, y, color.Black)
				}
			}
			pos += barW
		case "0":
			pos += barW
		case "_":
			for x := 0; x <= barW; x++ {
				for y := 0; y <= barH; y++ {
					img.Set(x+pos, y, color.RGBA{255, 0, 0, 255})
				}
			}
			pos += barW
		}
	}

	return img
}

func (this *Code93) getEncodedForPrint() string {
	var interCharSymb string
	if this.DebugPrint {
		interCharSymb = "_"
	} else {
		interCharSymb = "0"
	}
	checkSum := this.checksum()

	encoded := code93EncodeMap["SS"].encode
	encoded += interCharSymb

	for _, c := range this.msg + checkSum {
		ch := string(c)
		encoded += code93EncodeMap[ch].encode
		encoded += interCharSymb
	}

	encoded += code93EncodeMap["SS"].encode
	encoded += interCharSymb
	encoded += "1" + interCharSymb

	return encoded
}

func (this *Code93) checksum() string {
	var sumCCh, sumKCh string
	var sumC, sumK int
	weight := 1

	for i := len(this.msg) - 1; i >= 0; i-- {
		enCh := code93EncodeMap[string(this.msg[i])]
		sumC += enCh.value * weight
		weight += 1

		if weight >= 20 {
			weight = 0
		}
	}

	sumC = sumC % 47
	weight = 1

	findChar := func(sum int) string {
		r := ""
		for i, enCh := range code93EncodeMap {
			if enCh.value == sum {
				r = i
				break
			}
		}

		return r
	}

	sumCCh = findChar(sumC)
	afterSumC := this.msg + sumCCh

	for i := len(afterSumC) - 1; i >= 0; i-- {
		enCh := code93EncodeMap[string(afterSumC[i])]
		sumK += enCh.value * weight
		weight += 1

		if weight >= 15 {
			weight = 0
		}
	}

	sumK = sumK % 47
	sumKCh = findChar(sumK)

	return (sumCCh + sumKCh)
}

func main() {

}
