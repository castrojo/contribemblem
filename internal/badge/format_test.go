package badge

import "testing"

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{500, "500"},
		{1200, "1.2K"},
		{15000, "15.0K"},
		{1500000, "1.5M"},
		{2350000, "2.4M"},
	}

	for _, tt := range tests {
		result := FormatNumber(tt.input)
		if result != tt.expected {
			t.Errorf("FormatNumber(%d) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}
