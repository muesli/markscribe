package main

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

func humanized(t interface{}) string {
	switch v := t.(type) {
	case time.Time:
		return humanize.Time(v)
	default:
		return fmt.Sprintf("%v", t)
	}
}
