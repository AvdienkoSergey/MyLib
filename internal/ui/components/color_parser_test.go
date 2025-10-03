package components

import (
	"testing"
)

func TestNewColorParser(t *testing.T) {
	parser := NewColorParser()
	if parser == nil {
		t.Fatal("NewColorParser() returned nil")
	}
}

func TestColorParser_ParseToANSI_EmptyString(t *testing.T) {
	parser := NewColorParser()
	result, err := parser.ParseToANSI("", false)
	if err != nil {
		t.Errorf("Expected no error for empty color, got %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string for empty color, got '%s'", result)
	}
}

func TestColorParser_ParseToANSI_RGB(t *testing.T) {
	parser := NewColorParser()

	tests := []struct {
		color        string
		isBackground bool
		expected     string
	}{
		{"rgb(255,0,0)", false, "38;2;255;0;0"},
		{"rgb(255,0,0)", true, "48;2;255;0;0"},
		{"rgb(0,255,0)", false, "38;2;0;255;0"},
		{"rgb(0,0,255)", true, "48;2;0;0;255"},
		{"rgb(128,128,128)", false, "38;2;128;128;128"},
		{"rgb(255, 255, 255)", false, "38;2;255;255;255"},
	}

	for _, tt := range tests {
		result, err := parser.ParseToANSI(tt.color, tt.isBackground)
		if err != nil {
			t.Errorf("ParseToANSI(%s, %v) returned error: %v", tt.color, tt.isBackground, err)
		}
		if result != tt.expected {
			t.Errorf("ParseToANSI(%s, %v) = %s; expected %s", tt.color, tt.isBackground, result, tt.expected)
		}
	}
}

func TestColorParser_ParseToANSI_Hex(t *testing.T) {
	parser := NewColorParser()

	tests := []struct {
		color        string
		isBackground bool
		expected     string
	}{
		{"#FF0000", false, "38;2;255;0;0"},
		{"#FF0000", true, "48;2;255;0;0"},
		{"#00FF00", false, "38;2;0;255;0"},
		{"#0000FF", true, "48;2;0;0;255"},
		{"#808080", false, "38;2;128;128;128"},
		{"#FFFFFF", false, "38;2;255;255;255"},
		{"#000000", true, "48;2;0;0;0"},
	}

	for _, tt := range tests {
		result, err := parser.ParseToANSI(tt.color, tt.isBackground)
		if err != nil {
			t.Errorf("ParseToANSI(%s, %v) returned error: %v", tt.color, tt.isBackground, err)
		}
		if result != tt.expected {
			t.Errorf("ParseToANSI(%s, %v) = %s; expected %s", tt.color, tt.isBackground, result, tt.expected)
		}
	}
}

func TestColorParser_ParseToANSI_256Color(t *testing.T) {
	parser := NewColorParser()

	tests := []struct {
		color        string
		isBackground bool
		expected     string
	}{
		{"0", false, "38;5;0"},
		{"0", true, "48;5;0"},
		{"255", false, "38;5;255"},
		{"255", true, "48;5;255"},
		{"128", false, "38;5;128"},
		{"42", true, "48;5;42"},
	}

	for _, tt := range tests {
		result, err := parser.ParseToANSI(tt.color, tt.isBackground)
		if err != nil {
			t.Errorf("ParseToANSI(%s, %v) returned error: %v", tt.color, tt.isBackground, err)
		}
		if result != tt.expected {
			t.Errorf("ParseToANSI(%s, %v) = %s; expected %s", tt.color, tt.isBackground, result, tt.expected)
		}
	}
}

func TestColorParser_ParseToANSI_BasicColors(t *testing.T) {
	parser := NewColorParser()

	tests := []struct {
		color        string
		isBackground bool
		expected     string
	}{
		{"black", false, "30"},
		{"black", true, "40"},
		{"red", false, "31"},
		{"red", true, "41"},
		{"green", false, "32"},
		{"green", true, "42"},
		{"yellow", false, "33"},
		{"yellow", true, "43"},
		{"blue", false, "34"},
		{"blue", true, "44"},
		{"magenta", false, "35"},
		{"magenta", true, "45"},
		{"cyan", false, "36"},
		{"cyan", true, "46"},
		{"white", false, "37"},
		{"white", true, "47"},
		{"bright_black", false, "90"},
		{"bright_black", true, "100"},
		{"bright_red", false, "91"},
		{"bright_red", true, "101"},
		{"bright_green", false, "92"},
		{"bright_green", true, "102"},
		{"bright_yellow", false, "93"},
		{"bright_yellow", true, "103"},
		{"bright_blue", false, "94"},
		{"bright_blue", true, "104"},
		{"bright_magenta", false, "95"},
		{"bright_magenta", true, "105"},
		{"bright_cyan", false, "96"},
		{"bright_cyan", true, "106"},
		{"bright_white", false, "97"},
		{"bright_white", true, "107"},
	}

	for _, tt := range tests {
		result, err := parser.ParseToANSI(tt.color, tt.isBackground)
		if err != nil {
			t.Errorf("ParseToANSI(%s, %v) returned error: %v", tt.color, tt.isBackground, err)
		}
		if result != tt.expected {
			t.Errorf("ParseToANSI(%s, %v) = %s; expected %s", tt.color, tt.isBackground, result, tt.expected)
		}
	}
}

func TestColorParser_ParseToANSI_UnknownColor(t *testing.T) {
	parser := NewColorParser()

	tests := []string{
		"unknown",
		"invalid",
		"#ZZZ",
		"rgb()",
		"256",
		"300",
	}

	for _, color := range tests {
		result, err := parser.ParseToANSI(color, false)
		if err == nil {
			t.Errorf("ParseToANSI(%s, false) expected error, got result: %s", color, result)
		}
	}
}

func TestColorParser_parseRGB(t *testing.T) {
	parser := NewColorParser()

	tests := []struct {
		input    string
		r, g, b  int
	}{
		{"rgb(255,0,0)", 255, 0, 0},
		{"rgb(0,255,0)", 0, 255, 0},
		{"rgb(0,0,255)", 0, 0, 255},
		{"rgb(128,128,128)", 128, 128, 128},
		{"rgb(255, 255, 255)", 255, 255, 255},
		{"rgb( 0 , 0 , 0 )", 0, 0, 0},
	}

	for _, tt := range tests {
		r, g, b, err := parser.parseRGB(tt.input)
		if err != nil {
			t.Errorf("parseRGB(%s) returned error: %v", tt.input, err)
		}
		if r != tt.r || g != tt.g || b != tt.b {
			t.Errorf("parseRGB(%s) = (%d, %d, %d); expected (%d, %d, %d)", tt.input, r, g, b, tt.r, tt.g, tt.b)
		}
	}
}

func TestColorParser_parseHex(t *testing.T) {
	parser := NewColorParser()

	tests := []struct {
		input    string
		r, g, b  int
	}{
		{"#FF0000", 255, 0, 0},
		{"#00FF00", 0, 255, 0},
		{"#0000FF", 0, 0, 255},
		{"#808080", 128, 128, 128},
		{"#FFFFFF", 255, 255, 255},
		{"#000000", 0, 0, 0},
	}

	for _, tt := range tests {
		r, g, b, err := parser.parseHex(tt.input)
		if err != nil {
			t.Errorf("parseHex(%s) returned error: %v", tt.input, err)
		}
		if r != tt.r || g != tt.g || b != tt.b {
			t.Errorf("parseHex(%s) = (%d, %d, %d); expected (%d, %d, %d)", tt.input, r, g, b, tt.r, tt.g, tt.b)
		}
	}
}

func TestColorParser_rgbToANSI(t *testing.T) {
	parser := NewColorParser()

	tests := []struct {
		r, g, b      int
		isBackground bool
		expected     string
	}{
		{255, 0, 0, false, "38;2;255;0;0"},
		{255, 0, 0, true, "48;2;255;0;0"},
		{0, 255, 0, false, "38;2;0;255;0"},
		{0, 0, 255, true, "48;2;0;0;255"},
		{128, 128, 128, false, "38;2;128;128;128"},
	}

	for _, tt := range tests {
		result := parser.rgbToANSI(tt.r, tt.g, tt.b, tt.isBackground)
		if result != tt.expected {
			t.Errorf("rgbToANSI(%d, %d, %d, %v) = %s; expected %s", tt.r, tt.g, tt.b, tt.isBackground, result, tt.expected)
		}
	}
}

func TestColorParser_color256ToANSI(t *testing.T) {
	parser := NewColorParser()

	tests := []struct {
		color        string
		isBackground bool
		expected     string
	}{
		{"0", false, "38;5;0"},
		{"0", true, "48;5;0"},
		{"255", false, "38;5;255"},
		{"128", true, "48;5;128"},
	}

	for _, tt := range tests {
		result := parser.color256ToANSI(tt.color, tt.isBackground)
		if result != tt.expected {
			t.Errorf("color256ToANSI(%s, %v) = %s; expected %s", tt.color, tt.isBackground, result, tt.expected)
		}
	}
}

func TestColorParser_basicColorToANSI(t *testing.T) {
	parser := NewColorParser()

	tests := []struct {
		color        string
		isBackground bool
		expected     string
	}{
		{"red", false, "31"},
		{"red", true, "41"},
		{"blue", false, "34"},
		{"white", true, "47"},
		{"bright_red", false, "91"},
		{"bright_cyan", true, "106"},
	}

	for _, tt := range tests {
		result := parser.basicColorToANSI(tt.color, tt.isBackground)
		if result != tt.expected {
			t.Errorf("basicColorToANSI(%s, %v) = %s; expected %s", tt.color, tt.isBackground, result, tt.expected)
		}
	}
}

func TestColorParser_basicColorToANSI_Unknown(t *testing.T) {
	parser := NewColorParser()

	result := parser.basicColorToANSI("unknown", false)
	if result != "" {
		t.Errorf("basicColorToANSI('unknown', false) = %s; expected empty string", result)
	}

	result = parser.basicColorToANSI("invalid", true)
	if result != "" {
		t.Errorf("basicColorToANSI('invalid', true) = %s; expected empty string", result)
	}
}

func TestColorParser_parseRGB_Invalid(t *testing.T) {
	parser := NewColorParser()

	tests := []string{
		"rgb()",
		"rgb(255)",
		"rgb(255,0)",
		"rgb(a,b,c)",
	}

	for _, input := range tests {
		r, g, b, err := parser.parseRGB(input)
		if err == nil {
			t.Errorf("parseRGB(%s) succeeded with (%d, %d, %d); expected error", input, r, g, b)
		}
	}
}

func TestColorParser_parseHex_Invalid(t *testing.T) {
	parser := NewColorParser()

	tests := []string{
		"#",
		"#FFF",
		"#GGGGGG",
		"#12345",
	}

	for _, input := range tests {
		r, g, b, err := parser.parseHex(input)
		if err == nil {
			t.Errorf("parseHex(%s) succeeded with (%d, %d, %d); expected error", input, r, g, b)
		}
	}
}
