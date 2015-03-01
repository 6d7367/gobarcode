// implements http://en.wikipedia.org/wiki/Pharmacode
// 	n := 69
// 	f, _ := os.Create("pharmacode.png")
// 	b := linear.NewPharmacode(n)
// 	img := b.GetImage()
// 	png.Encode(f, img)
// 	f.Close()
package linear

import "fmt"
import "image"
import "image/color"

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

func (this *Pharmacode) GetImage() image.Image {
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

	return img
}
