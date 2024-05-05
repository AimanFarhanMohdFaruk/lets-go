package ui

import "time"

func displayTime(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
