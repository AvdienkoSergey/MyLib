package components

type UIComponent struct {
	Tag    string
	Text   string
	Width  uint16
	Height uint16
	Style  UIComponentStyle
}

type UIComponentBuilderI interface {
	SetText(val string) UIComponentBuilderI
	SetWidth(val uint16) UIComponentBuilderI
	SetHeight(val uint16) UIComponentBuilderI
	SetStyle(style UIComponentStyle) UIComponentBuilderI

	Build() UIComponent
}

type UIComponentBuilder struct {
	tag    string
	text   string
	width  uint16
	height uint16
	style  UIComponentStyle
}

func NewUIComponentBuilder() UIComponentBuilderI {
	return &UIComponentBuilder{}
}

func (c *UIComponentBuilder) SetTag(val string) UIComponentBuilderI {
	c.tag = val
	return c
}

func (c *UIComponentBuilder) SetText(val string) UIComponentBuilderI {
	c.text = val
	return c
}

func (c *UIComponentBuilder) SetWidth(val uint16) UIComponentBuilderI {
	c.width = val
	return c
}

func (c *UIComponentBuilder) SetHeight(val uint16) UIComponentBuilderI {
	c.height = val
	return c
}

func (c *UIComponentBuilder) SetStyle(style UIComponentStyle) UIComponentBuilderI {
	c.style = style
	return c
}

func (c *UIComponentBuilder) Build() UIComponent {
	return UIComponent{
		Tag:    c.tag,
		Text:   c.text,
		Width:  c.width,
		Height: c.height,
		Style:  c.style,
	}
}
