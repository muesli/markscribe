package main

import (
	"html/template"
)

func comment(s string) template.HTML {
	return template.HTML(s)
}
