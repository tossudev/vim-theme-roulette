package main


import (
    "fmt"
    "math"
)

func min3(a, b, c float64) float64 {
    return math.Min(math.Min(a, b), c)
}

func hueToRGB(h float64) (int, int, int) {
    kr := math.Mod(5+h*6, 6)
    kg := math.Mod(3+h*6, 6)
    kb := math.Mod(1+h*6, 6)

    rf := 1 - math.Max(min3(kr, 4-kr, 1), 0)
    gf := 1 - math.Max(min3(kg, 4-kg, 1), 0)
    bf := 1 - math.Max(min3(kb, 4-kb, 1), 0)

	r := int(rf*255.0)
	g := int(gf*255.0)
	b := int(bf*255.0)

    return r, g, b
}


func Rainbow(hue float64) string {
	r, g, b := hueToRGB(hue)

	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)

}
