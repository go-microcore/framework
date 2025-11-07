package framework

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"

	"github.com/aquasecurity/table"
)

//go:embed banner.txt
var banner string

func init() {
	fmt.Println(banner)
	fmt.Println()
	table := table.New(os.Stdout)
	table.AddRow("PID", strconv.Itoa(os.Getpid()))
	table.AddRow("PPID", strconv.Itoa(os.Getppid()))
	table.Render()
	fmt.Println()
}
