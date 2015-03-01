// implements http://en.wikipedia.org/wiki/RM4SCC
// 	msg := "BX11LT1A"
// 	f, _ := os.Create("rm4scc.png")
// 	b := post.NewRM4SCC(msg)
// 	img := b.GetImage()
// 	png.Encode(f, img)
// 	f.Close()
package post

import (
	"image"
	"image/color"
)

// t -tracker
// d - descentin
// a - ascenting
// f - full
var rm4sccEncodeMap map[string]string = map[string]string{
	"A":     "dada",
	"B":     "dfta",
	"C":     "tafd",
	"D":     "tfad",
	"E":     "tfft",
	"F":     "daad",
	"G":     "daft",
	"H":     "dfat",
	"I":     "atdf",
	"J":     "adtf",
	"K":     "adda",
	"L":     "fttf",
	"M":     "ftda",
	"N":     "fdta",
	"O":     "atfd",
	"P":     "adad",
	"Q":     "adft",
	"R":     "ftad",
	"S":     "ftft",
	"T":     "fdat",
	"U":     "aadd",
	"V":     "aftd",
	"W":     "afdt",
	"X":     "fatd",
	"Y":     "fadt",
	"Z":     "fftt",
	"0":     "ttff",
	"1":     "tdaf",
	"2":     "tdfa",
	"3":     "dtaf",
	"4":     "dtfa",
	"5":     "ddaa",
	"6":     "tadf",
	"7":     "tftf",
	"8":     "tfda",
	"9":     "datf",
	"START": "a",
	"STOP":  "f",
}

type rm4sccChecksumRow [6]string
type rm4sccChecsumMatrix [6]rm4sccChecksumRow

type RM4SCC struct {
	msg                 string
	BarHeight, BarWidth int
	DebugPrint          bool
}

func NewRM4SCC(msg string) *RM4SCC {
	b := new(RM4SCC)

	b.msg = msg
	b.BarHeight = 10
	b.BarWidth = 2
	b.DebugPrint = false

	return b
}

func (this *RM4SCC) GetImage() image.Image {
	encoded := this.getEncodedForPrint()

	pos := 0
	barH := this.BarHeight
	barW := this.BarWidth
	barWWide := barW * 2
	barHF := barH * 3

	imgH := barH * 3
	imgW := len(encoded) * barW * 3

	size := image.Rect(0, 0, imgW, imgH)
	img := image.NewRGBA(size)

	for _, c := range encoded {
		ch := string(c)
		switch ch {
		case "t":
			for x := 0; x <= barW; x++ {
				for y := barH; y <= barH*2; y++ {
					img.Set(x+pos, y, color.Black)
				}
			}
			pos += barW
		case "f":
			for x := 0; x <= barW; x++ {
				for y := 0; y <= barH*3; y++ {
					img.Set(x+pos, y, color.Black)
				}
			}
			pos += barW
		case "a":
			for x := 0; x <= barW; x++ {
				for y := 0; y <= barH*2; y++ {
					img.Set(x+pos, y, color.Black)
				}
			}
			pos += barW
		case "d":
			for x := 0; x <= barW; x++ {
				for y := barH; y <= barH*3; y++ {
					img.Set(x+pos, y, color.Black)
				}
			}
			pos += barW
		case "_":
			pos -= barW
			for x := 0; x <= barW; x++ {
				for y := barHF; y >= 0; y-- {
					img.Set(x+pos, y, color.RGBA{255, 0, 0, 255})
				}
			}
		default:
			pos += barW
		}

		pos += barWWide
	}

	return img
}

func (this *RM4SCC) getEncodedForPrint() string {
	var interCharSymb string

	if this.DebugPrint {
		interCharSymb = "_"
	} else {
		interCharSymb = ""
	}

	encoded := rm4sccEncodeMap["START"]
	encoded += interCharSymb
	for _, c := range this.msg {
		ch := string(c)
		encoded += rm4sccEncodeMap[ch]
		encoded += interCharSymb
	}

	encoded += rm4sccEncodeMap[this.checksum()]
	encoded += interCharSymb
	encoded += rm4sccEncodeMap["STOP"]
	encoded += interCharSymb

	return encoded
}

// чтобы получить символ для проверки на целостность необходимо вычислить два числа
// суммы верхних и нижних штрихов. У позиции штриха в кодируемом символе есть свой множитель,
// см. weighter. Полученные суммы верхних и нижних штрихов необходимо разделить на 6,
// полученный остаток является индексом линии и столбца в матрице, см. checksumDigit.
// Символ, находящийся по индексу и будет проверочным
func (this *RM4SCC) checksum() string {
	weighter := map[int]int{
		0: 0,
		1: 1,
		2: 2,
		3: 4,
	}

	checksumDigit := rm4sccChecsumMatrix{
		rm4sccChecksumRow{"0", "1", "2", "3", "4", "5"},
		rm4sccChecksumRow{"6", "7", "8", "9", "A", "B"},
		rm4sccChecksumRow{"C", "D", "E", "F", "G", "H"},
		rm4sccChecksumRow{"I", "J", "K", "L", "M", "N"},
		rm4sccChecksumRow{"O", "P", "Q", "R", "S", "T"},
		rm4sccChecksumRow{"U", "V", "W", "X", "Y", "Z"},
	}

	weightA := 0
	weightD := 0

	for _, c := range this.msg {
		ch := rm4sccEncodeMap[string(c)]
		tmpWeightA := 0
		tmpWeightD := 0

		for i, b := range ch {
			weight := 3 - i
			switch string(b) {
			case "a":
				tmpWeightA += weighter[weight] * 1
			case "d":
				tmpWeightD += weighter[weight] * 1
			case "f":
				tmpWeightA += weighter[weight] * 1
				tmpWeightD += weighter[weight] * 1
			}
		}

		weightA += tmpWeightA
		weightD += tmpWeightD
	}

	weightA = weightA % 6
	weightD = weightD % 6

	r := checksumDigit[weightA-1][weightD-1]

	return r
}
