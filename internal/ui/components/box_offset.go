package components

type Offset struct {
	X int
	Y int
}

func NewOffset(x int, y int) *Offset {
	return &Offset{
		X: x,
		Y: y,
	}
}

func CalculateOffset(
	cWidth int,
	cBorder string,
	cLayout string,
	cAlign string,
	cPaddingH int,
	cPaddingV int,
	cMarginH int,
	cMarginV int,
	cChildren []interface{}) (x, y int) {
	offset := NewOffset(0, 0)
	offset.X += CalculateAlignOffset(cWidth, cBorder, cPaddingH, cLayout, cAlign, cChildren)
	borderSize := GetBorderSize(cBorder)
	if borderSize.Width != 0 {
		offset.X += borderSize.Width
		offset.Y += borderSize.Height
	}
	if cPaddingH > 0 || cPaddingV > 0 {
		offset.X += cPaddingH
		offset.Y += cPaddingV
	}
	if cMarginH > 0 || cMarginV > 0 {
		offset.X += cMarginH
		offset.Y += cMarginV
	}
	return offset.X, offset.Y
}
