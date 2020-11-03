package main

import (
	"fmt"
	"reflect"
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

func reverse(s interface{}) interface{} {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}

	return s
}
