package framework

import (
	_ "embed"
	"fmt"
)

//go:embed banner.txt
var banner string

func init() {
	fmt.Println(banner)
	fmt.Println()
}
