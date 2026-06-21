package menu

import (
	"fmt"
	"time"
)

// formatTimeAgo renders t as a short relative duration ("5m ago", "2d3h
// ago", ...), falling back to an absolute date once it's more than a week
// old, or "Never" if t is the zero value.
func formatTimeAgo(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	d := time.Since(t)
	if d < time.Second {
		return "just now"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds ago", int(d.Minutes()), int(d.Seconds())-60*int(d.Minutes()))
	}
	if d < time.Hour*24 {
		return fmt.Sprintf("%dh%dm ago", int(d.Hours()), int(d.Minutes())-60*int(d.Hours()))
	}
	days := int(d.Hours() / 24)
	if days < 7 {
		return fmt.Sprintf("%dd%dh ago", days, int(d.Hours())-24*days)
	}
	return t.Format("01-02 15:04")
}
