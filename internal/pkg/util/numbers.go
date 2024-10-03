package util

import "fmt"

func FormatCount(num int) string {
	if num >= 1_000_000_000 {
		return fmt.Sprintf("%.2fB", float64(num)/1_000_000_000)
	} else if num >= 100_000_000 {
		return fmt.Sprintf("%.0fM", float64(num)/1_000_000)
	} else if num >= 10_000_000 {
		return fmt.Sprintf("%.1fM", float64(num)/1_000_000)
	} else if num >= 1_000_000 {
		return fmt.Sprintf("%.2fM", float64(num)/1_000_000)
	} else if num >= 100_000 {
		return fmt.Sprintf("%.0fk", float64(num)/1_000)
	} else if num >= 10_000 {
		return fmt.Sprintf("%.1fk", float64(num)/1_000)
	} else if num >= 1_000 {
		return fmt.Sprintf("%.2fk", float64(num)/1_000)
	} else {
		return fmt.Sprintf("%d", num)
	}
}
