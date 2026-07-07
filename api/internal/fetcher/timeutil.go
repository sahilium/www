package fetcher

import (
	"fmt"
	"time"
)

func timeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		m := int(diff.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	case diff < 24*time.Hour:
		h := int(diff.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	case diff < 30*24*time.Hour:
		d := int(diff.Hours() / 24)
		if d == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", d)
	case diff < 365*24*time.Hour:
		m := int(diff.Hours() / (24 * 30))
		if m == 1 {
			return "about 1 month ago"
		}
		return fmt.Sprintf("%d months ago", m)
	default:
		y := int(diff.Hours() / (24 * 365))
		if y == 1 {
			return "about 1 year ago"
		}
		return fmt.Sprintf("%d years ago", y)
	}
}
