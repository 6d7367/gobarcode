package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
)

// t -tracker
// d - descentin
// a - ascenting
// f - full
var encodeMap map[string]string = map[string]string{
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

type RM4SCC struct {
	msg                 string
	BarHeight, BarWidth int
	DebugPrint          bool
}

func NewRM4SCC(msg string) *RM4SCC {
	b := new(RM4SCC)

	b.msg = msg
	b.BarHeight = 25
	b.BarWidth = 2
	b.DebugPrint = false

	return b
}

func (this *RM4SCC) EncodeToPNG(w io.Writer) {
	encoded := this.getEncodedForPrint()

	pos := 0
	barH := this.BarHeight
	barW := this.BarWidth
	barWWide := barW * 2
	barHF := barH * 3

	imgH := barH * 3
	imgW := len(encoded) * barW * 3
	// imgW := 500

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

	png.Encode(w, img)
}

func (this *RM4SCC) getEncodedForPrint() string {
	var interCharSymb string

	if this.DebugPrint {
		interCharSymb = "_"
	} else {
		interCharSymb = ""
	}

	// checkDigit := this.checksum()

	encoded := encodeMap["START"]
	encoded += interCharSymb
	for _, c := range this.msg {
		ch := string(c)
		encoded += encodeMap[ch]
		encoded += interCharSymb
	}
	// encoded += encodeMap[checkDigit]
	// encoded += interCharSymb
	encoded += encodeMap["STOP"]
	encoded += interCharSymb

	return encoded
}

func (this *RM4SCC) checksum() string {
	var r int = 0

	return fmt.Sprint(r)
}

func main() {
	msg := "RM4SCC"
	f, _ := os.Create(msg + ".png")
	rm4scc := NewRM4SCC(msg)
	rm4scc.DebugPrint = true
	rm4scc.EncodeToPNG(f)
	f.Close()
}
