package main

import (
	"fmt"
	"os"

	"github.com/jamespwilliams/etymology"
)

func main() {
	f, _ := os.Open(os.Args[1])
	ety, err := etymology.New(f)
	fmt.Println(err)

	fmt.Println(ety.Lookup(etymology.Word{Word: "muscular", Language: "eng"}))
}
