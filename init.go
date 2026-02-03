package framework

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed banner.txt
var banner string

func init() {
	fmt.Fprintln(os.Stderr, banner)
	fmt.Fprintln(os.Stderr)
}
