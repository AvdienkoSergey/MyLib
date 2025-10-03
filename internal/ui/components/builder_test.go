package components

import (
	"testing"
)

func TestNewUIComponentBuilder(t *testing.T) {
	builder := NewUIComponentBuilder()
	if builder == nil {
		t.Fatal("NewUIComponentBuilder() returned nil")
	}
}

func TestUIComponentBuilder_SetText(t *testing.T) {
	builder := NewUIComponentBuilder()
	result := builder.SetText("test text")

	if result == nil {
		t.Fatal("SetText() returned nil")
	}

	component := result.Build()
	if component.Text != "test text" {
		t.Errorf("Expected Text to be 'test text', got '%s'", component.Text)
	}
}

func TestUIComponentBuilder_SetWidth(t *testing.T) {
	builder := NewUIComponentBuilder()
	result := builder.SetWidth(100)

	if result == nil {
		t.Fatal("SetWidth() returned nil")
	}

	component := result.Build()
	if component.Width != 100 {
		t.Errorf("Expected Width to be 100, got %d", component.Width)
	}
}

func TestUIComponentBuilder_SetHeight(t *testing.T) {
	builder := NewUIComponentBuilder()
	result := builder.SetHeight(50)

	if result == nil {
		t.Fatal("SetHeight() returned nil")
	}

	component := result.Build()
	if component.Height != 50 {
		t.Errorf("Expected Height to be 50, got %d", component.Height)
	}
}

func TestUIComponentBuilder_SetStyle(t *testing.T) {
	builder := NewUIComponentBuilder()
	style := UIComponentStyle{
		Color:      "red",
		Background: "blue",
	}

	result := builder.SetStyle(style)

	if result == nil {
		t.Fatal("SetStyle() returned nil")
	}

	component := result.Build()
	if component.Style.Color != "red" {
		t.Errorf("Expected Style.Color to be 'red', got '%s'", component.Style.Color)
	}
	if component.Style.Background != "blue" {
		t.Errorf("Expected Style.Background to be 'blue', got '%s'", component.Style.Background)
	}
}

func TestUIComponentBuilder_ChainedCalls(t *testing.T) {
	style := UIComponentStyle{
		Color:      "green",
		Background: "yellow",
		TextStyle: TextStyle{
			Bold:      true,
			Underline: true,
		},
	}

	component := NewUIComponentBuilder().
		SetText("chained text").
		SetWidth(200).
		SetHeight(100).
		SetStyle(style).
		Build()

	if component.Text != "chained text" {
		t.Errorf("Expected Text to be 'chained text', got '%s'", component.Text)
	}
	if component.Width != 200 {
		t.Errorf("Expected Width to be 200, got %d", component.Width)
	}
	if component.Height != 100 {
		t.Errorf("Expected Height to be 100, got %d", component.Height)
	}
	if component.Style.Color != "green" {
		t.Errorf("Expected Style.Color to be 'green', got '%s'", component.Style.Color)
	}
	if component.Style.Background != "yellow" {
		t.Errorf("Expected Style.Background to be 'yellow', got '%s'", component.Style.Background)
	}
	if !component.Style.TextStyle.Bold {
		t.Error("Expected TextStyle.Bold to be true")
	}
	if !component.Style.TextStyle.Underline {
		t.Error("Expected TextStyle.Underline to be true")
	}
}

func TestUIComponentBuilder_EmptyBuild(t *testing.T) {
	component := NewUIComponentBuilder().Build()

	if component.Text != "" {
		t.Errorf("Expected empty Text, got '%s'", component.Text)
	}
	if component.Width != 0 {
		t.Errorf("Expected Width to be 0, got %d", component.Width)
	}
	if component.Height != 0 {
		t.Errorf("Expected Height to be 0, got %d", component.Height)
	}
}

func TestUIComponentBuilder_PartialBuild(t *testing.T) {
	component := NewUIComponentBuilder().
		SetText("partial").
		SetWidth(50).
		Build()

	if component.Text != "partial" {
		t.Errorf("Expected Text to be 'partial', got '%s'", component.Text)
	}
	if component.Width != 50 {
		t.Errorf("Expected Width to be 50, got %d", component.Width)
	}
	if component.Height != 0 {
		t.Errorf("Expected Height to be 0, got %d", component.Height)
	}
}

func TestUIComponentBuilder_MultipleBuilds(t *testing.T) {
	builder := NewUIComponentBuilder().
		SetText("test").
		SetWidth(100)

	component1 := builder.Build()
	component2 := builder.Build()

	if component1.Text != component2.Text {
		t.Error("Multiple Build() calls should produce identical results")
	}
	if component1.Width != component2.Width {
		t.Error("Multiple Build() calls should produce identical results")
	}
}

func TestUIComponentBuilder_ComplexStyle(t *testing.T) {
	style := UIComponentStyle{
		Color:      "#FF5733",
		Background: "#33FF57",
		Flex: Flex{
			Direction: FlexDirectionRow,
			Justify:   FlexJustifyCenter,
			Align:     FlexAlignStretch,
			Wrap:      true,
			Gap:       10,
		},
		Border: Border{
			Style: BorderStyleRounded,
		},
		TextStyle: TextStyle{
			Bold:      true,
			Italic:    true,
			Underline: true,
			Blink:     false,
			Reverse:   false,
			Hidden:    false,
			Strike:    false,
			Dim:       false,
		},
	}

	component := NewUIComponentBuilder().
		SetText("complex component").
		SetWidth(300).
		SetHeight(150).
		SetStyle(style).
		Build()

	if component.Style.Flex.Direction != FlexDirectionRow {
		t.Errorf("Expected Flex.Direction to be 'row', got '%s'", component.Style.Flex.Direction)
	}
	if component.Style.Flex.Justify != FlexJustifyCenter {
		t.Errorf("Expected Flex.Justify to be 'center', got '%s'", component.Style.Flex.Justify)
	}
	if component.Style.Flex.Gap != 10 {
		t.Errorf("Expected Flex.Gap to be 10, got %d", component.Style.Flex.Gap)
	}
	if component.Style.Border.Style != BorderStyleRounded {
		t.Errorf("Expected Border.Style to be 'rounded', got '%s'", component.Style.Border.Style)
	}
}
