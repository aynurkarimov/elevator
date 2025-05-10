package main

import (
	"github.com/aynurkarimov/elevator/ui"
)

func main() {
	for {
		err := ui.Handler()

		if err != nil {
			panic(err)
		}
	}
}
