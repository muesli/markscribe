package main

import (
	"flag"
	"fmt"
	"os"
	"io/ioutil"

	"github.com/muesli/markscribe/internal"
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: markscribe [template]")
		os.Exit(1)
	}

	tplIn, err := ioutil.ReadFile(flag.Args()[0])
	if err != nil {
		fmt.Println("Can't read file:", err)
		os.Exit(1)
	}

	err = internal.New(tplIn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}