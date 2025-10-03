package components

import (
	"testing"
)

func TestNewUIComponentStyleBuilder(t *testing.T) {
	builder := NewUIComponentStyleBuilder()
	if builder == nil {
		t.Fatal("NewUIComponentStyleBuilder() returned nil")
	}
}

func TestUIComponentStyle_SetColor(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetColor("red").
		Build()

	if style.Color != "red" {
		t.Errorf("Expected Color to be 'red', got '%s'", style.Color)
	}
}

func TestUIComponentStyle_SetBackground(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetBackground("blue").
		Build()

	if style.Background != "blue" {
		t.Errorf("Expected Background to be 'blue', got '%s'", style.Background)
	}
}

func TestUIComponentStyle_SetFlex(t *testing.T) {
	flex := Flex{
		Direction: FlexDirectionRow,
		Justify:   FlexJustifyCenter,
		Align:     FlexAlignStart,
		Wrap:      true,
		Gap:       5,
	}

	style := NewUIComponentStyleBuilder().
		SetFlex(flex).
		Build()

	if style.Flex.Direction != FlexDirectionRow {
		t.Errorf("Expected Flex.Direction to be 'row', got '%s'", style.Flex.Direction)
	}
	if style.Flex.Justify != FlexJustifyCenter {
		t.Errorf("Expected Flex.Justify to be 'center', got '%s'", style.Flex.Justify)
	}
	if style.Flex.Align != FlexAlignStart {
		t.Errorf("Expected Flex.Align to be 'start', got '%s'", style.Flex.Align)
	}
	if !style.Flex.Wrap {
		t.Error("Expected Flex.Wrap to be true")
	}
	if style.Flex.Gap != 5 {
		t.Errorf("Expected Flex.Gap to be 5, got %d", style.Flex.Gap)
	}
}

func TestUIComponentStyle_SetBorder(t *testing.T) {
	border := Border{
		Style: BorderStyleRounded,
	}

	style := NewUIComponentStyleBuilder().
		SetBorder(border).
		Build()

	if style.Border.Style != BorderStyleRounded {
		t.Errorf("Expected Border.Style to be 'rounded', got '%s'", style.Border.Style)
	}
}

func TestUIComponentStyle_SetTextStyle(t *testing.T) {
	textStyle := TextStyle{
		Bold:      true,
		Italic:    true,
		Underline: false,
		Blink:     false,
	}

	style := NewUIComponentStyleBuilder().
		SetTextStyle(textStyle).
		Build()

	if !style.TextStyle.Bold {
		t.Error("Expected TextStyle.Bold to be true")
	}
	if !style.TextStyle.Italic {
		t.Error("Expected TextStyle.Italic to be true")
	}
	if style.TextStyle.Underline {
		t.Error("Expected TextStyle.Underline to be false")
	}
}

func TestUIComponentStyle_SetBold(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetBold(true).
		Build()

	if !style.TextStyle.Bold {
		t.Error("Expected TextStyle.Bold to be true")
	}
}

func TestUIComponentStyle_SetItalic(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetItalic(true).
		Build()

	if !style.TextStyle.Italic {
		t.Error("Expected TextStyle.Italic to be true")
	}
}

func TestUIComponentStyle_SetUnderline(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetUnderline(true).
		Build()

	if !style.TextStyle.Underline {
		t.Error("Expected TextStyle.Underline to be true")
	}
}

func TestUIComponentStyle_SetDim(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetDim(true).
		Build()

	if !style.TextStyle.Dim {
		t.Error("Expected TextStyle.Dim to be true")
	}
}

func TestUIComponentStyle_SetBlink(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetBlink(true).
		Build()

	if !style.TextStyle.Blink {
		t.Error("Expected TextStyle.Blink to be true")
	}
}

func TestUIComponentStyle_SetReverse(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetReverse(true).
		Build()

	if !style.TextStyle.Reverse {
		t.Error("Expected TextStyle.Reverse to be true")
	}
}

func TestUIComponentStyle_SetHidden(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetHidden(true).
		Build()

	if !style.TextStyle.Hidden {
		t.Error("Expected TextStyle.Hidden to be true")
	}
}

func TestUIComponentStyle_SetStrike(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetStrike(true).
		Build()

	if !style.TextStyle.Strike {
		t.Error("Expected TextStyle.Strike to be true")
	}
}

func TestUIComponentStyle_ChainedCalls(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetColor("#FF5733").
		SetBackground("#33FF57").
		SetBold(true).
		SetItalic(true).
		SetUnderline(true).
		Build()

	if style.Color != "#FF5733" {
		t.Errorf("Expected Color to be '#FF5733', got '%s'", style.Color)
	}
	if style.Background != "#33FF57" {
		t.Errorf("Expected Background to be '#33FF57', got '%s'", style.Background)
	}
	if !style.TextStyle.Bold {
		t.Error("Expected TextStyle.Bold to be true")
	}
	if !style.TextStyle.Italic {
		t.Error("Expected TextStyle.Italic to be true")
	}
	if !style.TextStyle.Underline {
		t.Error("Expected TextStyle.Underline to be true")
	}
}

func TestUIComponentStyle_AllTextStyles(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetBold(true).
		SetDim(true).
		SetItalic(true).
		SetUnderline(true).
		SetBlink(true).
		SetReverse(true).
		SetHidden(true).
		SetStrike(true).
		Build()

	if !style.TextStyle.Bold {
		t.Error("Expected TextStyle.Bold to be true")
	}
	if !style.TextStyle.Dim {
		t.Error("Expected TextStyle.Dim to be true")
	}
	if !style.TextStyle.Italic {
		t.Error("Expected TextStyle.Italic to be true")
	}
	if !style.TextStyle.Underline {
		t.Error("Expected TextStyle.Underline to be true")
	}
	if !style.TextStyle.Blink {
		t.Error("Expected TextStyle.Blink to be true")
	}
	if !style.TextStyle.Reverse {
		t.Error("Expected TextStyle.Reverse to be true")
	}
	if !style.TextStyle.Hidden {
		t.Error("Expected TextStyle.Hidden to be true")
	}
	if !style.TextStyle.Strike {
		t.Error("Expected TextStyle.Strike to be true")
	}
}

func TestUIComponentStyle_ComplexStyle(t *testing.T) {
	flex := Flex{
		Direction: FlexDirectionColumn,
		Justify:   FlexJustifySpaceBetween,
		Align:     FlexAlignCenter,
		Wrap:      false,
		Gap:       20,
	}

	border := Border{
		Style: BorderStyleDouble,
	}

	style := NewUIComponentStyleBuilder().
		SetColor("white").
		SetBackground("black").
		SetFlex(flex).
		SetBorder(border).
		SetBold(true).
		SetUnderline(true).
		Build()

	if style.Color != "white" {
		t.Errorf("Expected Color to be 'white', got '%s'", style.Color)
	}
	if style.Background != "black" {
		t.Errorf("Expected Background to be 'black', got '%s'", style.Background)
	}
	if style.Flex.Direction != FlexDirectionColumn {
		t.Errorf("Expected Flex.Direction to be 'column', got '%s'", style.Flex.Direction)
	}
	if style.Flex.Justify != FlexJustifySpaceBetween {
		t.Errorf("Expected Flex.Justify to be 'space-between', got '%s'", style.Flex.Justify)
	}
	if style.Flex.Align != FlexAlignCenter {
		t.Errorf("Expected Flex.Align to be 'center', got '%s'", style.Flex.Align)
	}
	if style.Flex.Wrap {
		t.Error("Expected Flex.Wrap to be false")
	}
	if style.Flex.Gap != 20 {
		t.Errorf("Expected Flex.Gap to be 20, got %d", style.Flex.Gap)
	}
	if style.Border.Style != BorderStyleDouble {
		t.Errorf("Expected Border.Style to be 'double', got '%s'", style.Border.Style)
	}
	if !style.TextStyle.Bold {
		t.Error("Expected TextStyle.Bold to be true")
	}
	if !style.TextStyle.Underline {
		t.Error("Expected TextStyle.Underline to be true")
	}
}

func TestUIComponentStyle_EmptyBuild(t *testing.T) {
	style := NewUIComponentStyleBuilder().Build()

	if style.Color != "" {
		t.Errorf("Expected empty Color, got '%s'", style.Color)
	}
	if style.Background != "" {
		t.Errorf("Expected empty Background, got '%s'", style.Background)
	}
	if style.TextStyle.Bold {
		t.Error("Expected TextStyle.Bold to be false")
	}
}

func TestFlexDirection_Constants(t *testing.T) {
	if FlexDirectionRow != "row" {
		t.Errorf("Expected FlexDirectionRow to be 'row', got '%s'", FlexDirectionRow)
	}
	if FlexDirectionColumn != "column" {
		t.Errorf("Expected FlexDirectionColumn to be 'column', got '%s'", FlexDirectionColumn)
	}
}

func TestFlexJustify_Constants(t *testing.T) {
	if FlexJustifyStart != "start" {
		t.Errorf("Expected FlexJustifyStart to be 'start', got '%s'", FlexJustifyStart)
	}
	if FlexJustifyCenter != "center" {
		t.Errorf("Expected FlexJustifyCenter to be 'center', got '%s'", FlexJustifyCenter)
	}
	if FlexJustifyEnd != "end" {
		t.Errorf("Expected FlexJustifyEnd to be 'end', got '%s'", FlexJustifyEnd)
	}
	if FlexJustifySpaceBetween != "space-between" {
		t.Errorf("Expected FlexJustifySpaceBetween to be 'space-between', got '%s'", FlexJustifySpaceBetween)
	}
}

func TestFlexAlign_Constants(t *testing.T) {
	if FlexAlignStart != "start" {
		t.Errorf("Expected FlexAlignStart to be 'start', got '%s'", FlexAlignStart)
	}
	if FlexAlignCenter != "center" {
		t.Errorf("Expected FlexAlignCenter to be 'center', got '%s'", FlexAlignCenter)
	}
	if FlexAlignEnd != "end" {
		t.Errorf("Expected FlexAlignEnd to be 'end', got '%s'", FlexAlignEnd)
	}
	if FlexAlignStretch != "stretch" {
		t.Errorf("Expected FlexAlignStretch to be 'stretch', got '%s'", FlexAlignStretch)
	}
}

func TestBorderStyle_Constants(t *testing.T) {
	if BorderStyleNone != "none" {
		t.Errorf("Expected BorderStyleNone to be 'none', got '%s'", BorderStyleNone)
	}
	if BorderStyleSolid != "solid" {
		t.Errorf("Expected BorderStyleSolid to be 'solid', got '%s'", BorderStyleSolid)
	}
	if BorderStyleDouble != "double" {
		t.Errorf("Expected BorderStyleDouble to be 'double', got '%s'", BorderStyleDouble)
	}
	if BorderStyleRounded != "rounded" {
		t.Errorf("Expected BorderStyleRounded to be 'rounded', got '%s'", BorderStyleRounded)
	}
}

func TestUIComponentStyle_OverwriteValues(t *testing.T) {
	style := NewUIComponentStyleBuilder().
		SetColor("red").
		SetColor("blue").
		SetBold(true).
		SetBold(false).
		Build()

	if style.Color != "blue" {
		t.Errorf("Expected Color to be 'blue' (overwritten), got '%s'", style.Color)
	}
	if style.TextStyle.Bold {
		t.Error("Expected TextStyle.Bold to be false (overwritten)")
	}
}
