package components

type BorderSize struct {
	Width  int // Добавка к ширине (лево + право)
	Height int // Добавка к высоте (верх + низ)
}

func GetBorderSize(borderType string) BorderSize {
	switch borderType {
	case "normal", "rounded", "thick", "double":
		return BorderSize{Width: 2, Height: 2}
	case "hidden", "":
		return BorderSize{Width: 0, Height: 0}
	default:
		return BorderSize{Width: 0, Height: 0}
	}
}
