package badge

import "fmt"

// FormatNumber formats large numbers with K/M suffixes
// 1200 → "1.2K", 1500000 → "1.5M"
func FormatNumber(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}
