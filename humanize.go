package main

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

func humanized(t interface{}) string {
	switch v := t.(type) {
	case time.Time:
		// flatten time to prevent updating README too often:
		v = time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location())

		if time.Since(v) <= time.Hour*24 {
			return "today"
		}

		return humanize.Time(v)
	default:
		return fmt.Sprintf("%v", t)
	}
}
