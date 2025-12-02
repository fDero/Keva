package main

import (
	"fmt"
	"os"

	"github.com/fDero/keva/core"
)

func main() {
	err := core.App.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
