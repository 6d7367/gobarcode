package gobarcode

func computeImgW(s string, barWideLen, barNarrowLen int) int {
	r := 0
	for _, c := range s {
		ch := string(c)
		switch ch {
		case "|":
			r += barNarrowLen
		case "▮":
			r += barWideLen
		case " ":
			r += barWideLen
		}
		r += barNarrowLen
	}

	return r
}
