package components

// FlexDirection определяет направление flex-контейнера
type FlexDirection string

const (
	FlexDirectionRow    FlexDirection = "row"
	FlexDirectionColumn FlexDirection = "column"
)

// FlexJustify определяет выравнивание по главной оси
type FlexJustify string

const (
	FlexJustifyStart        FlexJustify = "start"
	FlexJustifyCenter       FlexJustify = "center"
	FlexJustifyEnd          FlexJustify = "end"
	FlexJustifySpaceBetween FlexJustify = "space-between"
)

// FlexAlign определяет выравнивание по поперечной оси
type FlexAlign string

const (
	FlexAlignStart   FlexAlign = "start"
	FlexAlignCenter  FlexAlign = "center"
	FlexAlignEnd     FlexAlign = "end"
	FlexAlignStretch FlexAlign = "stretch"
)

// BorderStyle определяет стиль границы
type BorderStyle string

const (
	BorderStyleNone    BorderStyle = "none"
	BorderStyleSolid   BorderStyle = "solid"
	BorderStyleDouble  BorderStyle = "double"
	BorderStyleRounded BorderStyle = "rounded"
)

// TextStyle определяет стиль текста
type TextStyle struct {
	Bold      bool
	Dim       bool
	Italic    bool
	Underline bool
	Blink     bool
	Reverse   bool
	Hidden    bool
	Strike    bool
}

// Flex определяет flex-параметры
type Flex struct {
	Direction FlexDirection
	Justify   FlexJustify
	Align     FlexAlign
	Wrap      bool
	Gap       int
}

// Border определяет параметры границы
type Border struct {
	Style BorderStyle
}

// UIComponentStyle содержит все стили компонента
type UIComponentStyle struct {
	Color      string
	Background string
	Flex       Flex
	Border     Border
	TextStyle  TextStyle
}

type UIComponentStyleI interface {
	SetColor(val string) UIComponentStyleI
	SetBackground(val string) UIComponentStyleI
	SetFlex(val Flex) UIComponentStyleI
	SetBorder(val Border) UIComponentStyleI
	SetTextStyle(val TextStyle) UIComponentStyleI
	SetBold(val bool) UIComponentStyleI
	SetItalic(val bool) UIComponentStyleI
	SetUnderline(val bool) UIComponentStyleI
	SetDim(val bool) UIComponentStyleI
	SetBlink(val bool) UIComponentStyleI
	SetReverse(val bool) UIComponentStyleI
	SetHidden(val bool) UIComponentStyleI
	SetStrike(val bool) UIComponentStyleI

	Build() UIComponentStyle
}

func NewUIComponentStyleBuilder() UIComponentStyleI {
	return &UIComponentStyle{}
}

// SetColor устанавливает цвет текста
func (s *UIComponentStyle) SetColor(val string) UIComponentStyleI {
	s.Color = val
	return s
}

// SetBackground устанавливает цвет фона
func (s *UIComponentStyle) SetBackground(val string) UIComponentStyleI {
	s.Background = val
	return s
}

// SetFlex устанавливает flex-параметры
func (s *UIComponentStyle) SetFlex(val Flex) UIComponentStyleI {
	s.Flex = val
	return s
}

// SetBorder устанавливает параметры границы
func (s *UIComponentStyle) SetBorder(val Border) UIComponentStyleI {
	s.Border = val
	return s
}

// SetTextStyle устанавливает стиль текста
func (s *UIComponentStyle) SetTextStyle(val TextStyle) UIComponentStyleI {
	s.TextStyle = val
	return s
}

// SetBold устанавливает жирность текста
func (s *UIComponentStyle) SetBold(val bool) UIComponentStyleI {
	s.TextStyle.Bold = val
	return s
}

// SetItalic устанавливает курсив
func (s *UIComponentStyle) SetItalic(val bool) UIComponentStyleI {
	s.TextStyle.Italic = val
	return s
}

// SetUnderline устанавливает подчеркивание
func (s *UIComponentStyle) SetUnderline(val bool) UIComponentStyleI {
	s.TextStyle.Underline = val
	return s
}

// SetDim устанавливает тусклость текста
func (s *UIComponentStyle) SetDim(val bool) UIComponentStyleI {
	s.TextStyle.Dim = val
	return s
}

// SetBlink устанавливает мигание
func (s *UIComponentStyle) SetBlink(val bool) UIComponentStyleI {
	s.TextStyle.Blink = val
	return s
}

// SetReverse устанавливает инверсию цветов
func (s *UIComponentStyle) SetReverse(val bool) UIComponentStyleI {
	s.TextStyle.Reverse = val
	return s
}

// SetHidden устанавливает скрытие текста
func (s *UIComponentStyle) SetHidden(val bool) UIComponentStyleI {
	s.TextStyle.Hidden = val
	return s
}

// SetStrike устанавливает зачеркивание
func (s *UIComponentStyle) SetStrike(val bool) UIComponentStyleI {
	s.TextStyle.Strike = val
	return s
}

func (s *UIComponentStyle) Build() UIComponentStyle {
	return UIComponentStyle{
		Color:      s.Color,
		Background: s.Background,
		Flex:       s.Flex,
		Border:     s.Border,
		TextStyle:  s.TextStyle,
	}
}
