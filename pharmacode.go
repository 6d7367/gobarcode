// implements http://en.wikipedia.org/wiki/Pharmacode
package gobarcode

import "fmt"
import "io"
import "image"
import "image/color"
import "image/png"

type Pharmacode struct {
	n, BarWidth, BarHeight int
}

func NewPharmacode(n int) *Pharmacode {
	b := new(Pharmacode)
	b.n = n
	b.BarWidth = 2
	b.BarHeight = 50

	return b
}

// Example
// 	import "github.com/6d7367/gobarcode"
// 	...
// 	n := 69
// 	filename := fmt.Sprint(n) + ".png"
// 	f, _ := os.Create(filename)
// 	p := gobarcode.NewPharmacode(n)
// 	p.EncodeToPNG(f)
// 	f.Close()
func (this *Pharmacode) EncodeToPNG(w io.Writer) {
	encoded := fmt.Sprintf("%b", this.n+1)[1:]

	barH := this.BarHeight

	barWNarrow := this.BarWidth
	barWWide := 3 * barWNarrow
	imgH := barH
	imgW := len(encoded) * (barWWide * 2)

	size := image.Rect(0, 0, imgW, imgH)
	img := image.NewRGBA(size)

	pos := 0

	for _, c := range encoded {
		var offset int
		switch string(c) {
		case "1":
			offset = barWWide
		case "0":
			offset = barWNarrow
		}
		for x := 0; x <= offset; x++ {
			for y := 0; y <= barH; y++ {
				img.Set(x+pos, y, color.Black)
			}
		}

		pos += barWWide + offset

	}

	png.Encode(w, img)
}
